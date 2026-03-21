package router

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func NewRouter() *echo.Echo {

	e := echo.New()
	e.Use(middleware.RequestLogger())

	e.GET("/", func(c *echo.Context) error {
		return c.String(http.StatusOK, "healthy")
	})

	e.GET("/testIpv6Conn", func(c *echo.Context) error {

		success, err := ipv6_info.IpV6Google(c.Request().Context())
		if err != nil {
			return err
		}

		if success {
			return c.String(http.StatusOK, "working")
		}

		return c.String(http.StatusOK, "broken")
	})

	return e
}
