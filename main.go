package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/wladox/schedule-advisor/client"
	"github.com/wladox/schedule-advisor/config"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const (
	checkpointsUrl = "checkpoints"
	checkpointDetailsUrl = "checkpoints/details"
	checkpointDetailsSubtasksUrl = "subtasks"

	job_id = "3a095d9e884fcffafd803652ba1142e4"
	ingestionRateParams = "numBytesInPerSecond,numBytesOutPerSecond"

)

func main() {

	var env config.Config
	err := env.LoadFromENV()
	if err != nil {
		panic(errors.Wrap(err, "can't load config from env"))
	}

	flinkClient := client.NewClient(defaultClient(), env.FlinkMetricsAPIUrl)

	plan, err :=flinkClient.GetPlan(context.Background(), fmt.Sprintf("jobs/%s/plan", job_id))
	if err != nil {
		log.Fatal(err)
	}

	vertices := make(map[string]int, 0)

	for _, node := range plan.Nodes {
		vertices[node.Id] = node.Parallelism
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan bool)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		s := <-sig
		log.Println("received signal, shutting down ", s)
		cancel()
		log.Println("wait done. Shutting down http server.")
		done <- true
		os.Exit(0)
	}()

	ticker := time.NewTicker(5 * time.Second)
	var counter int
	for {
		select {
		case <-done:
			ticker.Stop()
			log.Println("Timer stopped")
			return
		case <-ticker.C:
			// collect ingestion rate metrics
			ingestionRateMetrics := make([]string, 0)
		 	counter += 5
			for taskId, parallelism := range vertices {
				for i := 0; i < parallelism; i++ {
					url := fmt.Sprintf("jobs/%s/vertices/%s/subtasks/%d/metrics?get=%s", job_id, taskId, i, ingestionRateParams)
					metrics, err := flinkClient.CollectVertexMetrics(ctx, url)
					if err != nil {
						log.Fatal(err)
					}
					for _, metric := range metrics {
						ingestionRateMetrics = append(ingestionRateMetrics, fmt.Sprintf("%v,%s,%d,%s,%s", counter, taskId, i, metric.Id, metric.Value))
					}
					fmt.Printf("vertexId: %s - subtaskIdx: %d - %s \n", taskId, i, metrics)
				}
			}
			err = writeMetrics(ingestionRateMetrics, "output/ingestion_rate.csv")
			ingestionRateMetrics = nil
			checkError("Could not write to file", err)

			// collect checkpoint metrics
			resp, _ := flinkClient.GetCheckpoints(ctx, fmt.Sprintf("jobs/%s/%s", job_id, checkpointsUrl))
			id := resp.Latest.Completed.Id
			metric, _ := flinkClient.GetCheckpointDetails(ctx, fmt.Sprintf("jobs/%s/%s/%d", job_id, checkpointDetailsUrl, id))

			operatorStateMetrics := make([]string, 0)
			for taskId, _ := range metric.Tasks {
				subtaskStats, _ := flinkClient.GetCheckpointSubtasks(ctx, fmt.Sprintf("jobs/%s/%s/%d/%s/%s", job_id, checkpointDetailsUrl, id,checkpointDetailsSubtasksUrl, taskId))
				for _, subtask := range subtaskStats.Subtasks {
					operatorStateMetrics = append(operatorStateMetrics, fmt.Sprintf("%d,%s,%d,%d", counter, taskId, subtask.Index, subtask.StateSize))
				}
			}
			err = writeMetrics(operatorStateMetrics, "output/operator_state.csv")
			operatorStateMetrics = nil
			checkError("Could not write to file", err)
		}
	}

}



// writeLines writes the lines to the given file.
func writeMetrics(lines []string, path string) error {

	s := strings.Split(path, "/")
	if len(s) > 0 {
		if _, err := os.Stat(s[0]); os.IsNotExist(err) {
			_ = os.Mkdir(s[0], 0766)
		}
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		checkError("Cannot open file", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for _, line := range lines {
		_, err = writer.WriteString(line + "\n")
		checkError("Cannot write to file", err)
	}
	return nil
}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}

func defaultClient() *http.Client {
	return &http.Client{
		Timeout: 30 * time.Second,
	}
}