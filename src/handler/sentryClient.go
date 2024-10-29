package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
)

type SentryClient struct {
	baseURL string
	token   string
}

func NewSentryClient(token string, baseURL string) *SentryClient {
	return &SentryClient{
		baseURL: baseURL,
		token:   token,
	}
}

func (c *SentryClient) getRequest(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch data: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (c *SentryClient) GetProjects(organizationSlug string, excludeFilter string, includeFiler string) ([]SentryProject, error) {
	url := fmt.Sprintf("%s/organizations/%s/projects/", c.baseURL, organizationSlug)
	body, err := c.getRequest(url)
	var projects []SentryProject
	if err != nil {
		return projects, err
	}

	if err := json.Unmarshal(body, &projects); err != nil {
		return projects, err
	}

	return projects, nil
}

func (c *SentryClient) IncludeProjectsByPattern(objects []SentryProject, includePattern string) ([]SentryProject, error) {
	includeRe, err := regexp.Compile(includePattern)
	if err != nil {
		return nil, fmt.Errorf("include regex compillation error: %v", err)
	}

	var filteredObjects []SentryProject
	for _, obj := range objects {
		if includeRe.MatchString(obj.Slug) {
			filteredObjects = append(filteredObjects, obj)
		}
	}
	return filteredObjects, nil
}

func (c *SentryClient) ExcludeProjectsByPattern(objects []SentryProject, excludePattern string) ([]SentryProject, error) {
	excludeRe, err := regexp.Compile(excludePattern)
	if err != nil {
		return nil, fmt.Errorf("exclude regex compillation error: %v", err)
	}

	var filteredObjects []SentryProject
	for _, obj := range objects {
		if !excludeRe.MatchString(obj.Slug) {
			filteredObjects = append(filteredObjects, obj)
		}
	}
	return filteredObjects, nil
}

func (c *SentryClient) GetIssues(organizationSlug string, projectSlug string) ([]SentryIssue, error) {
	url := fmt.Sprintf("%s/projects/%s/%s/issues/", c.baseURL, organizationSlug, projectSlug)
	body, err := c.getRequest(url)
	var issues []SentryIssue
	if err != nil {
		return issues, err
	}
	if err := json.Unmarshal(body, &issues); err != nil {
		return issues, err
	}

	return issues, nil
}
