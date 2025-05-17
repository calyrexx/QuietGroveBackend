package telegram

import (
	"context"
	"fmt"
	"github.com/Calyr3x/QuietGrooveBackend/internal/configuration"
	"github.com/Calyr3x/QuietGrooveBackend/internal/entities"
	"github.com/Calyr3x/QuietGrooveBackend/internal/pkg/errorspkg"
	"github.com/go-telegram/bot"
)

type TGNotifier struct {
	bot     *bot.Bot
	chatIDs []int
}

func NewTelegramNotifier(creds *configuration.TelegramBot) (*TGNotifier, error) {
	if creds == nil {
		return nil, errorspkg.NewErrConstructorDependencies("NewTelegramNotifier", "creds", "nil")
	}

	b, err := bot.New(creds.Token)
	if err != nil {
		return nil, err
	}

	return &TGNotifier{
		bot:     b,
		chatIDs: creds.AdminChatIDs,
	}, nil
}

func (n *TGNotifier) ReservationCreated(res entities.ReservationCreatedMessage) error {
	ctx := context.Background()
	text := fmt.Sprintf(
		"✅ *Новое бронирование*\n\n"+
			"🏠 Дом: %s\n"+
			"👤 Гость: %s, 📞: %s\n"+
			"📅 %s → %s\n"+
			"👥 Гостей: %d\n"+
			"💳 %d ₽",
		res.House,
		res.GuestName,
		res.GuestPhone,
		res.CheckIn.Format("02.01.2006"),
		res.CheckOut.Format("02.01.2006"),
		res.GuestsCount,
		res.TotalPrice,
	)

	for _, chatID := range n.chatIDs {
		println("Sending message to: ", chatID)
		_, err := n.bot.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    chatID,
			Text:      text,
			ParseMode: "Markdown",
		})
		if err != nil {
			continue
		}
	}

	return nil
}
