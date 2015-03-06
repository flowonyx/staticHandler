package staticHandler

import (
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fogcreek/logging"
)

var staticHandlerTags = []string{"noDowntime", "staticHandler.go"}

func hasErrorPage(root string, code int) bool {
	logging.VerboseWithTagsf(staticHandlerTags, "hasErrorPage(%s, %d)", root, code)

	errCode := strconv.Itoa(code)
	if _, err := os.Stat(filepath.Join(root, errCode+".html")); !os.IsNotExist(err) {
		return true
	}
	return false
}

func handleErrorPage(w http.ResponseWriter, r *http.Request, root string, code int) {
	logging.VerboseWithTagsf(staticHandlerTags, "handleErrorPage(%#v, %#v, %s, %d)", w, r, root, code)

	w.WriteHeader(code)
	sCode := strconv.Itoa(code)
	fullpath := filepath.Join(root, sCode+".html")
	if hasErrorPage(root, code) {
		f, err := os.Open(fullpath)
		if err != nil {
			logging.Errorf("Error opening %s: %v", fullpath, err)
			return
		}
		defer f.Close()
		io.Copy(w, f)
		return
	}

	w.Write([]byte(`
<!doctype html>
<html>
<head>
	<title>` + sCode + `: ` + http.StatusText(code) + `</title>
</head>
<body style="text-align:center">
<h1>` + sCode + `: ` + http.StatusText(code) + `</h1>
</body>
</html>`),
	)

}

// ServeFileOnly replies to the request with the contents of the named file but
// returns a 404 for a directory.
func ServeFileOnly(w http.ResponseWriter, r *http.Request, root string, name string) {

	fullpath := filepath.Join(root, name)

	logging.VerboseWithTagsf(staticHandlerTags, "ServeFileOnly(%s)", fullpath)

	info, err := os.Stat(fullpath)
	if err != nil {
		logging.DebugWithTagsf(staticHandlerTags, "Error from os.Stat(%s): %v", fullpath, err)
		handleErrorPage(w, r, root, http.StatusNotFound)
		return
	}

	if info.IsDir() {
		if _, err := os.Stat(filepath.Join(fullpath, "index.html")); os.IsNotExist(err) {
			logging.VerboseWithTagsf(staticHandlerTags, "index.html not found for this directory (%s): %v", fullpath, err)
			handleErrorPage(w, r, root, http.StatusNotFound)
			return
		}
	}

	logging.VerboseWithTagsf(staticHandlerTags, "Serving file:", fullpath)

	http.ServeFile(w, r, fullpath)
}

// FileOnlyHandler provides a static server for serving only files and not
// directory listings.
type FileOnlyHandler struct {
	Root        string
	StripPrefix string
}

func (fh *FileOnlyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upath := r.URL.Path
	if len(fh.StripPrefix) > 0 {
		upath = strings.TrimPrefix(upath, fh.StripPrefix)
	}
	if !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
		r.URL.Path = upath
	}
	ServeFileOnly(w, r, fh.Root, path.Clean(upath))
}

// NewFileOnlyHandler returns a static file handler that will only serve
// files for the given Root rather than directory listings.
func NewFileOnlyHandler(root string, stripPrefix string) *FileOnlyHandler {
	logging.VerboseWithTagsf(staticHandlerTags, "Setting FileOnlyHandler.Root to:", root)

	return &FileOnlyHandler{root, stripPrefix}
}
