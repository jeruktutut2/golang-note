package setup

import (
	"context"
	"golang-note/configuration"
	"golang-note/exception"
	"golang-note/middleware"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

func Echo(config *configuration.Configuration) (e *echo.Echo) {
	e = echo.New()
	e.Use(echomiddleware.Recover())
	e.HTTPErrorHandler = exception.CustomHTTPErrorHandler
	e.Use(middleware.SetGlobalRequestId)
	e.Use(middleware.SetGlobalTimeout(config.ApplicationTimeout))
	e.Use(middleware.SetGlobalRequestLog)
	return
}

func StartEcho(e *echo.Echo, config *configuration.Configuration) {
	go func() {
		if err := e.Start(":" + strconv.Itoa(int(config.ApplicationPort))); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()
}

func StopEcho(e *echo.Echo) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
	cancel()
	println(time.Now().String(), "server shutdown properly")
}
