package pushinator

import (
    "errors"
    "net/http"

    "github.com/go-resty/resty/v2"
)

type Client struct {
    token   string
    baseURL string
    http    *resty.Client
}

type Notification struct {
    ChannelID string `json:"channel_id"`
    Message   string `json:"content"`
}

func NewClient(token string) *Client {
    return NewClientWithHTTP(token, resty.New())
}

func NewClientWithHTTP(token string, httpClient *resty.Client) *Client {
    return &Client{
        token:   token,
        baseURL: "https://api.pushinator.com/api/v2",
        http:    httpClient,
    }
}

func (c *Client) SetBaseURL(baseURL string) {
    c.baseURL = baseURL
}

func (c *Client) SendNotification(channelId, message string) error {
    if c.token == "" {
        return errors.New("API token is required")
    }
    if channelId == "" {
        return errors.New("channel ID is required")
    }
    if message == "" {
        return errors.New("message is required")
    }

    notification := Notification{
        ChannelID: channelId,
        Message:   message,
    }

    resp, err := c.http.R().
        SetHeader("Authorization", "Bearer "+c.token).
        SetHeader("Content-Type", "application/json").
        SetBody(notification).
        Post(c.baseURL + "/notifications/send")

    if err != nil {
        return err
    }

    if resp.StatusCode() != http.StatusOK {
        return errors.New("failed to send notification: " + resp.Status())
    }

    return nil
}