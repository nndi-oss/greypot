package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/palantir/stacktrace"
)

type ExportType string

const (
	ExportType_HTML  ExportType = "html"
	ExportType_PDF   ExportType = "pdf"
	ExportType_PNG   ExportType = "png"
	ExportType_EXCEL ExportType = "excel"
)

type ReportBuilderClient struct {
	baseUrl string
}

func NewClient(baseUrl string) *ReportBuilderClient {
	return &ReportBuilderClient{baseUrl: baseUrl}
}

func (c *ReportBuilderClient) Get(reportId string, data interface{}, exportType ExportType) ([]byte, error) {
	url := fmt.Sprintf("http://%s/reports/export/%s/%s", c.baseUrl, exportType, reportId)
	body, err := json.Marshal(data)
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed to marshal: %+v to json", data)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed to POST: %s with data: %+v", url, data)
	}

	if resp.StatusCode != 200 {
		return nil, stacktrace.Propagate(err, "failed to POST: %s with data: %+v returned status: %d", url, data, resp.StatusCode)
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed to read response body")
	}

	var jsonResp response
	if err := json.Unmarshal(respBody, &jsonResp); err != nil {
		return nil, stacktrace.Propagate(err, "failed to unmarshal response body from json")
	}

	return []byte(jsonResp.Data), nil
}

func (c *ReportBuilderClient) GetBulk(reportId string, data interface{}, exportType ExportType) (*bulkResponse, error) {
	url := fmt.Sprintf("http://%s/reports/export/bulk/%s/%s", c.baseUrl, exportType, reportId)
	body, err := json.Marshal(data)
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed to marshal: %+v to json", data)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed to POST: %s with data: %+v", url, data)
	}

	if resp.StatusCode != 200 {
		return nil, stacktrace.Propagate(err, "failed to POST: %s with data: %+v returned status: %d", url, data, resp.StatusCode)
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed to read response body")
	}

	jsonResp := new(bulkResponse)
	if err := json.Unmarshal(respBody, jsonResp); err != nil {
		return nil, stacktrace.Propagate(err, "failed to unmarshal response body from json")
	}

	return jsonResp, nil
}

type response struct {
	Data     string `json:"data"`
	ReportId string `json:"reportId"`
}

type bulkResponse struct {
	ID       string     `json:"_id"`
	ReportId string     `json:"reportId"`
	Reports  []response `json:"reports"`
}
