package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/mmcdole/gofeed"
)

// TTL duration in seconds (25 Hours)
const TTLDuration = 25 * 60 * 60

func main() {
	if os.Getenv("AWS_LAMBDA_RUNTIME_API") != "" {
		// Running in AWS Lambda
		lambda.Start(handler)
	} else {
		// Running locally
		fmt.Println("Running locally...")
		err := handler(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
	}
}

func handler(ctx context.Context) error {
	// Initialize AWS DynamoDB client
	InitDB()

	var wg sync.WaitGroup
	var mu sync.Mutex
	newsMap := make(map[string][]NewsItem)

	for category, urls := range NewsFeeds {
		for _, url := range urls {
			wg.Add(1)
			go fetchNews(url, category, &wg, &mu, newsMap)
		}
	}

	wg.Wait()

	// Loop over the newsMap and insert into the database
	for category, items := range newsMap {
		for _, item := range items {
			err := InsertNews(item)
			if err != nil {
				fmt.Printf("Error inserting news item '%s' into the database: %v\n", item.Title, err)
			}
		}
		fmt.Printf("Processed %d items for category: %s\n", len(items), category)
	}

	return nil
}

func fetchNews(url string, category string, wg *sync.WaitGroup, mu *sync.Mutex, newsMap map[string][]NewsItem) {
	defer wg.Done()
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		fmt.Println("Error fetching news from URL:", url, err)
		return
	}

	now := time.Now()
	nowUnix := now.Unix()

	// Process each feed item
	for _, item := range feed.Items {
		// Create a unique ID based on the article URL
		newsID := generateNewsID(item.Link)

		// Parse the published date
		publishedAt := item.Published
		var timestamp int64

		// Try to parse the published date, fallback to current time if parsing fails
		if pubTime, err := time.Parse(time.RFC1123Z, publishedAt); err == nil {
			timestamp = pubTime.Unix() * 1000 // Convert to milliseconds
		} else {
			timestamp = nowUnix * 1000 // Use current time in milliseconds
		}

		// Create a single point (in a real app, you'd generate 10 points)
		point := NewsPoint{
			Text:        item.Title, // In a real app, this would be a summary point
			Description: item.Description,
			URL:         item.Link,
			Source:      feed.Title,
			PublishedAt: publishedAt,
		}

		// Create the news item
		newsItem := NewsItem{
			Category:  category,
			Timestamp: timestamp,
			NewsID:    newsID,
			Title:     item.Title,
			Points:    []NewsPoint{point}, // In a real app, this would have 10 points
			FetchedAt: nowUnix,
			TTL:       nowUnix + TTLDuration,
		}

		// Add to the map
		mu.Lock()
		newsMap[category] = append(newsMap[category], newsItem)
		mu.Unlock()
	}
}

// Generate a unique ID based on the article URL
func generateNewsID(url string) string {
	hash := md5.Sum([]byte(url))
	return hex.EncodeToString(hash[:])
}

// InsertNews inserts a news item into DynamoDB
func InsertNews(news NewsItem) error {
	// Check if the item already exists
	existingItem, err := getNewsItemByID(news.Category, news.NewsID)
	if err != nil && err.Error() != "NewsItem not found" {
		log.Printf("Error checking for existing news item: %v\n", err)
		return err
	}

	// If the item exists, we'll update it
	if existingItem != nil {
		log.Printf("News item with ID %s already exists, updating...\n", news.NewsID)
		// In a real app, you might want to merge or update fields
		// For simplicity, we'll just replace the entire item
	}

	// Marshal the news into DynamoDB map
	item, err := attributevalue.MarshalMap(news)
	if err != nil {
		log.Printf("Error marshaling news item: %v\n", err)
		return err
	}

	// Insert/update the news into the database
	_, err = DBClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String("Briefly-News"),
		Item:      item,
	})

	if err != nil {
		log.Printf("Couldn't add news item to Briefly-News table. Error: %v\n", err)
		return err
	}

	log.Printf("✅ Successfully inserted/updated news item '%s' into Briefly-News\n", news.Title)
	return nil
}

func getNewsItemByID(category, newsID string) (*NewsItem, error) {
	// Create a secondary index to query by newsID
	result, err := DBClient.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String("Briefly-News"),
		IndexName:              aws.String("NewsID-Index"), // You'll need to create this GSI
		KeyConditionExpression: aws.String("category = :category AND newsId = :newsId"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":category": &types.AttributeValueMemberS{Value: category},
			":newsId":   &types.AttributeValueMemberS{Value: newsID},
		},
		Limit: aws.Int32(1),
	})

	if err != nil {
		return nil, err
	}

	if len(result.Items) == 0 {
		return nil, fmt.Errorf("NewsItem not found")
	}

	var newsItem NewsItem
	err = attributevalue.UnmarshalMap(result.Items[0], &newsItem)
	if err != nil {
		return nil, err
	}

	return &newsItem, nil
}
