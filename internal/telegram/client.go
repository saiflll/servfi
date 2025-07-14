package telegram

import (
	"fmt"
	"log"

	

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var bot *tgbotapi.BotAPI


func InitBot() error {
	if BotToken == "" || ChatID == 0 {
		return nil
	}

	var err error
	bot, err = tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		return fmt.Errorf("gagal membuat instance bot Telegram: %w", err)
	}

	
	return nil
}


func SendAlert(messageText string) {
	if bot == nil || ChatID == 0 {
		return
	}

	go func(msgTxt string) {
		msg := tgbotapi.NewMessage(ChatID, msgTxt)
		msg.ParseMode = tgbotapi.ModeMarkdown
		msg.DisableWebPagePreview = true 

		if _, err := bot.Send(msg); err != nil {
			log.Printf("Error mengirim pesan Telegram: %v", err)
		}
	}(messageText)
}
