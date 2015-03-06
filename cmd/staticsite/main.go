/*
staticsite -d {directory path} -p {prefix} -d {another directory} -p {its prefix}

Allows quick static site or file serving sites.
*/
package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/flowonxy/staticHandler"
	"github.com/jawher/mow.cli"
)

func main() {

	app := cli.App("File Server", "Serves files from a list of directories")

	app.Spec = "[--port] (-dp)..."

	port := app.IntOpt("port", 8888, "Port under which to run")
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
		colonPort := fmt.Sprintf(":%d", *port)
		fmt.Println("Listening on port", *port)
		http.ListenAndServe(colonPort, nil)
	}

	app.Run(os.Args)

}
