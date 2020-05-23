package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

const (
	PORT  = ":80"
	KEY   = "blit.key"
	CRT   = "blit.crt"
	STORE = "https://p.mort.coffee"
	HTML  = `
<html>
<head>
    <meta charset="utf-8">
    <title>blit</title>
<style>
body {
    line-height: 1.5;
}
</style>
</head>
<body>
</body>
<pre id="man">
BLIT(1)                          User Commands                           BLIT(1)

<b>NAME</b>
	blit - upload pastes

<b>SYNOPSIS</b>
	curl --upload-file [FILE] https://blit.io
	echo "Hello World" | curl --upload-file - https://blit.io

<b>DESCRIPTION</b>
	Upload data to the server, get a URL back. The file will be deleted
	a while after the last read.

	<b>GET /, GET /index.html</b>
		Get this index page.

	<b>GET</b> /<u>ID</u>
		Fetch the raw file with the given ID.

	<b>GET</b> /<u>ID.EXT</u>
		Fetch the file with the given ID, displayed in a way
		appropriate for that file type.

	<b>PUT</b> /<u>PATH</u>
		Upload a file. If PATH has a file extension, the
		returned URL will have the same file extension.
		The URL of the uploaded file is returned.

<b>AUTHOR</b>
	Written by <a href="https://github.com/jmbi">jmbi</a>.

<b>WWW</b>
	<a href="https://github.com/jmbi/blit.io">https://github.com/jmbi/blit.io</a>

<b>REPORTING BUGS</b>
	Submit bug reports to
	<a href="https://github.com/jmbi/blit.io/issues">https://github.com/jmbi/blit.io/issues</a>.

blit 0.0.0                         May 2020                              BLIT(1)
</pre>
</html>
`
)

func main() {
	res, err := http.Get(STORE + "/favicon.ico")
	if err != nil {
		panic(err)
	}
	ico, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	c := &http.Client{}
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		uri := req.RequestURI
		if req.Method == "GET" {
			if uri == "/" || uri == "/index.html" {
				io.WriteString(w, HTML)
			} else if uri == "/favicon.ico" {
				w.Write(ico)
			} else {
				res, err := http.Get(STORE + "/" + uri[1:])
				if err != nil {
					fmt.Errorf("%v\n", err)
					return
				}
				dat, err := ioutil.ReadAll(res.Body)
				if err != nil {
					fmt.Errorf("%v\n", err)
					return
				}
				w.Header().Set("Content-Type", res.Header.Get("Content-Type"))
				w.Write(dat)
			}
		} else if req.Method == "PUT" {
			put, err := http.NewRequest(http.MethodPut, STORE+uri, req.Body)
			if err != nil {
				fmt.Errorf("%v\n", err)
				return
			}
			res, err := c.Do(put)
			if err != nil {
				fmt.Errorf("%v\n", err)
				return
			}
			out, err := ioutil.ReadAll(res.Body)
			if err != nil {
				fmt.Errorf("%v\n", err)
				return
			}
			uid := out[9+bytes.Index(out[8:], []byte("/")):]
			io.WriteString(w, fmt.Sprintf("https://blit.io/%s\n", uid))
		} else {
			io.WriteString(w, "404")
		}
	})
	fmt.Printf("listening on %s\n", PORT)
	panic(http.ListenAndServeTLS(PORT, CRT, KEY, nil))
}
