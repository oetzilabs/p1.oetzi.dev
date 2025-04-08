package states

import "p1/pkg/models"

type ClientState struct {
	Projects []*models.Project
	Servers  []*models.Server
	Brokers  []*models.Broker
	Metrics  *models.Metrics
}
