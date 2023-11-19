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
	"sync"
	"time"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func wsReadHandler(wg sync.WaitGroup, ws *websocket.Conn) error {
	defer wg.Done()

	// Handle WebSocket messages
	for {
		fmt.Println("Handling")
		// Read message from the client
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			return err
		}

		// Print received message
		log.Printf("Received: %s\n", msg)
	}
}
func wsWriteHandler(wg sync.WaitGroup, ws *websocket.Conn) error {
	defer wg.Done()

	for {
		// Write a new server message back to the user
		randomString := fmt.Sprintf("Server says: %s", random.String(10))
		var b bytes.Buffer
		w := io.Writer(&b)
		err := templates.ListItem(fmt.Sprintf(randomString)).Render(context.Background(), w)
		err = ws.WriteMessage(websocket.TextMessage, b.Bytes())
		if err != nil {
			wg.Done()
			return err
		}
		time.Sleep(3 * time.Second)
	}
}

func wsHandler(c echo.Context) error {

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

	wg := sync.WaitGroup{}
	wg.Add(2)
	go wsWriteHandler(wg, conn)
	go wsReadHandler(wg, conn)
	wg.Wait()

	return nil
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
	e.GET("/ws", wsHandler)

	e.Logger.Fatal(e.Start(":1323"))
}
