package http

import (
	"amazing_talker/pkg/delivery/http/view"
	"amazing_talker/pkg/model"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Register ...
func (h *Handler) Register(c echo.Context) error {
	ctx := c.Request().Context()

	var (
		req      view.RegisterReq
		iAccount model.IdentityAccount
	)

	if err := req.BindAndVerify(c); err != nil {
		return err
	}

	iAccount = req.ConvertToIdentityAccount()
	err := h.service.CreateIdentityAccount(ctx, &iAccount)
	if err != nil {
		return err
	}

	token, err := iAccount.CreateToken(h.cfg)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, view.RegisterResp{
		Token: token,
	})
}
