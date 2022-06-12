package models

type User struct {
	Id           int    `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
}

func MakeUser(
	iId int,
	iUsername string,
	iPasswordHash string,
) User {
	return User{
		Id:           iId,
		Username:     iUsername,
		PasswordHash: iPasswordHash,
	}
}
