package handlers

import (
	"archive/zip"
	"encoding/base64"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/nndi-oss/greypot/service"
	"github.com/sirupsen/logrus"
)

func ReportExportHandler(reportService service.ReportService, kind string) func(*fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		reportId := strings.TrimPrefix(ctx.Params("*"), "/")
		var body interface{}
		if err := ctx.BodyParser(&body); err != nil {
			logrus.Error(err)
			return ctx.Status(http.StatusInternalServerError).
				JSON(fiber.Map{
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
			return ctx.Status(http.StatusInternalServerError).
				JSON(fiber.Map{
					"err": err.Error(),
				})
		}

		return ctx.JSON(ExportResponse{
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

func BulkReportExportHandler(reportService service.ReportService, kind string) func(*fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		reportId := strings.TrimPrefix(ctx.Params("*"), "/")
		body := new(BulkExportRequest)
		if err := ctx.BodyParser(&body); err != nil {
			logrus.Error(err)
			return ctx.Status(http.StatusInternalServerError).
				JSON(fiber.Map{
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
				return ctx.Status(http.StatusInternalServerError).
					JSON(fiber.Map{
						"err": err.Error(),
					})
			}
			bulkResponse.Reports = append(bulkResponse.Reports, ExportResponse{
				ID:   fmt.Sprintf("%s:%s", reportId, entryReal.ID),
				Data: string(export),
				Type: kind,
			})
		}

		if ctx.Accepts("application/zip", "application/octet-stream") == "application/zip" {
			zipWriter := zip.NewWriter(ctx)
			defer zipWriter.Close()
			for _, report := range bulkResponse.Reports {
				reportFileName := strings.ReplaceAll(report.ID, "/", "_")
				reportFileName = strings.ReplaceAll(reportFileName, ":", "_")

				w, err := zipWriter.Create(filepath.Clean(fmt.Sprintf("%s.pdf", reportFileName)))
				if err != nil {
					logrus.Error(err)
					return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
						"message": "failed to process request",
					})
				}
				fileData, err := base64.StdEncoding.DecodeString(report.Data)
				if err != nil {
					logrus.Error(err)
					return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
						"message": "failed to process request",
					})
				}
				w.Write(fileData)
			}

			ctx.Type("zip")
			ctx.Set(fiber.HeaderContentDisposition, `attachment; filename="reports.zip"`)
			return ctx.SendStatus(200)
		}

		return ctx.JSON(bulkResponse)
	}
}
