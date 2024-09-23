# Vyper API Go SDK

![Vyper](https://images.vyper.trade/0000/vyper-header)

A Go SDK for interacting with the [Vyper API](https://build.vyper.trade/). This library allows developers to integrate Vyper's HTTP and WebSocket API into their Go applications with ease.

## Table of Contents

- [Vyper API Go SDK](#vyper-api-go-sdk)
  - [Table of Contents](#table-of-contents)
  - [Installation](#installation)
  - [Quick Start](#quick-start)
  - [Usage](#usage)
    - [Client Initialization](#client-initialization)
    - [REST API Example](#rest-api-example)
    - [WebSocket API Example](#websocket-api-example)
  - [API Documentation](#api-documentation)

## Installation

To install the Vyper API Go SDK, use `go get:`

```bash
go get github.com/Vyper-Terminal/vyper-client-go
```

## Quick Start

Here's a simple example to get you started:

```go
package main

import (
    "fmt"
    "log"

    vyperclientgo "github.com/Vyper-Terminal/vyper-client-go"
)

func main() {
    // Initialize the client with your API key
    client := vyperclientgo.NewVyperClient("your_api_key_here")

    // Get the list of chain IDs supported by Vyper
    chainIds, err := client.GetChainIds()
    if err != nil {
        log.Fatalf("Error fetching chain IDs: %v", err)
    }

    fmt.Println("Supported chain IDs:", chainIds)
}
```

## Usage

### Client Initialization

The `VyperClient` struct provides access to the RESTful API endpoints:

```go
// Create a client instance
client := vyperclientgo.NewVyperClient("your_api_key_here")
```

### REST API Example

Retrieve the market data for a specific token:

```go
package main

import (
    "fmt"
    "log"

    vyperclientgo "github.com/Vyper-Terminal/vyper-client-go"
)

func main() {
    client := vyperclientgo.NewVyperClient("your_api_key_here")

    // Fetch the All-Time High (ATH) data for a token
    tokenAth, err := client.GetTokenAth(1, "AVs9TA4nWDzfPJE9gGVNJMVhcQy3V9PGazuz33BfG2RA")
    if err != nil {
        log.Fatalf("Error fetching token ATH data: %v", err)
    }

    fmt.Printf("Market Cap USD: %f\n", tokenAth.MarketCapUSD)
    fmt.Printf("Timestamp: %s\n", tokenAth.Timestamp)
}
```

### WebSocket API Example

```go
package main

import (
    "fmt"
    "log"

    vyperclientgo "github.com/Vyper-Terminal/vyper-client-go"
)

func main() {
    wsClient := vyperclientgo.NewVyperWebsocketClient("your_api_key_here")

    // Define a message handler
    messageHandler := func(message interface{}) {
        fmt.Println("Received message:", message)
    }

    wsClient.SetMessageHandler(messageHandler)

    // Connect to the WebSocket and subscribe to token events
    err := wsClient.Connect(vyperclientgo.TokenEvents)
    if err != nil {
        log.Fatalf("Failed to connect to WebSocket: %v", err)
    }

    subscribeMessage := vyperclientgo.TokenSubscriptionMessage{
        Action: vyperclientgo.Subscribe,
        Types:  []vyperclientgo.SubscriptionType{vyperclientgo.PumpfunTokens},
    }

    err = wsClient.Subscribe(vyperclientgo.TokenEvents, subscribeMessage)
    if err != nil {
        log.Fatalf("Failed to subscribe to token events: %v", err)
    }

    fmt.Println("Subscribed to token events")

    err = wsClient.Listen()
    if err != nil {
        log.Fatalf("Error listening for WebSocket messages: %v", err)
    }
}
```

## API Documentation

For detailed information on the Vyper API, refer to the official documentation:

-   API Dashboard: [Vyper Dashboard](https://build.vyper.trade/)
-   API Documentation: [Vyper API Docs](ttps://docs.vyper.trade/)
