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
	iContext context.Context,
	iNamespace string,
	iOwnerKey models.PublicKey,
	iMinId int,
	iLimit int,
) ([]models.Material, error) {
	nodes, err := r.nodeRepository.FetchNodesByOwnerKey(
		iContext,
		iNamespace,
		iOwnerKey,
		iMinId,
		iLimit,
	)

	if err != nil {
		return []models.Material{}, err
	}

	if len(nodes) == 0 {
		return []models.Material{}, nil
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

func (r MaterialRepositorySql) fetchMaterialByNode(
	iContext context.Context,
	iNode models.Node,
) (models.Material, error) {
	row := r.db.QueryRowContext(iContext, `
		SELECT 
			name,
			quantity,
			unit
		FROM "material" 
		WHERE node_id = $1
	`, *iNode.Id)
	fmt.Println(123)

	var name string
	var quantity models.CustomDecimal
	var unit string

	err := row.Scan(
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
		iNode,
		name,
		quantity,
		unit,
	)
	return material, nil
}

func (r MaterialRepositorySql) FetchMaterialById(
	iContext context.Context,
	iId models.NodeId,
) (models.Material, error) {
	nodes, err := r.nodeRepository.FetchNodesById([]models.NodeId{models.NodeId(iId)})
	if err != nil {
		return models.Material{}, err
	}
	if len(nodes) == 0 {
		return models.Material{}, common.NotFound
	}
	node := nodes[0]
	material, err := r.fetchMaterialByNode(
		iContext,
		node,
	)
	if err != nil {
		return models.Material{}, err
	}
	return material, nil
}

func (r MaterialRepositorySql) FetchMaterialsByNodeId(
	iContext context.Context,
	iNamespace string,
	iMaterialNodeId map[string]bool,
) (map[models.NodeId]models.Material, error) {
	if len(iMaterialNodeId) == 0 {
		return map[models.NodeId]models.Material{}, nil
	}

	nodes, err := r.nodeRepository.FetchNodesByNodeId(
		iContext,
		iNamespace,
		iMaterialNodeId,
	)

	if err != nil {
		return map[models.NodeId]models.Material{}, err
	}

	if len(nodes) == 0 {
		return map[models.NodeId]models.Material{}, nil
	}

	materials := map[models.NodeId]models.Material{}
	for nodeId := range nodes {
		material, err := r.fetchMaterialByNode(iContext, nodes[nodeId])
		if err != nil {
			return map[models.NodeId]models.Material{}, err
		}
		materials[*material.Id] = material
	}

	return materials, nil
}

func (r MaterialRepositorySql) FetchMaterialsById(
	iContext context.Context,
	iIds map[models.NodeId]bool,
) (map[models.NodeId]models.Material, error) {
	if len(iIds) == 0 {
		return map[models.NodeId]models.Material{}, nil
	}

	ids := []models.NodeId{}
	for id := range iIds {
		ids = append(ids, id)
	}
	nodes, err := r.nodeRepository.FetchNodesById(ids)
	if err != nil {
		return map[models.NodeId]models.Material{}, err
	}

	if len(nodes) == 0 {
		return map[models.NodeId]models.Material{}, nil
	}

	nodeMap := map[models.NodeId]models.Node{}

	argString := []string{}
	arg := []interface{}{}
	count := 1
	for index := range nodes {
		nodeMap[*nodes[index].Id] = nodes[index]
		argString = append(argString, fmt.Sprintf("(node_id=$%d)", count))
		arg = append(arg, *nodes[index].Id)
		count++
	}

	query := `
		SELECT 
			node_id,
			name,
			quantity,
			unit
		FROM "material"
		WHERE 
	`

	query += strings.Join(argString, " OR ")
	fmt.Println("query")
	fmt.Println(query)
	response, err := r.db.QueryContext(
		iContext,
		query,
		arg...,
	)

	if err != nil {
		return map[models.NodeId]models.Material{}, nil
	}

	materials := map[models.NodeId]models.Material{}
	for response.Next() {
		nodeId := models.NodeId(0)
		name := ""
		quantity := models.CustomDecimal{}
		unit := ""
		err := response.Scan(
			&nodeId,
			&name,
			&quantity,
			&unit,
		)

		if err != nil {
			return map[models.NodeId]models.Material{}, nil
		}

		material := models.NewMaterial(
			nodeMap[nodeId],
			name,
			quantity,
			unit,
		)

		materials[*material.Id] = material
	}

	return materials, nil
}

func (r MaterialRepositorySql) upsertMaterialsByIds(
	iContext context.Context,
	iMaterials []models.Material,
) error {
	argString := []string{}
	arg := []interface{}{}
	count := 1
	query := `
		INSERT INTO "material" (
			node_id,
			name,
			quantity,
			unit
		) VALUES 
	`
	for i := range iMaterials {
		argString = append(argString, fmt.Sprintf("($%d, $%d, $%d, $%d)", count, count+1, count+2, count+3))
		arg = append(arg, []interface{}{
			*iMaterials[i].Id,
			iMaterials[i].Name,
			iMaterials[i].Quantity,
			iMaterials[i].Unit,
		}...)
		count += 4
	}

	query += strings.Join(argString, ",")
	query += ` 
		ON CONFLICT (node_id)
		DO UPDATE 
			SET 
				name=EXCLUDED.name,
				quantity=EXCLUDED.quantity,
				unit=EXCLUDED.unit
	`
	result, err := r.db.QueryContext(
		iContext,
		query,
		arg...,
	)

	if err != nil {
		return err
	}
	defer result.Close()
	return nil
}

func (r MaterialRepositorySql) UpsertMaterialsByIds(
	iContext context.Context,
	iMaterials map[models.NodeId]models.Material,
) (map[models.NodeId]models.Material, error) {
	if len(iMaterials) == 0 {
		return map[models.NodeId]models.Material{}, nil
	}

	nodes := map[models.NodeId]models.Node{}
	materialArray := make([]models.Material, 0, len(iMaterials))
	for id := range iMaterials {
		nodes[id] = iMaterials[id].Node
		materialArray = append(materialArray, iMaterials[id])
	}

	nodes, err := r.nodeRepository.UpsertNodesById(
		iContext,
		nodes,
	)

	if err != nil {
		return map[models.NodeId]models.Material{}, err
	}

	err = r.upsertMaterialsByIds(
		iContext,
		materialArray,
	)

	return iMaterials, nil
}

func (r MaterialRepositorySql) UpsertMaterialsByNodeIdsAndNamespace(
	iContext context.Context,
	iNamespace string,
	iMaterials map[string]models.Material,
) (map[string]models.Material, error) {
	if len(iMaterials) == 0 {
		return map[string]models.Material{}, nil
	}

	nodes := map[string]models.Node{}
	materialArray := make([]models.Material, 0, len(iMaterials))
	for id := range iMaterials {
		material := iMaterials[id]
		material.Node.Namespace = &iNamespace
		nodes[id] = material.Node
		materialArray = append(materialArray, material)
	}
	nodes, err := r.nodeRepository.UpsertNodesByNodeIdAndNamespace(
		iContext,
		iNamespace,
		nodes,
	)

	if err != nil {
		return map[string]models.Material{}, err
	}

	ret := map[string]models.Material{}
	for i := range materialArray {
		materialArray[i].Node = nodes[materialArray[i].NodeId]
		ret[materialArray[i].NodeId] = materialArray[i]
	}

	err = r.upsertMaterialsByIds(iContext, materialArray)

	if err != nil {
		return map[string]models.Material{}, err
	}

	return ret, nil
}
