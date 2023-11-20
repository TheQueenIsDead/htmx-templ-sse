package main

//go:generate templ generate templates

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log2 "github.com/labstack/gommon/log"
	"github.com/labstack/gommon/random"
	"htmx-temple-wss/templates"
	"io"
	"log"
	"net/http"
	"sync"
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

func wsReadHandler(logger echo.Logger, wg sync.WaitGroup, ws *websocket.Conn, msgChan chan string) error {
	logger.Debug("handling read operations for websocket")

	defer wg.Done()

	// Handle WebSocket messages

	for {
		logger.Debug("waiting to read from websocket")
		// Read message from the client

		var msg HtmxResponse
		err := ws.ReadJSON(&msg)
		if err != nil {
			logger.Error("could not marshal message into htmx template")
			logger.Error(err)
		}
		logger.Debug("websocket data received")

		// Print received message
		logger.Debug(fmt.Sprintf("received: %s", msg.ChatMessage))
		logger.Debug("adding read message to write channel")
		msgChan <- msg.ChatMessage
		logger.Debug("read message placed on write channel")

	}
}
func wsWriteHandler(logger echo.Logger, wg sync.WaitGroup, ws *websocket.Conn, msgChan chan string) (err error) {

	logger.Debug("handling write operations for websocket")
	defer wg.Done()

	serverMessageChan := make(chan string)
	go func() {
		for {
			// Write a new server message back to the user
			randomString := fmt.Sprintf("Server says: %s", random.String(10))
			var b bytes.Buffer
			w := io.Writer(&b)
			err = templates.ListItem(fmt.Sprintf(randomString)).Render(context.Background(), w)
			if err != nil {
				logger.Error(err)
			}
			logger.Debug(fmt.Sprintf("adding server message to channel: %s", randomString))
			serverMessageChan <- b.String()
			time.Sleep(3 * time.Second)
		}
	}()

	for {
		select {
		case serverMessage := <-serverMessageChan:
			logger.Debug(fmt.Sprintf("attempting to send server message: %s", serverMessage))
			err = ws.WriteMessage(websocket.TextMessage, []byte(serverMessage))
			if err != nil {
				logger.Error("error writing message to websocket")
				logger.Error(err)
				wg.Done()
				return
			}
		}
	}
	logger.Debug("broken out of handling write operations for websocket")
	return nil
}

func wsHandler(c echo.Context) error {

	c.Logger().Debug("handling incoming webhook request")

	var w http.ResponseWriter
	var r *http.Request
	w, r = c.Response(), c.Request()

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return err
	}
	defer conn.Close()

	msgChan := make(chan string, 1)

	wg := sync.WaitGroup{}
	wg.Add(2)
	go wsWriteHandler(c.Logger(), wg, conn, msgChan)
	go wsReadHandler(c.Logger(), wg, conn, msgChan)
	c.Logger().Debug("waiting on websocket r/w goroutines to complete")
	wg.Wait()
	c.Logger().Debug("websocket r/w goroutines completed")

	return nil
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
	e.GET("/ws", wsHandler)

	e.Logger.Fatal(e.Start(":1323"))
}
