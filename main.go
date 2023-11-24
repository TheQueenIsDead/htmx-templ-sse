package main

//go:generate templ generate templates

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/labstack/gommon/random"
	"htmx-temple-wss/templates"
	"io"
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

type Message struct {
	ChatMessage string `query:"chat_message"`
}

var (
	msgChan = make(chan string)
)

func messageHandler(c echo.Context) error {

	var message Message
	err := c.Bind(&message)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	// TODO: Actually return something
	return nil
}

func sseHandler(c echo.Context) error {
	// Set headers for SSE
	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")

	// Flush the headers to ensure the client receives them
	c.Response().Flush()

	go func() {
		// Infinite loop to keep the connection open
		for {
			// Write a new server message back to the user
			randomString := fmt.Sprintf("Server says: %s", random.String(10))
			var b bytes.Buffer
			w := io.Writer(&b)
			err := templates.IncomingMessage(fmt.Sprintf(randomString)).Render(context.Background(), w)
			if err != nil {
				c.Logger().Error(err)
			}
			c.Logger().Debug(fmt.Sprintf("adding server message to channel: %s", randomString))
			msgChan <- b.String()
			time.Sleep(1 * time.Second)

		}
	}()

	for {
		select {
		case serverMessage := <-msgChan:

			// Simulate some data to send to the client
			data := serverMessage
			// Send data to the client
			c.Logger().Debug(fmt.Sprintf("attempting to send server message: %s", serverMessage))
			_, err := c.Response().Write([]byte("data: " + data + "\n\n"))
			if err != nil {
				c.Logger().Error("error writing server sent event")
				c.Logger().Error(err)
				return err
			}
			c.Response().Flush()
			//// Introduce a delay (you can replace this with real-time updates)
			//<-time.After(1 * time.Second)

		}

	}

}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))

	e.Logger.SetLevel(log.DEBUG)

	// Index
	e.GET("/", func(c echo.Context) error {
		return templates.Index().Render(c.Request().Context(), c.Response().Writer)
	})

	// Message Creation
	e.POST("/message", func(c echo.Context) error {
		return templates.Index().Render(c.Request().Context(), c.Response().Writer)
	})

	// Server Sent Events
	e.GET("/sse", sseHandler)

	e.Logger.Fatal(e.Start(":1323"))
}
