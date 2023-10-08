package drive

import (
	typos "31arthur/drive-editor/models"
	"31arthur/drive-editor/pkg/domain"
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

//the main function called every 1 hour from the main.go file

func DriveAdapter(s *domain.APIServer) {
	// Load your service account JSON key file.
	keyFile := "pkg/adapter/gdrive/service-account-key.json"
	allFileDetails := []typos.GFile{}
	summaryFileDetails := []typos.SummaryFile{}
	updateFD := []typos.GFile{}
	index := 1
	done := make(chan struct{})

	fmt.Println("The Port Address:", s.ListenAddr)
	// Create a context and authenticate using the service account key file.
	ctx := context.Background()
	// withCredentialsFile function used to parse through the service account key
	service, err := drive.NewService(ctx, option.WithCredentialsFile(keyFile))
	if err != nil {
		log.Fatalf("Unable to create Drive service: %v", err)
	}

	// List the files in your Google Drive.
	files, err := listAllFiles(service, s)
	if err != nil {
		log.Fatalf("Unable to list files: %v", err)
	}

	// Print the list of files.
	// fmt.Println("Files in your Google Drive:")
	for _, file := range files {

		// return a new instance of GFile
		temp := HandleCreateFileField(*file, service)

		if temp.CreatedTime.After(HandleDriveTS(s)) {
			allFileDetails = append(allFileDetails, temp)

			// writing a go subroutine to seperately execute the summary extraction
			go func(temp typos.GFile) {
				summaryFileDetails = append(summaryFileDetails, HandleSummary(temp.ID, service))
				fmt.Printf("\nProcessed file no: %v / %v", index, len(files))
				// fmt.Printf(temp.Summary + "\n")
				index += 1
				//creating a label reference
				done <- struct{}{}
			}(temp)
		} else {
			updateFD = append(updateFD, temp)
		}
	}

	fmt.Println("Total files: ", len(files))
	//the collective Database Update functions
	DBUpdates(s, allFileDetails, updateFD)

	//through this it is waiting for all the summaries to be processed
	for range files {
		<-done
	}

	// this code will get executed only after the Summary Extraction subroutine is completed.
	if err := s.Store.UpdateSummary(summaryFileDetails); err != nil {
		fmt.Println("Update Summary: ", err)
	}
}

func listAllFiles(service *drive.Service, s *domain.APIServer) ([]*drive.File, error) {

	allFiles := []*drive.File{}
	pageToken := ""
	timeStamp := HandleDriveTS(s)
	// mimeType := "application/pdf"
	mimeType := "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	fmt.Println(timeStamp.Format(time.RFC3339))
	query := fmt.Sprintf("modifiedTime > '%s' and mimeType = '%s'", timeStamp.Format(time.RFC3339), mimeType)

	// Accessing the files in the directory
	for {
		query := service.Files.List().Q(query).Fields("files(id, name, createdTime, modifiedTime, labelInfo)").PageSize(1000).PageToken(pageToken)
		files, err := query.Do()
		if err != nil {
			return nil, err
		}
		for _, file := range files.Files {
			if !strings.Contains(file.Name, "~") {
				allFiles = append(allFiles, file)
			}
		}
		// This is used for pdf files, as there's no issue of temporary files
		// allFiles = append(allFiles, files.Files...)
		fmt.Printf("PageToken No : %v, files length: %v", pageToken, len(allFiles))
		if files.NextPageToken == "" {
			break
		}
		pageToken = files.NextPageToken
	}

	return allFiles, nil
}

func DBUpdates(s *domain.APIServer, allFileDetails []typos.GFile, updateFD []typos.GFile) {
	if err := s.Store.InsertFileData(allFileDetails); err != nil {
		fmt.Println("GDrive: ", err)
	}

	if err := s.Store.UpdateGDriveTS(); err != nil {
		fmt.Println("TS Drive: ", err)
	}

	if err := s.Store.UpdateFileDetails(updateFD); err != nil {
		fmt.Println("Updated Files in Drive: ", err)
	}

}
