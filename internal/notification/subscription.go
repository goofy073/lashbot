package notification

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"log/slog"
	"remnawave-tg-shop-bot/internal/config"
	"remnawave-tg-shop-bot/internal/database"
	"remnawave-tg-shop-bot/internal/handler"
	"remnawave-tg-shop-bot/internal/translation"
	"time"
)

type customerRepository interface {
	FindByExpirationRange(ctx context.Context, startDate, endDate time.Time) (*[]database.Customer, error)
	FindExpired(ctx context.Context, now time.Time) (*[]database.Customer, error)
	UpdateFields(ctx context.Context, id int64, updates map[string]interface{}) error
}

type tributeRepository interface {
	FindLatestActiveTributesByCustomerIDs(ctx context.Context, customerIDs []int64) (*[]database.Purchase, error)
}

type paymentProcessor interface {
	CreatePurchase(ctx context.Context, amount float64, months int, customer *database.Customer, invoiceType database.InvoiceType) (string, int64, error)
	ProcessPurchaseById(ctx context.Context, purchaseId int64) error
}

type SubscriptionService struct {
	customerRepository customerRepository
	purchaseRepository tributeRepository
	paymentService     paymentProcessor
	telegramBot        *bot.Bot
	tm                 *translation.Manager
	notify             func(context.Context, database.Customer) error
}

func NewSubscriptionService(customerRepository customerRepository,
	purchaseRepository tributeRepository,
	paymentService paymentProcessor,
	telegramBot *bot.Bot,
	tm *translation.Manager) *SubscriptionService {
	svc := &SubscriptionService{customerRepository: customerRepository, purchaseRepository: purchaseRepository, paymentService: paymentService, telegramBot: telegramBot, tm: tm}
	svc.notify = svc.sendNotification
	return svc
}

func (s *SubscriptionService) ProcessSubscriptionExpiration() error {
	ctx := context.Background()

	// 1. First, process users whose subscriptions have already expired
	if err := s.processExpiredSubscriptions(ctx); err != nil {
		slog.Error("Failed to process expired subscriptions", "error", err)
	}

	// 2. Then, process users whose subscriptions will expire soon
	customers, err := s.getCustomersWithExpiringSubscriptions()
	if err != nil {
		slog.Error("Failed to get customers with expiring subscriptions", "error", err)
		return err
	}

	slog.Info(fmt.Sprintf("Found %d customers with expiring subscriptions", len(*customers)))
	if len(*customers) == 0 {
		return nil
	}
	now := time.Now()

	customersIds := make([]int64, len(*customers))
	for i, customer := range *customers {
		customersIds[i] = customer.ID
	}

	latestActiveTributes, err := s.purchaseRepository.FindLatestActiveTributesByCustomerIDs(ctx, customersIds)
	if err != nil {
		slog.Error("Failed to query tribute purchases", "error", err)
		return err
	}

	customerIdTributes := make(map[int64]*database.Purchase, len(*latestActiveTributes))
	for i := range *latestActiveTributes {
		p := &(*latestActiveTributes)[i]
		customerIdTributes[p.CustomerID] = p
	}

	tributesProcessed := make(map[int64]bool, len(*latestActiveTributes))

	for _, customer := range *customers {
		daysUntilExpiration := s.getDaysUntilExpiration(now, *customer.ExpireAt)

		if p, ok := customerIdTributes[customer.ID]; ok {
			if daysUntilExpiration != 1 {
				continue
			}
			_, purchaseId, err := s.paymentService.CreatePurchase(ctx, p.Amount, p.Month, &customer, database.InvoiceTypeTribute)
			if err != nil {
				slog.Error("Failed to create tribute purchase", "error", err)
				continue
			}

			err = s.paymentService.ProcessPurchaseById(ctx, purchaseId)
			if err != nil {
				slog.Error("Failed to process tribute purchase", "error", err)
				continue
			}
			slog.Info("Tribute purchase processed successfully", "purchase_id", purchaseId)
			tributesProcessed[customer.ID] = true
		}
		if _, ok := tributesProcessed[customer.ID]; ok {
			continue
		}

		send := s.notify
		if send == nil {
			send = s.sendNotification
		}

		err := send(ctx, customer)
		if err != nil {
			slog.Error("Failed to send notification",
				"customer_id", customer.ID,
				"days_until_expiration", daysUntilExpiration,
				"error", err)
			continue
		}

		slog.Info("Notification sent successfully",
			"customer_id", customer.ID,
			"days_until_expiration", daysUntilExpiration)
	}

	slog.Info(fmt.Sprintf("Processed tributes customers %d with expiring subscriptions", len(tributesProcessed)))
	slog.Info(fmt.Sprintf("Sent notifications to %d customers with expiring subscriptions", len(*customers)-len(tributesProcessed)))
	return nil
}

func (s *SubscriptionService) processExpiredSubscriptions(ctx context.Context) error {
	now := time.Now()
	expiredCustomers, err := s.customerRepository.FindExpired(ctx, now)
	if err != nil {
		return err
	}

	slog.Info(fmt.Sprintf("Found %d customers with expired subscriptions that need to be kicked", len(*expiredCustomers)))
	for _, customer := range *expiredCustomers {
		// Kick user from channel
		_, err := s.telegramBot.BanChatMember(ctx, &bot.BanChatMemberParams{
			ChatID: config.GetPrivateChannelID(),
			UserID: customer.TelegramID,
		})
		if err != nil {
			slog.Error("Failed to ban chat member", "customer_id", customer.ID, "telegram_id", customer.TelegramID, "error", err)
		}

		// Unban them immediately so they can rejoin later if they purchase again
		_, err = s.telegramBot.UnbanChatMember(ctx, &bot.UnbanChatMemberParams{
			ChatID:       config.GetPrivateChannelID(),
			UserID:       customer.TelegramID,
			OnlyIfBanned: true,
		})
		if err != nil {
			slog.Error("Failed to unban chat member", "customer_id", customer.ID, "telegram_id", customer.TelegramID, "error", err)
		}

		// Clear subscription_link in db to mark them as kicked
		err = s.customerRepository.UpdateFields(ctx, customer.ID, map[string]interface{}{
			"subscription_link": nil,
		})
		if err != nil {
			slog.Error("Failed to clear subscription_link", "customer_id", customer.ID, "error", err)
		}

		// Notify the user that their subscription has expired
		_, err = s.telegramBot.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    customer.TelegramID,
			Text:      s.tm.GetText(customer.Language, "subscription_expired"),
			ParseMode: models.ParseModeHTML,
			ReplyMarkup: models.InlineKeyboardMarkup{
				InlineKeyboard: [][]models.InlineKeyboardButton{
					{s.tm.GetButton(customer.Language, "renew_subscription_button").InlineCallback(handler.CallbackBuy)},
				},
			},
		})
		if err != nil {
			slog.Error("Failed to send expiration notification", "customer_id", customer.ID, "error", err)
		}

		slog.Info("Successfully kicked expired customer", "customer_id", customer.ID, "telegram_id", customer.TelegramID)
	}
	return nil
}

func (s *SubscriptionService) getCustomersWithExpiringSubscriptions() (*[]database.Customer, error) {
	now := time.Now()
	endDate := now.AddDate(0, 0, 3)

	dbCustomers, err := s.customerRepository.FindByExpirationRange(context.Background(), now, endDate)
	if err != nil {
		return nil, err
	}

	return dbCustomers, nil
}

func (s *SubscriptionService) getDaysUntilExpiration(now time.Time, expireAt time.Time) int {
	nowDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	expireDate := time.Date(expireAt.Year(), expireAt.Month(), expireAt.Day(), 0, 0, 0, 0, expireAt.Location())

	duration := expireDate.Sub(nowDate)
	return int(duration.Hours() / 24)
}

func (s *SubscriptionService) sendNotification(ctx context.Context, customer database.Customer) error {
	expireDate := customer.ExpireAt.Format("02.01.2006")

	messageText := fmt.Sprintf(
		s.tm.GetText(customer.Language, "subscription_expiring"),
		expireDate,
	)

	_, err := s.telegramBot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    customer.TelegramID,
		Text:      messageText,
		ParseMode: models.ParseModeHTML,
		ReplyMarkup: models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{s.tm.GetButton(customer.Language, "renew_subscription_button").InlineCallback(handler.CallbackBuy)},
			},
		},
	})

	return err
}
