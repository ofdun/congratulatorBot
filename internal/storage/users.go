package storage

import (
	"database/sql"
	"fmt"
	"os"
)

func AddUserToMailing(id int64, time_ int64) error {
	dbInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"), os.Getenv("SSLMODE"))

	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return err
	}
	defer func() {
		if err = db.Close(); err != nil {
			return
		}
	}()

	query := "INSERT INTO users (id, mailing_time) VALUES ($1, $2)"
	if _, err = db.Exec(query, id, time_); err != nil {
		return err
	}

	return nil
}

func RemoveUserFromMailing(id int64) error {
	dbInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"), os.Getenv("SSLMODE"))

	db, err := sql.Open("postgres", dbInfo)
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

func GetIDsFromTime(time int) ([]int64, error) {
	dbInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"), os.Getenv("SSLMODE"))

	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = db.Close(); err != nil {
			return
		}
	}()

	query := "SELECT id FROM users WHERE mailing_time=$1"
	var Ids []int64
	result, err := db.Query(query, time)
	if err != nil {
		return nil, err
	}
	for result.Next() {
		var id int64
		if err = result.Scan(&id); err != nil {
			return nil, err
		}
		Ids = append(Ids, id)
	}
	if err != nil {
		return nil, err
	}

	return Ids, nil
}

func GetIfUserIsMailing(id int64) (bool, error) {
	dbInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"), os.Getenv("SSLMODE"))

	db, err := sql.Open("postgres", dbInfo)
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
