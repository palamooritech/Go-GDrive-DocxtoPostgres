package drive

import (
	typos "31arthur/drive-editor/models"
	"31arthur/drive-editor/pkg/domain"
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

func DriveAdapter(s *domain.APIServer) {
	// Load your service account JSON key file.
	keyFile := "pkg/adapter/gdrive/service-account-key.json"
	allFileDetails := []typos.GFile{}
	updateFD := []typos.GFile{}

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
		temp := HandleCreateFileField(*file)

		if temp.CreatedTime.After(HandleDriveTS(s)) {
			allFileDetails = append(allFileDetails, temp)
		} else {
			updateFD = append(updateFD, temp)
		}

	}
	fmt.Println("Total files: ", len(files))

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

func listAllFiles(service *drive.Service, s *domain.APIServer) ([]*drive.File, error) {

	allFiles := []*drive.File{}
	pageToken := ""
	timeStamp := HandleDriveTS(s)
	mimeType := "application/pdf"
	fmt.Println(timeStamp.Format(time.RFC3339))
	query := fmt.Sprintf("modifiedTime > '%s' and mimeType = '%s'", timeStamp.Format(time.RFC3339), mimeType)
	// List the files in your Google Drive.
	for {
		query := service.Files.List().Q(query).Fields("files(id, name, createdTime, modifiedTime, labelInfo)").PageSize(1000).PageToken(pageToken)
		files, err := query.Do()
		if err != nil {
			return nil, err
		}

		allFiles = append(allFiles, files.Files...)
		fmt.Printf("PageToken No : %v, files length: %v", pageToken, len(allFiles))
		if files.NextPageToken == "" {
			break
		}
		pageToken = files.NextPageToken
	}

	return allFiles, nil
}
