package main

// NewsPoint represents a single point in the 10-point news summary
type NewsPoint struct {
	Text        string `json:"text" dynamodbav:"text"`
	Description string `json:"description" dynamodbav:"description"`
	URL         string `json:"url" dynamodbav:"url"`
	Source      string `json:"source" dynamodbav:"source"`
	PublishedAt string `json:"publishedAt" dynamodbav:"publishedAt"`
}

// NewsItem represents a news item in the database
type NewsItem struct {
	Category  string      `json:"category" dynamodbav:"category"`   // Partition Key
	Timestamp int64       `json:"timestamp" dynamodbav:"timestamp"` // Sort Key
	NewsID    string      `json:"newsId" dynamodbav:"newsId"`
	Title     string      `json:"title" dynamodbav:"title"`
	Points    []NewsPoint `json:"points" dynamodbav:"points"`
	FetchedAt int64       `json:"fetchedAt" dynamodbav:"fetchedAt"`
	TTL       int64       `json:"ttl" dynamodbav:"ttl"`
}
