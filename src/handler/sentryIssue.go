package handler

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type SentryIssue struct {
	ID        uuid.UUID     `json:"redis_key"`        // uuid generated form issue fields
	UpdatedAt string        `json:"redis_updated_at"` // unix timestamp
	IssueId   string        `json:"id"`
	Level     string        `json:"level"`
	Status    string        `json:"status"`
	IssueType string        `json:"issueType"`
	Priority  string        `json:"priority"`
	Count     string        `json:"count"`
	Project   SentryProject `json:"project"`
}

type SentryIssues struct {
	Items map[uuid.UUID]SentryIssue
}

func (c *SentryIssue) setUUIDKey() {
	c.ID = uuid.NewSHA1(uuid.NameSpaceDNS, []byte(fmt.Sprintf("%s-%s-%s", "sentry", c.Project.Slug, c.IssueId)))
}

func (c *SentryIssue) setUpdatedAt() {
	c.UpdatedAt = fmt.Sprint(time.Now().Unix())
}
