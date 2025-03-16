package main

type News struct {
	Category  string `json:"category"`
	CreatedAt string `json:"created_at"`
	Data      []Data `json:"data"`
}

type Data struct {
	Title         string `json:"title"`
	Link          string `json:"link"`
	Source        string `json:"source"`
	PublishedData string `json:"published_data"`
	ImageUrl      string `json:"image_url,omitempty"`
}
