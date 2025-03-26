package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"

	"jirabot/internal/config" // Добавляем импорт конфига
)

type Bot struct {
	api       *tgbotapi.BotAPI
	channelID int64
}

// NewBot создает нового Telegram-бота
func NewBot(cfg *config.TelegramConfig) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		return nil, fmt.Errorf("bot init error: %w", err)
	}

	var chatID int64
	if _, err := fmt.Sscanf(cfg.ChannelID, "%d", &chatID); err != nil {
		return nil, fmt.Errorf("invalid channel ID: %w", err)
	}

	return &Bot{
		api:       bot,
		channelID: chatID,
	}, nil
}

// SendAlert отправляет уведомление в Telegram-канал
func (b *Bot) SendAlert(taskID, status, priority string, duration time.Duration) error {
	msg := tgbotapi.NewMessage(
		b.channelID,
		fmt.Sprintf("⚠️ *%s*\n- Status: %s\n- Priority: %s\n- Overdue: %v",
			taskID, status, priority, duration.Round(time.Minute)),
	)
	msg.ParseMode = "Markdown"

	_, err := b.api.Send(msg)
	return err
}
