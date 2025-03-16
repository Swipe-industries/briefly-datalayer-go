package main

type News struct {
	Category  string `json:"category" dynamodbav:"category"`     // Partition Key
	CreatedAt string `json:"created_at" dynamodbav:"created_at"` // Sort Key
	Data      []Data `json:"data" dynamodbav:"data"`             // News articles
}

type Data struct {
	Title         string `json:"title"`
	Link          string `json:"link"`
	Source        string `json:"source"`
	PublishedData string `json:"published_data"`
	ImageUrl      string `json:"image_url,omitempty"`
}
