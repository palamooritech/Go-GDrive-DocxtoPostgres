package storage

import (
	typos "31arthur/drive-editor/models"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateFileDetails(*typos.GFile) error
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=postgres dbname=postgres password=gobank sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) Init() error {
	return s.CreateFileDetailsTable()
}

func (s *PostgresStore) CreateFileDetailsTable() error {
	query := `create table if not exists filedetails (
		id varchar(255) primary key,
        file_name varchar(255),
		created_time timestamp,
		modified_time timestamp,
		touched boolean,
        case_number varchar(255),
        letter_type varchar(255),
        summary varchar(255),
		delivery_mode varchar(255),
        delivery_id varchar(255)
	)`

	_, err := s.db.Exec(query)
	// s.db.Close()
	return err
}

func (s *PostgresStore) CreateFileDetails(fileDetails *typos.GFile) error {
	query := `insert into filedetails(
		id ,
        file_name ,
		created_time ,
		modified_time ,
		touched,
        case_number,
        letter_type,
        summary,
		delivery_mode,
        delivery_id
	 ) values($1, $2, $3, $4,$5,$6,$7,$8,$9,$10)`

	resp, err := s.db.Query(query,
		fileDetails.ID,
		fileDetails.FileName,
		fileDetails.CreatedTime,
		fileDetails.ModifiedTime,
		fileDetails.Touched,
		fileDetails.CaseNumber,
		fileDetails.LetterType,
		fileDetails.Summary,
		fileDetails.DeliveryMode,
		fileDetails.DeliveryID,
	)

	if err != nil {
		log.Fatal(err)
		return err
	}

	fmt.Printf("%+v\n", resp)
	// s.db.Close()
	return nil
}
