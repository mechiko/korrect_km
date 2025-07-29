package templates

import (
	"html/template"
	"time"
)

// functions is a template.FuncMap providing custom date formatting functions for use in HTML templates.
var functions = template.FuncMap{
	"humanDate": func(t time.Time) string {
		if t.IsZero() {
			return ""
		}
		return t.Format("02 Jan 2006 at 15:04")
	},
	"formatDate": func(t time.Time) string {
		if t.IsZero() {
			return ""
		}
		return t.Format("2006-01-02")
	},
	"iterate": func(count int) []int {
		if count <= 0 {
			return []int{}
		}
		numbers := make([]int, count)
		for i := 0; i < count; i++ {
			numbers[i] = i
		}
		return numbers
	},
	"add": func(a, b int) int {
		return a + b
	},
	"sub": func(a, b int) int {
		return a - b
	},
	"mul": func(a, b int) int {
		return a * b
	},
	"min": func(a, b int) int {
		if a < b {
			return a
		}
		return b
	},
	"inc": func(i int) int {
		return i + 1
	},
	"noescape": func(str string) template.HTML {
		return template.HTML(str)
	},
}
