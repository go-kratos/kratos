package model

import "encoding/csv"

//CsvFormatter CsvFormatter
type CsvFormatter interface {
	GetFileName() string
	// ToCsv do not call flush
	ToCsv(writer *csv.Writer)
}

//ExportArgInterface export interface
type ExportArgInterface interface {
	// ExportFormat options: json, csv
	ExportFormat() string
}
