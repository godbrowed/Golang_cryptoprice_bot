package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings" // Додано для роботи з рядками

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/golang/freetype"
)

func getPrice(crypto string) (float64, error) {
	apiKey := "1451bd9877a2a3841315e390a6d096d0e863f854175d272b258b08ee0a79b63f"
	url := fmt.Sprintf("https://min-api.cryptocompare.com/data/price?fsym=%s&tsyms=USD", crypto)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Apikey %s", apiKey))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	var data map[string]float64
	err = json.Unmarshal(body, &data)
	if err != nil {
		return 0, err
	}
	price, exists := data["USD"]
	if !exists {
		return 0, fmt.Errorf("Не вдалося знайти валюту: %s", crypto)
	}

	return price, nil

}

func addPriceToImage(imagePath string, price float64) (image.Image, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Використовуємо тільки 2 змінні для jpeg.Decode
	img, err := jpeg.Decode(file)
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	rgba := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			rgba.Set(x, y, img.At(x, y))
		}
	}

	fontBytes, err := os.ReadFile("D:/IT/tgbots/Bymovement tg/DejaVuSans-Bold.ttf")
	if err != nil {
		return nil, err
	}
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return nil, err
	}

	ctx := freetype.NewContext()
	ctx.SetDPI(72)
	ctx.SetFont(font)
	ctx.SetFontSize(60)
	ctx.SetClip(rgba.Bounds())
	ctx.SetDst(rgba)
	ctx.SetSrc(image.NewUniform(color.RGBA{R: 138, G: 43, B: 255, A: 255})) // Трохи синій колір тексту

	priceText := fmt.Sprintf("$%.2f", price)
	pt := freetype.Pt(width/2-len(priceText)*7, height-140)
	_, err = ctx.DrawString(priceText, pt)
	if err != nil {
		return nil, err
	}

	return rgba, nil
}

func extractAmountAndCrypto(text string) (float64, string, error) {
	// Оновлений регулярний вираз для пошуку правильного формату, наприклад: "1 BTC"
	re := regexp.MustCompile(`^(\d+(\.\d+)?)\s*(\w+)$`)
	matches := re.FindStringSubmatch(text)
	if len(matches) < 4 {
		return 0, "", fmt.Errorf("invalid format")
	}
	amount, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return 0, "", err
	}
	crypto := matches[3]
	return amount, crypto, nil
}

func checkFile(imagePath string) {
	_, err := os.Stat(imagePath)
	if err != nil {
		log.Fatalf("File not found or error: %v", err)
	} else {
		fmt.Println("File exists")
	}
}

func main() {
	checkFile("images/image.jpg")

	bot, err := tgbotapi.NewBotAPI("7200898472:AAFXd0TPl7NEyu0qtGNmO5iwBT9eHNLfpFk")
	if err != nil {
		log.Fatal(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}

	for update := range updates {
		// Перевірка, чи є повідомлення і чи є текст
		if update.Message != nil && update.Message.Text != "" {
			chatID := update.Message.Chat.ID
			messageText := update.Message.Text

			// Перевірка, чи повідомлення містить правильний формат, наприклад "1 BTC"
			amount, crypto, err := extractAmountAndCrypto(messageText)
			if err != nil {
				continue
			}

			// Отримуємо ціну криптовалюти
			price, err := getPrice(crypto)
			if err != nil {
				continue
			}

			// Додавання ціни до зображення
			imagePath := "images/image.jpg"
			rgba, err := addPriceToImage(imagePath, price*amount)
			if err != nil {
				log.Fatal(err)
			}

			// Створюємо нове зображення
			outFile, err := os.Create("images/price_image.png")
			if err != nil {
				log.Fatal(err)
			}
			defer outFile.Close()

			// Кодуємо в PNG
			err = png.Encode(outFile, rgba)
			if err != nil {
				log.Fatal(err)
			}

			// Створення лінку на графік криптовалюти
			cryptoLink := fmt.Sprintf("https://www.tradingview.com/symbols/%sUSD/?exchange=CRYPTO", strings.ToUpper(crypto))

			// Створення текстового повідомлення для користувача
			msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("💰 Price: $%.2f", price*amount))
			msg.ParseMode = "MarkdownV2"

			// Створення кнопки для переходу до графіку
			photo := tgbotapi.NewPhotoUpload(chatID, "images/price_image.png")
			photo.Caption = msg.Text
			photo.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonURL("View Chart", cryptoLink),
				),
			)

			// Відправка фото з ціною і кнопкою
			_, err = bot.Send(photo)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
