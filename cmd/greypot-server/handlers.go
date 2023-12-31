package main

import (
	"archive/zip"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
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

		if c.Accepts("json", "application/json", "application/pdf") == "application/pdf" {
			tmpFileName, err := writeToTempFile(export)
			if err != nil {
				return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
					"message": "failed to process request",
				})
			}

			downloadFileName := request.Name
			if !strings.HasSuffix(downloadFileName, ".pdf") {
				downloadFileName = fmt.Sprintf("%s.pdf", downloadFileName)
			}

			return c.Download(tmpFileName, downloadFileName)
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
		request := new(BulkUploadTemplateRequest)
		if err := c.BodyParser(&request); err != nil {
			logrus.Error(err)
			return c.Status(http.StatusInternalServerError).
				JSON(fiber.Map{
					"err": err.Error(),
				})
		}

		reportId := strings.TrimSpace(request.Name)
		err := studioTemplateStore.Add(reportId, request.Template)
		defer studioTemplateStore.Remove(reportId)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"message":    "failed to upload template to store",
				"devMessage": err.Error(),
			})
		}

		bulkResponse := &handlers.BulkExportResponse{
			ID:       reportId,
			ReportID: reportId,
			Reports:  make([]handlers.ExportResponse, 0),
		}

		for _, entry := range request.Entries {
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

		if c.Accepts("application/zip", "application/octet-stream") == "application/zip" {
			// var zipBuffer bytes.Buffer
			zipWriter := zip.NewWriter(c)
			defer zipWriter.Close()
			for _, report := range bulkResponse.Reports {
				reportFileName := strings.ReplaceAll(report.ID, "/", "_")
				reportFileName = strings.ReplaceAll(reportFileName, ":", "_")

				w, err := zipWriter.Create(filepath.Clean(fmt.Sprintf("%s.pdf", reportFileName)))
				if err != nil {
					logrus.Error(err)
					return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
						"message": "failed to process request",
					})
				}
				fileData, err := base64.StdEncoding.DecodeString(report.Data)
				if err != nil {
					logrus.Error(err)
					return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
						"message": "failed to process request",
					})
				}
				w.Write(fileData)
			}

			c.Type("zip")
			c.Set(fiber.HeaderContentDisposition, `attachment; filename="reports.zip"`)
			return c.SendStatus(200)
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

		if c.Accepts("json", "application/json", "application/octet-stream") == "application/octet-stream" {
			tmpFileName, err := writeToTempFile(export)
			if err != nil {
				return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
					"message": "failed to process request",
				})
			}
			downloadFileName := reportId
			if !strings.HasSuffix(downloadFileName, ".xlsx") {
				downloadFileName = fmt.Sprintf("%s.xlsx", downloadFileName)
			}
			return c.Download(tmpFileName, downloadFileName)
		}

		return c.JSON(handlers.ExportResponse{
			ID:   reportId,
			Data: string(export),
			Type: "excel",
		})
	}
}

func writeToTempFile(base64EncodedFile []byte) (string, error) {
	tmpFile, err := os.CreateTemp(os.TempDir(), "greypot_")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()
	exportBytes, err := base64.StdEncoding.DecodeString(string(base64EncodedFile))
	if err != nil {
		return "", err
	}

	_, err = tmpFile.Write(exportBytes)
	if err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}
