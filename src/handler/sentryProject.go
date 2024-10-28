package handler

type SentryProject struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	Platform string `json:"platform"`
}
