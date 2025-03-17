package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/mmcdole/gofeed"
)

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
	newsMap := make(map[string]*News)

	for category, urls := range NewsFeeds {
		for _, url := range urls {
			wg.Add(1)
			go fetchNews(url, category, &wg, &mu, newsMap)
		}
	}

	wg.Wait()

	// Loop over the newsMap and insert into the database
	for _, n := range newsMap {
		err := InsertNews(*n)
		if err != nil {
			fmt.Println("Error inserting news into the database:", err)
			return err
		}
	}

	return nil
}

func fetchNews(url string, category string, wg *sync.WaitGroup, mu *sync.Mutex, newsMap map[string]*News) {
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
			Title:         item.Title,
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
	if existingNews, exists := newsMap[category]; exists {
		existingNews.Data = append(existingNews.Data, data...)
	} else {
		newsMap[category] = &News{
			Category:  category,
			CreatedAt: "latest",
			Data:      data,
		}
	}
	mu.Unlock()
}

func InsertNews(news News) error {
	// Marshal the news into JSON
	item, err := attributevalue.MarshalMap(news)
	if err != nil {
		log.Printf("Error marshaling news item: %v\n", err)
		return err
	}

	// Insert the news into the database
	_, err = DBClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String("Briefly-News"),
		Item:      item,
	})

	if err != nil {
		log.Printf("Couldn't add news item to Briefly-News table. Error: %v\n", err)
		return err
	}

	log.Println("✅ Successfully inserted news item into Briefly-News")
	return nil
}
