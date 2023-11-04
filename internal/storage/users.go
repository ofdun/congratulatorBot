package storage

import "database/sql"

type User struct {
	Id   int64
	Time string
}

func (u *User) AddUserToMailing() error {
	dbinfo := "host=localhost port=5432 user=postgres password=postgres dbname=users_bot sslmode=disable"

	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		return err
	}
	defer func() {
		if err = db.Close(); err != nil {
			return
		}
	}()

	query := "INSERT INTO users (id, mailing_time) VALUES ($1, $2)"
	if _, err = db.Exec(query, u.Id, u.Time); err != nil {
		return err
	}

	return nil
}

func RemoveUserFromMailing(id int64) error {
	dbinfo := "host=localhost port=5432 user=postgres password=postgres dbname=users_bot sslmode=disable"

	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		return err
	}
	defer func() {
		if err = db.Close(); err != nil {
			return
		}
	}()

	query := "DELETE FROM users WHERE id=$1"
	if _, err = db.Exec(query, id); err != nil {
		return err
	}

	return nil
}

func (u *User) GetMailingTime() (string, error) {
	dbinfo := "host=localhost port=5432 user=postgres password=postgres dbname=users_bot sslmode=disable"

	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		return "", err
	}
	defer func() {
		if err = db.Close(); err != nil {
			return
		}
	}()

	query := "SELECT mailing_time FROM users WHERE id=$1"
	result, err := db.Query(query, u.Id)
	for result.Next() {
		if err = result.Scan(&u.Time); err != nil {
			return "", err
		}
	}
	if err != nil {
		return "", err
	}

	return u.Time, nil
}

func (u *User) GetIDFromTime() (int64, error) {
	dbinfo := "host=localhost port=5432 user=postgres password=postgres dbname=users_bot sslmode=disable"

	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		return 0, err
	}
	defer func() {
		if err = db.Close(); err != nil {
			return
		}
	}()

	query := "SELECT id FROM users WHERE mailing_time=$1"
	result, err := db.Query(query, u.Time)
	for result.Next() {
		if err = result.Scan(&u.Id); err != nil {
			return 0, err
		}
	}
	if err != nil {
		return 0, err
	}

	return u.Id, nil
}

func GetIfUserIsMailing(id int64) (bool, error) {
	dbinfo := "host=localhost port=5432 user=postgres password=postgres dbname=users_bot sslmode=disable"

	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		return false, err
	}
	defer func() {
		if err = db.Close(); err != nil {
			return
		}
	}()

	query := "SELECT id FROM users WHERE id=$1"
	result, err := db.Query(query, id)
	for result.Next() {
		return true, nil
	}
	if err != nil {
		return false, err
	}

	return false, nil
}
