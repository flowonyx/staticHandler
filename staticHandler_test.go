package staticHandler

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServeFileWithFile(t *testing.T) {

	s := httptest.NewServer(NewFileOnlyHandler(".", ""))
	defer s.Close()

	res, err := http.Get(s.URL + "/staticHandler_test.go")
	if err != nil {
		t.Error(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Contains(body, []byte("TestServeFileWithFile")) {
		t.Error("Did not return correct file contents:", string(body))
	}

}

func TestServeFileWithDirectory(t *testing.T) {

	s := httptest.NewServer(NewFileOnlyHandler(".", ""))
	defer s.Close()

	res, err := http.Get(s.URL + "/")
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 404 {
		t.Error("Should have recieved 404!")
	}

}

func TestSetErrorPage(t *testing.T) {

	s := httptest.NewServer(NewFileOnlyHandler(".", ""))
	defer s.Close()

	SetErrorPage(404, "Testing 404!!!")

	res, err := http.Get(s.URL + "/")
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 404 {
		t.Error("Should have recieved 404!")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}
	if string(body) != "Testing 404!!!" {
		t.Error("Did not return correct file contents:", string(body))
	}

}

func TestErrorPageHandler(t *testing.T) {

	s := httptest.NewServer(NewFileOnlyHandler(".", ""))
	defer s.Close()

	ErrorPageHandler = func(w http.ResponseWriter, r *http.Request, code int) {
		w.WriteHeader(code)
		w.Write([]byte("Test custom error page handler."))
	}

	res, err := http.Get(s.URL + "/")
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 404 {
		t.Error("Should have recieved 404!")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}
	if string(body) != "Test custom error page handler." {
		t.Error("Did not return correct file contents:", string(body))
	}

}
