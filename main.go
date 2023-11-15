//package main
//
//import (
//	"fmt"
//	"github.com/labstack/echo/v4"
//	"github.com/labstack/echo/v4/middleware"
//	"golang.org/x/net/websocket"
//)
//
//func hello(c echo.Context) error {
//
//	websocket.Handler(func(ws *websocket.Conn) {
//		defer ws.Close()
//		for {
//			// Write
//			err := websocket.Message.Send(ws, "Hello, Client!")
//			if err != nil {
//				c.Logger().Error(err)
//				break
//			}
//
//			// Read
//			msg := ""
//			err = websocket.Message.Receive(ws, &msg)
//			if err != nil {
//				c.Logger().Error(err)
//				break
//			}
//			fmt.Printf("%s\n", msg)
//		}
//	}).ServeHTTP(c.Response(), c.Request())
//	return nil
//}
//
//func main() {
//	e := echo.New()
//	e.Use(middleware.Logger())
//	e.Use(middleware.Recover())
//	e.Static("/", "./templates")
//	e.GET("/ws", hello)
//	e.Logger.Fatal(e.Start(":1323"))
//}

package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/random"
	"golang.org/x/net/websocket"
	"net/http"
	"time"
)

func hello(c echo.Context) error {

	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()
		for {

			// Write
			randomString := random.String(10)
			err := websocket.Message.Send(ws, fmt.Sprintf("<div id=\"notifications\"> %s! </div>", randomString))
			if err != nil {
				c.Logger().Error(err)
				break
			}

			// Read
			msg := ""
			err = websocket.Message.Receive(ws, &msg)
			if err != nil {
				c.Logger().Error(err)
				break
			}

			fmt.Printf("%s\n", msg)
			time.Sleep(1 * time.Second)
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	//e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.Static("/", "./templates")
	e.GET("/ws", hello)

	e.Logger.Fatal(e.Start(":1323"))
}
