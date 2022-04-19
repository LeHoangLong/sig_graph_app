package repositories

import (
	"backend/internal/models"
	"database/sql"
	"fmt"
	"strings"
)

type MaterialRepositorySql struct {
	tx             *sql.Tx
	nodeRepository NodeRepositoryI
}

func MakeMaterialRepositorySql(
	iTx *sql.Tx,
	iNodeRepository NodeRepositoryI,
) MaterialRepositorySql {
	return MaterialRepositorySql{
		tx:             iTx,
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

	rows, err := r.tx.Query(`INSERT INTO "material"(
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
	fmt.Println(2)

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

	rows, err := r.tx.Query(statement, args...)
	if err != nil {
		return []models.Material{}, err
	}

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
