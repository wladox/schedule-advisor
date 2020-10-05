package flink

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/wladox/schedule-advisor/pkg/network"
	"io/ioutil"
	"log"
	"net/http"
)

type VertexMetric struct {
	Id    string `json:"id"`
	Value string `json:"value"`
}

type GraphNode struct {
	Id          string `json:"id"`
	Parallelism int    `json:"parallelism"`
}

type Jobs struct {
	Jobs []Job `json:"jobs"`
}

type Job struct {
	Id     string `json:"id"`
	Status string `json:"status"`
}

type JobPlanResponse struct {
	Plan JobPlan `json:"plan"`
}

type JobEntry struct {
	Jid   string `json:"jid"`
	State string `json:"state"`
}

type JobsOverview struct {
	Jobs []JobEntry
}

type JobPlan struct {
	Jid   string      `json:"jid"`
	Name  string      `json:"name"`
	Nodes []GraphNode `json:"nodes"`
}

type SubtaskResponse struct {
	Id                 int            `json:"id"`
	Status             string         `json:"status"`
	LatestAckTimestamp int64          `json:"latest_ack_timestamp"`
	StateSize          int            `json:"state_size"`
	Duration           int            `json:"end_to_end_duration"`
	NumSubtasks        int            `json:"num_subtasks"`
	Subtasks           []SubtaskStats `json:"subtasks"`
}

type SubtaskStats struct {
	Index        int    `json:"index"`
	AckTimestamp int64  `json:"ack_timestamp"`
	Duration     int    `json:"end_to_end_duration"`
	StateSize    int64  `json:"state_size"`
	Status       string `json:"status"`
}

type CheckpointConfig struct {
	Interval int `json:"interval"`
}

type CheckpointStats struct {
	Id          int    `json:"id"`
	Status      string `json:"status"`
	StateSize   int    `json:"state_size"`
	Duration    int    `json:"end_to_end_duration"`
	NumSubtasks int    `json:"num_subtasks"`
}

type CheckpointDetails struct {
	Id          int                        `json:"id"`
	Status      string                     `json:"status"`
	StateSize   int                        `json:"state_size"`
	Duration    int                        `json:"end_to_end_duration"`
	NumSubtasks int                        `json:"num_subtasks"`
	Tasks       map[string]CheckpointStats `json:"tasks"`
}

type CheckpointMetric struct {
	Latest struct {
		Completed CheckpointStats `json:"completed"`
	} `json:"latest"`
}

// Client to make the request to Flink API
type Client struct {
	doer network.HTTPDoer
	host string
}

// NewClient create new client
func NewClient(doer network.HTTPDoer, host string) *Client {
	return &Client{
		doer: doer,
		host: host,
	}
}

func (c *Client) CollectVertexMetrics(ctx context.Context, vertexMetricsUrl string) ([]VertexMetric, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", c.host, vertexMetricsUrl), nil)
	if err != nil {
		log.Fatal(err)
	}

	req = req.WithContext(ctx)
	res, err := c.doer.Do(req)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	metrics := make([]VertexMetric, 0)
	err = json.Unmarshal(body, &metrics)
	if err != nil {
		log.Fatal(err)
	}

	return metrics, nil
}

func (c *Client) GetCheckpointConfig(ctx context.Context, dest interface{}, url string) error {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", c.host, url), nil)
	if err != nil {
		return err
	}

	req = req.WithContext(ctx)
	res, err := c.doer.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, dest)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) GetCheckpoints(ctx context.Context, url string) (CheckpointMetric, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", c.host, url), nil)
	if err != nil {
		return CheckpointMetric{}, err
	}

	req = req.WithContext(ctx)
	res, err := c.doer.Do(req)
	if err != nil {
		return CheckpointMetric{}, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return CheckpointMetric{}, err
	}

	var metric CheckpointMetric
	err = json.Unmarshal(body, &metric)
	if err != nil {
		return CheckpointMetric{}, err
	}

	return metric, nil
}

func (c *Client) GetCheckpointDetails(ctx context.Context, url string) (CheckpointDetails, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", c.host, url), nil)
	if err != nil {
		return CheckpointDetails{}, err
	}

	req = req.WithContext(ctx)
	res, err := c.doer.Do(req)
	if err != nil {
		log.Fatal(err)
		return CheckpointDetails{}, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
		return CheckpointDetails{}, err
	}

	var metric CheckpointDetails
	err = json.Unmarshal(body, &metric)
	if err != nil {
		log.Fatal(err)
		return CheckpointDetails{}, err
	}

	return metric, nil
}

func (c *Client) GetCheckpointSubtasks(ctx context.Context, url string) (SubtaskResponse, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", c.host, url), nil)
	if err != nil {
		return SubtaskResponse{}, err
	}

	req = req.WithContext(ctx)
	res, err := c.doer.Do(req)
	if err != nil {
		return SubtaskResponse{}, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return SubtaskResponse{}, err
	}

	var metric SubtaskResponse
	err = json.Unmarshal(body, &metric)
	if err != nil {
		return SubtaskResponse{}, err
	}

	return metric, nil
}

func (c *Client) GetPlan(ctx context.Context, jobId string) (*JobPlan, error) {
	url := fmt.Sprintf("%s/jobs/%s/plan", c.host, jobId)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	res, err := c.doer.Do(req)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	if res.StatusCode != http.StatusInternalServerError {
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		var plan JobPlan
		err = json.Unmarshal(body, &plan)
		if err != nil {
			return nil, err
		}

		return &plan, nil
	}

	return &JobPlan{}, nil
}

func (c *Client) GetJob(ctx context.Context) (string, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", c.host, "jobs/overview"), nil)
	if err != nil {
		return "", err
	}

	req = req.WithContext(ctx)
	res, err := c.doer.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var overview JobsOverview
	err = json.Unmarshal(body, &overview)
	if err != nil {
		return "", err
	}

	if len(overview.Jobs) > 0 {
		return overview.Jobs[0].Jid, nil
	}

	return "", nil
}
