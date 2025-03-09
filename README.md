# Crypto Price Telegram Bot

A Telegram bot written in Go that fetches cryptocurrency prices and sends an image with the calculated value.

## Features
- Retrieves real-time cryptocurrency prices using the CryptoCompare API
- Processes user messages with cryptocurrency and amount (e.g., `1 BTC` or `10 ETH`)
- Generates an image with the price and sends it to the user
- Provides a link to the TradingView chart for the specified cryptocurrency

## Requirements
- Go 1.18+
- A Telegram bot token (get it from [BotFather](https://t.me/botfather))
- A CryptoCompare API key

## Installation
1. Clone the repository:
   ```sh
   git clone https://github.com/godbrowed/Golang_cryptoprice_bot.git
   cd Golang_cryptoprice_bot
   ```
2. Install dependencies:
   ```sh
   go mod tidy
   ```
3. Set up environment variables:
   ```sh
   export TELEGRAM_BOT_TOKEN="your-telegram-bot-token"
   export CRYPTOCOMPARE_API_KEY="your-api-key"
   ```
   Or create a `.env` file with:
   ```env
   TELEGRAM_BOT_TOKEN=your-telegram-bot-token
   CRYPTOCOMPARE_API_KEY=your-api-key
   ```
4. Run the bot:
   ```sh
   go run main.go
   ```

## Usage
- Send a message in the format `<amount> <crypto>` (e.g., `2 BTC`).
- The bot will fetch the current price and send an image with the calculated value.
- If you send `1 UTC` or `1 GMT`, the bot will reply with `Not Slava already calculated it above` and an upward-pointing emoji.

## Deployment
You can deploy the bot on a VPS or use Docker:
```sh
docker build -t crypto-bot .
docker run -d --env TELEGRAM_BOT_TOKEN=your-token --env CRYPTOCOMPARE_API_KEY=your-key crypto-bot
```

## Author
[GodBrowed](https://github.com/godbrowed)

