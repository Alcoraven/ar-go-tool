package arxlsx

type sheet struct {
	name    string
	titles  []string
	rowList []Row
}

func newSheet(name string) *sheet {
	return &sheet{
		name:    name,
		titles:  make([]string, 0),
		rowList: make([]Row, 0),
	}
}

type Sheet interface {
	AddRow(Row)                      // add a new row for sheet with Row
	AddRowWithValues(...interface{}) // add a new row for sheet with values
	CreateRow() Row                  // create a new row for sheet and return it
	GetSheetName() string            // return sheet.name of the sheet
	GetTitles() []string             // return sheet.titles of the sheet
	GetRows() []Row                  // return sheet.rowList of the sheet
}

func (s *sheet) AddRow(r Row) {
	s.rowList = append(s.rowList, r)
}

func (s *sheet) AddRowWithValues(values ...interface{}) {
	if len(values) == 0 {
		return
	}
	r := new(row)
	for _, v := range values {
		r.AddValue(v)
	}
	s.AddRow(r)
}

func (s *sheet) CreateRow() Row {
	r := new(row)
	return r
}

func (s *sheet) GetSheetName() string {
	return s.name
}

func (s *sheet) GetTitles() []string {
	return s.titles
}

func (s *sheet) GetRows() []Row {
	return s.rowList
}
