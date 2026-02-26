package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/carnex/gator/internal/database"
	"github.com/google/uuid"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		loggerUser, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return err
		}
		return handler(s, cmd, loggerUser)
	}

}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}
	name := cmd.Args[0]
	ctx := context.Background()
	_, err := s.db.GetUser(ctx, name)
	if err != nil {
		return err
	}
	errr := s.cfg.SetUser(name)
	if errr != nil {
		return fmt.Errorf("couldn't set current user: %w", errr)
	}

	fmt.Println("User switched successfully!")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}
	name := cmd.Args[0]
	ctx := context.Background()

	insertedUser, err := s.db.CreateUser(ctx, database.CreateUserParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: name})
	if err != nil {
		return err
	}
	s.cfg.SetUser(insertedUser.Name)
	fmt.Printf("%s was created\n", insertedUser.Name)

	return nil
}

func handlerReset(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage %s <name>", cmd.Name)
	}
	ctx := context.Background()
	err := s.db.DeleteUsers(ctx)
	if err != nil {
		return fmt.Errorf("couldn't delete users: %w", err)
	}
	fmt.Printf("DB users deleted\n")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}
	ctx := context.Background()
	users, err := s.db.GetUsers(ctx)
	if err != nil {
		return fmt.Errorf("unable to retrive users %w", err)
	}
	for _, user := range users {
		if user.Name == s.cfg.CurrentUserName {
			fmt.Printf("%s (current)\n", user.Name)
		} else {
			fmt.Printf("%s\n", user.Name)
		}
	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}
	time_between_requests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return nil
	}
	ticker := time.NewTicker(time_between_requests)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}
	ctx := context.Background()
	insertedFeed, err := s.db.CreateFeed(ctx, database.CreateFeedParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: cmd.Args[0], Url: cmd.Args[1], UserID: user.ID})
	if err != nil {
		return nil
	}
	feed, err := s.db.GetFeed(ctx, cmd.Args[1])
	if err != nil {
		return err
	}
	s.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), UserID: user.ID, FeedID: feed.ID})
	fmt.Printf(" %s added to database under user %s\n", insertedFeed.Url, user.Name)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}
	ctx := context.Background()
	feeds, err := s.db.GetFeeds(ctx)
	if err != nil {
		return err
	}
	for _, feed := range feeds {
		fmt.Printf("%s\n", feed)
	}
	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}
	ctx := context.Background()
	feed, err := s.db.GetFeed(ctx, cmd.Args[0])
	if err != nil {
		return err
	}
	newfeeds, err := s.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), UserID: user.ID, FeedID: feed.ID})
	fmt.Printf("Feed: %s succesfully follow for user: %s\n", newfeeds.FeedName, newfeeds.UserName)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	ctx := context.Background()
	following, err := s.db.GetFeedsFollowsForUser(ctx, user.ID)
	if err != nil {
		return err
	}
	for _, follows := range following {
		fmt.Printf("Following: %s\n", follows.Name_2)
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}
	ctx := context.Background()
	feed, err := s.db.GetFeed(ctx, cmd.Args[0])
	if err != nil {
		return err
	}
	s.db.Unfollow(ctx, database.UnfollowParams{UserID: user.ID, FeedID: feed.ID})
	return nil
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := 2
	if len(cmd.Args) == 1 {
		parsed, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			return err
		}
		limit = parsed
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{UserID: user.ID, Limit: int32(limit)})
	if err != nil {
		return nil
	}
	for _, post := range posts {
		fmt.Printf("Title: %s Description %s Url: %s published at: %s\n", post.Title, post.Description, post.Url, post.PublishedAt)
	}
	return nil
}
