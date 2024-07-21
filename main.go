package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/vartanbeno/go-reddit/v2/reddit"
)

const FILENAME = "subreddits.txt"

func main() {
	godotenv.Load()
	client := getAuthedClient()

	posts, err := getTopPosts(client, "golang")
	if err != nil {
		fmt.Println("Error in getTopPosts: ", err)
		return
	}

	for _, post := range posts {
		fmt.Printf("Post: %s\nURL: %s\n\n", post.Title, post.URL)
	}

	// Uncomment to trigger subscriptions
	// err = subscribeFromTextFile(client, FILENAME)
	// if err != nil {
	// 	fmt.Println("Error in subscribeFromTextFile: ", err)
	// 	return
	// }

	subredditInfo("rust")
}

func getAuthedClient() (client *reddit.Client) {
	id := os.Getenv("REDDIT_ID")
	secret := os.Getenv("REDDIT_SECRET")
	username := os.Getenv("REDDIT_USERNAME")
	password := os.Getenv("REDDIT_PASSWORD")

	credentials := reddit.Credentials{ID: id, Secret: secret, Username: username, Password: password}
	client, _ = reddit.NewClient(credentials)

	fmt.Println("Client created successfully.")

	return client
}

func getTopPosts(client *reddit.Client, subreddit string) (posts []*reddit.Post, err error) {
	posts, _, err = client.Subreddit.TopPosts(context.Background(), subreddit, &reddit.ListPostOptions{
		ListOptions: reddit.ListOptions{
			Limit: 5,
		},
		Time: "all",
	})

	if err != nil {
		fmt.Println("Error in TopPosts: ", err)
		return
	}

	fmt.Printf("Received %d posts.\n", len(posts))
	return
}

func subscribeFromTextFile(client *reddit.Client, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		subreddit := strings.TrimSpace(scanner.Text())
		if subreddit != "" {
			err := subscribeToSubreddit(client, subreddit)
			if err != nil {
				fmt.Printf("Failed to subscribe to %s: %v\n", subreddit, err)
			} else {
				fmt.Printf("Subscribed to %s\n", subreddit)
				count++
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	fmt.Printf("Subscribed to %d subreddits\n", count)

	return nil
}

func subscribeToSubreddit(client *reddit.Client, subreddit string) (err error) {
	result, err := client.Subreddit.Subscribe(context.Background(), subreddit)

	if err != nil {
		fmt.Println("Error in Subscribe: ", err)
		return
	}

	fmt.Println("Successfully subscribed to subreddit: ", result)
	return
}

func subredditInfo(subreddit string) (err error) {
	sr, _, err := reddit.DefaultClient().Subreddit.Get(context.Background(), subreddit)
	if err != nil {
		fmt.Println("Error in Get: ", err)
		return
	}

	fmt.Printf("%s was created on %s and has %d subscribers.\n", sr.NamePrefixed, sr.Created.Local(), sr.Subscribers)
	return
}
