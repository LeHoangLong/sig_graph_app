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

func (r MaterialRepositorySql) AddOrUpdateMaterial(
	iPublicKeyId int,
	iMaterial models.Material,
) error {

	_, err := r.db.Query(`INSERT INTO "material"(
			id,
			public_key_id,
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
		) 
		ON CONFLICT (id) DO UPDATE
		SET 
			public_key_id = $2,
			name = $3,
			quantity = $4,
			unit = $5,
			created_time = $6
		`,
		iMaterial.Id,
		iPublicKeyId,
		iMaterial.Name,
		iMaterial.Quantity.String(),
		iMaterial.Unit,
		iMaterial.CreatedTime.String(),
	)

	return err
}

func (r MaterialRepositorySql) FetchMaterials(
	iPublicKeyId int,
) ([]models.Material, error) {
	rows, err := r.db.Query(`SELECT 
			id,
			name,
			quantity,
			unit,
			created_time
		FROM "material" 
		WHERE public_key_id = $1
	`, iPublicKeyId)

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
