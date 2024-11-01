package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Message struct {
	Data string `json:"data"`
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/stream", streamHandler)

	e.Logger.Fatal(e.Start(":3000"))
}

func streamHandler(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Response().WriteHeader(http.StatusOK)

	enc := json.NewEncoder(c.Response())
	for i := 0; i < 5; i++ {
		message := Message{Data: fmt.Sprintf("data ke-%d", i+1)}
		fmt.Println("Mengirim:", message)
		fmt.Println()
		if err := enc.Encode(message); err != nil {
			return err
		}
		c.Response().Flush()
		time.Sleep(500 * time.Millisecond)
	}

	doneMessage := Message{Data: "done"}
	fmt.Println("Mengirim:", doneMessage)
	fmt.Println()
	if err := enc.Encode(doneMessage); err != nil {
		return err
	}
	c.Response().Flush()

	return nil
}
