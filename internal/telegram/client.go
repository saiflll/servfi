package telegram

import (
	"fmt"
	"log"

	// "net/http" // Tidak dibutuhkan lagi

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var bot *tgbotapi.BotAPI

// InitBot tidak perlu diubah.
func InitBot() error {
	if BotToken == "" || ChatID == 0 {
		return nil
	}

	var err error
	bot, err = tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		return fmt.Errorf("gagal membuat instance bot Telegram: %w", err)
	}

	// log.Printf("Bot Telegram berhasil diautorisasi sebagai %s", bot.Self.UserName)
	return nil
}

// SendAlert tidak perlu diubah.
func SendAlert(messageText string) {
	if bot == nil || ChatID == 0 {
		return
	}

	go func(msgTxt string) {
		msg := tgbotapi.NewMessage(ChatID, msgTxt)
		msg.ParseMode = tgbotapi.ModeMarkdown
		msg.DisableWebPagePreview = true // Menonaktifkan pratinjau link

		if _, err := bot.Send(msg); err != nil {
			log.Printf("Error mengirim pesan Telegram: %v", err)
		}
	}(messageText)
}
