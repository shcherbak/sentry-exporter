package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sentry-exporter/src/config"
	"sentry-exporter/src/handler"
	"sentry-exporter/src/metrics"
	"strconv"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
)

var (
	server           *http.Server
	shutdownWaiter1  sync.WaitGroup
	shutdownWaiter2  sync.WaitGroup
	issuesMiddleware *handler.SentryIssuesMiddleware
	buildVersion     string
)

func main() {

	ctx := context.Background()

	db, err := strconv.Atoi(config.GetConfig().REDIS_DBNO)
	if err != nil {
		log.Fatalf("Can't configure REDIS_DBNO: %v", err)
	}

	con := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", config.GetConfig().REDIS_ADDR, config.GetConfig().REDIS_PORT),
		DB:   db,
	})

	ttlSeconds, err := strconv.Atoi(config.GetConfig().TTL_SECONDS)
	if err != nil {
		log.Fatalf("Can't configure TTL_SECONDS: %v", err)
	}

	issuesMiddleware = &handler.SentryIssuesMiddleware{
		SentryClient: handler.NewSentryClient(
			config.GetConfig().AUTH_TOKEN,
			config.GetConfig().BASE_URL,
		),
		Context: ctx,
		Client:  con,
		TTL:     ttlSeconds,
	}

	grMax, err := strconv.Atoi(config.GetConfig().ROUTINE_MAX)
	if err != nil {
		log.Fatalf("Can't configure ROUTINE_MAX: %v", err)
	}
	ch := make(chan int, grMax)

	sleepSeconds, err := strconv.Atoi(config.GetConfig().SLEEP_SEC)
	if err != nil {
		log.Fatalf("Can't configure SLEEP_SEC: %v", err)
	}

	issuesCollector := metrics.SentryMetricsCollector{
		SentryIssues: &handler.SentryIssues{},
		Context:      ctx,
		Client:       con,
	}

	prometheus.MustRegister(issuesCollector)

	log.Printf("Creating server on port %s.", config.GetConfig().LISTEN_PORT)
	server = &http.Server{
		Addr:    fmt.Sprintf(":%s", config.GetConfig().LISTEN_PORT),
		Handler: nil,
	}
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/healthcheck", issuesMiddleware.HealthCheck)

	shutdownWaiter1.Add(1)
	ch <- 1
	go func() {
		defer func() { shutdownWaiter1.Done(); <-ch }()
		log.Printf("Starting http server. Build [%s]", buildVersion)
		server.ListenAndServe()
	}()

	for {
		projects, err := issuesMiddleware.SentryClient.GetProjects(config.GetConfig().ORGANIZATION_SLUG)
		if err != nil {
			log.Fatalf("Can't list projects: %v", err)
		}

		for i, p := range projects {
			fmt.Println(i, p.Slug)
		}

		fmt.Println("START main")
		for i, p := range projects {
			i := i
			p := p
			shutdownWaiter2.Add(1)
			ch <- 1
			go func() {
				defer func() { shutdownWaiter2.Done(); <-ch }()
				delay := time.Millisecond * time.Duration(50*(i+1))
				fmt.Printf("  START #%d (delay %v) %s\n", i, delay, p)
				issuesMiddleware.ImportIssueFromApiToRedis(config.GetConfig().ORGANIZATION_SLUG, p.Slug)
				time.Sleep(delay)
				fmt.Printf("  END #%d (delay %v) %s\n", i, delay, p)
			}()
		}

		shutdownWaiter2.Wait()
		fmt.Println("END main")
		time.Sleep(time.Duration(sleepSeconds) * time.Second)
	}

}
