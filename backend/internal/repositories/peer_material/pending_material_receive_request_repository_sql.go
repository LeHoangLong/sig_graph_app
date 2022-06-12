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

func MakePendingMaterialReceiveRequestRepositorySql(
	iDb *sql.DB,
) PendingMaterialReceiveRequestRepositorySql {
	return PendingMaterialReceiveRequestRepositorySql{
		db: iDb,
	}
}

func (r PendingMaterialReceiveRequestRepositorySql) CreatePendingReceiveMaterialRequest(
	iContext context.Context,
	iUser models.User,
	iToBeReceivedMaterial models.Material,
	iRelatedMaterials []models.Material,
	iOptions []models.SignatureOption,
	iSenderPublicKey string,
	iTransferTime models.CustomTime,
	iIsOutbound bool,
) (models.PendingMaterialReceiveRequest, error) {
	tx, err := r.db.BeginTx(iContext, nil)
	if err != nil {
		return models.PendingMaterialReceiveRequest{}, err
	}

	requestResponse := tx.QueryRowContext(
		iContext,
		`INSERT INTO "pending_receive_material_request" (
			main_node_id,
			transfer_time,
			sender_public_key,
			owner_id,
			is_outbound
		) VALUES (
			$1,
			$2,
			$3,
			$4,
			$5
		) RETURNING id
	`, *iToBeReceivedMaterial.Id, iTransferTime, iSenderPublicKey, iUser.Id, iIsOutbound)

	var pendingRequestId int
	err = requestResponse.Scan(&pendingRequestId)
	if err != nil {
		tx.Rollback()
		return models.PendingMaterialReceiveRequest{}, fmt.Errorf("cannot create pending request %w", err)
	}

	if len(iRelatedMaterials) > 0 {
		queryValString := make([]string, 0, len(iRelatedMaterials))
		queryVal := make([]interface{}, 0, len(iRelatedMaterials)*2)
		count := 0
		for _, material := range iRelatedMaterials {
			queryValString = append(queryValString, fmt.Sprintf("($%d, %d)", count, count+1))
			queryVal = append(queryVal, *&material.Id)
			queryVal = append(queryVal, pendingRequestId)
			count += 2
		}
		query := fmt.Sprintf(`
			INSERT INTO "node_from_peer" (
				node_id,
				pending_receive_material_request_id
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
	}

	if len(iOptions) > 0 {
		count := 1
		queryValString := make([]string, 0, len(iOptions))
		queryVal := make([]interface{}, 0, len(iOptions)*3)
		for _, option := range iOptions {
			queryValString = append(queryValString, fmt.Sprintf("($%d, $%d, $%d)", count, count+1, count+2))
			queryVal = append(queryVal, option.Signature)
			queryVal = append(queryVal, option.Id)
			queryVal = append(queryVal, pendingRequestId)
			count += 3
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
	}

	tx.Commit()

	ret := models.MakePendingMaterialReceiveRequest(
		pendingRequestId,
		iUser.Id,
		iToBeReceivedMaterial,
		iRelatedMaterials,
		iOptions,
		iSenderPublicKey,
		iTransferTime,
	)

	return ret, err
}
