package arxlsx

type sheet struct {
	name    string
	titles  []string
	rowList []Row
}

func newSheet(name string) *sheet {
	var s = &sheet{
		name:    name,
		titles:  make([]string, 0),
		rowList: make([]Row, 0),
	}
	return s
}

type Sheet interface {
	AddRow(Row)                      // add a new row for sheet with Row
	AddRowWithValues(...interface{}) // add a new row for sheet with values
	CreateRow() Row                  // create a new row for sheet and return it
	GetSheetName() string            // return sheet.name of the sheet
	SetTitles([]string)              // set sheet.titles
	GetTitles() []string             // return sheet.titles of the sheet
	GetRows() []Row                  // return sheet.rowList of the sheet
}

func (s *sheet) AddRow(r Row) {
	s.rowList = append(s.rowList, r)
}

func (s *sheet) AddRowWithValues(values ...interface{}) {
	r := newRow(values...)
	s.AddRow(r)
}

func (s *sheet) CreateRow() Row {
	r := newRow()
	s.AddRow(r)
	return r
}

func (s *sheet) GetSheetName() string {
	return s.name
}

func (s *sheet) SetTitles(titles []string) {
	s.titles = titles
}

func (s *sheet) GetTitles() []string {
	return s.titles
}

func (s *sheet) GetRows() []Row {
	return s.rowList
}
