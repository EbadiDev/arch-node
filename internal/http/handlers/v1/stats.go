package v1

import (
	"net/http"

	"github.com/cockroachdb/errors"
	"github.com/ebadidev/arch-node/pkg/xray"
	"github.com/labstack/echo/v4"
)

func StatsShow(x *xray.Xray) echo.HandlerFunc {
	return func(c echo.Context) error {
		stats, err := x.QueryStats()
		if err != nil {
			return errors.WithStack(err)
		}
		return c.JSON(http.StatusOK, stats)
	}
}
