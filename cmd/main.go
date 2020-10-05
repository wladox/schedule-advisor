package main

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	pb "github.com/wladox/schedule-advisor/internal/grpc"
	"github.com/wladox/schedule-advisor/pkg/app"
	"github.com/wladox/schedule-advisor/pkg/config"
	"github.com/wladox/schedule-advisor/pkg/flink"
	"github.com/wladox/schedule-advisor/pkg/forecast"
	"github.com/wladox/schedule-advisor/pkg/fs"
	"github.com/wladox/schedule-advisor/pkg/logger"
	"github.com/wladox/schedule-advisor/pkg/metrics"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	// load environment configuration
	var env config.Config
	err := env.LoadFromENV()
	if err != nil {
		panic(errors.Wrap(err, "can't load config from env"))
	}

	// setup logging
	log, err := logger.New(&env)
	if err != nil {
		panic(fmt.Sprintf("can't create logger. %s", err))
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	httpClient := defaultClient()

	// create client for Flink API
	flinkClient := flink.NewClient(httpClient, env.FlinkMetricsAPIUrl)

	// create client for prediction service
	predictionClient := forecast.NewClient(httpClient, env.PredictionServiceUrl)

	// create metrics exporter
	exporter := &fs.FileWriter{}

	// create storage
	storage := metrics.Storage{
		IngestionRate:      0,
		Bandwidth:          make(map[string]int, 0),
		OperatorStats:      make(map[string]int, 0),
		CheckpointInterval: 30,
	}

	metricsCh := make(chan flink.CheckpointConfig)

	// create new metric collector
	collector := metrics.NewMetricCollector(flinkClient, exporter, log, metricsCh, storage)

	calculator := forecast.NewCalculator(metricsCh, predictionClient, log, storage)
	calculator.Init()

	done := make(chan bool)

	go func() {
		ticker := time.NewTicker(time.Duration(env.Interval) * time.Millisecond)
		for {
			select {
			case <-done:
				ticker.Stop()
				log.Info("Timer stopped")
				return
			case <-ticker.C:

				err := collector.CollectDataflowInfo(ctx)
				if err != nil {
					log.WithError(err).Error("Metrics collection failed.")
				}
				log.Info("Dataflow metrics successfully collected.")

				err = collector.CollectNetworkInfo(ctx)
				if err != nil {
					log.WithError(err).Error("Metrics collection failed.")
				}
				log.Info("Network metrics successfully collected.")
			}
		}
	}()

	lis, err := net.Listen("tcp", ":5050")
	if err != nil {
		log.WithError(err).Fatal("failed to listen: ")
	}
	s := grpc.NewServer()

	registerShutdownHook(s, done, cancel)

	grpcServer := &app.GrpcServer{
		Calculator: calculator,
	}

	pb.RegisterSchedulerAdvisorServer(s, grpcServer)
	if err := s.Serve(lis); err != nil {
		log.WithError(err).Fatal("failed to serve: ")
	}
}

func registerShutdownHook(server *grpc.Server, done chan<- bool, cancelFunc context.CancelFunc) {
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		s := <-exit
		log.WithField("signal", s.String()).Info("received signal, shutting down ...")
		cancelFunc()

		log.Info("Shutting down cron job.")
		done <- true

		log.Info("Shutting down gRPC server.")
		server.Stop()

		os.Exit(0)
	}()
}

func defaultClient() *http.Client {
	return &http.Client{
		Timeout: 5 * time.Second,
	}
}
