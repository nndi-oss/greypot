package echo

import (
	"github.com/labstack/echo/v5"
	"github.com/nndi-oss/greypot"
	"github.com/nndi-oss/greypot/http/echo/handlers"
)

func UseGroup(app *echo.Group, s *greypot.Module) {
	app.GET("/reports/list", handlers.ReportListHandler(s.TemplateService))
	app.GET("/reports/preview/*", handlers.ReportPreviewHandler(s.TemplateService, s.TemplateEngine, s.ReportService))
	app.GET("/reports/render/*", handlers.ReportRenderHandlder(s.TemplateService))

	app.POST("/reports/export/html/*", handlers.ReportExportHandler(s.ReportService, "html"))
	app.POST("/reports/export/png/*", handlers.ReportExportHandler(s.ReportService, "png"))
	app.POST("/reports/export/pdf/*", handlers.ReportExportHandler(s.ReportService, "pdf"))

	app.POST("/reports/export/bulk/html/*", handlers.BulkReportExportHandler(s.ReportService, "html"))
	app.POST("/reports/export/bulk/png/*", handlers.BulkReportExportHandler(s.ReportService, "png"))
	app.POST("/reports/export/bulk/pdf/*", handlers.BulkReportExportHandler(s.ReportService, "pdf"))
}

func Use(app *echo.Echo, s *greypot.Module) {
	app.GET("/reports/list", handlers.ReportListHandler(s.TemplateService))
	app.GET("/reports/preview/*", handlers.ReportPreviewHandler(s.TemplateService, s.TemplateEngine, s.ReportService))
	app.GET("/reports/render/*", handlers.ReportRenderHandlder(s.TemplateService))

	app.POST("/reports/export/html/*", handlers.ReportExportHandler(s.ReportService, "html"))
	app.POST("/reports/export/png/*", handlers.ReportExportHandler(s.ReportService, "png"))
	app.POST("/reports/export/pdf/*", handlers.ReportExportHandler(s.ReportService, "pdf"))

	app.POST("/reports/export/bulk/html/*", handlers.BulkReportExportHandler(s.ReportService, "html"))
	app.POST("/reports/export/bulk/png/*", handlers.BulkReportExportHandler(s.ReportService, "png"))
	app.POST("/reports/export/bulk/pdf/*", handlers.BulkReportExportHandler(s.ReportService, "pdf"))
}
