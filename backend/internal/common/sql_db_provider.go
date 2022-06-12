package common

import (
	"database/sql"

	"go.uber.org/dig"
)

func ProvideSqlDb(
	iContainer *dig.Container,
) error {
	return iContainer.Provide(func(iConfig Config) (*sql.DB, error) {
		var err error
		connStr := ""
		connStr += " user=" + iConfig.DbUser
		connStr += " dbname=" + iConfig.DbName
		connStr += " password=" + iConfig.DbPassword
		connStr += " host=" + iConfig.DbHost
		connStr += " sslmode=" + iConfig.DbSslmode
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			return nil, err
		}
		return db, nil
	})
}
