# Private Channel Subscription Bot

A Telegram bot designed to automate the sale and management of paid subscriptions to private Telegram channels. 

This bot automatically generates one-time invite links for buyers upon successful payment, tracks their subscription expiration, and automatically kicks them from the channel when their subscription expires.

## Features

- **Automated Access Management**: Generates secure, single-use invite links for users upon successful payment.
- **Subscription Tracking**: Tracks user subscription expiration dates and automatically removes users from the private channel when their time runs out (they are immediately unbanned so they can buy again).
- **Multiple Payment Gateways**:
  - [YooKassa](https://yookassa.ru/developers/api)
  - [Crypto Pay](https://help.crypt.bot/crypto-pay-api)
  - Telegram Stars
  - Tribute
  - [Platega](https://platega.io): SBP, Russian cards, acquiring, international cards, and cryptocurrency
- **Tax Integration**: Automatic receipt generation via integration with "Мой Налог" (My Nalog) for self-employed individuals in Russia.
- **Notifications**: Users are automatically notified 3 days and 1 day before their subscription expires, and immediately upon expiration.
- **Localization**: Supports Russian and English interfaces.
- **Admin Features**: Whitelists/Blacklists for specific Telegram IDs.

## Prerequisites

- Docker and Docker Compose
- A registered Telegram Bot (via [@BotFather](https://t.me/BotFather))
- A private Telegram channel where the bot has administrator rights (specifically "Invite Users via Link" and "Ban Users" permissions).

## Installation & Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/your-repo/private-channel-bot.git
   cd private-channel-bot
   ```

2. **Configure Environment Variables**
   Copy the provided `.env.sample` file to `.env` and configure your credentials.
   ```bash
   cp .env.sample .env
   ```
   
   **Key Variables to set:**
   - `TELEGRAM_TOKEN`: Your bot token from @BotFather.
   - `PRIVATE_CHANNEL_ID`: The ID of your private channel (e.g., `-100123456789`). The bot must be an admin here.
   - `PRICE_1`, `PRICE_3`, `PRICE_6`, `PRICE_12`: Subscription prices in RUB for 1, 3, 6, and 12 months.
   - Configure your active payment gateways (CryptoPay, YooKassa, etc.).

3. **Database**
   The project uses PostgreSQL. By default, the provided `docker-compose.yaml` handles setting up the database. No manual configuration is needed unless you are running it without Docker.
   Database migrations run automatically on application startup.

4. **Start the Application**
   Run the bot and database using Docker Compose:
   ```bash
   docker-compose up -d --build
   ```

## Menu photo and pinned start message

The bot can show a photo on the start menu and all its inline-keyboard sections.
Put your own JPEG or PNG in `assets/menu.jpg` and set this in `.env`:

```dotenv
MENU_PHOTO_PATH=/assets/menu.jpg
```

- The file must be readable by the container's UID 1000. For a non-Docker installation, use a local path such as `MENU_PHOTO_PATH=./assets/menu.jpg`. Leave the value empty for text-only menus.
- Telegram photo limits: JPEG/PNG up to 10 MB; width + height at most 10000 pixels; aspect ratio at most 20. 

## Development

If you prefer to run the bot locally without Docker (requires a local PostgreSQL instance):

```bash
# Export env variables from your .env file
export $(cat .env | xargs)

# Run the bot
go run cmd/app/main.go
```
