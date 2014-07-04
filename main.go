package main

import (
	"code.google.com/p/go.exp/fsnotify"
	"code.google.com/p/go.net/websocket"
	"flag"
	"log"
	"net/http"
	"os"
)

type Event struct {
	Message string
	Done    chan error
}

type Entry struct {
	w          *fsnotify.Watcher
	file       string
	clients    []chan *Event
	Subscribed chan chan *Event
}

func NewEntry(w *fsnotify.Watcher, file string) *Entry {
	return &Entry{
		w:          w,
		file:       file,
		clients:    make([]chan *Event, 0, 10),
		Subscribed: make(chan chan *Event),
	}
}

func (entry *Entry) WatchStart() {
	if err := entry.w.Watch(entry.file); err != nil {
		log.Fatal(err)
	}
	for {
		select {
		case c, ok := <-entry.Subscribed:
			if !ok {
				return
			}
			entry.addClient(c)
		case ev := <-entry.w.Event:
			if ev.IsModify() {
				entry.notifyAll(ev.Name)
			}
		case err := <-entry.w.Error:
			log.Fatal(err)
		}
	}
}

func (entry *Entry) addClient(nc chan *Event) {
	for i, c := range entry.clients {
		if c == nil {
			entry.clients[i] = nc
			return
		}
	}
	entry.clients = append(entry.clients, nc)
}

func (entry *Entry) notifyAll(file string) {
	for i, c := range entry.clients {
		if c == nil {
			continue
		}
		e := &Event{file, make(chan error)}
		c <- e
		if err := <-e.Done; err != nil {
			log.Print(err)
			entry.clients[i] = nil
		}
	}
}

func (entry *Entry) WatchStop() {
	close(entry.Subscribed)
	entry.w.Close()
}

func (entry *Entry) Serve(ws *websocket.Conn) {
	c := make(chan *Event, 1)
	entry.Subscribed <- c
	for ev := range c {
		_, err := ws.Write([]byte(ev.Message))
		if err != nil {
			ev.Done <- err
			break
		}
		ev.Done <- nil
	}
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	w, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	entry := NewEntry(w, args[0])
	go entry.WatchStart()
	defer entry.WatchStop()

	page := NewPage(args[0])
	http.Handle("/page", http.Handler(page))
	http.Handle("/script", http.Handler(scriptContent))
	http.Handle("/event", websocket.Handler(entry.Serve))
	http.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir("."))))
	if err := http.ListenAndServe(":8888", nil); err != nil {
		log.Fatalf("ListenAndServe: %v", err)
	}
}
