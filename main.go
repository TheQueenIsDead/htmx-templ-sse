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

func hello(c echo.Context) error {

	var upgrader = websocket.Upgrader{}
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	for {
		mt, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return err
		}
		log.Printf("recv: %s", message)
		message = []byte("hello")

		randomString := random.String(10)
		var b bytes.Buffer
		//html := fmt.Sprintf("<ul hx-swap-oob='beforeend:#items'><li>%s</li></ul>", randomString)
		w := io.Writer(&b)
		err = templates.ListItem(randomString).Render(c.Request().Context(), w)
		if err != nil {
			c.Logger().Error("failed to template ListItem")
		}
		html := b.Bytes()

		err = ws.WriteMessage(mt, html)
		if err != nil {
			log.Println("write:", err)
			return err
		}
	}
	//
	//websocket.Handler(func(ws *websocket.Conn) {
	//	defer ws.Close()
	//	for {
	//
	//		// Write
	//		randomString := random.String(10)
	//
	//		// Partial Attempt
	//		// Return a list item wrapped in a list in order to mitigate interesting parsing behaviour when using OOB fragments
	//		// https://github.com/bigskysoftware/htmx/issues/1198#issuecomment-1763180864
	//		// May in time be resolved by
	//		var b bytes.Buffer
	//		//html := fmt.Sprintf("<ul hx-swap-oob='beforeend:#items'><li>%s</li></ul>", randomString)
	//		w := io.Writer(&b)
	//		err := templates.ListItem(randomString).Render(c.Request().Context(), w)
	//		if err != nil {
	//			c.Logger().Error("failed to template ListItem")
	//		}
	//		html := b.String()
	//
	//		err = websocket.Message.Send(ws, html)
	//		if err != nil {
	//			if errors.Is(err, syscall.EPIPE) {
	//				c.Logger().Error("connection broken")
	//				return
	//			}
	//			c.Logger().Error(err)
	//		}
	//
	//		//// Read
	//		//msg := ""
	//		//err = websocket.Message.Receive(ws, &msg)
	//		//if err != nil {
	//		//	if errors.Is(err, io.EOF) {
	//		//		log.Debug("No data to read")
	//		//	} else if errors.Is(err, syscall.EPIPE) {
	//		//		log.Debug("broken pipe")
	//		//	}
	//		//	c.Logger().Error(err)
	//		//}
	//
	//		//fmt.Printf("%s\n", msg)
	//		time.Sleep(3 * time.Second)
	//	}
	//}).ServeHTTP(c.Response(), c.Request())
	//return nil
}

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

		// Write message back to the client
		if err = conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Println("Error writing message:", err)
			break
		}

		randomString := random.String(10)
		var b bytes.Buffer
		//html := fmt.Sprintf("<ul hx-swap-oob='beforeend:#items'><li>%s</li></ul>", randomString)
		w := io.Writer(&b)
		err = templates.ListItem(randomString).Render(context.Background(), w)
		conn.WriteMessage(websocket.TextMessage, b.Bytes())
		if err != nil {
			return
		}

	}
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//e.GET("/test", test)

	e.GET("/", func(c echo.Context) error {
		return templates.Index().Render(c.Request().Context(), c.Response().Writer)
	})
	//e.GET("/ws", hello)
	// WebSocket endpoint
	e.GET("/ws", func(c echo.Context) error {
		wsHandler(c.Response(), c.Request())
		return nil
	})

	e.Logger.Fatal(e.Start(":1323"))
}
