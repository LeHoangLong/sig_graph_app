package repositories

import (
	"backend/internal/common"
	"backend/internal/models"
	"backend/internal/services/node_contract"
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"
)

type NodeRepositorySql struct {
	hasher node_contract.IdHasherI
	db     *sql.DB
}

func MakeNodeRepositorySql(
	iDb *sql.DB,
	iHasher node_contract.IdHasherI,
) NodeRepositorySql {
	return NodeRepositorySql{
		db:     iDb,
		hasher: iHasher,
	}
}

func (r NodeRepositorySql) CreateNode(iNode models.Node) (models.Node, error) {
	if iNode.Namespace == nil {
		return models.Node{}, common.InvalidArgument
	}
	args := []interface{}{
		iNode.NodeId,
		*iNode.OwnerPublicKey.Id,
		iNode.IsFinalized,
		time.Time(iNode.CreatedTime),
		iNode.Signature,
		iNode.Type,
		*iNode.Namespace,
	}
	previousNodeArrayStrings := []string{}
	nextNodeArrayStrings := []string{}

	count := 8

	for _, previous := range iNode.PreviousNodeHashedIds {
		args = append(args, previous)
		previousNodeArrayStrings = append(previousNodeArrayStrings, fmt.Sprintf("$%d", count))
		count++
	}

	for _, next := range iNode.NextNodeHashedIds {
		args = append(args, next)
		nextNodeArrayStrings = append(nextNodeArrayStrings, fmt.Sprintf("$%d", count))
		count++
	}

	statement := `
		INSERT INTO "node" (
			node_id,
			public_key_id,
			is_finalized,
			created_time,
			signature,
			type,
			namespace,
			previous_node_hashed_ids,
			next_node_hashed_ids
		) VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
	`

	statement += "ARRAY[" + strings.Join(previousNodeArrayStrings, ",") + "]::TEXT[],"
	statement += "ARRAY[" + strings.Join(nextNodeArrayStrings, ",") + "]::TEXT[]"
	statement += ") RETURNING id"

	result, err := r.db.Query(statement, args...)

	if err != nil {
		return models.Node{}, err
	}

	if !result.Next() {
		if err, ok := err.(*pq.Error); ok {
			if err.Code == "23505" {
				return models.Node{}, common.AlreadyExistsErr
			}
		}
		return models.Node{}, fmt.Errorf("no result from creating node")
	}

	var id models.NodeId
	result.Scan(&id)
	result.Close()

	if len(iNode.ChildrenIds) > 0 || len(iNode.ParentIds) > 0 {
		statement = `
			INSERT INTO "node_edge" (
				src_node_id,
				dst_node_id
			) VALUES (
		`
		argsStr := []string{}
		args = []interface{}{}
		for id, _ := range iNode.ParentIds {
			argsStr = append(argsStr, fmt.Sprintf("$%d, $%d", count, count+1))
			args = append(args, id, iNode.NodeId)
			count += 2
		}

		for id, _ := range iNode.ChildrenIds {
			argsStr = append(argsStr, fmt.Sprintf("$%d, $%d", count, count+1))
			args = append(args, iNode.NodeId, id)
			count += 2
		}
		statement += strings.Join(argsStr, "), (")
		statement += ") ON CONFLICT (src_node_id, dst_node_id) DO NOTHING"
		result, err = r.db.Query(statement, args...)
		if err != nil {
			return models.Node{}, err
		}
	}

	ret := iNode
	ret.Id = &id
	return ret, nil
}

func (r NodeRepositorySql) FetchChildrenIdOfNode(iId models.NodeId) (map[models.NodeId]bool, error) {
	result, err := r.db.Query(`
		SELECT dst_node_id FROM "node_edge" WHERE src_node_id=$1
	`, iId)

	if err != nil {
		return map[models.NodeId]bool{}, nil
	}

	ret := map[models.NodeId]bool{}
	for result.Next() {
		var childNodeId models.NodeId
		result.Scan(&childNodeId)
		ret[childNodeId] = true
	}

	return ret, nil
}

func (r NodeRepositorySql) FetchParentIdOfNode(iId models.NodeId) (map[models.NodeId]bool, error) {
	result, err := r.db.Query(`
		SELECT src_node_id FROM "node_edge" WHERE dst_node_id=$1
	`, iId)

	if err != nil {
		return map[models.NodeId]bool{}, nil
	}

	ret := map[models.NodeId]bool{}
	for result.Next() {
		var childNodeId models.NodeId
		result.Scan(&childNodeId)
		ret[childNodeId] = true
	}

	return ret, nil
}

func (r NodeRepositorySql) FetchNodesById(iId []models.NodeId) ([]models.Node, error) {
	if len(iId) == 0 {
		return []models.Node{}, nil
	}
	statement := `
		SELECT 
			n.id,
			n.node_id,
			n.is_finalized,
			n.previous_node_hashed_ids,
			n.next_node_hashed_ids,
			n.created_time,
			n.signature,
			n.type,
			n.namespace,
			pk.id,
			pk.value
		FROM "node" n
		INNER JOIN "public_key" pk
		ON n.public_key_id = pk.id AND (
	`
	argsString := []string{}
	args := []interface{}{}
	count := 1

	for _, id := range iId {
		argsString = append(argsString, fmt.Sprintf("(n.id = $%d)", count))
		args = append(args, id)
		count += 1
	}

	statement += strings.Join(argsString, " OR ")
	statement += ")"
	result, err := r.db.Query(statement, args...)

	if err != nil {
		return []models.Node{}, err
	}

	ret := []models.Node{}
	for result.Next() {
		node := models.Node{
			Id:                    new(models.NodeId),
			PreviousNodeHashedIds: map[string]bool{},
			NextNodeHashedIds:     map[string]bool{},
			ParentIds:             map[models.NodeId]bool{},
			ChildrenIds:           map[models.NodeId]bool{},
			Namespace:             new(string),
		}
		var publicKeyId models.PublicKeyId
		var publicKeyValue string
		var previousNodeHashedIds, nextNodeHashedIds []string
		err := result.Scan(
			node.Id,
			&node.NodeId,
			&node.IsFinalized,
			pq.Array(&previousNodeHashedIds),
			pq.Array(&nextNodeHashedIds),
			&node.CreatedTime,
			&node.Signature,
			&node.Type,
			node.Namespace,
			&publicKeyId,
			&publicKeyValue,
		)

		for _, hashedId := range previousNodeHashedIds {
			node.PreviousNodeHashedIds[hashedId] = true
		}

		for _, hashedId := range nextNodeHashedIds {
			node.NextNodeHashedIds[hashedId] = true
		}

		if err != nil {
			return []models.Node{}, err
		}

		node.OwnerPublicKey = models.MakePublicKey(&publicKeyId, publicKeyValue)
		ret = append(ret, node)
	}

	for idx, _ := range ret {
		childrenIds, err := r.FetchChildrenIdOfNode(*ret[idx].Id)
		if err != nil {
			return []models.Node{}, err
		}

		parentIds, err := r.FetchChildrenIdOfNode(*ret[idx].Id)
		if err != nil {
			return []models.Node{}, err
		}

		ret[idx].ParentIds = parentIds
		ret[idx].ChildrenIds = childrenIds
	}

	return ret, nil
}

func (r NodeRepositorySql) FetchNodesByOwnerKey(
	iContext context.Context,
	iNamespace string,
	iOwnerKey models.PublicKey,
	iMinId int,
	iLimit int,
) ([]models.Node, error) {
	result, err := r.db.Query(`
		SELECT 
			id,
			node_id,
			is_finalized,
			previous_node_hashed_ids,
			next_node_hashed_ids,
			created_time,
			signature,
			type,
			namespace
		FROM "node" 
		WHERE public_key_id = $1 AND id >= $2 AND namespace=$4
		LIMIT $3
	`, *iOwnerKey.Id, iMinId, iLimit, iNamespace)
	if err != nil {
		return []models.Node{}, err
	}
	defer result.Close()

	ret := []models.Node{}
	for result.Next() {
		node := models.Node{
			Id:                    new(models.NodeId),
			Namespace:             new(string),
			PreviousNodeHashedIds: map[string]bool{},
			NextNodeHashedIds:     map[string]bool{},
		}

		var previousNodeHashedIds, nextNodeHashedIds []string
		err := result.Scan(
			node.Id,
			&node.NodeId,
			&node.IsFinalized,
			pq.Array(&previousNodeHashedIds),
			pq.Array(&nextNodeHashedIds),
			&node.CreatedTime,
			&node.Signature,
			&node.Type,
			node.Namespace,
		)

		for i := range previousNodeHashedIds {
			node.PreviousNodeHashedIds[previousNodeHashedIds[i]] = true
		}

		for i := range nextNodeHashedIds {
			node.NextNodeHashedIds[nextNodeHashedIds[i]] = true
		}

		if err != nil {
			return []models.Node{}, err
		}

		node.OwnerPublicKey = iOwnerKey
		ret = append(ret, node)
	}
	return ret, nil
}

func (r NodeRepositorySql) FetchNodesByNodeId(
	iContext context.Context,
	iNamespace string,
	iId map[string]bool,
) (map[models.NodeId]models.Node, error) {
	if len(iId) == 0 {
		return map[models.NodeId]models.Node{}, nil
	}
	count := 2
	arg := []interface{}{iNamespace}
	argString := []string{}
	for id := range iId {
		argString = append(argString, fmt.Sprintf("(node_id=$%d)", count))
		arg = append(arg, id)
		count += 1
	}
	query := `
		SELECT 
			n.id,
			n.node_id,
			n.is_finalized,
			n.previous_node_hashed_ids,
			n.next_node_hashed_ids,
			n.created_time,
			n.signature,
			n.type,
			n.namespace,
			pk.id,
			pk.value
		FROM "node" n
		INNER JOIN "public_key" pk
			ON pk.id = n.public_key_id
		WHERE namespace=$1 AND 
	`
	query += strings.Join(argString, " OR ")
	result, err := r.db.QueryContext(
		iContext,
		query,
		arg...,
	)

	if err != nil {
		return map[models.NodeId]models.Node{}, nil
	}
	defer result.Close()

	ret := map[models.NodeId]models.Node{}
	for result.Next() {
		node := models.Node{
			Id:        new(models.NodeId),
			Namespace: new(string),
		}
		publicKeyId := models.PublicKeyId(0)
		publicKeyValue := ""

		result.Scan(
			node.Id,
			&node.NodeId,
			&node.IsFinalized,
			pq.Array(&node.PreviousNodeHashedIds),
			pq.Array(&node.NextNodeHashedIds),
			&node.CreatedTime,
			&node.Signature,
			&node.Type,
			node.Namespace,
			&publicKeyId,
			&publicKeyValue,
		)

		node.OwnerPublicKey = models.MakePublicKey(&publicKeyId, publicKeyValue)
		ret[*node.Id] = node
	}
	return ret, nil
}

func (r NodeRepositorySql) UpsertNodesById(
	iContext context.Context,
	iNodes map[models.NodeId]models.Node,
) (map[models.NodeId]models.Node, error) {
	if len(iNodes) == 0 {
		return map[models.NodeId]models.Node{}, nil
	}

	query := `
		INSERT INTO "node" (
			id,
			node_id,
			namespace,
			public_key_id,
			is_finalized,
			previous_node_hashed_ids,
			next_node_hashed_ids,
			created_time,
			signature,
			type
		) VALUES 
	`
	args := []interface{}{}
	argString := []string{}
	count := 1
	for id := range iNodes {
		nextNodeHashedIds := []string{}
		for id := range iNodes[id].NextNodeHashedIds {
			nextNodeHashedIds = append(nextNodeHashedIds, id)
		}

		previousNodeHashedIds := []string{}
		for id := range iNodes[id].PreviousNodeHashedIds {
			previousNodeHashedIds = append(previousNodeHashedIds, id)
		}

		if iNodes[id].Namespace == nil {
			return map[models.NodeId]models.Node{}, common.InvalidArgument
		}

		if iNodes[id].OwnerPublicKey.Id == nil {
			return map[models.NodeId]models.Node{}, common.InvalidArgument
		}

		args = append(args, []interface{}{
			id,
			iNodes[id].NodeId,
			*iNodes[id].Namespace,
			*iNodes[id].OwnerPublicKey.Id,
			iNodes[id].IsFinalized,
			pq.Array(previousNodeHashedIds),
			pq.Array(nextNodeHashedIds),
			iNodes[id].CreatedTime,
			iNodes[id].Signature,
			iNodes[id].Type,
		}...)
		counts := []interface{}{}
		for i := 0; i < 10; i++ {
			counts = append(counts, count+i)
		}
		argString = append(argString, fmt.Sprintf(`(
			$%d,
			$%d,
			$%d,
			$%d,
			$%d,
			$%d,
			$%d,
			$%d,
			$%d,
			$%d
		)`, counts...))
		count += 10
	}

	query += strings.Join(argString, " , ")
	query += `
		ON CONFLICT (id)
		DO UPDATE SET
			node_id=EXCLUDED.node_id,
			namespace=EXCLUDED.namespace,
			public_key_id=EXCLUDED.public_key_id,
			is_finalized=EXCLUDED.is_finalized,
			previous_node_hashed_ids=EXCLUDED.previous_node_hashed_ids,
			next_node_hashed_ids=EXCLUDED.next_node_hashed_ids,
			created_time=EXCLUDED.created_time,
			signature=EXCLUDED.signature,
			type=EXCLUDED.type
	`
	result, err := r.db.QueryContext(
		iContext,
		query,
		args...,
	)
	if err != nil {
		return map[models.NodeId]models.Node{}, err
	}
	defer result.Close()
	return iNodes, nil
}

func (r NodeRepositorySql) UpsertNodesByNodeIdAndNamespace(
	iContext context.Context,
	iNamespace string,
	iNodes map[string]models.Node,
) (map[string]models.Node, error) {
	if len(iNodes) == 0 {
		return map[string]models.Node{}, nil
	}

	query := `
		INSERT INTO "node" as n (
			id,
			node_id,
			namespace,
			public_key_id,
			is_finalized,
			previous_node_hashed_ids,
			next_node_hashed_ids,
			created_time,
			signature,
			type
		) VALUES 
	`
	args := []interface{}{}
	argString := []string{}
	count := 1
	for nodeId := range iNodes {
		if iNodes[nodeId].OwnerPublicKey.Id == nil {
			return map[string]models.Node{}, common.InvalidArgument
		}

		nextNodeHashedIds := []string{}
		for id := range iNodes[nodeId].NextNodeHashedIds {
			nextNodeHashedIds = append(nextNodeHashedIds, id)
		}

		previousNodeHashedIds := []string{}
		for id := range iNodes[nodeId].PreviousNodeHashedIds {
			previousNodeHashedIds = append(previousNodeHashedIds, id)
		}

		id := sql.NullInt32{
			Valid: iNodes[nodeId].Id == nil,
			Int32: 0,
		}
		if iNodes[nodeId].Id != nil {
			id.Int32 = int32(*iNodes[nodeId].Id)
		}
		args = append(args, []interface{}{
			id,
			nodeId,
			iNamespace,
			*iNodes[nodeId].OwnerPublicKey.Id,
			iNodes[nodeId].IsFinalized,
			pq.Array(previousNodeHashedIds),
			pq.Array(nextNodeHashedIds),
			iNodes[nodeId].CreatedTime,
			iNodes[nodeId].Signature,
			iNodes[nodeId].Type,
		}...)
		counts := []interface{}{}
		for i := 0; i < 10; i++ {
			counts = append(counts, count+i)
		}
		argString = append(argString, fmt.Sprintf(`(
			$%d,
			$%d,
			$%d,
			$%d,
			$%d,
			$%d,
			$%d,
			$%d,
			$%d,
			$%d
		)`, counts...))
		count += 10
	}

	query += strings.Join(argString, " , ")
	query += `
		ON CONFLICT (node_id, namespace)
		DO UPDATE SET
			id=coalesce(EXCLUDED.id, n.id),
			public_key_id=EXCLUDED.public_key_id,
			is_finalized=EXCLUDED.is_finalized,
			previous_node_hashed_ids=EXCLUDED.previous_node_hashed_ids,
			next_node_hashed_ids=EXCLUDED.next_node_hashed_ids,
			created_time=EXCLUDED.created_time,
			signature=EXCLUDED.signature,
			type=EXCLUDED.type
		RETURNING id, node_id
	`
	result, err := r.db.QueryContext(
		iContext,
		query,
		args...,
	)
	if err != nil {
		return map[string]models.Node{}, err
	}
	defer result.Close()
	ret := map[string]models.Node{}
	for result.Next() {
		id := models.NodeId(0)
		nodeId := ""
		err := result.Scan(
			&id,
			&nodeId,
		)

		if err != nil {
			return map[string]models.Node{}, err
		}

		node := iNodes[nodeId]
		node.Id = &id
		ret[nodeId] = node
	}
	return ret, nil
}
