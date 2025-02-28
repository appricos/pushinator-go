package pushinator

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

func TestSendNotification(t *testing.T) {
	tests := []struct {
		name           string
		token          string
		channelId      string
		message        string
		serverStatus   int
		serverResponse string
		wantErr        bool
		errMsg         string
	}{
		{
			name:         "successful notification",
			token:        "valid-token",
			channelId:    "channel-123",
			message:      "Hello, world!",
			serverStatus: http.StatusOK,
			wantErr:      false,
		},
		{
			name:         "missing token",
			token:        "",
			channelId:    "channel-123",
			message:      "Hello, world!",
			serverStatus: http.StatusOK,
			wantErr:      true,
			errMsg:       "API token is required",
		},
		{
			name:         "missing channel",
			token:        "valid-token",
			channelId:    "",
			message:      "Hello, world!",
			serverStatus: http.StatusOK,
			wantErr:      true,
			errMsg:       "channel ID is required",
		},
		{
			name:         "missing message",
			token:        "valid-token",
			channelId:    "channel-123",
			message:      "",
			serverStatus: http.StatusOK,
			wantErr:      true,
			errMsg:       "message is required",
		},
		{
			name:         "server error",
			token:        "valid-token",
			channelId:    "channel-123",
			message:      "Hello, world!",
			serverStatus: http.StatusBadRequest,
			wantErr:      true,
			errMsg:       "failed to send notification: 400 Bad Request",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/notifications/send", r.URL.Path, "Request should be sent to /notifications/send")
				assert.Equal(t, "Bearer "+tt.token, r.Header.Get("Authorization"))
				assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

				var payload map[string]string
				err := json.NewDecoder(r.Body).Decode(&payload)

				assert.NoError(t, err)
				assert.Equal(t, tt.channelId, payload["channel_id"], "channel_id should match input channelId")
				assert.Equal(t, tt.message, payload["content"], "content should match input message")

				w.WriteHeader(tt.serverStatus)
				if tt.serverResponse != "" {
					w.Write([]byte(tt.serverResponse))
				}
			}))

			defer server.Close()

			httpClient := resty.New()
			httpClient.SetTransport(server.Client().Transport)

			client := NewClientWithHTTP(tt.token, httpClient)
			client.SetBaseURL(server.URL)

			err := client.SendNotification(tt.channelId, tt.message)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNewClient(t *testing.T) {
	client := NewClient("test-token")
	assert.NotNil(t, client)
	assert.Equal(t, "test-token", client.token)
	assert.Equal(t, "https://api.pushinator.com/api/v2", client.baseURL)
	assert.NotNil(t, client.http)
}
