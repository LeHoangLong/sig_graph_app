package repositories

import (
	"backend/internal/models"
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type PeerProtocolRepositorySql struct {
	db *sql.DB
}

func MakePeerProtocolRepositorySql(
	iDb *sql.DB,
) PeerProtocolRepositorySql {
	return PeerProtocolRepositorySql{
		db: iDb,
	}
}

func (r PeerProtocolRepositorySql) FetchPeerProtocolByIds(
	iContext context.Context,
	iIds map[models.PeerProtocolId]bool,
) (map[models.PeerProtocolId]models.PeerProtocol, error) {
	if len(iIds) == 0 {
		return map[models.PeerProtocolId]models.PeerProtocol{}, nil
	}

	query := `
		SELECT 
			id,
			protocol,
			major_version,
			minor_version
		FROM "supported_peer_protocol"
		WHERE 
	`
	arg := []interface{}{}
	argString := []string{}
	count := 1
	for id := range iIds {
		arg = append(arg, id)
		argString = append(argString, fmt.Sprintf("(id=$%d)", count))
		count += 1
	}

	query += strings.Join(argString, " OR ")
	response, err := r.db.QueryContext(
		iContext,
		query,
		arg...,
	)
	if err != nil {
		return map[models.PeerProtocolId]models.PeerProtocol{}, err
	}
	defer response.Close()

	ret := map[models.PeerProtocolId]models.PeerProtocol{}
	for response.Next() {
		id := models.PeerProtocolId(0)
		protocol := ""
		majorVersion := 0
		minorVersion := 0

		err := response.Scan(
			&id,
			&protocol,
			&majorVersion,
			&minorVersion,
		)

		if err != nil {
			return map[models.PeerProtocolId]models.PeerProtocol{}, err
		}

		ret[id] = models.MakePeerProtocol(id, models.ProtocolName(protocol), majorVersion, minorVersion)
	}

	return ret, nil
}

func (r PeerProtocolRepositorySql) FilterSupportedProtocolsWithMajorVersion(
	iContext context.Context,
	iProtocols map[ProtocolWithMajorVersion]bool,
) (map[ProtocolWithMajorVersion][]models.PeerProtocol, error) {
	if len(iProtocols) == 0 {
		return map[ProtocolWithMajorVersion][]models.PeerProtocol{}, nil
	}

	argString := []string{}
	arg := []interface{}{}
	count := 1

	query := `
		SELECT 
			id,
			protocol,
			major_version,
			minor_version
		FROM "supported_peer_protocol"
		WHERE 
	`

	for protocol := range iProtocols {
		arg = append(arg, protocol.Protocol, protocol.MajorVersion)
		argString = append(argString, fmt.Sprintf("(protocol=$%d AND major_version=$%d)", count, count+1))
		count += 2
	}

	query += strings.Join(argString, " OR ")
	response, err := r.db.QueryContext(
		iContext,
		query,
		arg...,
	)

	if err != nil {
		return map[ProtocolWithMajorVersion][]models.PeerProtocol{}, nil
	}
	defer response.Close()

	ret := map[ProtocolWithMajorVersion][]models.PeerProtocol{}
	for response.Next() {
		id := models.PeerProtocolId(0)
		protocol := ""
		majorVersion := 0
		minorVersion := 0

		err := response.Scan(
			&id,
			&protocol,
			&majorVersion,
			&minorVersion,
		)

		if err != nil {
			return map[ProtocolWithMajorVersion][]models.PeerProtocol{}, err
		}

		protocolWithMajorVersion := MakeProtocolWithMajorVersion(protocol, majorVersion)
		ret[protocolWithMajorVersion] = append(ret[protocolWithMajorVersion], models.MakePeerProtocol(
			id,
			models.ProtocolName(protocol),
			majorVersion,
			minorVersion,
		))
	}

	return ret, nil

}
