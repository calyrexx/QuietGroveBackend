package telegram

import (
	"context"
	"fmt"
	"github.com/Calyr3x/QuietGrooveBackend/internal/configuration"
	"github.com/Calyr3x/QuietGrooveBackend/internal/pkg/errorspkg"
	"regexp"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/Calyr3x/QuietGrooveBackend/internal/entities"
	"github.com/Calyr3x/QuietGrooveBackend/internal/usecases"
)

type Adapter struct {
	bot          *bot.Bot
	adminChatIDs []int64
	verifSvc     *usecases.Verification
}

func NewAdapter(creds *configuration.TelegramBot) (*Adapter, error) {
	if creds == nil {
		return nil, errorspkg.NewErrConstructorDependencies("NewAdapter", "creds", "nil")
	}

	b, err := bot.New(creds.Token)
	if err != nil {
		return nil, err
	}
	return &Adapter{
		bot:          b,
		adminChatIDs: creds.AdminChatIDs,
	}, nil
}

func (a *Adapter) ReservationCreated(msg entities.ReservationCreatedMessage) error {
	ctx := context.Background()

	text := fmt.Sprintf(
		"✅ *Новое бронирование*\n"+
			"🏠 Дом: %s\n"+
			"👤 Гость: %s\n"+
			"📞 %s\n"+
			"📅 %s → %s\n"+
			"👥 %d гостей\n"+
			"💳 %d ₽",
		msg.House, msg.GuestName, msg.GuestPhone,
		msg.CheckIn.Format("02.01.2006"), msg.CheckOut.Format("02.01.2006"),
		msg.GuestsCount, msg.TotalPrice,
	)

	for _, chatID := range a.adminChatIDs {
		if _, err := a.bot.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    chatID,
			Text:      text,
			ParseMode: "Markdown",
		}); err != nil {
			return err
		}
	}
	return nil
}

func (a *Adapter) RegisterHandlers(ver *usecases.Verification) {
	a.verifSvc = ver

	re := regexp.MustCompile(`^\d{6}$`)
	a.bot.RegisterHandlerRegexp(
		bot.HandlerTypeMessageText,
		re,
		a.codeHandler,
	)
}

func (a *Adapter) codeHandler(ctx context.Context, b *bot.Bot, u *models.Update) {
	code := u.Message.Text
	tgID := u.Message.Chat.ID

	if err := a.verifSvc.Approve(ctx, code, tgID); err != nil {
		_, err = b.SendMessage(ctx, &bot.SendMessageParams{ChatID: tgID, Text: "❌ Код неверный или устарел"})
		if err != nil {
			return
		}
		return
	}

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{ChatID: tgID, Text: "✅ Личность подтверждена!"})
	if err != nil {
		return
	}
}

func (a *Adapter) Run(ctx context.Context) {
	a.bot.Start(ctx)
}
