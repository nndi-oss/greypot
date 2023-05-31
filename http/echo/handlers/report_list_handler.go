package handlers

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/nndi-oss/greypot/service"
	"github.com/sirupsen/logrus"
)

func ReportListHandler(tmplSrv service.TemplateService) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		templates, err := tmplSrv.ListTemplates()
		if err != nil {
			logrus.Error(err.Error())
			return ctx.JSON(http.StatusInternalServerError, responseMap{
				"err": err.Error(),
			})
		}

		return ctx.JSON(http.StatusOK, responseMap{
			"list": templates,
		})
	}
}
