package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nndi-oss/greypot/models"
	"github.com/nndi-oss/greypot/service"
	"github.com/nndi-oss/greypot/template/engine"
	"github.com/sirupsen/logrus"
)

func ReportPreviewHandler(templateService service.TemplateService,
	templateEngine engine.TemplateEngine,
	reportService service.ReportService,
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		reportId := strings.TrimPrefix(ctx.Param("reportId"), "/")
		data, err := extractData(ctx)
		if err != nil {
			logrus.Error(err)
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		templates, err := templateService.ListTemplates()
		if err != nil {
			logrus.Error(err)
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		png, err := reportService.ExportReportPng(reportId, data)
		if err != nil {
			logrus.Error(err)
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		rendered, err := templateEngine.Render(
			[]byte(previewTemplate),
			&models.TemplateContext{Data: map[string]interface{}{
				"reportId": reportId,
				"reports":  templates,
				"data":     data,
				"image":    string(png),
			}},
		)
		if err != nil {
			logrus.Error(err)
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.Data(http.StatusOK, "text/html; charset=utf-8", rendered)
	}
}

const previewTemplate = `
<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
	<link href="https://cdnjs.cloudflare.com/ajax/libs/jsoneditor/9.1.8/jsoneditor.min.css" rel="stylesheet" type="text/css">
    <script src="https://cdnjs.cloudflare.com/ajax/libs/jsoneditor/9.1.8/jsoneditor.min.js"></script>
	<style>
		
	</style>
</head>
<body>

<div style="display: flex; flex-direction: column; justify-content: center; align-items: center ">
	<h1 style="text-align: center">Preview Tool</h1>

	<div>
		<select id="reports">
		{{- range $c := .Values.reports }}
			{{- if (eq $c $.Values.reportId) }}
			<option selected="selected" value="{{ $c }}">{{ $c }}</option>
			{{- else }}
			<option value="{{ $c }}">{{ $c }}</option>
			{{- end }}
		{{- end }}
		</select>
		<button onclick="onGenerate()">Generate</button>
	</div>

	<div id="jsoneditor" style="width: 600px; height: 400px;"></div>
	<p>Result:</p>
	<div style="border: 2px solid; padding: 20px">
		<img src="data:image/png;base64, {{ .Values.image }}" alt="Report" />
	</div>

	<script>
		const container = document.getElementById("jsoneditor")
		const options = {
			mode: 'code',
			modes: ['code', 'form', 'text', 'tree', 'view', 'preview'], // allowed modes
	  	}
		const editor = new JSONEditor(container, options)
		const initialJson = {{ .Values.data }}
		editor.set(initialJson)

		function onGenerate() {
			const json = editor.get()
			const objJsonStr = JSON.stringify(json);
			const objJsonB64 = btoa(objJsonStr);

			var e = document.getElementById("reports");
			var reportId = e.value;

			window.location.href = "http://localhost:8080/reports/preview/" + reportId + "?d=" + objJsonB64;
		}
	</script>
</div>

</body>
</html>
`
