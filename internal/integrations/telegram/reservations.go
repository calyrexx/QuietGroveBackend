package telegram

import (
	"context"
	"fmt"
	"github.com/calyrexx/QuietGrooveBackend/internal/entities"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"strings"
)

func (a *Adapter) ReservationCreatedForAdmin(msg entities.ReservationCreatedMessage) error {
	ctx := context.Background()

	text := fmt.Sprintf(
		"✅ *Новое бронирование*\n"+
			"🏠 Дом: %s\n"+
			"👤 Гость: %s\n"+
			"📞 %s\n"+
			"📅 %s → %s\n"+
			"👥 %d гостей\n"+
			"💳 %d ₽\n",
		msg.HouseName, msg.GuestName, msg.GuestPhone,
		msg.CheckIn.Format("02.01.2006"), msg.CheckOut.Format("02.01.2006"),
		msg.GuestsCount, msg.TotalPrice,
	)

	if len(msg.Bathhouse) > 0 {
		text += "\n🔥 *Забронированы дополнительно:*\n"
		for _, bath := range msg.Bathhouse {
			fillOpt := ""
			if bath.FillOption != nil {
				fillOpt = "(" + *bath.FillOption + ")"
			}
			text += fmt.Sprintf(
				"• %s: %s с %s до %s %s\n",
				bath.Name,
				strings.ReplaceAll(bath.Date, "-", "."),
				bath.TimeFrom,
				bath.TimeTo,
				fillOpt,
			)
		}
	}

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

func (a *Adapter) ReservationCreatedForUser(msg entities.ReservationCreatedMessage, tgID int64) error {
	ctx := context.Background()

	text := fmt.Sprintf(
		"✅ *Ваше бронирование подтверждено!*\n"+
			"🏠 Дом: %s\n"+
			"📅 %s → %s\n"+
			"👥 %d гостей\n"+
			"💳 Стоимость проживания: %d ₽\n"+
			"📞 Наш номер для связи: +79867427283\n",
		msg.HouseName,
		msg.CheckIn.Format("02.01.2006"), msg.CheckOut.Format("02.01.2006"),
		msg.GuestsCount, msg.TotalPrice,
	)

	if len(msg.Bathhouse) > 0 {
		text += "\n🔥 *Забронированы дополнительно:*\n"
		for _, bath := range msg.Bathhouse {
			fillOpt := ""
			if bath.FillOption != nil {
				fillOpt = "(" + *bath.FillOption + ")"
			}
			text += fmt.Sprintf(
				"• %s: %s с %s до %s %s\n",
				bath.Name,
				strings.ReplaceAll(bath.Date, "-", "."),
				bath.TimeFrom,
				bath.TimeTo,
				fillOpt,
			)
		}
	}

	_, err := a.bot.SendMessage(ctx,
		&bot.SendMessageParams{
			ChatID:    tgID,
			Text:      text,
			ParseMode: "Markdown",
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func (a *Adapter) RemindUser(msg []entities.ReservationReminderNotification) error {
	ctx := context.Background()

	for _, m := range msg {
		text := fmt.Sprintf(
			"Уважаемый гость!\n"+
				"Ваше бронирование домика *%s* скоро начнётся!",
			m.HouseName,
		)

		_, err := a.bot.SendMessage(ctx,
			&bot.SendMessageParams{
				ChatID:    m.UserTgID,
				Text:      text,
				ParseMode: "Markdown",
				ReplyMarkup: &models.InlineKeyboardMarkup{
					InlineKeyboard: [][]models.InlineKeyboardButton{
						{
							{
								Text:         "Просмотреть бронирование 👀",
								CallbackData: fmt.Sprintf("view_resv_%s", m.UUID),
							},
						},
					},
				},
			},
		)
		if err != nil {
			a.logger.Error(err.Error())
		}
	}

	return nil
}

func (a *Adapter) myReservationsHandler(ctx context.Context, b *bot.Bot, u *models.Update) {
	var (
		tgID                  int64
		messageIdToDeleteBot  int
		messageIdToDeleteUser int
	)
	if u.CallbackQuery == nil {
		tgID = u.Message.Chat.ID
		messageIdToDeleteUser = u.Message.ID
	} else {
		tgID = u.CallbackQuery.From.ID
		messageIdToDeleteBot = u.CallbackQuery.Message.Message.ID
	}

	if messageIdToDeleteBot > 0 {
		_, err := a.bot.DeleteMessage(ctx, &bot.DeleteMessageParams{
			ChatID:    tgID,
			MessageID: messageIdToDeleteBot,
		})
		if err != nil {
			a.logger.Error(err.Error())
		}
	}

	if messageIdToDeleteUser > 0 {
		_, err := a.bot.DeleteMessage(ctx, &bot.DeleteMessageParams{
			ChatID:    tgID,
			MessageID: messageIdToDeleteUser,
		})
		if err != nil {
			a.logger.Error(err.Error())
		}
	}

	reservations, err := a.reservationSvc.GetByTelegramID(ctx, tgID)
	if err != nil {
		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: tgID,
			Text:   "⚠️ Не удалось получить бронирования. Попробуйте позже.",
		})
		if err != nil {
			a.logger.Error(err.Error())
		}
		return
	}

	if len(reservations) == 0 {
		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: tgID,
			Text:   "У вас пока нет бронирований.",
		})
		if err != nil {
			a.logger.Error(err.Error())
		}
		return
	}

	rows := make([][]models.InlineKeyboardButton, 0, len(reservations))
	for _, res := range reservations {
		text := fmt.Sprintf(
			"📅 %s → %s 🏠 %s",
			res.CheckIn.Format("02.01"),
			res.CheckOut.Format("02.01"),
			res.HouseName)
		btn := models.InlineKeyboardButton{
			Text:         text,
			CallbackData: fmt.Sprintf("view_resv_%s", res.UUID),
		}
		rows = append(rows, []models.InlineKeyboardButton{btn})
	}

	kb := &models.InlineKeyboardMarkup{
		InlineKeyboard: rows,
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      tgID,
		Text:        "Ваши бронирования:",
		ReplyMarkup: kb,
	})
	if err != nil {
		a.logger.Error(err.Error())
	}
}

func (a *Adapter) viewReservationCallback(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.CallbackQuery == nil {
		return
	}
	q := update.CallbackQuery

	uuid := strings.TrimPrefix(q.Data, "view_resv_")
	tgID := q.Message.Message.Chat.ID

	reservation, err := a.reservationSvc.GetDetailsByUUID(ctx, tgID, uuid)
	if err != nil {
		_, err = b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: q.ID,
			Text:            "Бронирование не найдено.",
			ShowAlert:       true,
		})
		if err != nil {
			a.logger.Error(err.Error())
		}
		return
	}

	var (
		statusMsg string
		canCancel bool
	)
	switch reservation.Status {
	case "confirmed":
		statusMsg = "Подтверждено ✅"
		canCancel = true
	case "cancelled":
		statusMsg = "Отменено ❌"
	case "checked_in":
		statusMsg = "В процессе ▶"
	case "checked_out":
		statusMsg = "Завершено ✅"
	}

	msg := fmt.Sprintf(
		"🏠 Дом: %s\n"+
			"📅 %s → %s\n"+
			"👥 %d гостей\n"+
			"💳 Стоимость проживания: %d₽\n"+
			"ℹ️ Статус: %s\n",
		reservation.HouseName,
		reservation.CheckIn.Format("02.01.2006"),
		reservation.CheckOut.Format("02.01.2006"),
		reservation.GuestsCount,
		reservation.TotalPrice,
		statusMsg,
	)

	if len(reservation.Bathhouse) > 0 {
		msg += "\n🔥 *Забронированы дополнительно*:\n"
		for _, bath := range reservation.Bathhouse {
			fillOpt := ""
			if bath.FillOptionName != nil {
				fillOpt = "(" + *bath.FillOptionName + ")"
			}
			msg += fmt.Sprintf("• %s: %s с %s до %s %s\n", bath.Name, bath.Date, bath.TimeFrom, bath.TimeTo, fillOpt)
		}
	}

	photo := &models.InputFileString{Data: reservation.ImageURL}

	_, err = b.SendPhoto(ctx, &bot.SendPhotoParams{
		ChatID:      tgID,
		Photo:       photo,
		Caption:     msg,
		ParseMode:   "Markdown",
		ReplyMarkup: a.buildReservationDetailKeyboard(uuid, canCancel),
	})
	if err != nil {
		a.logger.Error(err.Error())
		return
	}

	_, err = a.bot.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID:    tgID,
		MessageID: q.Message.Message.ID,
	})
	if err != nil {
		a.logger.Error(err.Error())
	}
}

func (a *Adapter) cancelReservationCallback(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.CallbackQuery == nil {
		return
	}
	q := update.CallbackQuery

	uuid := strings.TrimPrefix(q.Data, "cancel_resv_")
	tgID := q.Message.Message.Chat.ID

	reservation, err := a.reservationSvc.GetDetailsByUUID(ctx, tgID, uuid)
	if err != nil {
		_, err = b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: q.ID,
			Text:            "⚠ Бронирование не найдено.",
			ShowAlert:       true,
		})
		if err != nil {
			a.logger.Error(err.Error())
		}
		return
	}

	msg := fmt.Sprintf(
		"🏠 Дом: %s\n"+
			"📅 %s → %s\n"+
			"👥 %d гостей\n"+
			"💳 Стоимость проживания: %d₽\n"+
			"ℹ️ Статус: Отменено ❌\n",
		reservation.HouseName,
		reservation.CheckIn.Format("02.01.2006"),
		reservation.CheckOut.Format("02.01.2006"),
		reservation.GuestsCount,
		reservation.TotalPrice,
	)

	if len(reservation.Bathhouse) > 0 {
		msg += "\n🔥 *Забронированы дополнительно*:\n"
		for _, bath := range reservation.Bathhouse {
			fillOpt := ""
			if bath.FillOptionName != nil {
				fillOpt = "(" + *bath.FillOptionName + ")"
			}
			msg += fmt.Sprintf("• %s: %s с %s до %s %s\n", bath.Name, bath.Date, bath.TimeFrom, bath.TimeTo, fillOpt)
		}
	}

	err = a.reservationSvc.Cancel(ctx, tgID, uuid)
	if err != nil {
		_, err = b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: q.ID,
			Text:            "⚠️Не удалось отменить бронирование. Попробуйте позже.",
			ShowAlert:       true,
		})
		if err != nil {
			a.logger.Error(err.Error())
		}
		return
	}

	_, err = b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: q.ID,
		Text:            "Ваше бронирование отменено!",
		ShowAlert:       true,
	})
	if err != nil {
		a.logger.Error(err.Error())
	}

	_, err = b.EditMessageCaption(ctx, &bot.EditMessageCaptionParams{
		ChatID:      tgID,
		MessageID:   q.Message.Message.ID,
		Caption:     msg,
		ParseMode:   "Markdown",
		ReplyMarkup: a.buildReservationDetailKeyboard(uuid, false),
	})
	if err != nil {
		a.logger.Error(err.Error())
	}
}

func (a *Adapter) buildReservationDetailKeyboard(reservationUUID string, canCancel bool) *models.InlineKeyboardMarkup {
	var rows [][]models.InlineKeyboardButton

	if canCancel {
		rows = append(rows, []models.InlineKeyboardButton{
			{
				Text:         "Отменить бронирование ❌",
				CallbackData: fmt.Sprintf("cancel_resv_%s", reservationUUID),
			},
		})
	}
	rows = append(rows, []models.InlineKeyboardButton{
		{
			Text:         "⬅️ Назад",
			CallbackData: "my_reservations_back",
		},
	})
	return &models.InlineKeyboardMarkup{InlineKeyboard: rows}
}
