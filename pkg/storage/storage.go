package storage

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"strings"

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
	listQuery string

	//go:embed sql/create.sql
	createQuery string // createQuery

	//go:embed sql/update.sql
	updateQuery string

	//go:embed sql/delete.sql
	deleteQuery string

	//go:embed sql/if_exist.sql
	ifExistQuery string
)

type PrimaryKey int64

type Word struct {
	ID              PrimaryKey `json:"id"`
	Word            string     `json:"word"`
	Translation     string     `json:"translation"`
	Language        string     `json:"language"`
	ExampleSentence string     `json:"exampleSentence"`
}

type Update struct {
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
	Create(ctx context.Context, word Word) (PrimaryKey, error)
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

	if count == 0 {
		return customerr.ErrNoneOfSeedDataInserted
	}

	return nil
}

func (s *storage) List(ctx context.Context) ([]Word, error) {
	rows, err := s.db.QueryContext(ctx, listQuery)
	if err != nil {
		return nil, err
	}

	var words []Word
	for rows.Next() {
		var word Word
		exampleSentence := sql.NullString{}
		err = rows.Scan(&word.ID, &word.Word, &word.Translation, &word.Language, &exampleSentence)
		if err != nil {
			return nil, err
		}
		word.ExampleSentence = exampleSentence.String
		words = append(words, word)
	}

	if len(words) == 0 {
		return nil, customerr.ErrWordsNotFound
	}

	return words, nil
}

func (s *storage) Create(ctx context.Context, word Word) (PrimaryKey, error) {
	res, err := s.db.ExecContext(ctx, createQuery, word.Word, word.Translation, word.Language, word.ExampleSentence)
	if err != nil {
		if liteErr, ok := err.(*sqlite.Error); ok {
			if liteErr.Code() == sqlite3.P5_ConstraintUnique {
				return 0, customerr.ErrWordAlreadyExist
			}
		}
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return PrimaryKey(id), err
}

func (s *storage) Update(ctx context.Context, id PrimaryKey, update Update) error {
	row := s.db.QueryRowContext(ctx, ifExistQuery, id)

	var exist bool
	err := row.Scan(&exist)
	if err != nil {
		return err
	}

	if !exist {
		return customerr.ErrWordIDNotFound
	}

	var builder strings.Builder

	if update.Translation != nil {
		builder.WriteString(fmt.Sprintf("translation = '%s'", *update.Translation))
	}

	if update.Language != nil {
		if builder.Len() > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(fmt.Sprintf("language = '%s'", *update.Language))
	}

	if update.ExampleSentence != nil {
		if builder.Len() > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(fmt.Sprintf("example_sentence = '%s'", *update.ExampleSentence))
	}

	_, err = s.db.ExecContext(ctx, strings.Replace(updateQuery, "{setclause}", builder.String(), 1), id)
	return err
}

func (s *storage) Delete(ctx context.Context, id PrimaryKey) error {
	res, err := s.db.ExecContext(ctx, deleteQuery)
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
