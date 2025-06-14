package telegram

import (
	"context"
	"fmt"
	"github.com/calyrexx/QuietGrooveBackend/internal/entities"
	"github.com/go-telegram/bot"
	"strings"
)

func (a *Adapter) NewApplicationForEvent(res entities.NewApplication) error {
	ctx := context.Background()

	text := fmt.Sprintf(
		"🎉 *Новая заявка на мероприятие!*\n"+
			"👤 Имя: %s\n"+
			"📞 Телефон: %s\n"+
			"📅 Дата: %s\n"+
			"👥 Кол-во гостей: %d",
		res.Name,
		res.Phone,
		strings.ReplaceAll(res.CheckIn, "-", "."),
		res.GuestsCount,
	)

	for _, chatID := range a.adminChatIDs {
		_, err := a.bot.SendMessage(ctx,
			&bot.SendMessageParams{
				ChatID:    chatID,
				Text:      text,
				ParseMode: "Markdown",
			},
		)
		if err != nil {
			return err
		}
	}
	return nil
}
