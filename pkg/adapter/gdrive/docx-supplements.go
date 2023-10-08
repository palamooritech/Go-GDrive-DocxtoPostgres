package drive

import (
	"io"
	"log"
	"os"
	"strings"

	"github.com/gonfva/docxlib"
	"google.golang.org/api/drive/v3"
)

func downloadFile(id string, service *drive.Service) (string, error) {

	// Call the Drive API to download the file content. Returns a response stream.
	file, err := service.Files.Get(id).Download()
	if err != nil {
		log.Fatalf("Unable to download file: %v", err)
	}
	defer file.Body.Close()

	// Create a temporary file to store the downloaded content.
	tempFile, err := os.CreateTemp("", "temp-docx-*.docx")
	if err != nil {
		log.Fatalf("Unable to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Copy the downloaded content of file.body (response stream)
	// to the temporary file.
	// fmt.Println("the temporary file name: ", tempFile.Name())
	_, err = io.Copy(tempFile, file.Body)
	if err != nil {
		log.Fatalf("Unable to copy file content to temporary file: %v", err)
	}
	summary, err1 := Parsefile(tempFile.Name())
	if err1 != nil {
		log.Fatalf("Unable to parse file: %v %v", err1, summary)
	}
	// fmt.Println(summary)
	os.Remove(tempFile.Name())

	return summary, nil
}

func Parsefile(filename string) (string, error) {
	readFile, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	fileInfo, err := readFile.Stat()
	if err != nil {
		return "", err
	}
	size := fileInfo.Size()
	// uses the excellent gonfva pkg to parse .docx files.
	doc, err := docxlib.Parse(readFile, int64(size))
	if err != nil {
		return "", err
	}

	// for reference see the documentation of https://github.com/gonfva/docxlib/blob/master/main/main.go
	// paragraphs := []string{}
	holdingVal := ""
	for _, para := range doc.Paragraphs() {
		childPara := ""
		for _, child := range para.Children() {

			if child.Run != nil {
				childPara = childPara + "" + child.Run.Text.Text
			}

		}
		content, flag := SubCheck(childPara)
		if flag {
			// fmt.Println("the Summary:", content)
			return content, nil
		}
		if len(holdingVal) < len(childPara) {
			holdingVal = childPara
			if len(holdingVal) > 400 {
				break
			}
		}

		// paragraphs = append(paragraphs, childPara)
	}
	// fmt.Println("the largest", holdingVal)

	return holdingVal, nil

}

func SubCheck(para string) (string, bool) {

	if strings.Contains(para, "Sub:") {
		return strings.Replace(para, "Sub:", "", 1), true
	}
	return "", false
}
