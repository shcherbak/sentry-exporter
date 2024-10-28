package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var (
	ErrInsertFailed = errors.New("insert failed")
	ErrMarshlFailed = errors.New("marshal failed")
	ErrIDNotFound   = errors.New("id not found")
	ErrNoKeysFound  = errors.New("no keys found")
)

func (c *SentryIssue) RedisInsert(issue SentryIssue, ctx context.Context, con *redis.Client, ttl int) (uuid.UUID, error) {
	issue.setUUIDKey()
	issue.setUpdatedAt()
	srt, err := json.Marshal(issue)
	if err != nil {
		return issue.ID, ErrMarshlFailed
	}
	result, err := con.Set(ctx, fmt.Sprint(issue.ID), srt, time.Duration(ttl)*time.Second).Result()
	if err != nil {
		return issue.ID, ErrInsertFailed
	}
	log.Printf("Inserting %s: %s, %s, %s", issue.ID, issue.IssueId, issue.Project.Slug, result)
	return issue.ID, nil
}

func (c *SentryIssues) RedisRetrieve(id uuid.UUID, ctx context.Context, con *redis.Client) (SentryIssue, error) {
	result, err := con.Get(ctx, fmt.Sprint(id)).Result()
	if err != nil {
		return SentryIssue{}, ErrIDNotFound
	}
	issue := SentryIssue{}
	json.Unmarshal([]byte(result), &issue)
	return issue, nil
}

func (c *SentryIssues) RedisScan(ctx context.Context, con *redis.Client) (SentryIssues, error) {
	var issues = SentryIssues{
		Items: make(map[uuid.UUID]SentryIssue),
	}
	iter := con.Scan(ctx, 0, "*", 0).Iterator()
	for iter.Next(ctx) {
		uid, err := uuid.Parse(iter.Val())
		if err != nil {
			log.Printf("Can't parse uuid: %s, %v", iter.Val(), err)
			continue
		} else {
			issues.Items[uid], _ = c.RedisRetrieve(uid, ctx, con)
		}
	}
	if err := iter.Err(); err != nil {
		panic(err)
	}
	return issues, nil
}
