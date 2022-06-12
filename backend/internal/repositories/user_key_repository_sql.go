package repositories

import (
	"backend/internal/models"
	"context"
	"database/sql"
	"fmt"
)

type UserKeyRepositorySql struct {
	db *sql.DB
}

func MakeUserKeyRepositorySql(iDb *sql.DB) UserKeyRepositorySql {
	return UserKeyRepositorySql{
		db: iDb,
	}
}

func (r UserKeyRepositorySql) FetchUserKeyPairByUser(iUserId int) ([]models.UserKeyPair, error) {
	response, err := r.db.Query(`
		SELECT 
			public_key.id,
			public_key.value, 
			user_key.private_key, 
			user_key.is_default
		FROM "public_key" public_key
		INNER JOIN (
			SELECT public_key_id, is_default, private_key
			FROM "user_key"
			WHERE owner_id = $1
		) user_key
		ON public_key.id = user_key.public_key_id
	`, iUserId)

	if err != nil {
		return []models.UserKeyPair{}, err
	}

	ret := []models.UserKeyPair{}
	for response.Next() {
		var id int
		var publicKey, privateKey string
		var isDefault bool

		err := response.Scan(
			&id,
			&publicKey,
			&privateKey,
			&isDefault,
		)

		if err != nil {
			return []models.UserKeyPair{}, err
		}

		ret = append(ret, models.MakeUserKeyPair(
			models.MakePublicKey(
				&id,
				publicKey,
			),
			privateKey,
			isDefault,
		))
	}

	return ret, nil
}

func (r UserKeyRepositorySql) FetchDefaultUserKeyPair(iUserId int) (models.UserKeyPair, error) {
	response, err := r.db.Query(`
		SELECT 
			public_key.id,
			public_key.value, 
			user_key.private_key, 
			user_key.is_default
		FROM "public_key" public_key
		INNER JOIN (
			SELECT public_key_id, is_default, private_key
			FROM "user_key"
			WHERE owner_id = $1 AND is_default = TRUE
		) user_key
		ON public_key.id = user_key.public_key_id
	`, iUserId)

	if err != nil {
		return models.UserKeyPair{}, err
	}

	if !response.Next() {
		return models.UserKeyPair{}, fmt.Errorf("no default key found")
	}

	var id int
	var publicKey, privateKey string
	var isDefault bool

	err = response.Scan(
		&id,
		&publicKey,
		&privateKey,
		&isDefault,
	)

	if err != nil {
		return models.UserKeyPair{}, err
	}

	return models.MakeUserKeyPair(
		models.MakePublicKey(
			&id,
			publicKey,
		),
		privateKey,
		isDefault,
	), nil
}

func (r UserKeyRepositorySql) FetchPublicKeyByPeerId(
	iContext context.Context,
	iPeerId int,
) ([]models.PublicKey, error) {
	result, err := r.db.QueryContext(
		iContext,
		`
			SELECT
				public_key_id,
				value
			FROM "peer_key" peer_key
			INNER JOIN "public_key" public_key
			ON peer_key.owner_id=$1 AND peer_key.public_key_id = public_key.id
		`,
		iPeerId,
	)

	if err != nil {
		return []models.PublicKey{}, err
	}

	ret := []models.PublicKey{}
	for result.Next() {
		id := 0
		value := ""
		err := result.Scan(&id, &value)
		if err != nil {
			return []models.PublicKey{}, err
		}

		ret = append(ret, models.MakePublicKey(&id, value))
	}

	return ret, nil
}
