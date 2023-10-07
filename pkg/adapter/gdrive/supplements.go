package drive

import (
	typos "31arthur/drive-editor/models"
	"31arthur/drive-editor/pkg/domain"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"google.golang.org/api/drive/v3"
)

func dateFormatter(date string) time.Time {
	timestamp, err := time.Parse(time.RFC3339, date)
	if err != nil {
		log.Fatalf("This is an error: %v", err)
	}

	return timestamp

	// fmt.Println(reflect.TypeOf(timestamp))
}

func HandleCreateFileField(file drive.File) typos.GFile {
	fileDetails := typos.NewFile(
		file.Id,
		file.Name,
		dateFormatter(file.CreatedTime),
		dateFormatter(file.ModifiedTime),
		LetterIDExtract(file.Name),
		LetterTypeDeduce(file.Name),
	)

	return fileDetails

}

func LetterIDExtract(name string) string {

	pattern := `\[([0-9]+)\]`
	regex := regexp.MustCompile(pattern)

	matches := regex.FindAllStringSubmatch(name, -1)

	for _, match := range matches {
		if len(match) >= 2 {
			// Extract the matched value (e.g., "2023100601") from the submatch
			return match[1]
		}
	}

	return uuid.New().String()
}

func LetterTypeDeduce(name string) string {
	input := strings.ToLower(name)

	// Words to search for
	wordsToSearch := []string{"rti", "followup", "resolution", "affidavit", "legal notice"}

	result := ""
	// Iterate over the words and search for each one
	for _, word := range wordsToSearch {
		if strings.Contains(input, word) {
			if result == "followup" {
				// "followup" takes precedence, so we don't change the result
				continue
			}
			result = word
		}
	}

	if result == "" {
		result = "plea"
	}

	return result
}

func HandleDriveTS(s *domain.APIServer) time.Time {
	timestamp, err := s.Store.UseDriveTS()
	if err != nil {
		date := time.Date(2008, time.January, 1, 0, 0, 0, 0, time.UTC)
		return date
	}
	return timestamp
}
