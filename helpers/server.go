package helpers

import (
	"fmt"
	"github.com/spf13/afero"
	"net/http"
	"strings"
)

func FileServer(fs afero.Fs) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		var err error
		requestedFilename := strings.TrimPrefix(req.URL.Path, "/")
		fileData, err := afero.ReadFile(fs, requestedFilename)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			res.Write([]byte(fmt.Sprintf("Could not load file %s", requestedFilename)))
		}
		res.Write(fileData)
	})
}
