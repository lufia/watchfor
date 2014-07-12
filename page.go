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
	c.t.Execute(w, req.URL.Path[1:])
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

</body>
</html>
`

var scriptContent = ScriptContent(`
function getTargetPath() {
	return location.pathname.substr(1)
}
function getEventURL() {
	return "ws://" + location.host + "/event"
}

var Resource = function(path){
	this.path = path
}
Resource.prototype = {
	request: function(meth, type, onsuccess) {
		var xhr = new XMLHttpRequest()
		xhr.open(meth, '/files/' + this.path, true)
		if(type != '')
			xhr.responseType = type
		xhr.onload = function(e){
			if(this.status != 200)
				return
			ctype = this.getResponseHeader('Content-Type')
			onsuccess(ctype, this.response)
		}
		xhr.send()
	}
}

var ResourceController = function(r){
	this.r = r
	this.v = null
	this.ws = null
}
ResourceController.prototype = {
	bind: function(onbind){
		var c = this
		r.request('HEAD', '', function(type, data){
			c.v = onbind(type)
			c.v.createView()
			c.ws = new WebSocket(getEventURL(), ["event"])
			c.ws.onopen = function(){
				c.ws.onmessage = function(message){
					var m = JSON.parse(message.data)
					if(m.path != c.r.path)
						return
					c.refresh()
				}
			}
			c.refresh()
		})
	},
	refresh: function(){
		var c = this
		var type = this.v.getResponseType()
		this.r.request('GET', type, function(type, data){
			c.v.refreshView(data)
		})
	}
}

var DocumentView = function(){
}
DocumentView.prototype = {
	getResponseType: function(){
		return 'document'
	},
	createView: function(){
		var div = document.createElement('div')
		div.id = 'view'
		document.body.appendChild(div)
	},
	refreshView: function(data){
		var v = document.querySelector('#view')
		v.innerHTML = data.body.innerHTML
	}
}

var ImageView = function(){
}
ImageView.prototype = {
	getResponseType: function(){
		return 'blob'
	},
	createView: function(){
		var img = document.createElement('img')
		img.id = 'view'
		document.body.appendChild(img)
	},
	refreshView: function(data){
		var v = document.querySelector('#view')
		v.src = URL.createObjectURL(data)
	}
}

var TextView = function(){
}
TextView.prototype = {
	getResponseType: function(){
		return ''
	},
	createView: function(){
		var plain = document.createElement('pre')
		plain.id = 'view'
		document.body.appendChild(plain)
	},
	refreshView: function(data){
		var v = document.querySelector('#view')
		v.textContent = data
	}
}

window.onload = function(){
	path = getTargetPath()
	r = new Resource(path)
	ctlr = new ResourceController(r)
	ctlr.bind(function(type){
		if(/^text\/x?html/.test(type))
			return new DocumentView()
		else if(/^image\//.test(type))
			return new ImageView()
		else
			return new TextView()
	})
}
`)
