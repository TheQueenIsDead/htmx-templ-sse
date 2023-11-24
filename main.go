package main

//go:generate templ generate templates

import (
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log2 "github.com/labstack/gommon/log"
	"htmx-temple-wss/templates"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type HtmxResponse struct {
	ChatMessage string `json:"chat_message"`
	HEADERS     struct {
		HXRequest     string      `json:"HX-Request"`
		HXTrigger     string      `json:"HX-Trigger"`
		HXTriggerName interface{} `json:"HX-Trigger-Name"`
		HXTarget      string      `json:"HX-Target"`
		HXCurrentURL  string      `json:"HX-Current-URL"`
	} `json:"HEADERS"`
}

func sseHandler(c echo.Context) error {
	// Set headers for SSE
	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")

	// Flush the headers to ensure the client receives them
	c.Response().Flush()

	// Infinite loop to keep the connection open
	for {
		// Simulate some data to send to the client
		data := "Hello, SSE!"
		// Send data to the client
		c.Response().Write([]byte("data: " + data + "\n\n"))
		c.Response().Flush()
		// Introduce a delay (you can replace this with real-time updates)
		<-time.After(1 * time.Second)
	}
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Logger.SetLevel(log2.DEBUG)

	// Index
	e.GET("/", func(c echo.Context) error {
		return templates.Index().Render(c.Request().Context(), c.Response().Writer)
	})

	// WebSocket endpoint
	//e.GET("/ws", wsHandler)

	// sseHandler
	e.GET("/sse", sseHandler)

	e.Logger.Fatal(e.Start(":1323"))
}
