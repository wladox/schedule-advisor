package forecast

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/wladox/schedule-advisor/pkg/flink"
	"github.com/wladox/schedule-advisor/pkg/metrics"
	"time"
)

type Calculator struct {
	client  *ForecastClient
	metrics chan flink.CheckpointConfig
	storage metrics.Storage
	log     *log.Logger
}

type Placement struct {
	Assignments map[string]string
}

type Request struct {
	CurrentPlacement Placement
	NewPlacement     Placement
}

type CostsResponse struct {
	CurrentCosts int `json:"current_costs"`
	OptimalCosts int `json:"optimal_costs"`
}

func NewCalculator(metrics chan flink.CheckpointConfig, client *ForecastClient, log *log.Logger, store metrics.Storage) *Calculator {
	return &Calculator{
		client:  client,
		metrics: metrics,
		log:     log,
		storage: store,
	}
}

func (c *Calculator) Init() {
	go func() {
		for {
			cfg := <-c.metrics
			log.Info(cfg)
			close(c.metrics)
		}
	}()
}

func (c *Calculator) CalculateCosts(ctx context.Context, request Request) CostsResponse {

	// get checkpointing interval
	checkpointInterval := c.storage.GetCheckpointInterval()
	var lastCheckpoint int

	// get forecasts
	forecasts, err := c.client.GetForecasts(ctx, 30)
	if err != nil {
		c.log.WithError(err).Error("Failed fetching forecasts.")
		return CostsResponse{}
	}

	c.log.Info("Forecasts retrieved")

	maxBw := findMax(forecasts.Bandwidth)
	maxIRate := findMax(forecasts.IngestionRate)
	maxState := findMax(forecasts.OperatorState)

	maxDowntime := maxState * maxBw

	// get current bandwidth
	currentBw := c.storage.GetBandwidth("test")

	// calculate downtime costs
	downtimeCost := 1
	for taskId := range request.CurrentPlacement.Assignments {
		if request.CurrentPlacement.Assignments[taskId] != request.NewPlacement.Assignments[taskId] {
			// get current operator state
			size := c.storage.GetOperatorStats(taskId)
			downtimeCost += currentBw * size
		}
	}

	c.log.WithField("DowntimeCost", downtimeCost).Info("Downtime cost calculated")

	// get forecasts for state size and bandwidth
	currentTime := time.Now().Unix()

	downtimeCostNormalized := downtimeCost / int(maxDowntime)
	c.log.WithField("Downtime Normalized", downtimeCostNormalized).Info("Downtime calculated")

	// get current ingestion rate
	currentRate := c.storage.GetIngestionRate()

	ingestionRateNormalized := currentRate / int(maxIRate)
	c.log.WithField("Ingestion rate Normalized", downtimeCostNormalized).Info("Ingestion rate calculated")

	costs := (downtimeCostNormalized + ingestionRateNormalized) * (int(currentTime) - lastCheckpoint) / checkpointInterval

	c.log.WithField("CurrentCosts", costs).Info("Reconfiguration costs calculated")

	return CostsResponse{
		CurrentCosts: costs,
		OptimalCosts: costs,
	}
}

func findMax(slice []float64) float64 {
	max := 1.0
	if len(slice) > 0 {
		for _, v := range slice {
			if v >= max {
				max = v
			}
		}
	}
	return max
}
