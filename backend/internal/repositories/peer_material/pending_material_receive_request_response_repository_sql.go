package peer_material_repositories

import (
	"backend/internal/models"
	"database/sql"
)

type PendingMaterialReceiveRequestResponseRepositorySql struct {
	db *sql.DB
}

func MakePendingMaterialReceiveRequestResponseRepositorySql(
	iDb *sql.DB,
) PendingMaterialReceiveRequestResponseRepositorySql {
	return PendingMaterialReceiveRequestResponseRepositorySql{
		db: iDb,
	}
}

func (r PendingMaterialReceiveRequestResponseRepositorySql) SaveResponse(
	iRequestId int,
	iResponseId string,
) (models.PendingMaterialReceiveRequestResponse, error) {
	_, err := r.db.Query(`INSERT INTO "pending_receive_material_request_response" (
		request_id,
		response_id
	) VALUES (
		$1,
		$2
	)`, iRequestId, iResponseId)
	if err != nil {
		return models.PendingMaterialReceiveRequestResponse{}, err
	}

	return models.MakePendingMaterialReceiveRequestResponse(
		iRequestId,
		iResponseId,
	), nil
}
