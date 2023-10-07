package storage

import (
	typos "31arthur/drive-editor/models"
	"log"
	"time"

	"github.com/jackc/pgx"
)

type PGXStore struct {
	Pool *pgx.ConnPool
}

type PGXStorage interface {
	InsertFileData([]typos.GFile) error
	UseDriveTS() (time.Time, error)
	UpdateGDriveTS() error
	UpdateFileDetails([]typos.GFile) error
}

func NewPgxStore() (*PGXStore, error) {
	connConfig := pgx.ConnConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "gobank",
		Database: "postgres",
	}

	poolConfig := pgx.ConnPoolConfig{
		ConnConfig:     connConfig,
		MaxConnections: 5,
	}

	pool, err := pgx.NewConnPool(poolConfig)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &PGXStore{Pool: pool}, nil
}

func (p *PGXStore) Init() error {
	if err := p.CreatefileTable(); err != nil {
		return err
	}
	if err := p.CreateTimeStampTable(); err != nil {
		return err
	}
	return nil
}

func (p *PGXStore) CreatefileTable() error {
	query := `
		create table if not exists magnum(
			id varchar(255) primary key,
			file_name varchar(255),
			letter_id varchar(255),
			created_time timestamptz,
			modified_time timestamptz,
			touched boolean,
			case_number varchar(255),
			letter_type varchar(255),
			summary varchar(255),
			delivery_mode varchar(255),
			delivery_id varchar(255)
		);
		`
	_, err := p.Pool.Exec(query)
	return err
}

func (p *PGXStore) CreateTimeStampTable() error {
	query := `
		create table if not exists timestamps(
			id serial primary key,
			time_stamp timestamptz,
			purpose varchar(255) UNIQUE
		);
		`
	_, err := p.Pool.Exec(query)
	return err
}

func (p *PGXStore) InsertFileData(files []typos.GFile) error {
	// fmt.Println(files)
	// Start a transaction.
	tx, err := p.Pool.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Iterate over the files and insert them into the database.
	for _, file := range files {
		_, err := tx.Exec(`
            INSERT INTO magnum(
				id ,
				file_name ,
				letter_id,
				created_time ,
				modified_time ,
				touched,
				case_number,
				letter_type,
				summary,
				delivery_mode,
				delivery_id
			 ) values($1, $2, $3, $4,$5,$6,$7,$8,$9,$10,$11)`,
			file.ID,
			file.FileName,
			file.LID,
			file.CreatedTime,
			file.ModifiedTime,
			file.Touched,
			file.CaseNumber,
			file.LetterType,
			file.Summary,
			file.DeliveryMode,
			file.DeliveryID)

		if err != nil {
			return err
		}
	}

	// Commit the transaction to save the changes.
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (p *PGXStore) UpdateGDriveTS() error {

	tx, err := p.Pool.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	purpose := "gdrive-access"
	currentTime := time.Now().In(time.UTC)

	_, err1 := tx.Exec(`
			INSERT INTO timestamps (purpose,time_stamp)
			VALUES ($1, $2)
			ON CONFLICT (purpose)
			DO UPDATE SET time_stamp = excluded.time_stamp
        `, purpose, currentTime)

	if err1 != nil {
		return err1
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil

}

func (p *PGXStore) UseDriveTS() (time.Time, error) {

	var timeStamp time.Time

	err := p.Pool.QueryRow(`
		SELECT time_stamp FROM timestamps WHERE purpose = $1
		`, "gdrive-access").Scan(&timeStamp)

	if err != nil {
		return time.Time{}, err
	}

	return timeStamp, nil
}

func (p *PGXStore) UpdateFileDetails(files []typos.GFile) error {

	//start a transaction
	tx, err := p.Pool.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	for _, file := range files {
		_, err := tx.Exec(`
            UPDATE magnum
            SET 
			file_name = $1 ,
			letter_id = $2,
			modified_time = $3,
            WHERE id = $4
        `, file.FileName, file.LID, file.ModifiedTime, file.ID)

		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil

}
