package app

import (
	"github.com/vseriousv/price-bot/internal/price_monitoring_worker"
)

type priceMonitoringWorkerService struct {
	service *service
}

func NewPriceMonitoringWorkerService(service *service) *priceMonitoringWorkerService {
	return &priceMonitoringWorkerService{
		service: service,
	}
}

func (s *priceMonitoringWorkerService) Run() error {
	return price_monitoring_worker.StartWorker(s.service.Db, s.service.Config.TgToken)
}
