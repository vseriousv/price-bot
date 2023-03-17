package app

import "github.com/vseriousv/price-bot/internal/telegram"

type telegramService struct {
	service *service
}

func NewTelegramService(service *service) *telegramService {
	return &telegramService{
		service: service,
	}
}

func (s *telegramService) Run() error {
	return telegram.StartBot(s.service.Db, s.service.Config.TgToken)
}
