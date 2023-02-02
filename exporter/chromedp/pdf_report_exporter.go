package chromedp

import (
	"time"

	"github.com/nndi-oss/greypot/models"
	"github.com/palantir/stacktrace"
	"github.com/signintech/gopdf"
)

type pdfReportExporter struct {
	png *pngReportExporter
}

func NewPdfReportExporter(timeout time.Duration, vpHeight, vpWidth int) *pdfReportExporter {
	return &pdfReportExporter{png: NewPngReportExporter(timeout, vpHeight, vpWidth)}
}

func (pre *pdfReportExporter) Export(url string, renderedTemplate []byte) ([]byte, *models.PrintOptions, error) {
	png, opts, err := pre.png.Export(url, renderedTemplate)
	if err != nil {
		return nil, nil, stacktrace.Propagate(err, "failed to export pdf, image screenshot has a failure for url: %s", url)
	}

	pdf := gopdf.GoPdf{}
	defer pdf.Close()

	rect := gopdf.Rect{
		W: opts.PageWidth,
		H: opts.PageHeight,
	}
	rect.UnitsToPoints(gopdf.UnitPT)

	pdf.Start(gopdf.Config{PageSize: rect})

	pdf.AddPage()
	image, _ := gopdf.ImageHolderByBytes(png)
	if err := pdf.ImageByHolder(image, 0, 0, nil); err != nil {
		return nil, nil, stacktrace.Propagate(err, "failed to embed image into pdf for url : %s", url)
	}

	return pdf.GetBytesPdf(), opts, nil
}
