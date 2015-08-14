package staticHandler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

// FileOnlyHandler provides a static server for serving only files and not
// directory listings.
type FileOnlyHandler struct {
	Root        string
	StripPrefix string
}

// Keeps the HTML for each error response, if the user sets it.
var errorPages = make(map[int]string)

// SetErrorPage allows for directly setting the HTML for each error code.
func SetErrorPage(code int, html string) {
	errorPages[code] = html
}

func hasErrorPage(root string, code int) bool {

	errCode := strconv.Itoa(code)
	if _, err := os.Stat(filepath.Join(root, errCode+".html")); !os.IsNotExist(err) {
		return true
	}
	return false
}

func handleErrorPage(w http.ResponseWriter, r *http.Request, root string, code int) {

	w.WriteHeader(code)

	// If the HTML has been provide for this error, render it and return.
	if errorPages != nil {
		if html, ok := errorPages[code]; ok {
			w.Write([]byte(html))
			return
		}
	}

	// If the HTML has not been provided, look for an HTML page for this error.
	sCode := strconv.Itoa(code)
	fullpath := filepath.Join(root, sCode+".html")
	if hasErrorPage(root, code) {
		f, err := os.Open(fullpath)
		if err != nil {
			fmt.Errorf("Error opening %s: %v", fullpath, err)
			return
		}
		defer f.Close()
		io.Copy(w, f)
		return
	}

	// If no HTML or HTML page is available for this error, use this default HTML.
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
// returns a 404 for a directory unless an index.html file is found.
func ServeFileOnly(w http.ResponseWriter, r *http.Request, root string, name string) {

	fullpath := filepath.Join(root, name)

	info, err := os.Stat(fullpath)
	if err != nil {
		fmt.Errorf("Error from os.Stat(%s): %v", fullpath, err)
		handleErrorPage(w, r, root, http.StatusNotFound)
		return
	}

	if info.IsDir() {
		if _, err := os.Stat(filepath.Join(fullpath, "index.html")); os.IsNotExist(err) {
			handleErrorPage(w, r, root, http.StatusNotFound)
			return
		}
	}

	http.ServeFile(w, r, fullpath)
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
	return &FileOnlyHandler{root, stripPrefix}
}
