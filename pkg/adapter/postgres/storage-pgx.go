package storage

import (
	typos "31arthur/drive-editor/models"
	"database/sql"
	"log"
	"time"

	"github.com/jackc/pgx"
)

type PGXStore struct {
	Pool *pgx.ConnPool
}

// defining this interface and the functions of it, help in accessing
// these functions from outside of this package
type PGXStorage interface {
	InsertFileData([]typos.GFile) error
	UseDriveTS() (time.Time, error)
	UpdateGDriveTS() error
	UpdateFileDetails([]typos.GFile) error
	UpdateFileRequest(typos.EGFile) error
	UpdateSummary([]typos.SummaryFile) error
	AccessAll() []typos.GFile
}

// establishes the PGX variable, establishes the store and returns PGXStore.
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

// it is called to initiate creation of tables, if they are not already created
func (p *PGXStore) Init() error {
	if err := p.CreatefileTable(); err != nil {
		return err
	}
	if err := p.CreateTimeStampTable(); err != nil {
		return err
	}
	return nil
}

// for creating the main magnum file table
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
			summary varchar(512),
			delivery_mode varchar(255),
			delivery_id varchar(255),
			file_url varchar(255)
		);
		`
	_, err := p.Pool.Exec(query)
	return err
}

// for creating the timestamps table
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

// for creation of new rows, by inserting table. Handle an array of file details.
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
				delivery_id,
				file_url
			 ) values($1, $2, $3, $4,$5,$6,$7,$8,$9,$10,$11,$12)`,
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
			file.DeliveryID,
			file.FURL)

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

//updates the timestamp field, so that it can check with the modified timestamp

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

// for accessing the timestamp
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

// this is used to update the fields of the file, and if some
// fields are empty, they are left with the previous values
func (p *PGXStore) UpdateFileRequest(file typos.EGFile) error {

	sqlStatement := `
		UPDATE magnum
		SET
			case_number = COALESCE($2, case_number),			
			letter_type = COALESCE($3, letter_type),
			summary = COALESCE($4, summary),
			delivery_mode = COALESCE($5, delivery_mode),
			delivery_id = COALESCE($6, delivery_id),
			touched = $7
		WHERE
			id = $1
	`
	_, err := p.Pool.Exec(
		sqlStatement,
		file.ID,
		sql.NullString{String: file.CaseNumber, Valid: file.CaseNumber != ""},
		sql.NullString{String: file.LetterType, Valid: file.LetterType != ""},
		sql.NullString{String: file.Summary, Valid: file.Summary != ""},
		sql.NullString{String: file.DeliveryMode, Valid: file.DeliveryMode != ""},
		sql.NullString{String: file.DeliveryID, Valid: file.DeliveryID != ""},
		true)

	if err != nil {
		return err
	}

	return nil
}

// the greatest use of concurrency by me, I
// could drag the time for server initialization from more than 10-15 minutes
// to mere seconds. This updates the summary of the field individually.
func (p *PGXStore) UpdateSummary(files []typos.SummaryFile) error {

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
			summary = $1 
            WHERE id = $2
        `, file.Summary, file.ID)

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

// gives all the fields and their values in the table
func (p *PGXStore) AccessAll() []typos.GFile {

	tx, err := p.Pool.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	query := "SELECT * FROM magnum"
	rows, err := tx.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	magnumRows := []typos.GFile{}

	for rows.Next() {
		row := typos.GFile{}
		err := rows.Scan(
			&row.ID,
			&row.LID,
			&row.FileName,
			&row.CreatedTime,
			&row.ModifiedTime,
			&row.Touched,
			&row.CaseNumber,
			&row.LetterType,
			&row.Summary,
			&row.DeliveryMode,
			&row.DeliveryID,
			&row.FURL)
		if err != nil {
			log.Fatal(err)
		}

		magnumRows = append(magnumRows, row)
	}
	return magnumRows
}
