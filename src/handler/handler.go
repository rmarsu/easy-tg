package handler

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type bot interface {
	Send(c tgbotapi.Chattable) error
	WaitForMessage(upd *tgbotapi.Update) *tgbotapi.Update
}

type Handler func(sender bot, upd *tgbotapi.Update)
