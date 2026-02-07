package discord

import (
	"context"
	"net/http"
)

// EditApplicationParams are the parameters for editing the current application.
type EditApplicationParams struct {
	CustomInstallURL               *string  `json:"custom_install_url,omitempty"`
	Description                    *string  `json:"description,omitempty"`
	RoleConnectionsVerificationURL *string  `json:"role_connections_verification_url,omitempty"`
	InteractionsEndpointURL        *string  `json:"interactions_endpoint_url,omitempty"`
	Flags                          *int     `json:"flags,omitempty"`
	Icon                           *string  `json:"icon,omitempty"`
	CoverImage                     *string  `json:"cover_image,omitempty"`
	Tags                           []string `json:"tags,omitempty"`
}

// GetCurrentApplication returns the current application object.
func (c *Client) GetCurrentApplication(ctx context.Context) (*Application, error) {
	app := new(Application)
	err := c.doRequest(ctx, http.MethodGet, "/applications/@me", nil, app)
	if err != nil {
		return nil, err
	}
	return app, nil
}

// EditCurrentApplication edits properties of the current application.
func (c *Client) EditCurrentApplication(ctx context.Context, params *EditApplicationParams) (*Application, error) {
	app := new(Application)
	err := c.doRequest(ctx, http.MethodPatch, "/applications/@me", params, app)
	if err != nil {
		return nil, err
	}
	return app, nil
}
