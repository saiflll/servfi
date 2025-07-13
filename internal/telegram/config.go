package telegram

import (
	"log"
	"os"
	"strconv"
)

var (
	BotToken string
	ChatID   int64
)

func LoadConfig() {
	BotToken = os.Getenv("TELEGRAM_BOT_TOKEN")
	if BotToken == "" {
		log.Println("Peringatan: Environment variable TELEGRAM_BOT_TOKEN tidak diatur. Notifikasi Telegram akan dinonaktifkan.")
	}

	chatIDStr := os.Getenv("TELEGRAM_CHAT_ID")
	if chatIDStr == "" {
		log.Println("Peringatan: Environment variable TELEGRAM_CHAT_ID tidak diatur. Notifikasi Telegram akan dinonaktifkan.")
	} else {
		var err error
		ChatID, err = strconv.ParseInt(chatIDStr, 10, 64)
		if err != nil {
			log.Printf("Peringatan: TELEGRAM_CHAT_ID '%s' tidak valid. Harus berupa integer. Notifikasi Telegram akan dinonaktifkan. Error: %v", chatIDStr, err)
			BotToken = ""
		}
	}
}
