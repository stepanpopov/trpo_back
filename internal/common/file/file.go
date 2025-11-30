package file

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// CreateFile creates a new file in the specified directory with a unique name
// based on the hash of the file's content and the provided extension.
//
// The function calculates the hash of the file's content to generate a unique
// filename, checks if a file with the same name already exists, and creates
// the file if it does not exist. The file's content is written to the newly
// created file. The file's read position is restored to its original position
// after the operation.
//
// Parameters:
//   - file: An io.ReadSeeker representing the file to be saved. The file's read
//     position will be restored to its original position after the operation.
//   - extension: A string representing the file extension to be appended to the
//     generated filename.
//   - dirPath: A string representing the directory path where the file will be saved.
//
// Returns:
//   - filename: A string representing the generated filename (including the extension).
//   - path: A string representing the full path of the created file.
//   - error: An error if the file creation process fails at any step.
func CreateFile(file io.ReadSeeker, extension string, dirPath string) (filename string, path string, _ error) {
	beginPosition, err := file.Seek(0, io.SeekCurrent) // save begin position
	if err != nil {
		return "", "", fmt.Errorf("can't do file seek: %w", err)
	}

	filename, err = FileHash(file, extension)
	if err != nil {
		return "", "", fmt.Errorf("can't get fil hash: %w", err)
	}

	path = filepath.Join(dirPath, filename)

	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		newFD, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0o644)
		if err != nil {
			return "", "", fmt.Errorf("can't create file to save avatar: %w", err)
		}
		defer newFD.Close()

		if _, err := io.Copy(newFD, file); err != nil {
			return "", "", fmt.Errorf("can't write sent avatar to file: %w", err)
		}
	} else if err != nil {
		return "", "", fmt.Errorf("can't check file stat: %w", err)
	}

	if _, err = file.Seek(beginPosition, io.SeekStart); err != nil { // go back to beginPosition
		return "", "", fmt.Errorf("can't do file seek: %w", err)
	}

	return filename, path, nil
}

// FileHash generates a unique filename based on the SHA-256 hash of the file's content
// and appends the provided file extension.
//
// The function reads the content of the file to compute its hash, ensuring that the
// generated filename is unique for the given content. The file's read position is
// restored to its original position after the operation.
//
// Parameters:
//   - file: An io.ReadSeeker representing the file whose hash will be computed.
//   - extension: A string representing the file extension to be appended to the hash.
//
// Returns:
//   - A string representing the generated filename (hash + extension).
//   - An error if the file's hash could not be computed or if there was an issue
//     reading/seeking the file.
func FileHash(file io.ReadSeeker, extension string) (string, error) {
	beginPosition, err := file.Seek(0, io.SeekCurrent) // save begin position
	if err != nil {
		return "", fmt.Errorf("can't do file seek: %w", err)
	}

	// Create standard filename
	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", fmt.Errorf("can't write sent avatar to hasher: %w", err)
	}
	newFileName := hex.EncodeToString(hasher.Sum(nil))

	if _, err = file.Seek(beginPosition, io.SeekStart); err != nil { // go back to beginPosition
		return "", fmt.Errorf("can't do file seek: %w", err)
	}

	return newFileName + extension, nil
}
