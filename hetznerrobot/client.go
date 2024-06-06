package hetznerrobot

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type HetznerRobotClient struct {
	username string
	password string
	url      string
}

func NewHetznerRobotClient(username string, password string, url string) HetznerRobotClient {
	return HetznerRobotClient{
		username: username,
		password: password,
		url:      url,
	}
}

func (c *HetznerRobotClient) makeAPICall(ctx context.Context, method string, uri string, body io.Reader) ([]byte, error) {
	tflog.Debug(ctx, "requesting Hetzner webservice", map[string]interface{}{
		"uri":    uri,
		"method": method,
		"body":   body,
	})

	r, err := http.NewRequestWithContext(ctx, method, uri, body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	r.SetBasicAuth(c.username, c.password)

	client := http.Client{}

	response, err := client.Do(r)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}

	defer response.Body.Close()

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	tflog.Debug(ctx, "got hetzner webservice response", map[string]interface{}{
		"status": response.StatusCode,
		"body":   bytes,
	})

	if response.StatusCode > 400 {
		return nil, fmt.Errorf("hetzner webservice response status %d: %s", response.StatusCode, bytes)
	}

	return bytes, nil
}
