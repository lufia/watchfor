package main

import (
	"io"
	"net/http"
)

type PageContent string

func (c PageContent) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	io.WriteString(w, string(c))
}

type ScriptContent string

func (c ScriptContent) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/javascript")
	io.WriteString(w, string(c))
}

var pageContent = PageContent(`<!doctype html>
<html>
<head>
<meta charset="utf-8">
<script src="/script"></script>
<title>test</title>
</head>
<body>

<img id="view" src="/files/test.jpg">

</body>
</html>
`)

var scriptContent = ScriptContent(`
console.log("new")
var ws = new WebSocket("ws://localhost:8888/event", ["event"])
ws.onopen = function(){
	ws.onmessage = function(message){
		var xhr = new XMLHttpRequest()
		xhr.open('GET', '/files/' + message.data, true)
		xhr.responseType = 'blob'
		xhr.onload = function(e){
			if(this.status != 200)
				return
			document.querySelector('#view').src = window.URL.createObjectURL(this.response)
		}
		xhr.send()
	}
}
`)
