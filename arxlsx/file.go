package arxlsx

import (
	"errors"
	"fmt"
	"github.com/tealeg/xlsx"
	"runtime/debug"
)

type file struct {
	filePath string
	sheets   []Sheet
}

type File interface {
	CreateSheet(sheetName string) (Sheet, error) // add a sheet for xlsx file and return it
	AddSheet(Sheet)                              // add a sheet for xlsx file

	Export() error // export xlsx file
}

func (xf *file) AddSheet(s Sheet) {
	xf.sheets = append(xf.sheets, s)
}

func (xf *file) CreateSheet(sheetName string) (Sheet, error) {
	if sheetName == "" {
		sheetName = fmt.Sprintf("Sheet%d", len(xf.sheets)+1)
	} else {
		for _, s := range xf.sheets {
			if s.GetSheetName() == sheetName {
				return nil, errors.New("sheet name already exists")
			}
		}
	}

	s := newSheet(sheetName)

	xf.sheets = append(xf.sheets, s)
	return s, nil
}

func (xf *file) ChangePath(newPath string) {
	xf.filePath = newPath
}

func (xf *file) Export() error {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println("recovered from panic:", e)
			fmt.Println(string(debug.Stack()))
		}
	}()
	xlsxFile := xlsx.NewFile()

	for _, s := range xf.sheets {
		ns, e := xlsxFile.AddSheet(s.GetSheetName())
		if e != nil {
			panic(e)
		}

		for _, r := range s.GetRows() {
			nr := ns.AddRow()
			for _, value := range r.GetValues() {
				cell := nr.AddCell()
				switch value.(type) {
				case int:
					cell.SetInt(value.(int))
				case int32:
					cell.SetInt64(int64(value.(int32)))
				case int64:
					if v := value.(int64); v > 999999999999999 {
						cell.SetString(fmt.Sprintf("%v", value))
					} else {
						cell.SetInt64(value.(int64))
					}
				case uint:
					cell.SetInt64(int64(value.(uint)))
				case uint32:
					cell.SetInt64(int64(value.(uint32)))
				case uint64:
					if v := value.(uint64); v > 999999999999999 {
						cell.SetString(fmt.Sprintf("%v", value))
					} else {
						cell.SetInt64(int64(value.(uint64)))
					}
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
		if titles := s.GetTitles(); len(titles) > 0 {
			rowTitle, _ := ns.AddRowAtIndex(0)
			for _, title := range titles {
				cell := rowTitle.AddCell()
				cell.Value = title
			}
		}
	}

	if err := xlsxFile.Save(xf.filePath); err != nil {
		panic(err)
	}

	return nil
}
