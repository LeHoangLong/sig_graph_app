package repositories

import (
	"backend/internal/common"
	"backend/internal/models"
	"context"
	"database/sql"
)

type PeerRepositorySql struct {
	db *sql.DB
}

func MakePeerRepositorySql(iDb *sql.DB) PeerRepositorySql {
	return PeerRepositorySql{
		db: iDb,
	}
}

func (r PeerRepositorySql) FetchPeers(iUserId models.UserId) ([]models.Peer, error) {
	response, err := r.db.Query(`
		SELECT 
			peer_key_peer.peer_id, 
			peer_key_peer.alias, 
			public_key.id, 
			public_key.value 
		FROM "public_key" public_key
		INNER JOIN (
			SELECT peer.id as peer_id, peer.alias, peer_key.public_key_id FROM "peer_key" peer_key
			INNER JOIN (
				SELECT id, alias, owner_id 
				FROM "peer" 
				WHERE owner_id = $1
			) peer ON peer_key.peer_id = peer.id
		) peer_key_peer ON public_key.id = peer_key_peer.public_key_id
	`, iUserId)

	if err != nil {
		return []models.Peer{}, err
	}

	peerMap := map[models.PeerId]models.Peer{}
	for response.Next() {
		var peerId models.PeerId
		var publicKeyId models.PublicKeyId
		var peerAlias, publicKeyValue string

		response.Scan(
			&peerId,
			&peerAlias,
			&publicKeyId,
			&publicKeyValue,
		)

		publicKey := models.MakePublicKey(&publicKeyId, publicKeyValue)

		if peer, ok := peerMap[peerId]; ok {
			peer.PublicKey = append(peer.PublicKey, publicKey)
		} else {
			peerMap[peerId] = models.MakePeer(peerId, peerAlias, []models.PublicKey{publicKey})
		}
	}

	ret := []models.Peer{}
	for _, v := range peerMap {
		ret = append(ret, v)
	}

	return ret, nil
}

func (r PeerRepositorySql) FetchPeerEndPoints(
	iContext context.Context,
	iPeerId models.PeerId,
) ([]models.PeerEndpoint, error) {
	result, err := r.db.QueryContext(iContext, `
		SELECT 
			pe.id,
			pe.url,
			pr.id,
			pr.protocol,
			pr.major_version,
			pr.minor_version
		FROM "peer_endpoint" pe
		INNER JOIN "supported_peer_protocol" pr
		ON pe.peer_id = $1 AND pe.protocol_id = pr.id
	`, iPeerId)

	if err != nil {
		return []models.PeerEndpoint{}, err
	}

	ret := []models.PeerEndpoint{}
	for result.Next() {
		var endpointId int
		var endpointUrl string
		protocolId := models.PeerProtocolId(0)
		protocolName := models.ProtocolName("")
		protocolMajor, protocolMinor := 0, 0

		err := result.Scan(
			&endpointId,
			&endpointUrl,
			&protocolId,
			&protocolName,
			&protocolMajor,
			&protocolMinor,
		)

		if err != nil {
			return []models.PeerEndpoint{}, err
		}

		endpoint := models.MakePeerEndpoint(
			endpointId,
			endpointUrl,
			models.MakePeerProtocol(
				protocolId,
				protocolName,
				protocolMajor,
				protocolMinor,
			),
		)
		ret = append(ret, endpoint)
	}

	return ret, nil
}

func (r PeerRepositorySql) FetchPeerByKeyId(
	iContext context.Context,
	iPublicKeyId models.PublicKeyId,
) (models.Peer, error) {
	response, err := r.db.QueryContext(
		iContext,
		`
			SELECT 
				peer.id,
				peer.alias,
				public_key.id,
				public_key.value
			FROM  (
				SELECT 
					peer.id,
					peer.alias
				FROM "peer" peer
				INNER JOIN "peer_key" peer_key
				ON peer.id = peer_key.peer_id
				AND peer_key.public_key_id=$1
			) peer
			INNER JOIN "peer_key" peer_key
				ON peer.id = peer_key.peer_id
			INNER JOIN "public_key" public_key
				ON public_key.id = peer_key.public_key_id
		`,
		iPublicKeyId,
	)
	if err != nil {
		return models.Peer{}, err
	}
	defer response.Close()

	peerId := models.PeerId(0)
	peerAlias := ""
	keys := []models.PublicKey{}

	for response.Next() {
		keyId := models.PublicKeyId(0)
		keyValue := ""
		err := response.Scan(
			&peerId,
			&peerAlias,
			&keyId,
			&keyValue,
		)

		if err != nil {
			return models.Peer{}, err
		}

		key := models.MakePublicKey(&keyId, keyValue)
		keys = append(keys, key)
	}
	if len(keys) == 0 {
		/// no peer with the specified key id
		return models.Peer{}, common.NotFound
	}

	peer := models.MakePeer(peerId, peerAlias, keys)
	return peer, nil
}
