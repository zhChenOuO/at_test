package errors

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// HTTPErrorHandlerForEcho responds error response according to given error.
func HTTPErrorHandlerForEcho(err error, c echo.Context) {
	if err == nil {
		return
	}

	echoErr, ok := err.(*echo.HTTPError)
	if ok {
		_ = c.JSON(echoErr.Code, echoErr)
		return
	}

	causeErr := errors.Cause(err)
	_err, ok := causeErr.(*_error)
	if !ok || _err == nil {
		_ = c.JSON(http.StatusInternalServerError, ErrInternalError)
		return
	}

	_ = c.JSON(_err.Status, GetHTTPError(c, _err))
}

//ErrMiddleware provide error middleware
func ErrMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err != nil {
			logFields := map[string]interface{}{}

			// 紀錄 Request 資料
			req := c.Request()
			{
				logFields["requestMethod"] = req.Method
				logFields["requestURL"] = req.URL.String()
			}

			// 紀錄 Response 資料
			resp := c.Response()
			resp.After(func() {
				logFields["responseStatus"] = resp.Status
				// 根據狀態碼用不同等級來紀錄
				logger := log.Ctx(req.Context()).With().Fields(logFields).Logger()
				switch {
				case resp.Status >= http.StatusInternalServerError:
					logger.Error().Msgf("%+v", err)
				case resp.Status >= http.StatusBadRequest:
					logger.Warn().Msgf("%+v", err)
				default:
					logger.Debug().Msgf("%+v", err)
				}
			})
		}
		return err
	}

}

// NotFoundHandlerForEcho responds not found response.
func NotFoundHandlerForEcho(c echo.Context) error {
	return c.JSON(http.StatusNotFound, ErrPageNotFound)
}
