package internal

import (
	"context"
	"database/sql"
	"log"
	"telecart/pkg/svc"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

type SqlLiteStore struct {
	conn *sql.DB
}

func NewSqlLiteStore(dns string) (svc.Store, error) {
	var err error

	db, err := sql.Open("sqlite3", dns)
	if err != nil {
		return nil, errors.Wrap(err, "open db connection: "+dns)
	}

	return &SqlLiteStore{conn: db}, nil
}

func (s *SqlLiteStore) Init() error {
	query := `
CREATE TABLE IF NOT EXISTS messages(
id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, 
name TEXT,
message TEXT
)`

	if _, err := s.conn.Exec(query); err != nil {
		return errors.Wrap(err, "exec query: "+query)
	}

	return nil
}

func (s *SqlLiteStore) Save(ctx context.Context, msg *svc.Message) error {
	stmp, err := s.conn.Prepare("INSERT INTO messages (name, message) VALUES (?, ?)")
	if err != nil {
		return errors.Wrap(err, "prepare stmp")
	}

	result, err := stmp.Exec(msg.Name, msg.Message)
	if err != nil {
		return errors.Wrap(err, "db exec")
	}

	log.Println(result.LastInsertId())
	return nil
}

func (s *SqlLiteStore) Close() error {
	return s.conn.Close()
}
