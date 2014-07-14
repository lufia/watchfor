package main

import (
	"code.google.com/p/go.net/websocket"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

func usage() {
	fmt.Fprintf(os.Stderr, `Usage: watchfor [options] [dir]

options:
	-a addr (default :8080)
		http address
	-s ext
		source file extension. (e.g. .dot)
	-t ext
		target file extension. (e.g. .png)
	-c command
		command for convert source to target
`)
	os.Exit(1)
}

func main() {
	addr := flag.String("a", ":8080", "http address")
	cmd := flag.String("c", "", "cook command")
	src := flag.String("s", "", "source extension")
	dest := flag.String("t", "", "target extension")
	flag.Usage = usage
	flag.Parse()
	rule := &Rule{
		Cmd:     Command(*cmd),
		SrcExt:  FileExt(*src),
		DestExt: FileExt(*dest),
	}

	args := flag.Args()
	if len(args) == 0 {
		args = []string{"."}
	}
	entry, err := NewEntry(args[0], rule)
	if err != nil {
		log.Fatal(err)
	}
	go entry.eventLoop()

	dir := args[0]
	page := NewPage(dir)
	http.Handle("/", AvoidCache(http.Handler(page)))
	http.Handle("/script", AvoidCache(http.Handler(scriptContent)))
	http.Handle("/event", websocket.Handler(entry.Serve))
	f := AvoidCache(http.FileServer(http.Dir(".")))
	http.Handle("/files/", http.StripPrefix("/files/", f))
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatalf("ListenAndServe: %v", err)
	}
}
