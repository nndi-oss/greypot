+++
date = "2023-02-02-10T06:43:48+02:00"
title = "Advanced: Implementating a Custom Template Engine"
draft = false
weight = 50
description = "Advanced: Implementating a Custom Template Engine"
toc = true
bref = "Advanced.TemplateEngine"
+++


```go
package main

import (
	"bytes"
	"fmt"

	"github.com/PuerkitoBio/goquery"
	"github.com/nndi-oss/greypot/models"
	"github.com/nndi-oss/greypot/template/engine"
	"github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
)

type html2ExcelTemplateEngine struct {
	pongo engine.TemplateEngine
}

func NewHtml2ExcelTemplateEngine() *html2ExcelTemplateEngine {
	return &html2ExcelTemplateEngine{
		pongo: engine.NewDjangoTemplateEngine(),
	}
}

func (pte *html2ExcelTemplateEngine) Render(templateContent []byte, ctx *models.TemplateContext) ([]byte, error) {
	out, err := pte.pongo.Render(templateContent, ctx)
	if err != nil {
		return nil, err
	}

	return pte.ConvertHtmlToExcel(out)
}

func (pte *html2ExcelTemplateEngine) ConvertHtmlToExcel(htmlData []byte) ([]byte, error) {
	var outExcel bytes.Buffer
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(htmlData))
	if err != nil {
		return nil, err
	}
	st := &state{}

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			logrus.WithError(err)
			return
		}
	}()

	doc.Find("table").Each(func(i int, tbl *goquery.Selection) {
		st.NumTablesFound = st.NumTablesFound + 1
		if st.CurrentRowNum == 0 {
			st.IsFirstRow = true
		}

		sheetName := tbl.AttrOr("data-sheet-name", fmt.Sprintf("Sheet %d", st.NumTablesFound))
		_, err := f.NewSheet(sheetName)
		if err != nil {
			logrus.WithError(err)
			return
		}

		for idx, row := range toExcelSheetData(tbl) {
			cell, err := excelize.CoordinatesToCellName(1, idx+1)
			if err != nil {
				return
			}
			f.SetSheetRow(sheetName, cell, &row)
		}

		st.CurrentRowNum++
	})

	// Set active sheet of the workbook.
	f.SetActiveSheet(st.NumTablesFound)
	_, err = f.WriteTo(&outExcel)
	if err != nil {
		return nil, err
	}
	return outExcel.Bytes(), err
}

type state struct {
	TableNameInferred string
	IsFirstRow        bool
	CurrentRowNum     int
	NumTablesFound    int
	NumRowsProcessed  int
}

func toExcelSheetData(tbl *goquery.Selection) [][]any {
	sheetData := make([][]any, 0)
	bodyRows := tbl.Find("tbody > tr")
	if bodyRows != nil {
		if bodyRows.Length() > 0 {
			bodyRows.Each(func(i int, tr *goquery.Selection) {
				tds := tr.Children()
				tdsExcelRow := make([]any, 0)

				tds.Each(func(i int, td *goquery.Selection) {
					tdsExcelRow = append(tdsExcelRow, td.Text())
				})

				sheetData = append(sheetData, tdsExcelRow)
			})
		}
	}
	return sheetData
}

```