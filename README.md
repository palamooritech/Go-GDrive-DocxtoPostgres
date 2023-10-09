# GoDocxToPostgres

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

**GoDocxToPostgres** is a powerful Go application designed to simplify the management of .docx files stored in Google Drive by efficiently streaming them into a PostgreSQL database. This versatile tool provides a seamless solution for organizing and associating your document files with various cases, projects, or applications.

## Features

- **Google Drive Integration**: Connects seamlessly with Google Drive API to access and retrieve .docx files.
- **PostgreSQL Database**: Streamlines file content, metadata, and associations into a PostgreSQL database for easy retrieval and manipulation.
- **File Management**: Efficiently handles .docx files, including their content, metadata, and organization.
- **Case Association**: Enables the association of documents with specific cases or projects, enhancing document management.
- **Scalability**: Designed for scalability, making it suitable for both small-scale and large-scale document management.

## Prerequisites

Before you can use this project, you'll need to set up a few things:

1. **Create a Google Service Account:**
   - Go to the [Google Cloud Console](https://console.cloud.google.com/).
   - Create a new project.
   - Navigate to the "IAM & Admin" section.
   - Create a new Service Account.
   - Download the JSON key file for your Service Account.

2. **Share Your Google Drive Folder:**
   - Share the Google Drive folder containing your documents with the email address of your Service Account.

3. **Place the Service Account Key File:**
   - Place the downloaded `service-account-key.json` file in the `pkg/adapter/gdrive` folder of this project.

4. **Configure Mimetype and Customize Code in pkg/adapter/gdrive/supplements.go:**
   - Adjust the Mimetype as per your file requirements in `supplements.go` and `docx-supplements.go`.
   - Modify the code in these files to suit your specific needs.

5. **Dockerize PostgreSQL and pgAdmin4:**
   - Dockerize your PostgreSQL database and pgAdmin4 in separate containers.

6. **Running the Project:**
   - Once everything is set up, you can run the project from the commands in the 'Makefile':
     ```
     make run
     ```

## Usage

- Your Go application will now be able to access and process documents from Google Drive and store them in the PostgreSQL database.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.


## Contributions

Contributions are welcome! Whether you're fixing a bug, improving the documentation, or adding new features, your contributions help make GoDocxToPostgres better for everyone. See our Contribution Guidelines for more details.


      *
     / \
    /___\
