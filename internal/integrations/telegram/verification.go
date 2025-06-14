package telegram

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (a *Adapter) verificationHandler(ctx context.Context, b *bot.Bot, u *models.Update) {
	code := u.Message.Text
	tgID := u.Message.Chat.ID

	if len(code) != 6 {
		_, err := b.SendMessage(ctx,
			&bot.SendMessageParams{
				ChatID: tgID,
				Text:   "⚠️ Код должен содержать 6 цифр",
			},
		)
		if err != nil {
			a.logger.Error(err.Error())
		}
		return
	}

	if err := a.verifSvc.Approve(ctx, code, tgID); err != nil {
		_, err = b.SendMessage(ctx,
			&bot.SendMessageParams{
				ChatID: tgID,
				Text:   "❌ Код неверный или устарел",
			},
		)
		if err != nil {
			a.logger.Error(err.Error())
		}
		return
	}

	_, err := b.SendMessage(ctx,
		&bot.SendMessageParams{
			ChatID: tgID,
			Text:   "✅ Личность подтверждена!",
		},
	)
	if err != nil {
		a.logger.Error(err.Error())
	}
}
