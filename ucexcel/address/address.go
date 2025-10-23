package address

import (
	"github.com/xuri/excelize/v2"
)

type address struct {
	row int
	col int
}

// утилиты для адресации к ячейкам
func New(row int, col int) *address {
	return &address{
		row: row,
		col: col,
	}
}

// преобразуем адрес 1 в букву адреса ексель А
func (a *address) Address() string {
	// excelize expects 1-based col/row
	name, _ := excelize.CoordinatesToCellName(a.col+1, a.row)
	return name
}

func (a *address) NextCol() string {
	a.col += 1
	out := a.Address()
	return out
}

func (a *address) NextRow() string {
	a.row += 1
	a.col = 0
	out := a.Address()
	return out
}

func (a *address) MoveTo(row int, col int) string {
	a.row = row
	a.col = col
	out := a.Address()
	return out
}
func (a *address) AddCol(col int) string {
	a.col += col
	out := a.Address()
	return out
}
func (a *address) AddRow(row int) string {
	a.row += row
	out := a.Address()
	return out
}

func (a *address) ShiftCol(col int) string {
	name, _ := excelize.CoordinatesToCellName(a.col+1+col, a.row)
	return name
}

func (a *address) Range(row int, col int) string {
	name, _ := excelize.CoordinatesToCellName(a.col+1+col, a.row+row)
	return name
}

func (a *address) Row() int {
	return a.row
}

func (a *address) Col() int {
	return a.col
}
