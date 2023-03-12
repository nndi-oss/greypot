package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/nndi-oss/greypot"
	"github.com/nndi-oss/greypot/http/fiber/handlers"
	"github.com/sirupsen/logrus"
)

func generatePDF(studioModule *greypot.Module, studioTemplateStore *inmemoryTemplateRepository) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		request := new(UploadTemplateRequest)
		err := c.BodyParser(&request)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"message":    "failed to parse request body",
				"devMessage": err.Error(),
			})
		}
		reportId := strings.TrimSpace(request.Name)
		err = studioTemplateStore.Add(reportId, request.Template)
		defer studioTemplateStore.Remove(reportId)

		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"message":    "failed to upload template to store",
				"devMessage": err.Error(),
			})
		}

		export, err := studioModule.ReportService.ExportReportPdf(reportId, request.Data)
		if err != nil {
			logrus.Error(err)
			return c.Status(http.StatusInternalServerError).
				JSON(fiber.Map{
					"err": err.Error(),
				})
		}

		return c.JSON(handlers.ExportResponse{
			ID:   reportId,
			Data: string(export),
			Type: "pdf",
		})
	}
}

func generateBulkPDF(studioModule *greypot.Module, studioTemplateStore *inmemoryTemplateRepository) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		reportId := strings.TrimPrefix(c.Params("*"), "/")
		body := new(handlers.BulkExportRequest)
		if err := c.BodyParser(&body); err != nil {
			logrus.Error(err)
			return c.Status(http.StatusInternalServerError).
				JSON(fiber.Map{
					"err": err.Error(),
				})
		}

		bulkResponse := &handlers.BulkExportResponse{
			ID:       body.ID,
			ReportID: reportId,
			Reports:  make([]handlers.ExportResponse, 0),
		}

		for _, entry := range body.Entries {
			entryReal := entry
			reportData := entryReal.Data
			var export []byte

			export, err := studioModule.ReportService.ExportReportPdf(reportId, reportData)
			if err != nil {
				logrus.Error(err)
				return c.Status(http.StatusInternalServerError).
					JSON(fiber.Map{
						"err": err.Error(),
					})
			}
			bulkResponse.Reports = append(bulkResponse.Reports, handlers.ExportResponse{
				ID:   fmt.Sprintf("%s:%s", reportId, entryReal.ID),
				Data: string(export),
				Type: "pdf",
			})
		}

		return c.JSON(bulkResponse)
	}
}

func generateExcel(studioModule *greypot.Module, studioTemplateStore *inmemoryTemplateRepository) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		reportId := strings.TrimPrefix(c.Params("*"), "/")
		var body interface{}
		if err := c.BodyParser(&body); err != nil {
			logrus.Error(err)
			return c.Status(http.StatusInternalServerError).
				JSON(fiber.Map{
					"err": err.Error(),
				})
		}

		html2Excel := NewHtml2ExcelTemplateEngine()
		mod := greypot.NewModule(
			greypot.WithTemplatesRepository(studioTemplateStore),
			greypot.WithTemplateEngine(html2Excel),
		)
		export, err := mod.ReportService.ExportReportHtml(reportId, body)
		if err != nil {
			logrus.Error(err)
			return c.Status(http.StatusInternalServerError).
				JSON(fiber.Map{
					"err": err.Error(),
				})
		}

		return c.JSON(handlers.ExportResponse{
			ID:   reportId,
			Data: string(export),
			Type: "excel",
		})
	}
}
