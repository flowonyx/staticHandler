package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"bitbucket.org/unclespike/fileserver"

	"github.com/jawher/mow.cli"
)

func main() {

	app := cli.App("File Server", "Serves files from a list of directories")

	app.Spec = "(-dp)..."

	dirs := app.StringsOpt("d dir", []string{}, "Directory to serve.")
	prefixes := app.StringsOpt("p prefix", []string{}, "Prefix under which to serve files from this directory.")

	app.Action = func() {
		if len(*dirs) != len(*prefixes) {
			fmt.Println("You need to specify a prefix for each directory")
			os.Exit(1)
		}
		for i, dir := range *dirs {
			prefix := (*prefixes)[i]
			if !strings.HasPrefix(prefix, "/") {
				prefix = "/" + prefix
			}
			if !strings.HasSuffix(prefix, "/") {
				prefix += "/"
			}
			fmt.Printf("Serving %s at %s\n", dir, prefix)
			static := staticHandler.NewFileOnlyHandler(dir, prefix)
			http.Handle(prefix, static)
		}
		// http.Handle("/", http.FileServer(http.Dir(".")))
		http.ListenAndServe(":9994", nil)
	}

	app.Run(os.Args)

}
