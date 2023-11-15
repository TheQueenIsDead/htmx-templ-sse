package main

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/random"
	"golang.org/x/net/websocket"
	"net/http"
	"syscall"
	"time"
)

func hello(c echo.Context) error {

	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()
		for {

			// Write
			randomString := random.String(10)

			// Partial Attempt
			// Return a list item wrapped in a list in order to mitigate interesting parsing behaviour when using OOB fragments
			// https://github.com/bigskysoftware/htmx/issues/1198#issuecomment-1763180864
			// May in time be resolved by
			html := fmt.Sprintf("<ul hx-swap-oob='beforeend:#items'><li>%s</li></ul>", randomString)

			err := websocket.Message.Send(ws, html)
			if err != nil {
				if errors.Is(err, syscall.EPIPE) {
					c.Logger().Error("connection broken")
					return
				}
				c.Logger().Error(err)
			}

			//// Read
			//msg := ""
			//err = websocket.Message.Receive(ws, &msg)
			//if err != nil {
			//	if errors.Is(err, io.EOF) {
			//		log.Debug("No data to read")
			//	} else if errors.Is(err, syscall.EPIPE) {
			//		log.Debug("broken pipe")
			//	}
			//	c.Logger().Error(err)
			//}

			//fmt.Printf("%s\n", msg)
			time.Sleep(3 * time.Second)
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.Static("/", "./templates")
	e.GET("/ws", hello)

	e.Logger.Fatal(e.Start(":1323"))
}
