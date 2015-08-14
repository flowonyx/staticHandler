# staticHandler
A layer on top of the stdlib `http.ServeFile` that eliminates directory browsing. Also includes a command line server for quickly serving files from multiple directories under chosen prefixes.

When you want to serve files from a directory without an index.html file and you don't want the directory to be browsable, this allows you to easily do that.

Just `go get github.com/flowonyx/staticHandler`.

For an example of how to use this library, look at `cmd/staticsite/main.go`.

Basically, `staticHandler.NewFileOnlyHandler(root, stripPrefix)` returns an http.Handler. The `stripPrefix` parameter is the string prefix you want to strip off when serving requests. You can leave it as an empty string if you don't want to strip any prefix.

Also of note, is that it will look for a custom error page when there is an error ({error code}.html). It also allows setting the HTML for each error code through the `SetErrorPage` function. You can also set the `ErrorPageHandler` to a custom function for serving error pages.
