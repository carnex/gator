package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/carnex/gator/internal/database"
	"github.com/google/uuid"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedurl string) (*RSSFeed, error) {
	output := RSSFeed{}
	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "GET", feedurl, nil)
	if err != nil {
		return &output, err
	}
	req.Header.Set("User-Agent", "gator")
	resp, err := client.Do(req)
	if err != nil {
		return &output, err
	}
	defer resp.Body.Close()
	reader, err := io.ReadAll(resp.Body)
	if err != nil {
		return &output, err
	}
	xml.Unmarshal(reader, &output)
	output.Channel.Title = html.UnescapeString(output.Channel.Title)
	output.Channel.Description = html.UnescapeString(output.Channel.Description)
	for i, item := range output.Channel.Item {
		output.Channel.Item[i].Title = html.UnescapeString(item.Title)
		output.Channel.Item[i].Description = html.UnescapeString(item.Description)
	}
	return &output, err
}

func scrapeFeeds(s *state) error {
	ctx := context.Background()
	nextFeed, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		return err
	}
	err = s.db.MarkFeedFetched(ctx, database.MarkFeedFetchedParams{ID: nextFeed.ID, LastFetchedAt: sql.NullTime{Time: time.Now(), Valid: true}, UpdatedAt: time.Now()})
	if err != nil {
		return err
	}
	feed, err := fetchFeed(ctx, nextFeed.Url)
	if err != nil {
		return err
	}
	for _, item := range feed.Channel.Item {
		pubtime, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			return err
		}
		post, err := s.db.CreatePosts(ctx, database.CreatePostsParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Title: item.Title, Url: item.Link, Description: item.Description, PublishedAt: pubtime, FeedID: nextFeed.ID})
		if err != nil {
			return err
		}
		fmt.Printf("Succesfully save %s\n", post.Title)
	}

	return nil
}
