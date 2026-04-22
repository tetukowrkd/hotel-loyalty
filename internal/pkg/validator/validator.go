package validator

import (
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"hotel-loyalty/internal/pkg/errors"

	v "github.com/go-playground/validator/v10"
)

var validate = v.New()

func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

// 🔥 helper untuk extract error
func ParseValidationError(err error) map[string]string {
	errors := make(map[string]string)

	if errs, ok := err.(v.ValidationErrors); ok {
		for _, e := range errs {
			switch e.Field() {
			case "Name":
				errors["name"] = "name is required"
			case "Email":
				errors["email"] = "invalid email format"
			case "Password":
				errors["password"] = "password must be at least 6 characters"
			}
		}
	}

	return errors
}

func ValidateImage(file multipart.File, filename string, size int64) error {

	// size (2MB)
	if size > 2<<20 {
		return errors.NewBadRequest("File too large (max 2MB)")
	}

	// extension
	ext := strings.ToLower(filepath.Ext(filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return errors.NewBadRequest("Invalid file extension")
	}

	// MIME check
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil {
		return errors.ErrInternal
	}

	contentType := http.DetectContentType(buffer)

	if contentType != "image/jpeg" && contentType != "image/png" {
		return errors.NewBadRequest("Invalid image format")
	}

	// reset file pointer
	file.Seek(0, io.SeekStart)

	return nil
}
