package db

import (
	"github.com/skilld-labs/dbr"
	"time"
)

type User struct {
	ID        int
	Name      string
	Username  string
	Email     string
	State     string
	CreatedAt time.Time
}

func (db *DbAPI) GetUserByID(id int) (User, error) {
	user := User{}
	_, err := db.Db.Select("id", "name", "username", "email", "state", "created_at").From("users").Where(dbr.Eq("id", id)).Load(&user)
	return user, err
}
