package arxlsx

import (
	"fmt"
	"github.com/tealeg/xlsx"
)

type row struct {
	values []interface{}
}

type Row interface {
	AddValue(val interface{})
}

func (r *row) AddValue(val interface{}) {
	r.values = append(r.values, val)
}

type sheet struct {
	name   string
	titles []string
	rows   [][]interface{}
}

type Sheet interface {
	AddRow(...interface{})

	getSheetName() string
	getTitles() []string
	getRows() [][]interface{}
}

func (s *sheet) AddRow(values ...interface{}) {
	r := make([]interface{}, 0)
	for _, v := range values {
		r = append(r, v)
	}
	s.rows = append(s.rows, r)
}

func (s *sheet) getSheetName() string {
	return s.name
}

func (s *sheet) getTitles() []string {
	return s.titles
}

func (s *sheet) getRows() [][]interface{} {
	return s.rows
}

type file struct {
	filePath string
	sheets   []Sheet
}

type File interface {
	AddSheet(sheet2 Sheet)
	Export() error
}

func (xf *file) AddSheet(s Sheet) {
	xf.sheets = append(xf.sheets, s)
}

func (xf *file) Export() error {
	xlsxFile := xlsx.NewFile()

	for _, s := range xf.sheets {
		newSheet, e := xlsxFile.AddSheet(s.getSheetName())
		if e != nil {
			return e
		}

		var cell *xlsx.Cell
		var cellRow *xlsx.Cell
		for _, values := range s.getRows() {
			r := newSheet.AddRow()
			for _, value := range values {
				cell = r.AddCell()
				switch value.(type) {
				case int:
					cell.SetInt(value.(int))
				case int32:
					cell.SetInt64(int64(value.(int32)))
				case int64:
					cell.SetInt64(value.(int64))
				case uint32:
					cell.SetValue(value.(uint32))
				case uint64:
					cell.SetValue(value.(uint64))
				case float64:
					cell.SetFloat(value.(float64))
				case string:
					cell.SetString(value.(string))
				default:
					cell.SetString(fmt.Sprintf("%v", value))
				}
			}
		}

		// 更新表头
		if titles := s.getTitles(); len(titles) > 0 {
			rowTitle, _ := newSheet.AddRowAtIndex(0)
			for _, title := range titles {
				cellRow = rowTitle.AddCell()
				cellRow.Value = title
			}
		}
	}

	if err := xlsxFile.Save(xf.filePath); err != nil {
		return err
	}

	return nil
}

func NewRow() Row {
	return &row{
		values: make([]interface{}, 0),
	}
}

func NewSheet(sheetName string) Sheet {
	return &sheet{
		name:   sheetName,
		titles: make([]string, 0),
		rows:   make([][]interface{}, 0),
	}
}

func NewFile(filePath string) File {
	return &file{
		filePath: filePath,
		sheets:   make([]Sheet, 0),
	}
}
