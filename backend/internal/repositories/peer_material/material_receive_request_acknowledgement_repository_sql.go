package peer_material_repositories

import (
	"backend/internal/common"
	"backend/internal/models"
	"context"
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

func (r MaterialReceiveRequestAcknowledgementRepositorySql) SaveOutboundAcknowledgement(
	iRequestId models.MaterialReceiveRequestId,
	iResponseId models.MaterialReceiveRequestResponseId,
) (models.MaterialReceiveRequestAcknowledgement, error) {
	_, err := r.db.Query(`INSERT INTO "outbound_receive_material_request_acknowledgement" (
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

func (r MaterialReceiveRequestAcknowledgementRepositorySql) FetchOutboundAcknowledementByRequestId(
	iContext context.Context,
	iRequestId models.MaterialReceiveRequestId,
) (models.MaterialReceiveRequestAcknowledgement, error) {
	response := r.db.QueryRowContext(
		iContext,
		`
			SELECT 
				request_id,
				response_id
			FROM "outbound_receive_material_request_acknowledgement"
			WHERE response_id=$1
		`, iRequestId,
	)

	requestId := models.MaterialReceiveRequestId(0)
	responseId := models.MaterialReceiveRequestResponseId("")
	err := response.Scan(&requestId, &responseId)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.MaterialReceiveRequestAcknowledgement{}, common.NotFound
		}
		return models.MaterialReceiveRequestAcknowledgement{}, err
	}

	acknowledgement := models.MakeMaterialReceiveRequestAcknowledgement(
		requestId,
		responseId,
	)

	return acknowledgement, nil
}

func (r MaterialReceiveRequestAcknowledgementRepositorySql) FetchOutboundAcknowledementByResponseId(
	iContext context.Context,
	iResponseId models.MaterialReceiveRequestResponseId,
) (models.MaterialReceiveRequestAcknowledgement, error) {
	response := r.db.QueryRowContext(
		iContext,
		`
			SELECT 
				request_id,
				response_id
			FROM "outbound_receive_material_request_acknowledgement"
			WHERE response_id=$1
		`, iResponseId,
	)

	requestId := models.MaterialReceiveRequestId(0)
	responseId := models.MaterialReceiveRequestResponseId("")
	err := response.Scan(&requestId, &responseId)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.MaterialReceiveRequestAcknowledgement{}, common.NotFound
		}
		return models.MaterialReceiveRequestAcknowledgement{}, err
	}

	acknowledgement := models.MakeMaterialReceiveRequestAcknowledgement(
		requestId,
		responseId,
	)

	return acknowledgement, nil
}

func (r MaterialReceiveRequestAcknowledgementRepositorySql) SaveInboundAcknowledgement(
	iRequestId models.MaterialReceiveRequestId,
	iResponseId models.MaterialReceiveRequestResponseId,
) (models.MaterialReceiveRequestAcknowledgement, error) {
	_, err := r.db.Query(`INSERT INTO "inbound_receive_material_request_acknowledgement" (
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

func (r MaterialReceiveRequestAcknowledgementRepositorySql) FetchInboundAcknowledementByRequestId(
	iContext context.Context,
	iRequestId models.MaterialReceiveRequestId,
) (models.MaterialReceiveRequestAcknowledgement, error) {
	response := r.db.QueryRowContext(
		iContext,
		`
			SELECT 
				request_id,
				response_id
			FROM "inbound_receive_material_request_acknowledgement"
			WHERE response_id=$1
		`, iRequestId,
	)

	requestId := models.MaterialReceiveRequestId(0)
	responseId := models.MaterialReceiveRequestResponseId("")
	err := response.Scan(&requestId, &responseId)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.MaterialReceiveRequestAcknowledgement{}, common.NotFound
		}
		return models.MaterialReceiveRequestAcknowledgement{}, err
	}

	acknowledgement := models.MakeMaterialReceiveRequestAcknowledgement(
		requestId,
		responseId,
	)

	return acknowledgement, nil
}

func (r MaterialReceiveRequestAcknowledgementRepositorySql) FetchInboundAcknowledementByResponseId(
	iContext context.Context,
	iResponseId models.MaterialReceiveRequestResponseId,
) (models.MaterialReceiveRequestAcknowledgement, error) {
	response := r.db.QueryRowContext(
		iContext,
		`
			SELECT 
				request_id,
				response_id
			FROM "inbound_receive_material_request_acknowledgement"
			WHERE response_id=$1
		`, iResponseId,
	)

	requestId := models.MaterialReceiveRequestId(0)
	responseId := models.MaterialReceiveRequestResponseId("")
	err := response.Scan(&requestId, &responseId)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.MaterialReceiveRequestAcknowledgement{}, common.NotFound
		}
		return models.MaterialReceiveRequestAcknowledgement{}, err
	}

	acknowledgement := models.MakeMaterialReceiveRequestAcknowledgement(
		requestId,
		responseId,
	)

	return acknowledgement, nil
}
