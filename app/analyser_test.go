package app

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestHandleAnalyseWebPageGeneralSuccess tests the end-point to return success.
func TestHandleAnalyseWebPageGeneralSuccess(t *testing.T) {
	r, _ := http.NewRequest("POST", "/analyse", bytes.NewBuffer([]byte(`{"url":"https://www.google.com"}`)))
	w := httptest.NewRecorder()

	HandleAnalyseWebPage(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
}

// TestHandleAnalyseWebPageNotSuccess tests the end-point when the url returns a non 200 response.
func TestHandleAnalyseWebPageNotSuccess(t *testing.T) {
	r, _ := http.NewRequest("POST", "/analyse", bytes.NewBuffer([]byte(`{"url":"https://www.buymeacoffee.com/freeformatter"}`)))
	w := httptest.NewRecorder()

	HandleAnalyseWebPage(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// TestHandleAnalyseWebPageInvalidParams tests the end-point when the paremeters can not be decoded.
func TestHandleAnalyseWebPageInvalidParams(t *testing.T) {
	r, _ := http.NewRequest("POST", "/analyse", bytes.NewBuffer([]byte(`{"url":https://www.google.com}`)))
	w := httptest.NewRecorder()

	HandleAnalyseWebPage(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestHandleAnalyseWebPageNoURLParams tests the end-point when the url request parameter is not provided.
func TestHandleAnalyseWebPageNoURLParams(t *testing.T) {
	r, _ := http.NewRequest("POST", "/analyse", bytes.NewBuffer([]byte(`{"not_url":"not a url parameter"}`)))
	w := httptest.NewRecorder()

	HandleAnalyseWebPage(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestHandleAnalyseWebPageInvalidURLValue tests the end-point when the url rvalue provided is invalid and has no protocol specified.
func TestHandleAnalyseWebPageInvalidURLValue(t *testing.T) {
	r, _ := http.NewRequest("POST", "/analyse", bytes.NewBuffer([]byte(`{"url":"invalid url sent"}`)))
	w := httptest.NewRecorder()

	HandleAnalyseWebPage(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetHTMLVersion(t *testing.T) {
	htmlVersion := getHTMLVersion(`
		<!DOCTYPE html>
		<html>
		<head>
		<title>The document title</title>
		</head>
		<body>
		<h1>Main heading</h1>
		<p>A paragraph.</p>
		</body>
		</html>`)

	assert.Equal(t, "HTML5", htmlVersion)
}

func TestGetHTMLVersionXHTML1(t *testing.T) {
	htmlVersion := getHTMLVersion(`
	<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN"
	"http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd">
	<html>
	<head>
	<title>The document title</title>
	</head>
	<body>
	<h1>Main heading</h1>
	<p>A paragraph.</p>
	</body>
	</html>`)

	assert.Equal(t, "XHTML", htmlVersion)
}

func TestGetHTMLVersionHTML4(t *testing.T) {
	htmlVersion := getHTMLVersion(`<!DOCTYPE html PUBLIC "-//W3C//DTD HTML 4.01//EN"
	"http://www.w3.org/TR/html4/strict.dtd">
	<html>
	<head>
	<title>The document title</title>
	</head>
	<body>
	<h1>Main heading</h1>
	<p>A paragraph.</p>
	</body>
	</html>`)

	assert.Equal(t, "HTML 4.01", htmlVersion)
}

func TestGetHTMLVersionHTML3(t *testing.T) {
	htmlVersion := getHTMLVersion(`<!DOCTYPE html PUBLIC "-//W3C//DTD HTML 3.2 Final//EN">
	<html>
	<head>
	<title>The document title</title>
	</head>
	<body>
	<h1>Main heading</h1>
	<p>A paragraph.</p>
	</body>
	</html>`)

	assert.Equal(t, "HTML 3.2", htmlVersion)
}

func TestGetHTMLVersionHTML2(t *testing.T) {
	htmlVersion := getHTMLVersion(`<!DOCTYPE html PUBLIC "-//IETF//DTD HTML 2.0//EN">
	<html>
	<head>
	<title>The document title</title>
	</head>
	<body>
	<h1>Main heading</h1>
	<p>A paragraph.</p>
	</body>
	</html>`)

	assert.Equal(t, "HTML 2.0", htmlVersion)
}

func TestGetHTMLVersionHTML1(t *testing.T) {
	htmlVersion := getHTMLVersion(`
	<html>
	<head>
	<title>The document title</title>
	</head>
	<body>
	<h1>Main heading</h1>
	<p>A paragraph.</p>
	</body>
	</html>`)

	assert.Equal(t, "HTML 1.0", htmlVersion)
}
