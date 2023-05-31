package handlers

import (
	"encoding/base64"
	"encoding/json"

	"github.com/labstack/echo/v5"
	"github.com/palantir/stacktrace"
)

func extractData(ctx echo.Context) (interface{}, error) {
	encodedData := ctx.QueryParam("d")
	if encodedData == "" {
		return map[string]interface{}{}, nil
	}

	strData, err := base64.StdEncoding.DecodeString(encodedData)
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed to decode base64 data")
	}

	var jsonData interface{}
	if err := json.Unmarshal(strData, &jsonData); err != nil {
		return nil, stacktrace.Propagate(err, "failed to unmarshal data to json")
	}

	return jsonData, nil
}
