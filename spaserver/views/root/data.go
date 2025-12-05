package root

import "korrectkm/domain"

// https://symbl.cc/ru/unicode-table/#basic-latin
func (t *page) PageData() (interface{}, error) {
	return struct{ Copyright string }{Copyright: "\u00a9 ООО \u00abНЕВАКОД\u00bb"}, nil
	// return domain.ModelFooter{Copyright: "\u00a9 ООО \u00abНЕВАКОД\u00bb"}
}

func (t *page) InitData(_ domain.Apper) (interface{}, error) {
	return struct{ Copyright string }{Copyright: "\u00a9 ООО \u00abНЕВАКОД\u00bb"}, nil
}
