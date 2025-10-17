package selfdb

import (
	"database/sql"
	"fmt"
	"korrectkm/repo/selfdb/migrations"
	"path/filepath"

	"github.com/pressly/goose/v3"
	"github.com/upper/db/v4/adapter/sqlite"
)

// проверка наличия создание и миграция
// вызывается однажды при старте программы
func (r *DbSelf) Check() (err error) {
	defer func() {
		if rr := recover(); rr != nil {
			err = fmt.Errorf("%s check panic %v", modError, rr)
		}
	}()

	if !r.dbInfo.Exists {
		uri := r.dbInfo.SqliteUri(filepath.Join(r.dbInfo.Path, r.dbInfo.File))
		uri.Options["mode"] = "rwc"
		r.dbSession, err = sqlite.Open(uri)
	} else {
		r.dbSession, err = r.dbInfo.Connect()
	}
	if err != nil {
		r.dbInfo.Exists = false
		return fmt.Errorf("%s prepareSelf ошибка подключения к БД %w", modError, err)
	}

	db, ok := r.dbSession.Driver().(*sql.DB)
	if !ok {
		r.dbInfo.Exists = false
		r.Close()
		return fmt.Errorf("%s prepareSelf ошибка получения *sql.DB %T", modError, r.dbSession.Driver())
	}
	dialect := r.dbInfo.Driver
	switch dialect {
	case "sqlite":
		if err := r.makeMigrationsSqlite(db); err != nil {
			r.dbInfo.Exists = false
			r.Close()
			return fmt.Errorf("%s %w", modError, err)
		}
	case "mssql":
		uri := r.dbInfo.MssqlUri().String()
		if err := r.makeMigrationsMs(uri); err != nil {
			r.dbInfo.Exists = false
			r.Close()
			return fmt.Errorf("%s %w", modError, err)
		}
	default:
		r.dbInfo.Exists = false
		r.Close()
		return fmt.Errorf("%s ошибка драйвера %s", modError, dialect)
	}
	// пробуем получить версию миграции
	if r.version, err = goose.GetDBVersion(db); err != nil {
		r.dbInfo.Exists = false
		r.Close()
		return fmt.Errorf("%s %w", modError, err)
	}

	r.dbInfo.Exists = true
	return nil
}

func (r *DbSelf) makeMigrationsSqlite(DB *sql.DB) error {
	goose.SetBaseFS(migrations.Sqlite)
	if err := goose.SetDialect("sqlite3"); err != nil {
		return err
	}
	if err := goose.Up(DB, "sqlite"); err != nil {
		return err
	}
	return nil
}

func (r *DbSelf) makeMigrationsMs(uri string) error {
	goose.SetBaseFS(migrations.Mssql)
	if err := goose.SetDialect("mssql"); err != nil {
		return fmt.Errorf("failed to set MSSQL dialect: %w", err)
	}
	dbGoose, err := goose.OpenDBWithDriver("mssql", uri)
	if err == nil {
		return fmt.Errorf("failed to open MSSQL connection: %w", err)
	}
	defer dbGoose.Close()
	if err := goose.Up(dbGoose, "mssql"); err != nil {
		return fmt.Errorf("failed to run MSSQL migrations: %w", err)
	}
	return nil
}
