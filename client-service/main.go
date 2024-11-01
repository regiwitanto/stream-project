package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"

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

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET},
	}))

	e.GET("/consume", consumeStreamHandler)

	e.Logger.Fatal(e.Start(":3001"))
}

func consumeStreamHandler(c echo.Context) error {
	resp, err := http.Get("http://localhost:3000/stream")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return echo.NewHTTPError(resp.StatusCode, "Failed to connect to stream service")
	}

	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	c.Response().WriteHeader(http.StatusOK)

	reader := bufio.NewReader(resp.Body)
	decoder := json.NewDecoder(reader)

	for {
		var message Message
		if err := decoder.Decode(&message); err != nil {
			break
		}

		fmt.Println("Menerima:", message.Data)
		fmt.Println()

		formattedChunk := fmt.Sprintf("data: %s\n\n", message.Data)
		c.Response().Write([]byte(formattedChunk))
		c.Response().Flush()
	}

	return nil
}
