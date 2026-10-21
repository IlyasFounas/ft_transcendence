package mapillary

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

// The server has no Mapillary token: it can't ask whether an image still exists
var ErrNoToken = errors.New("MAPILLARY_ACCESS_TOKEN is not set")

const graphURL = "https://graph.mapillary.com/"

// Asks the Mapillary API whether an image can still be served (used by the hub when a player can't load a location)
type Checker struct {
	token	string
	baseURL	string
	client	*http.Client
}

func NewChecker() (*Checker, error) {
	token := os.Getenv("MAPILLARY_ACCESS_TOKEN")
	if token == "" {
		return nil, ErrNoToken
	}
	return &Checker{token: token, baseURL: graphURL, client: &http.Client{Timeout: 5 * time.Second}}, nil
}

// Tells whether the image mid still exists on Mapillary:
//   - (true, nil): it exists, the problem is on the player's side (connection, browser...)
//   - (false, nil): Mapillary says this id can't be served anymore (deleted, made private...): the location is broken
//   - (false, err): no answer that can be trusted (network, token refused, quota, Mapillary down): nothing can be concluded
func (c *Checker) ImageExists(ctx context.Context, mid int64) (bool, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		c.baseURL+strconv.FormatInt(mid, 10)+"?fields=id", nil)
	if err != nil {
		return false, err
	}
	// the token goes in a header, not in the URL
	// the errors of the http client contain the URL, so the token would end up in the logs
	req.Header.Set("Authorization", "OAuth "+c.token)
	resp, err := c.client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		return true, nil
	}
	var body struct {
		Error struct {
			Message	string	`json:"message"`
			Code	int		`json:"code"`
		} `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return false, fmt.Errorf("mapillary answered %d", resp.StatusCode)
	}
	// code 100: this id isn't an image the API can serve
	if body.Error.Code == 100 {
		return false, nil
	}
	return false, fmt.Errorf("mapillary answered %d (code %d): %s", resp.StatusCode, body.Error.Code, body.Error.Message)
}