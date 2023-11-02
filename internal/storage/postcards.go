package storage

import (
	"database/sql"
	"fmt"
	"telegramBot/internal/model"
)

type Database struct {
	host      string
	user      string
	password  string
	port      string
	sslMode   string
	tableName string
	dbname    string
}

func NewDatabase() *Database {
	return &Database{
		host:     "localhost",
		port:     "5432",
		user:     "postgres",
		password: "postgres",
		dbname:   "users_bot",
		sslMode:  "disable",
	}
}

type PostcardsPostgresStorage struct {
	db *Database
}

func NewPostcardsPostgresStorage(db *Database) *PostcardsPostgresStorage {
	return &PostcardsPostgresStorage{db: db}
}

// GetPostcardsFromStorage returns all postcards from storage
func (p *PostcardsPostgresStorage) GetPostcardsFromStorage() ([][2]string, error) {
	var postcards [][2]string
	var name2path = [2]string{}

	dbInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		p.db.host, p.db.port, p.db.user, p.db.password, p.db.dbname, p.db.sslMode)

	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return postcards, err
	}
	defer func() {
		if err = db.Close(); err != nil {
			return
		}
	}()

	query := "SELECT * FROM postcards"

	result, err := db.Query(query)
	if err != nil {
		return postcards, err
	}
	defer func() {
		if err = result.Close(); err != nil {
			return
		}
	}()

	for result.Next() {
		if err = result.Scan(&name2path[0], &name2path[1]); err != nil {
			return nil, err
		}
		postcards = append(postcards, name2path)
	}

	return postcards, err
}

func (p *PostcardsPostgresStorage) AddPostcardToStorage(postcard *model.Postcard) error {

	dbInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		p.db.host, p.db.port, p.db.user, p.db.password, p.db.dbname, p.db.sslMode)

	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return err
	}
	defer func() {
		if err = db.Close(); err != nil {
			return
		}
	}()

	query := "INSERT INTO postcards (Name, path) VALUES ($1, $2)"

	if _, err = db.Exec(query, postcard.Name, postcard.Path); err != nil {
		return err
	}

	return nil
}

func (p *PostcardsPostgresStorage) RemovePostcardFromStorage(postcard *model.Postcard) error {
	dbInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		p.db.host, p.db.port, p.db.user, p.db.password, p.db.dbname, p.db.sslMode)

	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return err
	}
	defer func() {
		if err = db.Close(); err != nil {
			return
		}
	}()

	query := "DELETE FROM postcards WHERE Name=$1"
	if _, err = db.Exec(query, postcard.Name); err != nil {
		return err
	}
	return nil
}
