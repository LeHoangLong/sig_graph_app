package repositories

import (
	"backend/internal/models"
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type PeerKeyRepositorySql struct {
	db *sql.DB
}

func MakePeerKeyRepositorySql(
	iDb *sql.DB,
) PeerKeyRepositorySql {
	return PeerKeyRepositorySql{
		db: iDb,
	}
}

func (r PeerKeyRepositorySql) CreateOrFetchPeerKeysByValue(
	iContext context.Context,
	iOwner models.User,
	iPeerKeys []string,
) ([]models.PeerKey, error) {
	if len(iPeerKeys) == 0 {
		return []models.PeerKey{}, nil
	}

	argStr := []string{}
	arg := []interface{}{iOwner.Id}
	count := 2
	for _, key := range iPeerKeys {
		argStr = append(argStr, fmt.Sprintf("(public_key.value=$%d)", count))
		arg = append(arg, key)
		count++
	}

	query := `
		SELECT 
			public_key.id, 
			public_key.value
		FROM "peer_key" peer_key
		INNER JOIN "public_key" public_key
		ON peer_key.public_key_id = public_key.id
			AND peer_key.owner_id=$1
			AND 
	`

	query = query + strings.Join(argStr, " OR ")

	result, err := r.db.Query(query, arg...)
	if err != nil {
		return []models.PeerKey{}, err
	}
	defer result.Close()

	existingKeysMap := map[string]models.PeerKey{}
	for result.Next() {
		keyId := 0
		keyVal := ""
		err := result.Scan(&keyId, &keyVal)
		if err != nil {
			return []models.PeerKey{}, err
		}

		existingKeysMap[keyVal] = models.MakePeerKey(models.MakePublicKey(&keyId, keyVal))
	}

	keysToInsert := []string{}
	for _, key := range iPeerKeys {
		if _, ok := existingKeysMap[key]; !ok {
			keysToInsert = append(keysToInsert, key)
		}
	}

	argStr = []string{}
	arg = []interface{}{}
	count = 1
	newKeysMap := map[string]models.PeerKey{}
	if len(keysToInsert) > 0 {
		for _, key := range keysToInsert {
			argStr = append(argStr, fmt.Sprintf("($%d)", count))
			arg = append(arg, key)
			count++
		}

		query = `
			INSERT INTO "public_key" (value) VALUES 
		`

		query += strings.Join(argStr, ", ")
		query += " RETURNING id"
		result, err := r.db.Query(query, arg...)
		if err != nil {
			return []models.PeerKey{}, err
		}

		newPublicKeyIds := []int{}
		for result.Next() {
			publicKeyId := 0
			err := result.Scan(&publicKeyId)
			if err != nil {
				return []models.PeerKey{}, err
			}
			newPublicKeyIds = append(newPublicKeyIds, publicKeyId)
		}

		argStr = []string{}
		arg = []interface{}{iOwner.Id}
		count = 2
		for _, id := range newPublicKeyIds {
			argStr = append(argStr, fmt.Sprintf("($%d, $1)", count))
			arg = append(arg, id)
			count++
		}
		query = `
			INSERT INTO "peer_key" (public_key_id, owner_id) VALUES 
		`
		query += strings.Join(argStr, ", ")
		_, err = r.db.Query(query, arg...)
		if err != nil {
			return []models.PeerKey{}, err
		}

		for index := range keysToInsert {
			newKey := models.MakePeerKey(models.MakePublicKey(&newPublicKeyIds[index], keysToInsert[index]))
			newKeysMap[newKey.Value] = newKey
		}
	}

	ret := []models.PeerKey{}
	for _, key := range iPeerKeys {
		if existingKey, ok := existingKeysMap[key]; ok {
			ret = append(ret, existingKey)
		} else if newKey, ok := newKeysMap[key]; ok {
			ret = append(ret, newKey)
		} else {
			return []models.PeerKey{}, fmt.Errorf("Something wrong, no key matched")
		}
	}

	return ret, nil
}
