package discord

import (
	"context"
	"fmt"
	"net/http"
)

// GetUser returns a user object for the given user ID.
func (c *Client) GetUser(ctx context.Context, userID Snowflake) (*User, error) {
	user := new(User)
	route := fmt.Sprintf("/users/%s", userID)
	err := c.doRequest(ctx, http.MethodGet, route, nil, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetCurrentUser returns the user object of the requester's account.
func (c *Client) GetCurrentUser(ctx context.Context) (*User, error) {
	user := new(User)
	err := c.doRequest(ctx, http.MethodGet, "/users/@me", nil, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
