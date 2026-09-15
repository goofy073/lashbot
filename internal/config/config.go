package config

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type config struct {
	menuPhotoPath                                             string
	telegramToken                                             string
	price1, price3, price6, price12                           int
	starsPrice1, starsPrice3, starsPrice6, starsPrice12       int
	defaultLanguage                                           string
	databaseURL                                               string
	cryptoPayURL, cryptoPayToken                              string
	botURL                                                    string
	yookasaURL, yookasaShopId, yookasaSecretKey, yookasaEmail string
	moynalogURL, moynalogUsername, moynalogPassword           string
	feedbackURL                                               string
	channelURL                                                string
	serverStatusURL                                           string
	supportURL                                                string
	tosURL                                                    string
	isYookasaEnabled                                          bool
	isCryptoEnabled                                           bool
	isTelegramStarsEnabled                                    bool
	isMoynalogEnabled                                         bool
	adminTelegramId                                           int64
	privateChannelID                                          int64
	miniApp                                                   string
	enableAutoPayment                                         bool
	healthCheckPort                                           int
	tributeWebhookUrl, tributeAPIKey, tributePaymentUrl       string
	yookasaWebhookUrl                                         string
	plategaMerchantId, plategaSecret, plategaWebhookUrl       string
	isPlategaSBPEnabled, isPlategaCardsEnabled                bool
	isPlategaAcquiringEnabled, isPlategaWorldwideEnabled      bool
	isPlategaCryptoEnabled                                    bool
	isWebAppLinkEnabled                                       bool
	daysInMonth                                               int
	blockedTelegramIds                                        map[int64]bool
	whitelistedTelegramIds                                    map[int64]bool
	requirePaidPurchaseForStars                               bool
}

var conf config

func MenuPhotoPath() string {
	return conf.menuPhotoPath
}

func DefaultLanguage() string {
	return conf.defaultLanguage
}

func GetTributeWebHookUrl() string {
	return conf.tributeWebhookUrl
}

func GetTributeAPIKey() string {
	return conf.tributeAPIKey
}

func GetTributePaymentUrl() string {
	return conf.tributePaymentUrl
}

func GetYookasaWebHookUrl() string {
	return conf.yookasaWebhookUrl
}

func PlategaMerchantId() string {
	return conf.plategaMerchantId
}

func PlategaSecret() string {
	return conf.plategaSecret
}

func GetPlategaWebHookUrl() string {
	return conf.plategaWebhookUrl
}

func IsPlategaSBPEnabled() bool {
	return conf.isPlategaSBPEnabled
}

func IsPlategaCardsEnabled() bool {
	return conf.isPlategaCardsEnabled
}

func IsPlategaAcquiringEnabled() bool {
	return conf.isPlategaAcquiringEnabled
}

func IsPlategaWorldwideEnabled() bool {
	return conf.isPlategaWorldwideEnabled
}

func IsPlategaCryptoEnabled() bool {
	return conf.isPlategaCryptoEnabled
}

func IsPlategaEnabled() bool {
	return conf.isPlategaSBPEnabled || conf.isPlategaCardsEnabled ||
		conf.isPlategaAcquiringEnabled || conf.isPlategaWorldwideEnabled ||
		conf.isPlategaCryptoEnabled
}

func GetMiniAppURL() string {
	return conf.miniApp
}

func GetBlockedTelegramIds() map[int64]bool {
	return conf.blockedTelegramIds
}

func GetWhitelistedTelegramIds() map[int64]bool {
	return conf.whitelistedTelegramIds
}

func FeedbackURL() string {
	return conf.feedbackURL
}

func ChannelURL() string {
	return conf.channelURL
}

func ServerStatusURL() string {
	return conf.serverStatusURL
}

func SupportURL() string {
	return conf.supportURL
}

func TosURL() string {
	return conf.tosURL
}

func YookasaEmail() string {
	return conf.yookasaEmail
}

func Price1() int {
	return conf.price1
}

func Price3() int {
	return conf.price3
}

func Price6() int {
	return conf.price6
}

func Price12() int {
	return conf.price12
}

func DaysInMonth() int {
	return conf.daysInMonth
}

func Price(month int) int {
	switch month {
	case 1:
		return conf.price1
	case 3:
		return conf.price3
	case 6:
		return conf.price6
	case 12:
		return conf.price12
	default:
		return conf.price1
	}
}

func StarsPrice(month int) int {
	switch month {
	case 1:
		return conf.starsPrice1
	case 3:
		return conf.starsPrice3
	case 6:
		return conf.starsPrice6
	case 12:
		return conf.starsPrice12
	default:
		return conf.starsPrice1
	}
}

func TelegramToken() string {
	return conf.telegramToken
}

func TelegramProxyURL() string {
	return envStringDefault("TELEGRAM_PROXY_URL", "")
}

func MoynalogProxyURL() string {
	return envStringDefault("MOYNALOG_PROXY_URL", "")
}

func DadaBaseUrl() string {
	return conf.databaseURL
}

func CryptoPayUrl() string {
	return conf.cryptoPayURL
}

func CryptoPayToken() string {
	return conf.cryptoPayToken
}

func BotURL() string {
	return conf.botURL
}

func SetBotURL(botURL string) {
	conf.botURL = botURL
}

func YookasaUrl() string {
	return conf.yookasaURL
}

func YookasaShopId() string {
	return conf.yookasaShopId
}

func YookasaSecretKey() string {
	return conf.yookasaSecretKey
}

func IsCryptoPayEnabled() bool {
	return conf.isCryptoEnabled
}

func IsYookasaEnabled() bool {
	return conf.isYookasaEnabled
}

func IsTelegramStarsEnabled() bool {
	return conf.isTelegramStarsEnabled
}

func RequirePaidPurchaseForStars() bool {
	return conf.requirePaidPurchaseForStars
}

func GetAdminTelegramId() int64 {
	return conf.adminTelegramId
}

func GetPrivateChannelID() int64 {
	return conf.privateChannelID
}

func GetHealthCheckPort() int {
	return conf.healthCheckPort
}

func IsWepAppLinkEnabled() bool {
	return conf.isWebAppLinkEnabled
}

func MoynalogUrl() string {
	return conf.moynalogURL
}

func MoynalogUsername() string {
	return conf.moynalogUsername
}

func MoynalogPassword() string {
	return conf.moynalogPassword
}

func IsMoynalogEnabled() bool {
	return conf.isMoynalogEnabled
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Panicf("env %q not set", key)
	}
	return v
}

func mustEnvInt(key string) int {
	v := mustEnv(key)
	i, err := strconv.Atoi(v)
	if err != nil {
		log.Panicf("invalid int in %q: %v", key, err)
	}
	return i
}

func envIntDefault(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		log.Panicf("invalid int in %q: %v", key, err)
	}
	return i
}

func envStringDefault(key string, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

func envBool(key string) bool {
	return os.Getenv(key) == "true"
}

func InitConfig() {
	if os.Getenv("DISABLE_ENV_FILE") != "true" {
		if err := godotenv.Load(".env"); err != nil {
			log.Println("No .env loaded:", err)
		}
	}
	var err error
	conf.adminTelegramId, err = strconv.ParseInt(os.Getenv("ADMIN_TELEGRAM_ID"), 10, 64)
	if err != nil {
		panic("ADMIN_TELEGRAM_ID .env variable not set")
	}

	conf.privateChannelID, err = strconv.ParseInt(os.Getenv("PRIVATE_CHANNEL_ID"), 10, 64)
	if err != nil {
		panic("PRIVATE_CHANNEL_ID .env variable not set or invalid")
	}

	conf.telegramToken = mustEnv("TELEGRAM_TOKEN")
	conf.menuPhotoPath = strings.TrimSpace(os.Getenv("MENU_PHOTO_PATH"))

	conf.isWebAppLinkEnabled = func() bool {
		isWebAppLinkEnabled := os.Getenv("IS_WEB_APP_LINK") == "true"
		return isWebAppLinkEnabled
	}()

	conf.miniApp = envStringDefault("MINI_APP_URL", "")

	conf.defaultLanguage = envStringDefault("DEFAULT_LANGUAGE", "ru")

	conf.daysInMonth = envIntDefault("DAYS_IN_MONTH", 30)

	conf.healthCheckPort = envIntDefault("HEALTH_CHECK_PORT", 8080)

	conf.enableAutoPayment = envBool("ENABLE_AUTO_PAYMENT")

	conf.price1 = mustEnvInt("PRICE_1")
	conf.price3 = mustEnvInt("PRICE_3")
	conf.price6 = mustEnvInt("PRICE_6")
	conf.price12 = mustEnvInt("PRICE_12")

	conf.isTelegramStarsEnabled = envBool("TELEGRAM_STARS_ENABLED")
	if conf.isTelegramStarsEnabled {
		conf.starsPrice1 = envIntDefault("STARS_PRICE_1", conf.price1)
		conf.starsPrice3 = envIntDefault("STARS_PRICE_3", conf.price3)
		conf.starsPrice6 = envIntDefault("STARS_PRICE_6", conf.price6)
		conf.starsPrice12 = envIntDefault("STARS_PRICE_12", conf.price12)
	}

	conf.requirePaidPurchaseForStars = envBool("REQUIRE_PAID_PURCHASE_FOR_STARS")

	conf.databaseURL = mustEnv("DATABASE_URL")

	conf.isCryptoEnabled = envBool("CRYPTO_PAY_ENABLED")
	if conf.isCryptoEnabled {
		conf.cryptoPayURL = mustEnv("CRYPTO_PAY_URL")
		conf.cryptoPayToken = mustEnv("CRYPTO_PAY_TOKEN")
	}

	conf.isYookasaEnabled = envBool("YOOKASA_ENABLED")
	if conf.isYookasaEnabled {
		conf.yookasaURL = mustEnv("YOOKASA_URL")
		conf.yookasaShopId = mustEnv("YOOKASA_SHOP_ID")
		conf.yookasaSecretKey = mustEnv("YOOKASA_SECRET_KEY")
		conf.yookasaEmail = mustEnv("YOOKASA_EMAIL")
		conf.yookasaWebhookUrl = os.Getenv("YOOKASA_WEBHOOK_URL")
	}

	if envBool("PLATEGA_ENABLED") {
		conf.plategaMerchantId = mustEnv("PLATEGA_MERCHANT_ID")
		conf.plategaSecret = mustEnv("PLATEGA_SECRET")
		conf.plategaWebhookUrl = mustEnv("PLATEGA_WEBHOOK_URL")
		conf.isPlategaSBPEnabled = envBool("PLATEGA_SBP_ENABLED")
		conf.isPlategaCardsEnabled = envBool("PLATEGA_CARDS_ENABLED")
		conf.isPlategaAcquiringEnabled = envBool("PLATEGA_ACQUIRING_ENABLED")
		conf.isPlategaWorldwideEnabled = envBool("PLATEGA_WORLDWIDE_ENABLED")
		conf.isPlategaCryptoEnabled = envBool("PLATEGA_CRYPTO_ENABLED")
	}

	conf.serverStatusURL = os.Getenv("SERVER_STATUS_URL")
	conf.supportURL = os.Getenv("SUPPORT_URL")
	conf.feedbackURL = os.Getenv("FEEDBACK_URL")
	conf.channelURL = os.Getenv("CHANNEL_URL")
	conf.tosURL = os.Getenv("TOS_URL")

	conf.tributeWebhookUrl = os.Getenv("TRIBUTE_WEBHOOK_URL")
	if conf.tributeWebhookUrl != "" {
		conf.tributeAPIKey = mustEnv("TRIBUTE_API_KEY")
		conf.tributePaymentUrl = mustEnv("TRIBUTE_PAYMENT_URL")
	}

	conf.blockedTelegramIds = func() map[int64]bool {
		v := os.Getenv("BLOCKED_TELEGRAM_IDS")
		if v != "" {
			ids := strings.Split(v, ",")
			var blockedMap = make(map[int64]bool)
			for _, idStr := range ids {
				id, err := strconv.ParseInt(strings.TrimSpace(idStr), 10, 64)
				if err != nil {
					panic(fmt.Sprintf("invalid telegram ID in BLOCKED_TELEGRAM_IDS: %v", err))
				}
				blockedMap[id] = true
			}
			slog.Info("Loaded blocked telegram IDs", "count", len(blockedMap))
			return blockedMap
		} else {
			slog.Info("No blocked telegram IDs specified")
			return map[int64]bool{}
		}
	}()

	conf.whitelistedTelegramIds = func() map[int64]bool {
		v := os.Getenv("WHITELISTED_TELEGRAM_IDS")
		if v != "" {
			ids := strings.Split(v, ",")
			var whitelistedMap = make(map[int64]bool)
			for _, idStr := range ids {
				id, err := strconv.ParseInt(strings.TrimSpace(idStr), 10, 64)
				if err != nil {
					panic(fmt.Sprintf("invalid telegram ID in WHITELISTED_TELEGRAM_IDS: %v", err))
				}
				whitelistedMap[id] = true
			}
			slog.Info("Loaded whitelisted telegram IDs", "count", len(whitelistedMap))
			return whitelistedMap
		} else {
			slog.Info("No whitelisted telegram IDs specified")
			return map[int64]bool{}
		}
	}()

	conf.isMoynalogEnabled = envBool("MOYNALOG_ENABLED")
	if conf.isMoynalogEnabled {
		conf.moynalogURL = envStringDefault("MOYNALOG_URL", "https://moynalog.ru/api/v1")
		conf.moynalogUsername = mustEnv("MOYNALOG_USERNAME")
		conf.moynalogPassword = mustEnv("MOYNALOG_PASSWORD")
	}
}
