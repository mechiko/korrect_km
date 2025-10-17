package app

import (
	"fmt"
	"html/template"
	"os"
	"reflect"
	"strings"
)

func FilePackageContents(filename string) string {
	b, err := os.ReadFile(filename)
	if err != nil {
		panic(fmt.Errorf("%s %s", err.Error(), filename))
	}
	return string(b)
}

var funcMapHtml = template.FuncMap{
	// The name "inc" is what the function will be called in the template text.
	"inc": func(i int) int {
		return i + 1
	},
	"noescape": func(str string) template.HTML {
		return template.HTML(str)
	},
}

// func getDirectory() string {
// 	_, file, _, ok := runtime.Caller(1)
// 	if ok {
// 		return path.Dir(file)
// 	}

// 	return ""
// }

// возвращает строку содержащую список полей структуры через запятую
// s pointer struct
func ListFieldStruct(s interface{}) (out string) {
	outArray := make([]string, 0)
	t := reflect.TypeOf(s)
	kind := t.Kind()
	if kind != reflect.Pointer {
		return ""
	}
	elemS := t.Elem()
	if kind != reflect.Struct {
		return ""
	}
	for i := 0; i < elemS.NumField(); i++ {
		field := elemS.Field(i)
		outArray = append(outArray, field.Name)
	}
	out = strings.Join(outArray, ",")
	return out
}

func ArrayFieldStruct(s interface{}) (out []string) {
	out = make([]string, 0)
	t := reflect.TypeOf(s)
	kind := t.Kind()
	if kind != reflect.Pointer {
		return out
	}
	elemS := t.Elem()
	if kind != reflect.Struct {
		return out
	}
	for i := 0; i < elemS.NumField(); i++ {
		field := elemS.Field(i)
		out = append(out, field.Name)
	}
	return out
}

// возвращает строку содержащую значение поля по имени
// s pointer struct
func FieldByNameStruct(s interface{}, name string) (out string) {
	t := reflect.ValueOf(s)
	kind := t.Kind()
	if kind != reflect.Pointer {
		return ""
	}
	elemS := t.Elem()
	if kind != reflect.Struct {
		return ""
	}
	out = elemS.FieldByName(name).String()
	return out
}
