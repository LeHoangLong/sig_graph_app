package repositories

import (
	"backend/internal/common"
	"backend/internal/models"
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type MaterialRepositorySql struct {
	db             *sql.DB
	nodeRepository NodeRepositoryI
}

func MakeMaterialRepositorySql(
	iDb *sql.DB,
	iNodeRepository NodeRepositoryI,
) MaterialRepositorySql {
	return MaterialRepositorySql{
		db:             iDb,
		nodeRepository: iNodeRepository,
	}
}

func (r MaterialRepositorySql) AddMaterial(
	iMaterial models.Material,
) (models.Material, error) {
	node, err := r.nodeRepository.CreateNode(iMaterial.Node)
	if err != nil {
		return models.Material{}, err
	}

	rows, err := r.db.Query(`INSERT INTO "material"(
			node_id,
			name,
			quantity,
			unit
		) VALUES (
			$1,
			$2,
			$3,
			$4
		)
		`,
		node.Id,
		iMaterial.Name,
		iMaterial.Quantity.String(),
		iMaterial.Unit,
	)
	rows.Close()
	if err != nil {
		return models.Material{}, err
	}

	createdMaterial := iMaterial
	createdMaterial.Node = node

	return createdMaterial, err
}

func (r MaterialRepositorySql) FetchMaterialsByOwner(
	iOwnerKey models.PublicKey,
	iMinId int,
	iLimit int,
) ([]models.Material, error) {
	nodes, err := r.nodeRepository.FetchNodesByOwnerKey(iOwnerKey, iMinId, iLimit)
	if err != nil {
		return []models.Material{}, err
	}
	statement := `SELECT 
			name,
			quantity,
			unit
		FROM "material" 
		WHERE 
	`

	count := 1
	argsClause := []string{}
	args := []interface{}{}
	for _, node := range nodes {
		argsClause = append(argsClause, fmt.Sprintf("(node_id = $%d)", count))
		count++
		args = append(args, node.Id)
	}

	if len(nodes) > 0 {
		statement += strings.Join(argsClause, " OR ")
	}

	fmt.Println("statement")
	fmt.Println(statement)

	rows, err := r.db.Query(statement, args...)
	if err != nil {
		return []models.Material{}, err
	}
	defer rows.Close()

	ret := []models.Material{}
	i := 0
	for rows.Next() {
		var material models.Material
		rows.Scan(
			&material.Name,
			&material.Quantity,
			&material.Unit,
		)

		material.Node = nodes[i]
		ret = append(ret, material)
		i += 1
	}

	return ret, nil
}

func (r MaterialRepositorySql) FetchMaterialById(
	iContext context.Context,
	iId int,
) (models.Material, error) {
	nodes, err := r.nodeRepository.FetchNodesById([]int{iId})
	if err != nil {
		return models.Material{}, err
	}
	if len(nodes) == 0 {
		return models.Material{}, common.NotFound
	}
	node := nodes[0]
	row := r.db.QueryRowContext(iContext, `
		SELECT 
			name,
			quantity,
			unit
		FROM "material" 
		WHERE node_id = $1
	`, *node.Id)

	var name string
	var quantity models.CustomDecimal
	var unit string

	err = row.Scan(
		&name,
		&quantity,
		&unit,
	)
	if err == sql.ErrNoRows {
		return models.Material{}, common.NotFound
	} else if err != nil {
		return models.Material{}, err
	}

	material := models.NewMaterial(
		node,
		name,
		quantity,
		unit,
	)
	return material, nil
}
