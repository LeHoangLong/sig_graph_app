package repositories

import (
	"backend/internal/models"
	graph_id_service "backend/internal/services/graph_id"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"
)

type NodeRepositorySql struct {
	hasher graph_id_service.IdHasher
	tx     *sql.Tx
}

func MakeNodeRepositorySql(
	iTx *sql.Tx,
) NodeRepositorySql {
	return NodeRepositorySql{
		tx: iTx,
	}
}

func (r NodeRepositorySql) CreateNode(iNode models.Node) (models.Node, error) {
	args := []interface{}{
		iNode.NodeId,
		*iNode.OwnerPublicKey.Id,
		iNode.IsFinalized,
		time.Time(iNode.CreatedTime),
		iNode.Signature,
		iNode.Type,
	}
	previousNodeArrayStrings := []string{}
	nextNodeArrayStrings := []string{}

	count := 6

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
			previous_node_hashed_ids,
			next_node_hashed_ids
		) VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
	`

	statement += "ARRAY[" + strings.Join(previousNodeArrayStrings, ",") + "]::TEXT[],"
	statement += "ARRAY[" + strings.Join(nextNodeArrayStrings, ",") + "]::TEXT[]"
	statement += ") RETURNING id"

	result, err := r.tx.Query(statement, args...)

	if err != nil {
		return models.Node{}, err
	}

	if !result.Next() {
		return models.Node{}, fmt.Errorf("no result from creating node")
	}

	var id int
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
		result, err = r.tx.Query(statement, args...)
		if err != nil {
			return models.Node{}, err
		}
	}

	ret := iNode
	ret.Id = &id
	return ret, nil
}

func (r NodeRepositorySql) FetchChildrenIdOfNode(iId int) (map[int]bool, error) {
	result, err := r.tx.Query(`
		SELECT dst_node_id FROM "node_edge" WHERE src_node_id=$1
	`, iId)

	if err != nil {
		return map[int]bool{}, nil
	}

	ret := map[int]bool{}
	for result.Next() {
		var childNodeId int
		result.Scan(&childNodeId)
		ret[childNodeId] = true
	}

	return ret, nil
}

func (r NodeRepositorySql) FetchParentIdOfNode(iId int) (map[int]bool, error) {
	result, err := r.tx.Query(`
		SELECT src_node_id FROM "node_edge" WHERE dst_node_id=$1
	`, iId)

	if err != nil {
		return map[int]bool{}, nil
	}

	ret := map[int]bool{}
	for result.Next() {
		var childNodeId int
		result.Scan(&childNodeId)
		ret[childNodeId] = true
	}

	return ret, nil
}

func (r NodeRepositorySql) FetchNodesById(iId []int) ([]models.Node, error) {
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
	result, err := r.tx.Query(statement, args...)

	if err != nil {
		return []models.Node{}, err
	}

	ret := []models.Node{}
	index := 0
	for result.Next() {
		node := models.Node{}
		var publicKeyId int
		var publicKeyValue string

		result.Scan(
			node.Id,
			&node.NodeId,
			&node.IsFinalized,
			pq.Array(&node.PreviousNodeHashedIds),
			pq.Array(&node.NextNodeHashedIds),
			&node.CreatedTime,
			&node.Signature,
			&node.Type,
			&publicKeyId,
			&publicKeyValue,
		)
		node.OwnerPublicKey = models.MakePublicKey(&publicKeyId, publicKeyValue)
		if index > len(iId) || *node.Id != iId[index] {

		}
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

func (r NodeRepositorySql) FetchNodesByOwnerKey(iOwnerKey models.PublicKey, iMinId int, iLimit int) ([]models.Node, error) {
	result, err := r.tx.Query(`
		SELECT 
			id,
			node_id,
			is_finalized,
			previous_node_hashed_ids,
			next_node_hashed_ids,
			created_time,
			signature,
			type,
		FROM "node" 
		WHERE public_key_id = $1 AND id > $2
		LIMIT $3
	`, iOwnerKey.Id, iMinId, iLimit)
	if err != nil {
		return []models.Node{}, err
	}

	ret := []models.Node{}
	for result.Next() {
		node := models.Node{}
		result.Scan(
			&node.Id,
			&node.NodeId,
			&node.IsFinalized,
			pq.Array(&node.PreviousNodeHashedIds),
			pq.Array(&node.NextNodeHashedIds),
			&node.CreatedTime,
			&node.Signature,
			&node.Type,
		)

		node.OwnerPublicKey = iOwnerKey
		ret = append(ret, node)
	}
	return ret, nil
}
