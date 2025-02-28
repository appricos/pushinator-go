# pushinator-go

# Pushinator Go Client

A Go package that enables developers to send push notifications seamlessly through the Pushinator API.

## Installation

Install the package using `go get`:

```sh
go get github.com/appricos/pushinator-go
```

## Usage

### Initializing the Client

To start using the `pushinator-go` package, create a client instance by passing your API token:

```go
package main

import (
	pushinator "github.com/appricos/pushinator-go"
)

func main() {
	client := pushinator.NewClient("PUSHINATOR_API_TOKEN")
}
```

### Sending Notifications

To send a notification to a specific channel, use the `SendNotification` method. Provide your channel ID and the notification content as arguments:

```go
package main

import (
	"fmt"
	"log"

	pushinator "github.com/appricos/pushinator-go"
)

func main() {
	client := pushinator.NewClient("PUSHINATOR_API_TOKEN")

	err := client.SendNotification("PUSHINATOR_CHANNEL_ID", "Hello from Go! 🚀")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Println("✅ Notification sent successfully!")
}
```