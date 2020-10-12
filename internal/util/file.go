// Package util exposes utility functions
package util

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// known extensions
var textExt = map[string]bool{
	".txt":  true,
	".json": true,
	".yaml": true,
}

// IsTextFile checks if the file located at the provided path contains text.
func IsTextFile(fp string) (bool, error) {
	log.Debug().
		Str("filepath", fp).
		Msg("Checking if file is text.")

	// if the extension is known, we're good.
	if isText, ok := textExt[path.Ext(fp)]; ok {
		return isText, nil
	}

	log.Debug().
		Str("filepath", fp).
		Msg("Unknown file extension. Verifying content.")

	absFp, _ := filepath.Abs(fp)

	// the extension is not known; read an initial chunk of the file and check if it looks like text
	f, err := os.Open(absFp)
	if err != nil {
		return false, errors.Wrap(err, "Unable to open file")
	}

	var buf [1024]byte
	n, err := f.Read(buf[0:])
	if err != nil {
		return false, errors.Wrap(err, "Unable to read file chunk to determine if its text")
	}

	return strings.Contains(http.DetectContentType(buf[0:n]), "text/"), f.Close()
}
