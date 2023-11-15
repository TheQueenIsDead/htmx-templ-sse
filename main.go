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
			err := websocket.Message.Send(ws, fmt.Sprintf("<div id=\"notifications\" hx-swap-oob=\"afterend\" > %s! </div>", randomString))
			if err != nil {
				if errors.Is(err, syscall.EPIPE) {
					c.Logger().Error("connection broken")
					return
				}
				c.Logger().Error(err)
				//panic(err)
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
			//	//panic(err)
			//}

			//fmt.Printf("%s\n", msg)
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
