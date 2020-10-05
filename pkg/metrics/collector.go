package metrics

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/wladox/schedule-advisor/pkg/api"
	"github.com/wladox/schedule-advisor/pkg/flink"
	"google.golang.org/grpc"
	"strconv"
	"time"
)

const (
	checkpointsUrl               = "checkpoints"
	checkpointsConfigUrl         = "checkpoints/config"
	checkpointDetailsUrl         = "checkpoints/details"
	checkpointDetailsSubtasksUrl = "subtasks"
	ingestionRateParams          = "numBytesInPerSecond,numBytesOutPerSecond"
)

type Archiver interface {
	Save(data []string, dest string) error
}

type Output struct {
	Time         string
	TaskId       string
	SubtaskIndex int
	Name         string
	Value        string
}

type BandwidthOutput struct {
	Time  string
	Link  string
	Value string
}

type MetricCollector struct {
	client  *flink.Client
	writer  Archiver
	log     *log.Logger
	calc    chan<- flink.CheckpointConfig
	jobId   string
	storage Storage
}

func NewMetricCollector(client *flink.Client, writer Archiver, log *log.Logger, calc chan flink.CheckpointConfig, store Storage) *MetricCollector {
	return &MetricCollector{
		client:  client,
		writer:  writer,
		log:     log,
		calc:    calc,
		storage: store,
	}
}

func (c *MetricCollector) getCheckpointConfig(ctx context.Context) error {

	url := fmt.Sprintf("jobs/%s/%s", c.jobId, checkpointsConfigUrl)

	cfg := flink.CheckpointConfig{}
	err := c.client.GetCheckpointConfig(ctx, cfg, url)
	if err != nil {
		return err
	}

	log.WithField("Checkpoint Interval", cfg.Interval).Info("Checkpoint config retrieved.")
	c.calc <- cfg
	return nil
}

func (c *MetricCollector) CollectDataflowInfo(ctx context.Context) error {
	jobId, err := c.client.GetJob(ctx)
	if err != nil {
		c.log.WithError(err).Error("Can't find a job")
		return err
	}

	if len(jobId) == 0 {
		c.log.Info("No job is running.")
		return nil
	}

	if len(c.jobId) == 0 {
		c.jobId = jobId
		err = c.getCheckpointConfig(ctx)
		if err != nil {
			c.log.WithError(err).Error("Checkpoint info could not be retrieved.")
			return err
		}
	}

	plan, err := c.client.GetPlan(ctx, jobId)
	if err != nil {
		c.log.WithError(err).Error("Could not retrieve job plan")
		return err
	}

	if len(plan.Nodes) == 0 {
		c.log.Info("Job plan is empty.")
		return nil
	}

	vertices := make(map[string]int, 0)

	for _, node := range plan.Nodes {
		vertices[node.Id] = node.Parallelism
	}

	// INGESTION RATE
	ingestionRateMetrics, err := c.collectIngestionRate(ctx, jobId, vertices)
	if err != nil {
		c.log.WithError(err).Error("Ingestion rate could not be retrieved")
		return err
	}

	err = c.writer.Save(convertToOutput(ingestionRateMetrics), "output/ingestion_rate.csv")
	if err != nil {
		c.log.WithError(err).Error("Could not write to file")
		return err
	}

	// OPERATOR STATE
	checkpointMetrics, err := c.collectCheckpointMetrics(ctx, jobId)
	if err != nil {
		c.log.WithError(err).Error("Checkpoint metrics could not be retrieved")
		return err
	}

	err = c.writer.Save(convertToOutput(checkpointMetrics), "output/operator_state.csv")
	if err != nil {
		c.log.WithError(err).Error("Could not write to file")
		return err
	}

	return nil
}

func (c *MetricCollector) CollectNetworkInfo(ctx context.Context) error {
	conn := openGrpcConnection()
	if conn != nil {
		defer conn.Close()
	}

	log.Info("Connected to grpc server.")

	meshClient := api.NewResourceMeshRPCClient(conn)
	all, err := meshClient.GetAllNodes(ctx, &api.Empty{})
	if err != nil {
		log.WithError(err).Fatal("Connection to api master failed.")
	}

	output := make([]BandwidthOutput, 0)

	now := time.Now()

	for _, node := range all.Nodes {
		for _, info := range node.Info.PeerInformation {
			if info.MetricName == "Bandwidth" && len(info.Peer) > 0 {
				log.WithField("From", node.Id.Id).WithField("To", info.Peer).Infof("Bandwidth: %f", info.Metric)
				bo := BandwidthOutput{
					Time:  now.String(),
					Link:  fmt.Sprintf("%s_%s", node.Id.Id, info.Peer),
					Value: fmt.Sprintf("%d", int(info.Metric)),
				}
				output = append(output, bo)
			}
		}
	}

	err = c.writer.Save(convertBandwidthToOutput(output), "output/bandwidth.csv")
	if err != nil {
		c.log.WithError(err).Error("Could not write to file")
		return err
	}

	return nil
}

func (c *MetricCollector) collectIngestionRate(ctx context.Context, jobId string, vertices map[string]int) ([]Output, error) {
	ingestionRateMetrics := make([]Output, 0)

	for taskId, parallelism := range vertices {
		for i := 0; i < parallelism; i++ {
			url := fmt.Sprintf("jobs/%s/vertices/%s/subtasks/%d/metrics?get=%s", jobId, taskId, i, ingestionRateParams)
			metrics, err := c.client.CollectVertexMetrics(ctx, url)
			if err != nil {
				c.log.WithError(err).Error("Vertex metrics could not be collected")
			}
			for _, metric := range metrics {
				output := Output{
					Time:         time.Now().String(),
					TaskId:       taskId,
					SubtaskIndex: i,
					Name:         metric.Id,
					Value:        metric.Value,
				}
				ingestionRateMetrics = append(ingestionRateMetrics, output)
			}
		}
	}

	return ingestionRateMetrics, nil
}

func (c *MetricCollector) collectCheckpointMetrics(ctx context.Context, jobId string) ([]Output, error) {
	cp, _ := c.client.GetCheckpoints(ctx, fmt.Sprintf("jobs/%s/%s", jobId, checkpointsUrl))

	checkpointId := cp.Latest.Completed.Id

	cpDetails, _ := c.client.GetCheckpointDetails(ctx, fmt.Sprintf("jobs/%s/%s/%d", jobId, checkpointDetailsUrl, checkpointId))

	ts := time.Now()

	metrics := make([]Output, 0)
	for taskId, _ := range cpDetails.Tasks {
		subtaskStats, _ := c.client.GetCheckpointSubtasks(ctx, fmt.Sprintf("jobs/%s/%s/%d/%s/%s", jobId, checkpointDetailsUrl, checkpointId, checkpointDetailsSubtasksUrl, taskId))

		for _, subtask := range subtaskStats.Subtasks {
			output := Output{
				Time:         ts.String(),
				TaskId:       taskId,
				SubtaskIndex: subtask.Index,
				Name:         "state_size",
				Value:        strconv.Itoa(int(subtask.StateSize)),
			}
			metrics = append(metrics, output)
		}
	}

	return metrics, nil
}

func openGrpcConnection() *grpc.ClientConn {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial("0.0.0.0:2020", grpc.WithInsecure())
	if err != nil {
		log.WithError(err).Error("Could not open connection to grpc server")
		return nil
	}
	return conn
}

func convertToOutput(stats []Output) []string {
	lines := make([]string, 0)
	for _, stat := range stats {
		line := fmt.Sprintf("%d,%s,%d,%s,%s", time.Now().Unix(), stat.TaskId, stat.SubtaskIndex, stat.Name, stat.Value)
		lines = append(lines, line)
	}
	return lines
}

func convertBandwidthToOutput(stats []BandwidthOutput) []string {
	lines := make([]string, 0)
	for _, stat := range stats {
		line := fmt.Sprintf("%d,%s,%s", time.Now().Unix(), stat.Link, stat.Value)
		lines = append(lines, line)
	}
	return lines
}
