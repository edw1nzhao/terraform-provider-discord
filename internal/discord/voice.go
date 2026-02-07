package discord

import (
	"context"
	"net/http"
)

// ListVoiceRegions returns an array of voice region objects.
func (c *Client) ListVoiceRegions(ctx context.Context) ([]*VoiceRegion, error) {
	var regions []*VoiceRegion
	err := c.doRequest(ctx, http.MethodGet, "/voice/regions", nil, &regions)
	if err != nil {
		return nil, err
	}
	return regions, nil
}
