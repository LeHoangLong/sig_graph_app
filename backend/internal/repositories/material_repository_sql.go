package repositories

import (
	"backend/internal/models"
	"database/sql"
)

type MaterialRepositorySql struct {
	db *sql.DB
}

func MakeMaterialRepositorySql(
	iDb *sql.DB,
) MaterialRepositorySql {
	return MaterialRepositorySql{
		db: iDb,
	}
}

func (r MaterialRepositorySql) AddMaterial(iUserId int, iMaterial models.Material) error {
	_, err := r.db.Query(`INSERT INTO "material"(
			user_id,
			id,
			name,
			quantity,
			unit,
			created_time
		) VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6
		)`,
		iUserId,
		iMaterial.Id,
		iMaterial.Name,
		iMaterial.Quantity.String(),
		iMaterial.Unit,
		iMaterial.CreatedTime.String(),
	)

	return err
}

func (r MaterialRepositorySql) FetchMaterials(iUserId int) ([]models.Material, error) {
	rows, err := r.db.Query(`SELECT 
			id,
			name,
			quantity,
			unit,
			created_time
		FROM "material" 
		WHERE user_id = $1
	`, iUserId)

	if err != nil {
		return []models.Material{}, err
	}

	ret := []models.Material{}
	for rows.Next() {
		var material models.Material
		rows.Scan(
			&material.Id,
			&material.Name,
			&material.Quantity,
			&material.Unit,
			&material.CreatedTime,
		)

		ret = append(ret, material)
	}

	return ret, nil
}
