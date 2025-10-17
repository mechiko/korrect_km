package application

import (
	"fmt"
	"strconv"
	"time"
	_ "time/tzdata"
)

var moscowTZ *time.Location

func init() {
	var err error
	moscowTZ, err = time.LoadLocation("Europe/Moscow")
	if err != nil {
		panic(fmt.Errorf("failed to load Moscow timezone: %w", err))
	}
}

func (a *Application) InitDateKv() {
	t := time.Now().In(moscowTZ)
	_, m, _ := t.Date()
	switch (int(m) - 1) / 3 {
	case 0:
		a.startTime = time.Date(t.Year(), 1, 1, 1, 0, 0, 0, moscowTZ)
	case 1:
		a.startTime = time.Date(t.Year(), 4, 1, 1, 0, 0, 0, moscowTZ)
	case 2:
		a.startTime = time.Date(t.Year(), 7, 1, 1, 0, 0, 0, moscowTZ)
	case 3:
		a.startTime = time.Date(t.Year(), 10, 1, 1, 0, 0, 0, moscowTZ)
	}
	a.endTime = a.startTime.AddDate(0, 3, -1)
}

func (a *Application) InitDateMn() {
	t := time.Now().In(moscowTZ)
	_, m, _ := t.Date()
	a.startTime = time.Date(t.Year(), m, 1, 1, 0, 0, 0, moscowTZ)
	a.endTime = a.startTime.AddDate(0, 1, -1)
	a.period = a.MnCalc()
}

func (a *Application) NowDateString() string {
	n := time.Now()
	return fmt.Sprintf("%4d.%02d.%02d %02d:%02d:%02d", n.Local().Year(), n.Local().Month(), n.Local().Day(), n.Local().Hour(), n.Local().Minute(), n.Local().Second())
}

func (a *Application) StartDateString() string {
	return fmt.Sprintf("%4d.%02d.%02d", a.startTime.Local().Year(), a.startTime.Local().Month(), a.startTime.Local().Day())
}

func (a *Application) EndDateString() string {
	return fmt.Sprintf("%4d.%02d.%02d", a.endTime.Local().Year(), a.endTime.Local().Month(), a.endTime.Local().Day())
}

func (a *Application) SetStartDate(d time.Time) {
	a.startTime = d
}

func (a *Application) SetEndDate(d time.Time) {
	a.endTime = d
}

func (a *Application) StartDate() time.Time {
	return a.startTime
}

func (a *Application) EndDate() time.Time {
	return a.endTime
}

func (a *Application) KvartalCalc() string {
	_, m, _ := a.startTime.Date()
	return strconv.Itoa(((int(m) - 1) / 3) + 1)
}

func (a *Application) MnCalc() (out string) {
	_, m, _ := a.startTime.Date()
	switch m {
	case 1:
		out = "янв"
	case 2:
		out = "фев"
	case 3:
		out = "мар"
	case 4:
		out = "апр"
	case 5:
		out = "май"
	case 6:
		out = "июн"
	case 7:
		out = "июл"
	case 8:
		out = "авг"
	case 9:
		out = "сен"
	case 10:
		out = "окт"
	case 11:
		out = "ноя"
	case 12:
		out = "дек"
	}
	return out
}

func (a *Application) SetKvartal(k int) {
	t := time.Now().In(moscowTZ)
	switch k {
	case 1:
		a.startTime = time.Date(t.Year(), 1, 1, 1, 0, 0, 0, moscowTZ)
	case 2:
		a.startTime = time.Date(t.Year(), 4, 1, 1, 0, 0, 0, moscowTZ)
	case 3:
		a.startTime = time.Date(t.Year(), 7, 1, 1, 0, 0, 0, moscowTZ)
	case 4:
		a.startTime = time.Date(t.Year(), 10, 1, 1, 0, 0, 0, moscowTZ)
	}
	a.endTime = a.startTime.AddDate(0, 3, -1)
}

func (a *Application) SetPeriod(y int, p string) {
	switch p {
	case "янв":
		a.startTime = time.Date(y, 1, 1, 1, 0, 0, 0, moscowTZ)
		a.endTime = a.startTime.AddDate(0, 1, -1)
	case "фев":
		a.startTime = time.Date(y, 2, 1, 1, 0, 0, 0, moscowTZ)
		a.endTime = a.startTime.AddDate(0, 1, -1)
	case "мар":
		a.startTime = time.Date(y, 3, 1, 1, 0, 0, 0, moscowTZ)
		a.endTime = a.startTime.AddDate(0, 1, -1)
	case "1кв":
		a.startTime = time.Date(y, 1, 1, 1, 0, 0, 0, moscowTZ)
		a.endTime = a.startTime.AddDate(0, 3, -1)
	case "апр":
		a.startTime = time.Date(y, 4, 1, 1, 0, 0, 0, moscowTZ)
		a.endTime = a.startTime.AddDate(0, 1, -1)
	case "май":
		a.startTime = time.Date(y, 5, 1, 1, 0, 0, 0, moscowTZ)
		a.endTime = a.startTime.AddDate(0, 1, -1)
	case "июн":
		a.startTime = time.Date(y, 6, 1, 1, 0, 0, 0, moscowTZ)
		a.endTime = a.startTime.AddDate(0, 1, -1)
	case "2кв":
		a.startTime = time.Date(y, 4, 1, 1, 0, 0, 0, moscowTZ)
		a.endTime = a.startTime.AddDate(0, 3, -1)
	case "июл":
		a.startTime = time.Date(y, 7, 1, 1, 0, 0, 0, moscowTZ)
		a.endTime = a.startTime.AddDate(0, 1, -1)
	case "авг":
		a.startTime = time.Date(y, 8, 1, 1, 0, 0, 0, moscowTZ)
		a.endTime = a.startTime.AddDate(0, 1, -1)
	case "сен":
		a.startTime = time.Date(y, 9, 1, 1, 0, 0, 0, moscowTZ)
		a.endTime = a.startTime.AddDate(0, 1, -1)
	case "3кв":
		a.startTime = time.Date(y, 7, 1, 1, 0, 0, 0, moscowTZ)
		a.endTime = a.startTime.AddDate(0, 3, -1)
	case "окт":
		a.startTime = time.Date(y, 10, 1, 1, 0, 0, 0, moscowTZ)
		a.endTime = a.startTime.AddDate(0, 1, -1)
	case "ноя":
		a.startTime = time.Date(y, 11, 1, 1, 0, 0, 0, moscowTZ)
		a.endTime = a.startTime.AddDate(0, 1, -1)
	case "дек":
		a.startTime = time.Date(y, 12, 1, 1, 0, 0, 0, moscowTZ)
		a.endTime = a.startTime.AddDate(0, 1, -1)
	case "4кв":
		a.startTime = time.Date(y, 10, 1, 1, 0, 0, 0, moscowTZ)
		a.endTime = a.startTime.AddDate(0, 3, -1)
	}
	a.period = p
}

func (a *Application) Period() string {
	return a.period
}
