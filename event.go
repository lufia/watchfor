package main

import (
	"code.google.com/p/go.exp/fsnotify"
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"log"
	"time"
)

type Event struct {
	Path string     `json:"path"`
	Done chan error `json:"-"`
}

// 監視ディレクトリひとつを表す
type Entry struct {
	w          *fsnotify.Watcher
	dir        string
	clients    []chan *Event
	Subscribed chan chan *Event
	rule       *Rule
}

func NewEntry(dir string, rule *Rule) (*Entry, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	if err := w.Watch(dir); err != nil {
		return nil, err
	}
	return &Entry{
		w:          w,
		dir:        dir,
		clients:    make([]chan *Event, 0, 10),
		Subscribed: make(chan chan *Event),
		rule:       rule,
	}, nil
}

func (entry *Entry) eventLoop() {
	que := make(map[string]*fsnotify.FileEvent)
	for {
		select {
		case c := <-entry.Subscribed:
			entry.addClient(c)
		case ev := <-entry.w.Event:
			if ev.IsDelete() {
				continue
			}
			que[ev.Name] = ev
		case err := <-entry.w.Error:
			log.Fatal(err)
		case <-time.After(100 * time.Millisecond):
			for key, ev := range que {
				entry.notifyAll(ev.Name)
				delete(que, key)
			}
		}
	}
}

func (entry *Entry) addClient(nc chan *Event) {
	for i, c := range entry.clients {
		// websocketの通信でエラーになった場合は、
		// notifyAll()がnilにしているので再利用する。
		if c == nil {
			entry.clients[i] = nc
			return
		}
	}
	entry.clients = append(entry.clients, nc)
}

func (entry *Entry) notifyAll(file string) {
	target, err := entry.rule.Eval(file)
	if err == ErrNotCovered {
		// 監視対象外なので何もしない
		return
	}
	if err != nil {
		log.Printf("Eval(%v) = %v", file, err)
		return
	}
	for i, c := range entry.clients {
		if c == nil {
			continue
		}
		event := &Event{target, make(chan error)}
		c <- event
		if err := <-event.Done; err != nil {
			log.Print(err)
			entry.clients[i] = nil
		}
	}
}

func (entry *Entry) Serve(ws *websocket.Conn) {
	c := make(chan *Event, 1)
	entry.Subscribed <- c
	fout := json.NewEncoder(ws)
	for ev := range c {
		log.Print(ev)
		if err := fout.Encode(ev); err != nil {
			ev.Done <- err
			break
		}
		ev.Done <- nil
	}
}
