package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/demmydemon/serve/serve"
	"golang.org/x/exp/slices"
)

type allowedHosts []string

func (a *allowedHosts) Set(value string) error {
	*a = append(*a, value)
	return nil
}

func (a *allowedHosts) String() string {
	return strings.Join(*a, ",")
}

func parseCommandline() (int, []string, allowedHosts) {
	port := 80
	files := []string{}
	allowed := allowedHosts{}

	flag.IntVar(&port, "port", 8181, "Specify a port to listen on")
	flag.Var(&allowed, "allow", "What host is allowed to access the files")

	flag.Parse()

	for _, name := range flag.Args() {
		files = append(files, filepath.Base(filepath.Clean(name)))
	}
	return port, files, allowed
}

func filesIn(dir string) []string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	files := []string{}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		files = append(files, entry.Name())
	}
	slices.Sort(files)
	return files
}

func main() {
	port, files, allow := parseCommandline()
	cwd, err := os.Getwd()

	log.SetFlags(log.Ltime)

	if err != nil {
		log.Fatal(err)
	}
	if len(files) == 0 {
		log.Printf("No files specified, using all files in %s\n", cwd)
		files = filesIn(cwd)
	}
	if len(allow) == 0 {
		log.Fatal("You did not -allow anyone to access.")
	}
	log.Printf("Allowed: %s", allow.String())
	log.Fatal(serve.Begin(cwd, port, files, allow))
}
