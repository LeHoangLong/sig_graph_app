package peer_material_repositories

import (
	"backend/internal/models"
	"database/sql"
)

type MaterialReceiveRequestAcknowledgementRepositorySql struct {
	db *sql.DB
}

func MakeMaterialReceiveRequestAcknowledgementRepositorySql(
	iDb *sql.DB,
) MaterialReceiveRequestAcknowledgementRepositorySql {
	return MaterialReceiveRequestAcknowledgementRepositorySql{
		db: iDb,
	}
}

func (r MaterialReceiveRequestAcknowledgementRepositorySql) SaveAcknowledgement(
	iRequestId models.MaterialReceiveRequestId,
	iResponseId models.MaterialReceiveRequestResponseId,
) (models.MaterialReceiveRequestAcknowledgement, error) {
	_, err := r.db.Query(`INSERT INTO "receive_material_request_acknowledgement" (
		request_id,
		response_id
	) VALUES (
		$1,
		$2
	)`, iRequestId, iResponseId)
	if err != nil {
		return models.MaterialReceiveRequestAcknowledgement{}, err
	}

	return models.MakeMaterialReceiveRequestAcknowledgement(
		iRequestId,
		iResponseId,
	), nil
}
