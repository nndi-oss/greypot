package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/nndi-oss/greypot/service"
	"github.com/sirupsen/logrus"
)

func ReportExportHandler(reportService service.ReportService, kind string) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		reportId := strings.TrimPrefix(ctx.PathParam("*"), "/")
		var body interface{}
		if err := ctx.Bind(&body); err != nil {
			logrus.Error(err)
			return ctx.JSON(http.StatusInternalServerError, responseMap{
				"err": err.Error(),
			})
		}

		var export []byte
		var err error
		switch kind {
		case "html":
			export, err = reportService.ExportReportHtml(reportId, body)
		case "pdf":
			export, err = reportService.ExportReportPdf(reportId, body)
		case "png":
			export, err = reportService.ExportReportPng(reportId, body)
		}

		if err != nil {
			logrus.Error(err)
			return ctx.JSON(http.StatusInternalServerError, responseMap{
				"err": err.Error(),
			})
		}

		return ctx.JSON(http.StatusOK, ExportResponse{
			ID:   reportId,
			Data: string(export),
			Type: kind,
		})
	}
}

type ExportResponse struct {
	ID   string `json:"reportId"`
	Data string `json:"data"`
	Type string `json:"type"`
}

type BulkExportEntry struct {
	ID   string `json:"_id"`
	Data any    `json:"data"`
}

type BulkExportRequest struct {
	ID      string            `json:"_id"`
	Entries []BulkExportEntry `json:"entries"`
}

type BulkExportResponse struct {
	ID       string           `json:"_id"`
	ReportID string           `json:"reportId"`
	Reports  []ExportResponse `json:"reports"`
}

func BulkReportExportHandler(reportService service.ReportService, kind string) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		reportId := strings.TrimPrefix(ctx.PathParam("*"), "/")
		body := new(BulkExportRequest)
		if err := ctx.Bind(&body); err != nil {
			logrus.Error(err)
			return ctx.JSON(http.StatusInternalServerError, responseMap{
				"err": err.Error(),
			})
		}

		bulkResponse := &BulkExportResponse{
			ID:       body.ID,
			ReportID: reportId,
			Reports:  make([]ExportResponse, 0),
		}

		for _, entry := range body.Entries {
			entryReal := entry
			reportData := entryReal.Data
			var export []byte
			var err error
			switch kind {
			case "html":
				export, err = reportService.ExportReportHtml(reportId, reportData)
			case "pdf":
				export, err = reportService.ExportReportPdf(reportId, reportData)
			case "png":
				export, err = reportService.ExportReportPng(reportId, reportData)
			}
			if err != nil {
				logrus.Error(err)
				return ctx.JSON(http.StatusInternalServerError, responseMap{
					"err": err.Error(),
				})
			}
			bulkResponse.Reports = append(bulkResponse.Reports, ExportResponse{
				ID:   fmt.Sprintf("%s:%s", reportId, entryReal.ID),
				Data: string(export),
				Type: kind,
			})
		}

		return ctx.JSON(http.StatusOK, bulkResponse)
	}
}
