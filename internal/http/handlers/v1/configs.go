package v1

import (
	"fmt"
	"net/http"

	"github.com/ebadidev/arch-node/internal/utils"
	"github.com/ebadidev/arch-node/pkg/xray"
	"github.com/labstack/echo/v4"
)

func ConfigsStore(x *xray.Xray) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		var config xray.Config
		if err = c.Bind(&config); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "Cannot parse the request body.",
			})
		}
		if err = config.Validate(); err != nil {
			return c.JSON(http.StatusUnprocessableEntity, map[string]string{
				"message": fmt.Sprintf("Validation error: %v", err.Error()),
			})
		}

		if c.Request().Header.Get("X-App-Name") != "Arch-Manager" {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": fmt.Sprintf("Unknown client."),
			})
		}

		for _, i := range config.Inbounds {
			isFree := utils.PortFree(i.Port)
			if i.Tag != "api" && i.Tag != "remote" && !isFree {
				return c.JSON(http.StatusUnprocessableEntity, map[string]string{
					"message": fmt.Sprintf("The port '%s.%d' is already in use", i.Tag, i.Port),
				})
			}
			if i.Tag == "remote" && !isFree {
				currentInbound := x.Config().FindInbound("remote")
				if currentInbound == nil || currentInbound.Port != i.Port {
					return c.JSON(http.StatusUnprocessableEntity, map[string]string{
						"message": fmt.Sprintf("The port '%s.%d' is already in use", i.Tag, i.Port),
					})
				}
			}
			if i.Tag == "api" {
				if i.Port, err = utils.FreePort(); err != nil {
					return c.JSON(http.StatusUnprocessableEntity, map[string]string{
						"message": fmt.Sprintf("API inbound port failed, err: %v", err.Error()),
					})
				}
			}
		}

		x.SetConfig(&config)

		go x.Restart()

		return c.JSON(http.StatusOK, map[string]string{
			"message": "The configs stored successfully.",
		})
	}
}
