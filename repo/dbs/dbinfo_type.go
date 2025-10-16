package dbs

import "korrectkm/config"

// host
// user
// password
// name имя файла без пути
// file полное имя файла и путь
// name имя БД
// driver драйвер
// Exists файл существует и только
type DbInfo struct {
	Host       string
	User       string
	Pass       string
	File       string
	Name       string
	Driver     string
	Connection string
	Exists     bool // только для sqlite делает поиск файла
}

func NewDbInfo(cfg config.DatabaseConfiguration) *DbInfo {
	return &DbInfo{
		Host:       cfg.Host,
		User:       cfg.User,
		Pass:       cfg.Pass,
		File:       cfg.File,
		Name:       cfg.DbName,
		Driver:     cfg.Driver,
		Connection: cfg.Connection,
	}
}
