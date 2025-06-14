package telegram

import (
	"context"
	"github.com/calyrexx/QuietGrooveBackend/internal/configuration"
	"github.com/calyrexx/QuietGrooveBackend/internal/pkg/errorspkg"
	"github.com/calyrexx/QuietGrooveBackend/internal/usecases"
	"github.com/calyrexx/zeroslog"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"log/slog"
	"regexp"
)

const tgBot string = "telegramBot"

type Adapter struct {
	bot            *bot.Bot
	logger         *slog.Logger
	adminChatIDs   []int64
	verifSvc       *usecases.Verification
	reservationSvc *usecases.Reservation
}

func NewAdapter(creds *configuration.TelegramBot, logger *slog.Logger) (*Adapter, error) {
	if creds == nil {
		return nil, errorspkg.NewErrConstructorDependencies("NewAdapter", "creds", "nil")
	}
	if logger == nil {
		return nil, errorspkg.NewErrConstructorDependencies("NewAdapter", "logger", "nil")
	}

	newLogger := logger.With(zeroslog.ServiceKey, tgBot)

	b, err := bot.New(creds.Token)
	if err != nil {
		return nil, err
	}
	return &Adapter{
		bot:          b,
		logger:       newLogger,
		adminChatIDs: creds.AdminChatIDs,
	}, nil
}

func (a *Adapter) RegisterHandlers(ver *usecases.Verification, res *usecases.Reservation) {
	a.verifSvc = ver
	a.reservationSvc = res

	onlyDigits := regexp.MustCompile(`^\d+$`)

	a.bot.RegisterHandlerMatchFunc(
		func(u *models.Update) bool {
			return u.Message != nil && onlyDigits.MatchString(u.Message.Text)
		},
		a.verificationHandler,
	)

	a.bot.RegisterHandler(
		bot.HandlerTypeMessageText,
		"🏡 Мои бронирования",
		bot.MatchTypeExact,
		a.myReservationsHandler,
	)

	a.bot.RegisterHandler(
		bot.HandlerTypeCallbackQueryData,
		"view_resv_",
		bot.MatchTypePrefix,
		a.viewReservationCallback,
	)

	a.bot.RegisterHandler(
		bot.HandlerTypeMessageText,
		"/start",
		bot.MatchTypeExact,
		a.startHandler,
	)

}

func (a *Adapter) startHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	tgID := update.Message.Chat.ID

	replyMarkup := &models.ReplyKeyboardMarkup{
		Keyboard: [][]models.KeyboardButton{
			{
				{Text: "🏡 Мои бронирования"},
			},
		},
		ResizeKeyboard:  true,
		OneTimeKeyboard: false,
	}

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      tgID,
		Text:        "Добро пожаловать! Выберите действие:",
		ReplyMarkup: replyMarkup,
	})
	if err != nil {
		a.logger.Error(err.Error())
	}
}

func (a *Adapter) Run(ctx context.Context) {
	a.bot.Start(ctx)
}
