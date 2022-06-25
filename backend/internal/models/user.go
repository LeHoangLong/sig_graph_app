package models

type UserId int
type User struct {
	Id           UserId `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
}

func MakeUser(
	iId UserId,
	iUsername string,
	iPasswordHash string,
) User {
	return User{
		Id:           iId,
		Username:     iUsername,
		PasswordHash: iPasswordHash,
	}
}
