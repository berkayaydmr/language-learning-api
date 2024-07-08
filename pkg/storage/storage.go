package storage

import (
	"context"
	"database/sql"
	_ "embed"

	customerr "github.com/berkayaydmr/language-learning-api/pkg/error"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

var (
	//go:embed sql/create_tables.sql
	createTablesQuery string

	//go:embed sql/seed_data.sql
	seedDataQuery string

	//go:embed sql/list.sql
	listWordsQuery string

	//go:embed sql/create.sql
	createWordQuery string

	//go:embed sql/update.sql
	updateWordQuery string

	//go:embed sql/delete.sql
	deleteWordQuery string
)

type PrimaryKey int64

type Word struct {
	ID              PrimaryKey `json:"id"`
	Word            string     `json:"word"`
	Translation     string     `json:"translation"`
	Language        string     `json:"language"`
	ExampleSentence string     `json:"exampleSentence"`
}

func (w *Word) Scan(query *sql.Rows) error {
	return query.Scan(&w.ID, &w.Word, &w.Translation, &w.Language, &w.ExampleSentence)
}

type Update struct {
	Word            *string `json:"word"`
	Translation     *string `json:"translation"`
	Language        *string `json:"language"`
	ExampleSentence *string `json:"exampleSentence"`
}

type Storage interface {
	Open(ctx context.Context, dsn string) error
	Close() error
	CreateTables(ctx context.Context) error
	SeedData(ctx context.Context) error

	List(ctx context.Context) ([]Word, error)
	Create(ctx context.Context, word Word) (*PrimaryKey, error)
	Update(ctx context.Context, id PrimaryKey, update Update) error
	Delete(ctx context.Context, id PrimaryKey) error
}

type storage struct {
	db *sql.DB
}

func New() Storage {
	return &storage{}
}

func (s *storage) Open(ctx context.Context, dsn string) error {
	conn, err := sql.Open("sqlite", dsn)
	if err != nil {
		return err
	}

	s.db = conn
	return nil
}

func (s *storage) Close() error {
	return s.db.Close()
}

func (s *storage) CreateTables(ctx context.Context) error {
	_, err := s.db.Exec(createTablesQuery)

	if err != nil {
		if liteErr, ok := err.(*sqlite.Error); ok {
			// TODO: find table exist error code
			if liteErr.Code() == sqlite3.P4_TABLE {
				return nil
			}
		}
	}

	return err
}

func (s *storage) SeedData(ctx context.Context) error {
	result, err := s.db.ExecContext(ctx, seedDataQuery)
	if err != nil {
		return err
	}

	count, err := result.LastInsertId()
	if err != nil {
		return err
	}

	if count < 0 {
		return customerr.ErrNoneOfSeedDataInserted
	}

	return nil
}

func (s *storage) List(ctx context.Context) ([]Word, error) {
	rows, err := s.db.QueryContext(ctx, listWordsQuery)
	if err != nil {
		return nil, err
	}

	words := []Word{}
	for rows.Next() {
		w := Word{}
		err = w.Scan(rows)
		if err != nil {
			return nil, err
		}

		words = append(words, w)
	}

	if len(words) == 0 {
		return nil, customerr.ErrWordsNotFound
	}

	return words, nil
}

func (s *storage) Create(ctx context.Context, word Word) (*PrimaryKey, error) {
	res, err := s.db.ExecContext(ctx, createWordQuery, word.Word, word.Translation, word.Language, word.ExampleSentence)
	if err != nil {
		if liteErr, ok := err.(*sqlite.Error); ok {
			if liteErr.Code() == sqlite3.P5_ConstraintUnique {
				return nil, customerr.ErrWordAlreadyExist
			}
		}
		return nil, err
	}

	id, err := res.LastInsertId()
	primaryKeyId := PrimaryKey(id)
	return &primaryKeyId, err
}

func (s *storage) Update(ctx context.Context, id PrimaryKey, update Update) error {
	_, err := s.db.ExecContext(ctx, updateWordQuery, update.Word, update.Language, update.Translation, update.ExampleSentence, id)

	return err
}

func (s *storage) Delete(ctx context.Context, id PrimaryKey) error {
	res, err := s.db.ExecContext(ctx, deleteWordQuery)
	if err != nil {
		return err
	}

	r, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if r == 0 {
		return customerr.ErrWordIDNotFound
	}

	return nil
}
