package repositories

import (
	"backend/internal/models"
	"context"
	"database/sql"
)

type UserEndpointRepositorySql struct {
	db                     *sql.DB
	peerProtocolRepository PeerProtocolRepositoryI
}

func MakeUserEndpointRepositorySql(
	iDb *sql.DB,
	iPeerProtocolRepository PeerProtocolRepositoryI,
) UserEndpointRepositorySql {
	return UserEndpointRepositorySql{
		db:                     iDb,
		peerProtocolRepository: iPeerProtocolRepository,
	}
}

func (r UserEndpointRepositorySql) FetchUserEndpointByUserId(
	iContext context.Context,
	iUserId models.UserId,
) ([]models.UserEndpoint, error) {
	response, err := r.db.QueryContext(
		iContext,
		`
			SELECT
				id,
				user_id,
				url,
				protocol_id
			FROM "user_endpoint"
			WHERE user_id=$1
		`,
		iUserId,
	)

	if err != nil {
		return []models.UserEndpoint{}, nil
	}
	defer response.Close()

	type UserEndpoint struct {
		Id         models.UserEndpointId
		UserId     models.UserId
		Url        string
		ProtocolId models.PeerProtocolId
	}

	tempUserEndpoints := map[models.UserEndpointId]UserEndpoint{}
	for response.Next() {
		userEndpoint := UserEndpoint{}
		err := response.Scan(
			&userEndpoint.Id,
			&userEndpoint.UserId,
			&userEndpoint.Url,
			&userEndpoint.ProtocolId,
		)

		if err != nil {
			return []models.UserEndpoint{}, nil
		}

		tempUserEndpoints[userEndpoint.Id] = userEndpoint
	}

	protocolIds := map[models.PeerProtocolId]bool{}
	for id := range tempUserEndpoints {
		protocolIds[tempUserEndpoints[id].ProtocolId] = true
	}

	protocols, err := r.peerProtocolRepository.FetchPeerProtocolByIds(
		iContext,
		protocolIds,
	)

	if err != nil {
		return []models.UserEndpoint{}, nil
	}

	ret := []models.UserEndpoint{}
	for id := range tempUserEndpoints {
		if protocol, ok := protocols[tempUserEndpoints[id].ProtocolId]; ok {
			/// should always be ok
			userEndpoint := models.MakeUserEndpoint(
				id,
				tempUserEndpoints[id].UserId,
				tempUserEndpoints[id].Url,
				protocol,
			)
			ret = append(ret, userEndpoint)
		}
	}

	return ret, nil
}
