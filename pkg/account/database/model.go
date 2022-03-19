package database

type User struct {
	UserID         uint64 `json:"user_id"`
	UserName       string `json:"user_name"`
	HashedPassword []byte `json:"password"`
}
