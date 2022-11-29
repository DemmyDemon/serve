package serve

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Server struct {
	Files      []string
	IndexCache []byte
	Allowed    []string
	Dir        string
	Port       int
}

func Begin(cwd string, port int, files []string, allow []string) error {
	server := Server{
		Files:   files,
		Dir:     cwd,
		Port:    port,
		Allowed: allow,
	}

	return http.ListenAndServe(fmt.Sprintf(":%d", port), server)
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	remote := strings.SplitN(r.RemoteAddr, ":", 2)[0]
	req := fmt.Sprintf("[%s] %s %s:", remote, r.Method, r.RequestURI)

	if !s.Allow(remote) {
		log.Println(req, "Remote not allowed")
		w.WriteHeader(http.StatusForbidden)
		w.Header().Add("Content-Type", "text/plain")
		w.Write([]byte(`Not allowed!`))
		return
	}

	if r.Method != "GET" {
		log.Println(req, `GET only, please.`)
		w.WriteHeader(http.StatusForbidden)
		w.Header().Add("Content-Type", "text/plain")
		w.Write([]byte(`GET only!`))
	}

	if r.RequestURI == "/" {
		log.Println(req, `Serving Index`)
		errList := s.Index(w)
		for _, err := range errList {
			log.Println(req, "Index error:", err)
		}
		return
	}

	fileName := strings.TrimPrefix(r.RequestURI, "/")
	for _, candidate := range s.Files {
		if candidate == fileName {
			log.Println(req, "Serving file")
			err := s.ServeFile(w, fileName)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Header().Add("Content-Type", "text/plain")
				w.Write([]byte(err.Error()))
			}
			return
		}
	}

	log.Println(req, "File not whitelisted")
	w.WriteHeader(http.StatusNotFound)
	w.Header().Add("Content-Type", "text/plain")
	w.Write([]byte("Nope"))
}

func (s Server) Allow(remote string) bool {
	for _, candidate := range s.Allowed {
		if remote == candidate {
			return true
		}
	}
	return false
}

func (s Server) ServeFile(w http.ResponseWriter, fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		log.Printf("Error opening %s: %s\n", fileName, err)
		return err
	}
	defer file.Close()

	w.WriteHeader(http.StatusOK)

	buffer := make([]byte, 1024)
	for {
		read, err := file.Read(buffer)
		if read > 0 {
			_, err := w.Write(buffer[:read])
			if err != nil {
				return err
			}
		}
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
}

func (s Server) Index(w http.ResponseWriter) []string {
	w.Header().Add("Content-Type", "text/html")

	index := bytes.Buffer{}
	index.WriteString(`<!DOCTYPE html><html><head><title>Serve</title><style>`)

	index.WriteString(` table { border: 1px solid black; width: 100%; border-collapse: collapse; } `)

	index.WriteString(` thead tr {background-color: #FFE; line-height: 1.5em; }`)
	index.WriteString(` thead th { border: 1px solid black; } `)

	index.WriteString(` tbody td { text-align: center; border: 1px solid black; padding: 0.25em; } `)
	index.WriteString(` tbody td:nth-child(1) { text-align: left; } `)
	index.WriteString(` tbody tr:nth-child(odd) { background-color: #EFE; } `)
	index.WriteString(` tbody tr:nth-child(even)  { background-color: #EEF; } `)

	index.WriteString(`</style></head><body>`)

	index.WriteString(`<table id="files"><thead><tr><th colspan="3">`)
	index.WriteString(s.Dir)
	index.WriteString(`</th></tr><tr><th>Name</th><th>Size (KiB)</th><th>Modified</th></tr></thead><tbody>`)
	errList := []string{}
	for _, fileName := range s.Files {

		info, err := os.Stat(fileName)
		if err != nil {
			errList = append(errList, err.Error())
			continue
		}

		index.WriteString(`<tr>`)

		index.WriteString(`<td><a href="/`)
		index.WriteString(fileName)
		index.WriteString(`">`)
		index.WriteString(fileName)
		index.WriteString(`</a></td>`)

		index.WriteString(`<td>`)
		index.WriteString(fmt.Sprintf("%0.2f", float32(info.Size())/1024.0))
		index.WriteString(`</td>`)

		index.WriteString(`<td>`)
		index.WriteString(info.ModTime().Format(`2006-01-02 15:04:05`))
		index.WriteString(`</td>`)

		index.WriteString(`</tr>`)
	}
	index.WriteString(`</tbody></table>`)

	index.WriteRune('\n')
	index.WriteString(`<!-- Generated `)
	index.WriteString(time.Now().Format(`2006-01-02 15:04:05`))
	index.WriteString(` -->`)
	index.WriteRune('\n')
	if len(errList) > 0 {
		index.WriteString("<!-- Errors encountered:\n")
		for _, err := range errList {
			index.WriteRune('\t')
			index.WriteString(err)
			index.WriteRune('\n')
		}
		index.WriteString("-->\n")
	}

	index.WriteString(`</body></html>`)

	w.Write(index.Bytes())

	return errList
}
