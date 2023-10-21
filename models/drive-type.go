package typos

import (
	"time"
)

type GFile struct {
	ID           string    `json:"id"`
	LID          string    `json:"lid"`
	FileName     string    `json:"file_name"`
	CreatedTime  time.Time `json:"created_time"`
	ModifiedTime time.Time `json:"modified_time"`
	Touched      bool      `json:"touched"`
	CaseNumber   string    `json:"case_number"`
	LetterType   string    `json:"letter_type"`
	Summary      string    `json:"summary"`
	DeliveryMode string    `json:"delivery_mode"`
	DeliveryID   string    `json:"delivery_id"`
	FURL         string    `json:"file_url"`
}

type EGFile struct {
	ID           string `json:"id"`
	LID          string `json:"lid"`
	FileName     string `json:"file_name"`
	Touched      bool   `json:"touched"`
	CaseNumber   string `json:"caseNumber"`
	LetterType   string `json:"letterType"`
	Summary      string `json:"summary"`
	DeliveryMode string `json:"deliveryMode"`
	DeliveryID   string `json:"deliveryID"`
}

type SummaryFile struct {
	ID      string `json:"id"`
	Summary string `json:"summary:"`
}

type SearchData struct {
	Keyword string `json:"keyword"`
}

func NewFile(
	id string, fileName string, createdTime time.Time,
	modifiedTime time.Time, letterID string, letterType string,
	fileURL string) GFile {
	//creates a new instance of GFile to avoid immutability issues
	temp := new(GFile)
	temp.ID = id
	temp.LID = letterID
	temp.FileName = fileName
	temp.CreatedTime = createdTime
	temp.ModifiedTime = modifiedTime
	temp.Touched = false
	temp.CaseNumber = ""
	temp.LetterType = letterType
	temp.Summary = ""
	temp.DeliveryMode = ""
	temp.DeliveryID = ""
	temp.FURL = fileURL
	return *temp
}

func NewSummaryFile(id string, summary string) SummaryFile {
	temp := new(SummaryFile)
	temp.ID = id
	temp.Summary = summary
	return *temp
}

func AccessFileRow(id string, letterID string, fileName string,
	createdTime time.Time, modifiedTime time.Time, touched bool,
	caseNumber string, letterType string, summary string,
	deliveryMode string, deliveryID string, fileURL string) GFile {

	temp := new(GFile)
	temp.ID = id
	temp.LID = letterID
	temp.FileName = fileName
	temp.CreatedTime = createdTime
	temp.ModifiedTime = modifiedTime
	temp.Touched = touched
	temp.CaseNumber = caseNumber
	temp.LetterType = letterType
	temp.Summary = summary
	temp.DeliveryMode = deliveryMode
	temp.DeliveryID = deliveryID
	temp.FURL = fileURL

	return *temp
}
