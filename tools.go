package toolkit

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const randomStringSource = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_+"

// Tools is the type used to instantiate this module.Any variable of this type will have ccess
// to all the methods with the reciver *Tools
type Tools struct {
	MaxFileSize  int
	allowedTypes []string
}

func (t *Tools) RandomString(n int) string {
	s, r := make([]rune, n), []rune(randomStringSource)
	for i := range s {
		p, _ := rand.Prime(rand.Reader, len(r))
		x, y := p.Uint64(), uint64(len(r))
		s[i] = r[x%y]
	}
	return string(s)
}

type UploadFile struct {
	NewFileName      string
	OriginalFileName string
	FileSize         int64
}

func (t *Tools) UploadFile(r *http.Request, uploadDir string, rename ...bool) ([]*UploadFile, error) {
	renameFile := true
	if len(rename) > 0 {
		renameFile = rename[0]
	}
	var uploadFiles []*UploadFile
	if t.MaxFileSize == 0 {
		t.MaxFileSize = 1024 * 1024 * 1024
	}
	err := r.ParseMultipartForm(int64(t.MaxFileSize))
	if err != nil {
		return nil, errors.New("the upload file  is too big")
	}
	for _, fHeadeers := range r.MultipartForm.File {
		for _, hdr := range fHeadeers {
			uploadFiles, err = func(uploadFiles []*UploadFile) ([]*UploadFile, error) {
				var uploadFile UploadFile
				infile, err := hdr.Open()
				if err != nil {
					return nil, err
				}
				defer infile.Close()

				buff := make([]byte, 512)
				_, err = infile.Read(buff)
				if err != nil {
					return nil, err
				}
				//check to see if the file type is permitted
				allowed := false
				fileType := http.DetectContentType(buff)
				// allowedTypes := []string("image/jpeg","image/png","image/gif")
				if len(t.allowedTypes) > 0 {
					for _, x := range t.allowedTypes {
						if strings.EqualFold(fileType, x) {
							allowed = true
						}
					}

				} else {
					allowed = true
				}
				if !allowed {
					return nil, errors.New("the upload file type is not permitted")
				}
				_, err = infile.Seek(0, 0)
				if err != nil {
					return nil, err
				}
				if renameFile {
					uploadFile.NewFileName = fmt.Sprintf("%s%s", t.RandomString(25), filepath.Ext(hdr.Filename))
				} else {
					uploadFile.NewFileName = hdr.Filename
				}
				var outfile *os.File
				defer outfile.Close()
				if outfile, err = os.Create(filepath.Join(uploadDir, uploadFile.NewFileName)); err != nil {
					return nil, err
				} else {
					fileSize, err := io.Copy(outfile, infile)
					if err != nil {
						return nil, err
					}
					uploadFile.FileSize = fileSize
					uploadFiles = append(uploadFiles, &uploadFile)
					return uploadFiles, nil
				}
			}(uploadFiles)
			if err != nil {
				return uploadFiles, err
			}
		}
	}
	return uploadFiles, nil
}
