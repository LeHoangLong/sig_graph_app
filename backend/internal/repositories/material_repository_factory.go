package repositories

import (
	"context"
	"database/sql"
)

type MaterialRepositoryFactory struct {
	db *sql.DB
}

func MakeMaterialRepositoryFactory(
	iDb *sql.DB,
) MaterialRepositoryFactory {
	return MaterialRepositoryFactory{
		db: iDb,
	}
}

type MaterialRepositoryHandler func(iRepository MaterialRepositoryI) error

func (f MaterialRepositoryFactory) GetRepository(
	iContext context.Context,
	iHandler MaterialRepositoryHandler,
) error {
	tx, err := f.db.BeginTx(iContext, nil)
	if err != nil {
		return nil
	}

	nodeRepository := MakeNodeRepositorySql(tx)
	materialRepository := MakeMaterialRepositorySql(tx, nodeRepository)

	err = iHandler(materialRepository)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
