package slack

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/slack-go/slack"
)

type SlackClient struct {
	client      *slack.Client
	user        *slack.User
	displayName string
}

type StdLibClient struct {
	apiToken string
	client   *http.Client
}

func NewSlackClient(apiToken string, userId string, displayName string) (*SlackClient, error) {
	client := slack.New(apiToken, slack.OptionDebug(true))
	user, err := client.GetUserInfo(userId)
	if err != nil {
		return nil, err
	}
	return &SlackClient{client: client, user: user, displayName: displayName}, nil
}

func NewStdLibClient(apiToken string) *StdLibClient {
	return &StdLibClient{
		apiToken: apiToken,
		client:   &http.Client{},
	}
}

func (c *StdLibClient) UpdateDisplayName(userID, displayName string) error {
	// https://api.slack.com/methods/users.profile.set#profile_updates

	data := url.Values{
		"profile": {`{"display_name": "` + displayName + `"}`},
	}

	req, err := http.NewRequest("POST", "https://slack.com/api/users.profile.set", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+c.apiToken)

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (c *SlackClient) UpdateStatus(statusText, statusEmoji string) error {
	fmt.Printf("Slack client UpdateStatus called with statusText=%s, statusEmoji=%s\n", statusText, statusEmoji)

	// Update the real name (display name) of the user
	err := c.client.SetUserRealNameContextWithUser(context.Background(), c.user.ID, c.displayName)
	if err != nil {
		fmt.Printf("Error setting real name: %v\n", err)
		return err
	}

	// Update the custom status
	err = c.client.SetUserCustomStatus(statusText, statusEmoji, 0)
	if err == nil {
		fmt.Printf("Updated status for %s (%s, %s): %s %s\n", c.user.Profile.RealName, c.user.ID, c.user.Profile.Email, statusEmoji, statusText)
	} else {
		fmt.Printf("Error setting custom status: %v\n", err)
	}
	return err
}

func (c *SlackClient) UnsetStatus() error {
	err := c.client.UnsetUserCustomStatus()
	if err != nil {
		fmt.Printf("Error unsetting custom status: %v\n", err)
	} else {
		fmt.Printf("Unset custom status for %s (%s, %s)\n", c.user.Profile.RealName, c.user.ID, c.user.Profile.Email)
	}
	return err
}
