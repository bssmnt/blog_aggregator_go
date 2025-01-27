package commands

import (
	"blog_aggregator_go/internal/config"
	"blog_aggregator_go/internal/database"
	"blog_aggregator_go/internal/rss"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"time"
)

type State struct {
	Db  *database.Queries
	Cfg *config.Config
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	CommandNames map[string]func(*State, Command) error
}

func MiddlewareLoggedIn(handler interface{}) func(*State, Command) error {
	return func(s *State, cmd Command) error {
		if s.Cfg.CurrentUserName == "" {
			return errors.New("no user is currently logged in")
		}

		user, err := s.Db.GetUser(context.Background(), s.Cfg.CurrentUserName)
		if err != nil {
			return err
		}

		switch h := handler.(type) {
		case func(s *State, cmd Command) error:
			return h(s, cmd)
		case func(s *State, cmd Command, user database.User) error:
			return h(s, cmd, user)
		default:
			return fmt.Errorf("invalid handler type: %T", h)
		}
	}
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.CommandNames[name] = f
}

func (c *Commands) Run(s *State, cmd Command) error {
	if handler, exists := c.CommandNames[cmd.Name]; exists {
		return handler(s, cmd)
	} else {
		return errors.New("command not found")
	}
}

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return errors.New("please specify a username: login <username>")
	}

	username := cmd.Args[0]

	_, err := s.Db.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("user %s not found", username)
	}

	err = s.Cfg.SetUser(username)
	if err != nil {
		return err
	}
	fmt.Println("user set to:", username)
	return nil
}

func HandlerRegister(s *State, cmd Command) error {

	if len(cmd.Args) != 1 {
		return errors.New("please specify a username: register <username>")
	}
	username := cmd.Args[0]
	ctx := context.Background()
	params := database.CreateUserParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: username}

	user, err := s.Db.CreateUser(ctx, params)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" {
				return fmt.Errorf("user %s already exists", username)
			}
		}
		return err
	}

	err = s.Cfg.SetUser(username)
	if err != nil {
		return err
	}

	err = config.Save(*s.Cfg)
	if err != nil {
		return err
	}

	fmt.Println("user set to:", user)

	return nil
}

func HandlerReset(s *State, cmd Command) error {
	err := s.Db.DeleteAllData(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func HandlerUsers(s *State, cmd Command) error {
	allUsers, err := s.Db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	if len(allUsers) == 0 {
		fmt.Println("no users found")
	}

	currentUser := s.Cfg.CurrentUserName
	for _, user := range allUsers {
		if user == currentUser {
			fmt.Printf("* %s (current)\n", currentUser)
		} else {
			fmt.Printf("* %s\n", user)
		}
	}
	return nil
}

func HandlerAgg(s *State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("expected 1 argument (time_between_reqs), got %d", len(cmd.Args))
	}

	timeBetweenReqs, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("invalid duration format: %v", err)
	}

	fmt.Printf("Collecting feeds every %v\n", timeBetweenReqs)

	ticker := time.NewTicker(timeBetweenReqs)
	for ; ; <-ticker.C {
		if err := HandlerScrapeFeeds(s, cmd); err != nil {
			fmt.Printf("Error scraping feeds: %v\n", err)
			continue
		}
	}
}

func HandlerAddFeed(s *State, cmd Command, user database.User) error {

	if len(cmd.Args) != 2 {
		return errors.New("please specify a feed name and url: add <feed name> <url>")
	}

	ctx := context.Background()

	currentUserID := user.ID

	feedName := cmd.Args[0]
	feedURL := cmd.Args[1]

	params := database.CreateFeedParams{ID: uuid.New(), Name: feedName, Url: feedURL, UserID: currentUserID}
	err := s.Db.CreateFeed(ctx, params)
	if err != nil {
		return err
	}

	feedID := params.ID
	feedUserID := params.UserID

	now := time.Now()
	_, err = s.Db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    currentUserID,
		FeedID:    feedID,
	})

	fmt.Printf("%s added, from URL %s, feed id: %v, added by %s\n", feedName, feedURL, feedID, feedUserID)

	return nil
}

func HandlerFeeds(s *State, cmd Command) error {
	ctx := context.Background()

	feeds, err := s.Db.GetFeeds(ctx)
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		fmt.Printf("%s added, from URL %s, added by %s\n", feed.FeedName, feed.FeedUrl, feed.UsersName)
	}
	return nil
}

func HandlerFollow(s *State, cmd Command, user database.User) error {
	ctx := context.Background()

	currentUserID := user.ID

	if len(cmd.Args) != 1 {
		return errors.New("please provide url")
	}

	url := cmd.Args[0]

	feed, err := s.Db.GetFeedByURL(ctx, url)
	if err != nil {
		return err
	}

	now := time.Now()

	feedFollow, err := s.Db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    currentUserID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return err
	}

	fmt.Printf("following '%v' as '%v'\n", feedFollow.FeedName, feedFollow.UserName)

	return nil
}

func HandlerFollowing(s *State, cmd Command, user database.User) error {
	ctx := context.Background()

	feedFollows, err := s.Db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		return err
	}

	fmt.Println("currently following:")
	for _, follow := range feedFollows {
		fmt.Printf("- %s\n", follow.FeedName)
	}

	return nil
}

func HandlerUnfollow(s *State, cmd Command, user database.User) error {
	ctx := context.Background()

	if len(cmd.Args) != 1 {
		return errors.New("please provide valid syntax: unfollow <feed name>")
	}

	feedURL := cmd.Args[0]

	feed, err := s.Db.GetFeedByURL(ctx, feedURL)
	if err != nil {
		return errors.New("feed not found, please check the url")
	}

	err = s.Db.UnfollowFeed(ctx, database.UnfollowFeedParams{
		feed.ID,
		user.ID,
	})

	if err != nil {
		return errors.New("feed unfollow failed, are you definitely following this feed")
	}

	return nil
}

func HandlerScrapeFeeds(s *State, cmd Command) error {
	ctx := context.Background()

	feed, err := s.Db.GetNextFeedToFetch(ctx)
	if err != nil {
		return err
	}

	err = s.Db.MarkFeedFetched(ctx, feed.ID)
	if err != nil {
		return err
	}

	rssFeed, err := rss.FetchFeed(ctx, feed.Url)
	if err != nil {
		return err
	}

	fmt.Printf("\nFetching from: %s\n", rssFeed.Channel.Title)
	for _, item := range rssFeed.Channel.Item {
		pubDate, err := time.Parse("Mon, 02 Jan 2006 15:04:05 -0700", item.PubDate)
		if err != nil {
			return err
		}

		_, err = s.Db.GetPostByURL(ctx, database.GetPostByURLParams{
			Url:    item.Link,
			FeedID: feed.ID,
		})
		if err == nil {
			continue
		}

		fmt.Printf("[%s] %s\n",
			time.Now().Format("15:04:05"),
			item.Title)

		params := database.InsertPostParams{
			ID:          uuid.New(),
			Title:       item.Title,
			Url:         item.Link,
			FeedID:      feed.ID,
			PublishedAt: pubDate,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		err = s.Db.InsertPost(ctx, params)
		if err != nil {
			return err
		}

	}

	return nil
}
