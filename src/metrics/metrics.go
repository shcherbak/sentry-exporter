package metrics

import (
	"context"
	"fmt"
	"sentry-exporter/src/handler"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
)

var (
	ID                     = "id"
	IssueId                = "issue_id"
	Level                  = "level"
	Status                 = "status"
	IssueType              = "issue_type"
	Priority               = "priority"
	ProjectSlug            = "project_slug"
	sentryIssueMetricValue = float64(-1)
	sentryIssueMetricDesc  = prometheus.NewDesc(
		"sentry_issue_events",
		"count of sentry issues",
		[]string{
			ID,
			IssueId,
			Level,
			Status,
			IssueType,
			Priority,
			ProjectSlug,
		}, nil,
	)
)

type SentryMetricsCollector struct {
	SentryIssues *handler.SentryIssues
	Context      context.Context
	Client       *redis.Client
}

func (cc SentryMetricsCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(cc, ch)
}

func (cc SentryMetricsCollector) Collect(ch chan<- prometheus.Metric) {
	issues, _ := cc.SentryIssues.RedisScan(cc.Context, cc.Client)
	for id, issue := range issues.Items {
		f, err := strconv.ParseFloat(issue.Count, 64)
		if err == nil {
			sentryIssueMetricValue = f
		}
		ch <- prometheus.MustNewConstMetric(
			sentryIssueMetricDesc,
			prometheus.CounterValue,
			sentryIssueMetricValue,
			fmt.Sprint(id),
			fmt.Sprint(issue.IssueId),
			fmt.Sprint(issue.Level),
			fmt.Sprint(issue.Status),
			fmt.Sprint(issue.IssueType),
			fmt.Sprint(issue.Priority),
			fmt.Sprint(issue.Project.Slug),
		)

	}
}
