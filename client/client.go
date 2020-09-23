package client

import (
	"context"
	"encoding/json"
	"fmt"

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
	StateSize    int    `json:"state_size"`
	Status       string `json:"status"`
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

type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client to make the request to voucher-service
type Client struct {
	doer HTTPDoer
	host string
}

// NewClient create new client to voucher service to fetch details
func NewClient(doer HTTPDoer, host string) *Client {
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

func (c *Client) GetCheckpoints(ctx context.Context, url string) (CheckpointMetric, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", c.host, url), nil)
	if err != nil {
		log.Fatal(err)
	}

	req = req.WithContext(ctx)
	res, err := c.doer.Do(req)
	if err != nil {
		log.Fatal(err)
		return CheckpointMetric{}, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
		return CheckpointMetric{}, err
	}

	var metric CheckpointMetric
	err = json.Unmarshal(body, &metric)
	if err != nil {
		log.Fatal(err)
		return CheckpointMetric{}, err
	}

	return metric, nil
}

func (c *Client) GetCheckpointDetails(ctx context.Context, url string) (CheckpointDetails, error) {
	fmt.Println(fmt.Sprintf("%s/%s", c.host, url))
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", c.host, url), nil)
	if err != nil {
		log.Fatal(err)
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
	fmt.Println(fmt.Sprintf("%s/%s", c.host, url))
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", c.host, url), nil)
	if err != nil {
		log.Fatal(err)
	}

	req = req.WithContext(ctx)
	res, err := c.doer.Do(req)
	if err != nil {
		log.Fatal(err)
		return SubtaskResponse{}, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
		return SubtaskResponse{}, err
	}

	var metric SubtaskResponse
	err = json.Unmarshal(body, &metric)
	if err != nil {
		log.Fatal(err)
		return SubtaskResponse{}, err
	}

	return metric, nil
}

func (c *Client) GetPlan(ctx context.Context, url string) (*JobPlan, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", c.host, url), nil)
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

	var planResp JobPlanResponse
	err = json.Unmarshal(body, &planResp)
	if err != nil {
		log.Fatal(err)
	}

	return &planResp.Plan, nil
}
