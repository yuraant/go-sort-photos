package main

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/url"
	"path"

	"github.com/zserge/webview"
)

var indexHTML = `
<!doctype html>
<html>
	<head>
		<meta http-equiv="X-UA-Compatible" content="IE=edge">
	</head>
	<body>
	<div id="open" style="position:absolute;top:50%;left:50%;transform:translate(-50%, -50%);">
		<button onclick="external.invoke('opendir')">Open directory</button>
	</div>
	<div id="test">
	 <ul id="list">
	 </ul>
	</div>
	<div id="rename" style="position:absolute;left:50%;transform:translate(-50%, -50%);">
	</div>	
	</body>
</html>`

func handleRPC(w webview.WebView, data string) {
	switch {
	case data == "rename":
		log.Println("open", w.Dialog(webview.DialogTypeAlert, webview.DialogFlagInfo, "Rename", ""))
	case data == "opendir":
		openningDir := w.Dialog(webview.DialogTypeOpen, webview.DialogFlagDirectory, "Open directory", "")
		pathes := readDir(openningDir)
		w.InjectCSS(string(MustAsset("static/css/styles.css")))
		w.Eval(string(MustAsset("static/js/Sortable.js")))
		w.Eval(`document.getElementById("open").remove();`)
		w.Eval(`function putIMG(path, id) {var ul = document.getElementById("list");var img = document.createElement('img');img.src = path;var li = document.createElement("li"); li.setAttribute('data-id', id); li.appendChild(img);ul.appendChild(li);}`)
		var i int
		for _, v := range pathes {
			i++
			b, err := ioutil.ReadFile(v)
			if err != nil {
				log.Fatalln(err)
			}
			w.Eval(fmt.Sprintf(`putIMG("data:image/jpeg;base64, %s", %d)`, template.JSEscapeString(base64.StdEncoding.EncodeToString(b)), i))
		}
		w.Eval(`var sort = document.getElementById('list'); var sortable = Sortable.create(sort, {group: "list"});`)
		w.Eval(`var el = document.getElementById('rename'); el.innerHTML = "<button onclick=\"external.invoke('rename')\">Rename</button>"`)
	}
}

func readDir(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	var data []string
	for _, f := range files {
		if path.Ext(f.Name()) == ".jpg" {

			pathF := dir + "/" + f.Name()
			data = append(data, pathF)
		}

	}
	return data
}

func main() {
	w := webview.New(webview.Settings{
		Width:  1200,
		Height: 800,
		Title:  "text",
		URL:    "data:text/html," + url.PathEscape(indexHTML),
		ExternalInvokeCallback: handleRPC,
		Debug:     true,
		Resizable: true,
	})
	defer w.Exit()

	w.Run()
}
