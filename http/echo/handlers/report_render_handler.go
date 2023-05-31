package handlers

import (
	"bytes"
	"net/http"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/nndi-oss/greypot/service"
	"github.com/sirupsen/logrus"
)

func ReportRenderHandlder(templateService service.TemplateService) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		reportId := strings.TrimPrefix(ctx.PathParam("*"), "/")
		data, err := extractData(ctx)
		if err != nil {
			logrus.Error(err.Error())
			return ctx.JSON(http.StatusInternalServerError, responseMap{
				"err": err.Error(),
			})
		}

		html, err := templateService.RenderTemplate(reportId, data)
		if err != nil {
			logrus.Error(err.Error())
			return ctx.JSON(http.StatusInternalServerError, responseMap{
				"err": err.Error(),
			})
		}

		r := bytes.NewReader(html)
		return ctx.Stream(http.StatusOK, "text/html;charset=utf-8", r)
	}
}
