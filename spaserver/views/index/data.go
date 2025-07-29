package index

// https://symbl.cc/ru/unicode-table/#basic-latin
func (t *page) PageData() interface{} {
	return struct{ Copyright string }{Copyright: "\u00a9 ООО \u00abНЕВАКОД\u00bb"}
	// return domain.ModelFooter{Copyright: "\u00a9 ООО \u00abНЕВАКОД\u00bb"}
}

func (t *page) InitData() interface{} {
	return struct{ Copyright string }{Copyright: "\u00a9 ООО \u00abНЕВАКОД\u00bb"}
}
