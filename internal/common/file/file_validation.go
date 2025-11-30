package file

import (
	"errors"
	"io"
	"net/http"
)

// CheckMimeType verifies the MIME type of a file and checks if it matches any of the provided valid types.
//
// This function reads the first 512 bytes of the file to determine its MIME type using
// http.DetectContentType. It then compares the detected MIME type against the list of
// allowed types provided as arguments.
//
// The file's read position is restored to its original position after the MIME type check.
//
// Parameters:
//   - file: An io.ReadSeeker representing the file to check. The file's read position will
//     be restored to its original position after the operation.
//   - correctTypes: A variadic list of strings representing the allowed MIME types.
//
// Returns:
//   - A string representing the detected MIME type of the file.
//   - An error if the file's MIME type could not be determined or if there was an issue
//     reading/seeking the file.
func CheckMimeType(file io.ReadSeeker, correctTypes ...string) (string, error) {
	curPosition, err := file.Seek(0, io.SeekCurrent) // save current position
	if err != nil {
		return "", err
	}

	var fileHeader [512]byte
	if _, err := file.Read(fileHeader[:]); err != nil {
		return "", err
	}

	if _, err = file.Seek(curPosition, io.SeekStart); err != nil { // go back to curPosition
		return "", err
	}

	fileType := http.DetectContentType(fileHeader[:])
	for _, correctType := range correctTypes {
		if correctType == fileType {
			return fileType, nil
		}
	}

	return fileType, errors.New("incorrect type")
}
