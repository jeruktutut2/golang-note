package globals

import "github.com/labstack/echo/v4"

type SSEClient struct {
	Id      int64
	Context echo.Context
}

var SSEClients []SSEClient
