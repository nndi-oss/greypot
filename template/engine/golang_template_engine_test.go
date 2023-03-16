package engine

import (
	"testing"

	"github.com/nndi-oss/greypot/models"
	"github.com/stretchr/testify/assert"
)

var golangTemplateRepo = NewGolangTemplateEngine()

func TestGolangTemplateRepository_Render(t *testing.T) {
	assertTemplate(t,
		"{{ .Data.name }}",
		map[string]interface{}{"name": "matan"},
		"matan",
	)

	assertTemplate(t,
		"hello {{ .Data.name }}",
		map[string]interface{}{"name": "matan"},
		"hello matan",
	)

	assertTemplate(t, `{{- range .Data.items }}
{{ . }}
{{- end }}
`,
		map[string]interface{}{"items": []int{1, 2, 3, 4}},
		"\n1\n2\n3\n4\n",
	)
}

func assertTemplate(t *testing.T, template string, data interface{}, result string) {
	rendered, err := golangTemplateRepo.Render([]byte(template), &models.TemplateContext{Data: data})
	assert.Nil(t, err)
	assert.EqualValues(t, result, rendered)
}
