package repositories

import (
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

func (r PeerRepositorySql) FetchPeers(iUserId int) ([]models.Peer, error) {
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
			) peer ON peer_key.owner_id = peer.id
		) peer_key_peer ON public_key.id = peer_key_peer.public_key_id
	`, iUserId)

	if err != nil {
		return []models.Peer{}, err
	}

	peerMap := map[int]models.Peer{}
	for response.Next() {
		var peerId, publicKeyId int
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
	iPeerId int,
) ([]models.PeerEndpoint, error) {
	result, err := r.db.QueryContext(iContext, `
		SELECT 
			pe.id,
			pe.url,
			pr.protocol
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
		var endpointUrl, endpointProtocol string

		err := result.Scan(
			&endpointId,
			&endpointUrl,
			&endpointProtocol,
		)

		if err != nil {
			return []models.PeerEndpoint{}, err
		}

		endpoint := models.MakePeerEndpoint(
			endpointId,
			endpointUrl,
			models.PeerProtocol(endpointProtocol),
		)
		ret = append(ret, endpoint)
	}

	return ret, nil
}
