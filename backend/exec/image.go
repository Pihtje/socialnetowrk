package exec

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/google/uuid"
)

var imagesDir string

func init() {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	imagesDir = filepath.Join(basepath, "../images/")
}

// Returns the file extension based on the request headers "Content-Type".
func getFileExtensionByMimeType(fileHeader *multipart.FileHeader) string {
	contentType := fileHeader.Header.Get("Content-Type")
	switch contentType {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	default:
		return ""
	}
}

// Saves an image and returns the filename.
func SaveImage(file io.Reader, subdirectory string, prefix string, extension string) (string, error) {
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	uniqueFilename := prefix + uuid.New().String() + extension
	path := filepath.Join(imagesDir, subdirectory, uniqueFilename)

	if _, err := os.Stat(filepath.Join(imagesDir, subdirectory)); os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Join(imagesDir, subdirectory), 0755)
		if err != nil {
			return "", err
		}
	}

	if err := os.WriteFile(path, fileBytes, 0644); err != nil {
		return "", err
	}

	relativePath := filepath.Join(subdirectory, uniqueFilename)
	fullPath := GetFullImagePath(relativePath)
	return fullPath, nil
}

// Saves an image file to the "images/user_images" directory.
func SaveUserAvatar(file io.Reader, fileHeader *multipart.FileHeader) (string, error) {
	extension := getFileExtensionByMimeType(fileHeader)
	return SaveImage(file, "user_images", "user_pic_", extension)
}

// Renames an image file within the "images" directory.
func RenameImage(oldPath string, newPath string) error {
	return os.Rename(oldPath, newPath)
}

// Removes an image file from the "images" directory.
func RemoveExistingUserAvatar(userID int64) error {
	user, err := GetUser("users.user_id", fmt.Sprint(userID))
	if err != nil {
		fmt.Println("RemoveExistingUserAvatar err 1: ", err)
		return err
	}

	if !user[0].ImageUrl.Valid {
		fmt.Println("RemoveExistingUserAvatar: nothing to remove, returning early")
		return nil
	}

	prefix := fmt.Sprintf("user_pic_%d_", userID)
	dir := filepath.Join(imagesDir, "user_images")

	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), prefix) {
			fmt.Println("RemoveExistingUserAvatar: file.Name(): ", file.Name(), ", prefix: ", prefix)
			err = os.Remove(filepath.Join(dir, file.Name()))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Saves an image file to the "images/post_images" directory.
func SavePostImages(files []*multipart.FileHeader, postID int) ([]string, error) {
	var savedPaths []string
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			return nil, err
		}
		defer file.Close()

		extension := getFileExtensionByMimeType(fileHeader)
		filenamePrefix := fmt.Sprintf("post_pic_%d_", postID)
		savedPath, err := SaveImage(file, "post_images", filenamePrefix, extension)
		if err != nil {
			return nil, err
		}
		savedPaths = append(savedPaths, savedPath)
	}
	return savedPaths, nil
}

// Inserts the image path into the "images" table.
func SaveImagePath(postID int, imagePath string) (int64, error) {
	statement, err := Db.Prepare("INSERT INTO images (post_id, image_URL) VALUES (?, ?);")
	if err != nil {
		fmt.Println("Error preparing statement: ", err)
		return 0, err
	}
	defer statement.Close()

	res, err := statement.Exec(postID, imagePath)
	if err != nil {
		fmt.Println("Error executing statement: ", err)
		return 0, err
	}

	imageID, err := res.LastInsertId()
	if err != nil {
		fmt.Println("Error retrieving last insert id: ", err)
		return 0, err
	}

	return imageID, nil
}

// Saves an image file to the "images/comment_images" directory.
func SaveCommentImages(files []*multipart.FileHeader, commentID int) ([]string, error) {
	var savedPaths []string
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			return nil, err
		}
		defer file.Close()

		extension := getFileExtensionByMimeType(fileHeader)
		filenamePrefix := fmt.Sprintf("comment_pic_%d_", commentID)
		savedPath, err := SaveImage(file, "comment_images", filenamePrefix, extension)
		if err != nil {
			return nil, err
		}
		savedPaths = append(savedPaths, savedPath)
	}
	return savedPaths, nil
}

// Inserts the comment image path into the "images" table.
func SaveCommentImagePath(commentID int, imagePath string) (int64, error) {
	query := "INSERT INTO images (comment_id, image_URL) VALUES (?, ?);"

	statement, err := Db.Prepare(query)
	if err != nil {
		fmt.Println("Error preparing statement: ", err)
		return 0, err
	}
	defer statement.Close()

	res, err := statement.Exec(commentID, imagePath)
	if err != nil {
		fmt.Println("Error executing statement: ", err)
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		fmt.Println("Error getting last insert ID: ", err)
		return 0, err
	}

	return id, nil
}

// Function to prepend the imagesDir to the relative path.
func GetFullImagePath(relativePath string) string {
	const imagesDir = "images"
	return filepath.Join(imagesDir, relativePath)
}
