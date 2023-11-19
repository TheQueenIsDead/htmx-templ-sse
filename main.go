package main

//go:generate templ generate templates

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/random"
	"htmx-temple-wss/templates"
	"io"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer conn.Close()

	// Handle WebSocket messages
	for {
		fmt.Println("Handling")
		// Read message from the client
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		// Print received message
		log.Printf("Received: %s\n", msg)

		randomString := fmt.Sprintf("Server says: %s", random.String(10))
		var b bytes.Buffer
		w := io.Writer(&b)

		// Write received msg back to user
		err = templates.ListItem(fmt.Sprintf("You said: %s", msg)).Render(context.Background(), w)
		err = conn.WriteMessage(websocket.TextMessage, b.Bytes())
		if err != nil {
			return
		}

		// Write a new server message back to the user
		b.Reset()
		err = templates.ListItem(fmt.Sprintf(randomString)).Render(context.Background(), w)
		err = conn.WriteMessage(websocket.TextMessage, b.Bytes())
		if err != nil {
			return
		}

	}
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Index
	e.GET("/", func(c echo.Context) error {
		return templates.Index().Render(c.Request().Context(), c.Response().Writer)
	})

	// WebSocket endpoint
	e.GET("/ws", func(c echo.Context) error {
		wsHandler(c.Response(), c.Request())
		return nil
	})

	e.Logger.Fatal(e.Start(":1323"))
}
