package main

import (
	"crypto"
	"encoding/hex"
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pkg/browser"
)

var version = "0.2.0"
var changePath = "/__changes_d49689c9-665a-4d1b-9e50-05c7b8f764d7" // Arbitrary path that is very unlikely to be used by a real file.

func main() {
	address := flag.String("address", "localhost:3000", "Address of the server")
	folder := flag.String("folder", ".", "Folder that should be served")
	autoRefresh := flag.Bool("auto-refresh", true, "Inject JavaScript into HTML pages that automatically reloads when a file in the specified folder changes.")
	open := flag.Bool("open", true, "Opens the server address in the default browser on startup")
	flag.Parse()

	handler := http.FileServer(http.Dir(*folder))

	if *autoRefresh {
		http.HandleFunc(changePath, func(w http.ResponseWriter, r *http.Request) {
			hash := crypto.SHA1.New()
			err := filepath.Walk(*folder, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				hash.Write([]byte(path))
				hash.Write([]byte("#"))
				hash.Write([]byte(info.ModTime().String()))
				hash.Write([]byte("\n"))
				return err
			})
			if err != nil {
				log.Panic(err)
			}

			w.Header().Add("Cache-Control", "no-store, must-revalidate")

			w.Write([]byte(hex.EncodeToString(hash.Sum(nil))))
		})

		handler = JsInjectionInterceptor{handler}
	}

	http.Handle("/", handler)

	println("\x1b[34mserve " + version + "\x1b[0m")

	if *open {
		var url = "http://" + *address
		log.Print("Opening \"" + url + "\"...")
		browser.OpenURL(url)
	}

	log.Print("Listening on " + *address + " to serve \"" + *folder + "\" with auto-refresh \"" + strconv.FormatBool(*autoRefresh) + "\"...")
	err := http.ListenAndServe(*address, nil)
	if err != nil {
		log.Panic(err)
	}
}

type JsInjectionInterceptor struct {
	delegate http.Handler
}

func (i JsInjectionInterceptor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	inMemoryResponseWriter := &InMemoryResponseWriter{
		header:     http.Header{},
		data:       []byte{},
		statusCode: 200,
	}

	// Write response from the delegate to the in-memory response writer:
	i.delegate.ServeHTTP(inMemoryResponseWriter, r)

	// Copy headers:
	for k, v := range inMemoryResponseWriter.header {
		for _, vi := range v {
			w.Header().Add(k, vi)
		}
	}

	w.Header().Del("Content-Length") // Remove Content-Length, because it'd be wrong after injection.

	w.Header().Add("Expires", "0")

	w.WriteHeader(inMemoryResponseWriter.statusCode)

	if strings.HasPrefix(inMemoryResponseWriter.Header().Get("Content-Type"), "text/html") {
		js := []byte(`<script>
		async function getContent() {
			return await (await fetch("` + changePath + `")).text()
		}

		let lastContent = null;

		async function checkChanges() {
			const newContent = await getContent()
			if (lastContent == null) {
				lastContent = newContent
			}
			if (newContent != lastContent) {
				location.reload()
			}
		}

		setInterval(() => { checkChanges() }, 1000)
		</script>`)

		inMemoryResponseWriter.data = append(inMemoryResponseWriter.data, js...)
	}

	w.Write(inMemoryResponseWriter.data)
}

type InMemoryResponseWriter struct {
	header     http.Header
	data       []byte
	statusCode int
}

func (i *InMemoryResponseWriter) Header() http.Header {
	return i.header
}

func (i *InMemoryResponseWriter) Write(data []byte) (int, error) {
	i.data = append(i.data, data...)
	return len(data), nil
}

func (i *InMemoryResponseWriter) WriteHeader(statusCode int) {
	i.statusCode = statusCode
}
