package handler

import (
	"context"
	"io"
	"log"
	"net/http"

	"github.com/redis/go-redis/v9"
)

type SentryIssuesMiddleware struct {
	SentryClient *SentryClient
	Context      context.Context
	Client       *redis.Client
	TTL          int
}

func (c *SentryIssuesMiddleware) ImportIssueFromApiToRedis(organization string, project string) {
	sentryIssues, err := c.SentryClient.GetIssues(organization, project)
	if err != nil {
		log.Fatalf("Can't set sentryIssues: %v", err)
	}
	for _, i := range sentryIssues {
		i.RedisInsert(i, c.Context, c.Client, c.TTL)
	}
}

func (s *SentryIssuesMiddleware) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "OK")
}
