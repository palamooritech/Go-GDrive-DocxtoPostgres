package drive

import (
	typos "31arthur/drive-editor/models"
	"31arthur/drive-editor/pkg/domain"
	"log"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"google.golang.org/api/drive/v3"
)

// formats the date with the timezone value
func dateFormatter(date string) time.Time {
	timestamp, err := time.Parse(time.RFC3339, date)
	if err != nil {
		log.Fatalf("This is an error: %v", err)
	}

	return timestamp

	// fmt.Println(reflect.TypeOf(timestamp))
}

// it is used for creation og a row of GFile type and return a new instance
func HandleCreateFileField(file drive.File, service *drive.Service) typos.GFile {
	file_url := "https://drive.google.com/file/d/" + file.Id + "/view?usp=sharing"
	// file_url := file.WebViewLink
	// fmt.Println("File Url", file_url)
	fileDetails := typos.NewFile(
		file.Id,
		strings.TrimSuffix(file.Name, ".docx"),
		dateFormatter(file.CreatedTime),
		dateFormatter(file.ModifiedTime),
		LetterIDExtract(file.Name),
		LetterTypeDeduce(file.Name),
		file_url,
	)
	return fileDetails

}

// extracts the letter id from the name of the file => []
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

// deduce the type of the letter based on the name of the file
func LetterTypeDeduce(name string) string {
	input := strings.ToLower(name)

	// Words to search for
	wordsToSearch := []string{"rti", "followup", "resolution", "affidavit", "legal notice", "follow up"}

	result := ""
	// Iterate over the words and search for each one
	for _, word := range wordsToSearch {
		if strings.Contains(input, word) {
			if result == "followup" || result == "follow up" {
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

// takes the timestamp of the last modified value
func HandleDriveTS(s *domain.APIServer) time.Time {
	timestamp, err := s.Store.UseDriveTS()
	if err != nil {
		date := time.Date(2008, time.January, 1, 0, 0, 0, 0, time.UTC)
		return date
	}
	return timestamp
}

// extracts the summary of the file from the body of the docx file.
func HandleSummary(id string, service *drive.Service) typos.SummaryFile {

	summary, err := downloadFile(id, service)
	if err != nil {
		log.Fatalf("Error downloading .docx file: %v", err)
	}

	return typos.NewSummaryFile(id, truncateString(summary, 500))
}

func truncateString(s string, maxLength int) string {
	if utf8.RuneCountInString(s) > maxLength {
		// Convert the string to a rune slice
		runes := []rune(s)

		// Truncate to 500 characters
		truncatedString := string(runes[:maxLength])
		return truncatedString
	} else {
		// If the string is already 500 characters or shorter, keep it as is
		return s
	}
}
