package arxlsx

type row struct {
	values []interface{}
}

func (r *row) addValue(value interface{}) {
	r.values = append(r.values, value)
}

func newRow(values ...interface{}) Row {
	r := new(row)
	if len(values) > 0 {
		for _, v := range values {
			r.addValue(v)
		}
	}
	return r
}

type Row interface {
	AddValue(vals ...interface{})
	GetValues() []interface{}
}

func (r *row) AddValue(vals ...interface{}) {
	for _, v := range vals {
		r.addValue(v)
	}
}

func (r *row) GetValues() []interface{} {
	return r.values
}
