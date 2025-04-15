package arxlsx

func NewRow(values ...interface{}) Row {
	//r := row{
	//	values: make([]interface{}, 0),
	//}
	//if len(values) > 0 {
	//	r.AddValue(values...)
	//}
	//return &r
	return newRow(values...)
}

func NewSheet(sheetName string) Sheet {
	return &sheet{
		name:    sheetName,
		titles:  make([]string, 0),
		rowList: make([]Row, 0),
	}
}

func NewFile(filePath string) File {
	return &file{
		filePath: filePath,
		sheets:   make([]Sheet, 0),
	}
}
