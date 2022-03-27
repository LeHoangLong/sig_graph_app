package repositories

import (
	"backend/internal/models"
	"database/sql"
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
			user_key.id, 
			public_key.id,
			public_key.value, 
			user_key.private_key, 
			user_key.is_default
		FROM "public_key" public_key
		INNER JOIN (
			SELECT id, public_key_id, is_default, private_key
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
		var id, publicKeyId int
		var publicKey, privateKey string
		var isDefault bool

		err := response.Scan(
			&id,
			&publicKeyId,
			&publicKey,
			&privateKey,
			&isDefault,
		)

		if err != nil {
			return []models.UserKeyPair{}, err
		}

		ret = append(ret, models.MakeUserKeyPair(
			id,
			models.MakePublicKey(
				publicKeyId,
				publicKey,
			),
			privateKey,
			isDefault,
		))
	}

	return ret, nil
}
