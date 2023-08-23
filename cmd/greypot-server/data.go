package main

import "github.com/nndi-oss/greypot/http/fiber/handlers"

type UploadTemplateRequest struct {
	Name     string
	Template string
	Data     any
}

type BulkUploadTemplateRequest struct {
	Name     string
	Template string
	Data     any
	Entries  []handlers.BulkExportEntry
}
