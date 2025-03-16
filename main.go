package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/mmcdole/gofeed"
)

func main() {
	// Initialize DynamoDB connection
	InitDB()

	// Use DBClient from the config
	fmt.Println("DynamoDB client is ready to use:", DBClient)

	var wg sync.WaitGroup
	var mu sync.Mutex
	var news []News

	for category, urls := range NewsFeeds {
		for _, url := range urls {
			wg.Add(1)
			go fetchNews(url, category, &wg, &mu, &news)
		}
	}

	wg.Wait()

	// Convert to JSON
	newsJSON, err := json.MarshalIndent(news, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling news to JSON:", err)
		return
	}

	fmt.Println(string(newsJSON))

	// Insert news into DynamoDB
	for _, n := range news {
		err := insertNewsIntoDynamoDB(n)
		if err != nil {
			fmt.Println("Error inserting news into DynamoDB:", err)
		}
	}
}

func fetchNews(url string, category string, wg *sync.WaitGroup, mu *sync.Mutex, news *[]News) {
	defer wg.Done()
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		fmt.Println("Error fetching news from URL:", url, err)
		return
	}

	var data []Data
	for _, item := range feed.Items {
		newItem := Data{
			Title:         item.Description,
			Link:          item.Link,
			Source:        feed.Title,
			PublishedData: item.Published,
		}

		// Check if the image is present
		if len(item.Enclosures) > 0 {
			if item.Image != nil && item.Image.URL != "" {
				newItem.ImageUrl = item.Image.URL
			} else if len(item.Enclosures) > 0 && item.Enclosures[0].URL != "" {
				newItem.ImageUrl = item.Enclosures[0].URL
			} else if item.Extensions != nil {
				for _, ext := range item.Extensions {
					if ext["image"] != nil && len(ext["image"]) > 0 {
						if ext["image"][0].Attrs["url"] != "" {
							newItem.ImageUrl = ext["image"][0].Attrs["url"]
							break
						}
					}
				}
			}
		}
		data = append(data, newItem)
	}

	mu.Lock()
	*news = append(*news, News{
		Category:  category,
		CreatedAt: time.Now().Format(time.RFC3339),
		Data:      data,
	})
	mu.Unlock()
}

func insertNewsIntoDynamoDB(news News) error {
	item, err := attributevalue.MarshalMap(news)
	if err != nil {
		return fmt.Errorf("failed to marshal news: %w", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String("Briefly-News"),
		Item:      item,
	}

	_, err = DBClient.PutItem(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to put item into DynamoDB: %w", err)
	}

	return nil
}
