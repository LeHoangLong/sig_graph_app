package repositories

import (
	"backend/internal/common"
	"backend/internal/models"
	"context"
	"database/sql"
)

type UserRepositorySql struct {
	db *sql.DB
}

func MakeUserRepositorySql(
	iDb *sql.DB,
) UserRepositorySql {
	return UserRepositorySql{
		db: iDb,
	}
}

func (r UserRepositorySql) CreateUser(username string, passwordHash string) (*models.User, error) {
	return nil, nil
}

func (r UserRepositorySql) GetUser(username string) (*models.User, error) {
	return nil, nil
}

func (r UserRepositorySql) GetUserById(iContext context.Context, iId models.UserId) (models.User, error) {
	result := r.db.QueryRowContext(
		iContext,
		`
			SELECT 
				id, 
				username,
				password_hash,
				salt
			FROM "user" u 
			WHERE id=$1
		`,
		iId,
	)

	id := 0
	var username, passwordHash, salt string
	err := result.Scan(
		&id,
		&username,
		&passwordHash,
		&salt,
	)

	if err != nil {
		return models.User{}, err
	}

	user := models.MakeUser(models.UserId(id), username, passwordHash)
	return user, nil
}

func (r UserRepositorySql) FindUserWithPublicKey(iPublicKey string) (models.User, error) {
	result := r.db.QueryRow(`
		SELECT 
			u.id, 
			u.username,
			u.password_hash,
			u.salt
		FROM "user" u 
		JOIN "user_key" k 
		ON u.id = k.owner_id 
		JOIN "public_key" pk 
		ON pk.value = $1 AND pk.id = k.public_key_id
	`, iPublicKey)

	id := 0
	var username, passwordHash, salt string
	err := result.Scan(
		&id,
		&username,
		&passwordHash,
		&salt,
	)

	if err == sql.ErrNoRows {
		return models.User{}, common.NotFound
	}

	if err != nil {
		return models.User{}, err
	}

	ret := models.MakeUser(
		models.UserId(id),
		username,
		passwordHash,
	)
	return ret, nil
}

func (r UserRepositorySql) DoesUserExist(username string) (bool, error) {
	return false, nil
}
