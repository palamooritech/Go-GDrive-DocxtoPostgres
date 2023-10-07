package typos

import "time"

type GFile struct {
	ID           string    `json:"id"`
	LID          string    `json:"lid"`
	FileName     string    `json:"fileName"`
	CreatedTime  time.Time `json:"createdTime"`
	ModifiedTime time.Time `json:"modifiedTime"`
	Touched      bool      `json:"touched"`
	CaseNumber   string    `json:"caseNumber"`
	LetterType   string    `json:"letterType"`
	Summary      string    `json:"summary"`
	DeliveryMode string    `json:"deliveryMode"`
	DeliveryID   string    `json:"deliveryID"`
}

func NewFile(
	id string, fileName string, createdTime time.Time,
	modifiedTime time.Time, letterID string, letterType string) GFile {
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
	return *temp
}
