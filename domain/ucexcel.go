package domain

type ExcelAddress interface {
	Address() string
	NextCol() string
	NextRow() string
	MoveTo(row int, col int) string
	ShiftCol(col int) string
	Range(row int, col int) string
	AddCol(int) string
	AddRow(int) string
	Row() int
	Col() int
}
