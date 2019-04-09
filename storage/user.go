package storage

import (
	"github.com/jmoiron/sqlx"
)

type User struct {
	UserName     string `db:"user_name"`
	PasswordHash string `db:"password_hash"`
}

// LoginUser check user name and password
func LoginUser(userName string) (User, error) {
	var usr User
	err := sqlx.Get(db, &usr, `
		select user_name,
		password_hash
		from users
		where user_name=$1`,
		userName,
	)
	if err != nil {
		return usr, err
	}
	return usr, nil
}

func UpdateUserPassword(usr User) error {
	_, err := db.Exec(`
		update users
		set password_hash=$2
		where user_name=$1`,
		usr.UserName,
		usr.PasswordHash,
	)
	if err != nil {
		return err
	}

	return nil
}
