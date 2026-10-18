package httpclient

import (
    "net/http"
    "encoding/json"
    "io"
    "fmt"
)

type Client struct {
    *http.Client
}

func New() *Client {
    return &Client {
        Client: &http.Client{},
    }
}

func (c *Client) GetJSON(url string, out any) error {
    resp, err := c.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return fmt.Errorf(
            "GET %s: %d: %s",
            url,
            resp.StatusCode,
            body,
        )
    }

    if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return err
	}
    return nil
}