package peer_material_repositories

import (
	"backend/internal/models"
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type PendingMaterialReceiveRequestRepositorySql struct {
	db *sql.DB
}

func (r PendingMaterialReceiveRequestRepositorySql) createPendingReceiveMaterialRequest(
	iContext context.Context,
	iNodeId string,
	iInfo models.PendingMaterialInfo,
	iOptions []models.SignatureOption,
) (models.PendingMaterialReceiveRequest, error) {
	tx, err := r.db.BeginTx(iContext, nil)
	if err != nil {
		return models.PendingMaterialReceiveRequest{}, err
	}

	requestResponse := tx.QueryRowContext(
		iContext,
		`INSERT INTO "pending_receive_material_request" (
			current_node_id
		) VALUES (
			$1
		) RETURNING id
	`, iNodeId)

	var pendingRequestId int
	err = requestResponse.Scan(&pendingRequestId)
	if err != nil {
		tx.Rollback()
		return models.PendingMaterialReceiveRequest{}, fmt.Errorf("cannot create pending request %w", err)
	}

	queryValString := make([]string, 0, len(iOptions))
	queryVal := make([]interface{}, 0, len(iOptions)*3)
	for _, option := range iOptions {
		queryValString = append(queryValString, "(?, ?, ?)")
		queryVal = append(queryVal, option.Signature)
		queryVal = append(queryVal, option.Id)
		queryVal = append(queryVal, pendingRequestId)
	}

	query := fmt.Sprintf(`
		INSERT INTO "signature_option" (
			signature,
			new_node_id,
			pending_request_id
		) VALUES %s
	`, strings.Join(queryValString, ","))

	_, err = tx.ExecContext(
		iContext,
		query,
		queryVal...,
	)
	if err != nil {
		tx.Rollback()
		return models.PendingMaterialReceiveRequest{}, fmt.Errorf("cannot create signature options %w", err)
	}

	_, err = tx.ExecContext(
		iContext,
		`
			INSERT INTO "pending_material_info" (
				name,
				quantity,
				unit,
				pending_request_id
			) VALUES (
				$1,
				$2,
				$3,
				$4
			)
		`,
		iInfo.Name,
		iInfo.Quantity.String(),
		iInfo.Unit,
		pendingRequestId,
	)
	if err != nil {
		tx.Rollback()
		return models.PendingMaterialReceiveRequest{}, fmt.Errorf("cannot create pending material info %w", err)
	}

	tx.Commit()

	return models.PendingMaterialReceiveRequest{}, err
}
