package repositories

import (
	"backend/internal/services/node_contract"
	"context"
	"database/sql"
)

type MaterialRepositoryFactory struct {
	db     *sql.DB
	hasher node_contract.IdHasherI
}

func MakeMaterialRepositoryFactory(
	iDb *sql.DB,
	iHasher node_contract.IdHasherI,
) MaterialRepositoryFactory {
	return MaterialRepositoryFactory{
		db:     iDb,
		hasher: iHasher,
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

	nodeRepository := MakeNodeRepositorySql(f.db, f.hasher)
	materialRepository := MakeMaterialRepositorySql(f.db, nodeRepository)

	err = iHandler(materialRepository)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
