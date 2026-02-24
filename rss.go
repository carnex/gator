package main

import (
	"context"
	"encoding/xml"
	"html"
	"io"
	"net/http"
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
