package main

import (
	"html/template"
	"io"
	"net/http"
)

type PageContent struct {
	t    *template.Template
	file string
}

func NewPage(file string) *PageContent {
	t := template.Must(template.New(file).Parse(pageTemplate))
	return &PageContent{t, file}
}

func (c PageContent) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	c.t.Execute(w, c.file)
}

type ScriptContent string

func (c ScriptContent) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/javascript")
	io.WriteString(w, string(c))
}

var pageTemplate = `<!doctype html>
<html>
<head>
<meta charset="utf-8">
<script src="/script"></script>
<title>test</title>
</head>
<body>

<img id="view" src="/files/{{ . }}">

</body>
</html>
`

var scriptContent = ScriptContent(`
var ws = new WebSocket("ws://localhost:8888/event", ["event"])
var task = null
ws.onopen = function(){
	ws.onmessage = function(message){
		if(task != null)
			clearTimeout(task)
		task = setTimeout(function(){
			var xhr = new XMLHttpRequest()
			xhr.open('GET', '/files/' + message.data, true)
			xhr.responseType = 'blob'
			xhr.onload = function(e){
				if(this.status != 200)
					return
				document.querySelector('#view').src = window.URL.createObjectURL(this.response)
			}
			xhr.send()
		}, 500)
	}
}
`)
