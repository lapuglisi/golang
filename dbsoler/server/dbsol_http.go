package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/lapuglisi/dbsoler/coredb"
)

const (
	HttpServerAddress string = ":5120"
	MtsolHomeNotFound string = "A variável MTSOL_HOME não está definida."
)

func StartServer() (err error) {

	// Define os handlers
	http.HandleFunc("/", httpHandler)
	http.HandleFunc("/view", datViewHandler)

	return http.ListenAndServe(HttpServerAddress, nil)
}

func httpHandler(response http.ResponseWriter, req *http.Request) {

	beginBody := []byte(`
	<html>
	<head>
	  <title>DBSoler Suave</title>
	<head>
	<body style='font-family: verdana; font-size: 8pt'>`)

	response.Write(beginBody)

	mtsolHome := os.Getenv("MTSOL_HOME")
	if len(mtsolHome) == 0 {
		response.Write([]byte(MtsolHomeNotFound))
		return
	}

	// Starting from MTSOL_HOME, get dat/ directory
	datDir := mtsolHome + "/dat"

	_, err := os.Stat(datDir)
	if err != nil {
		response.Write([]byte(err.Error()))
		return
	}

	// Verifica se parametros dir esta sendo usado
	subDir := req.FormValue("dir")
	if len(subDir) > 0 {
		datDir = fmt.Sprintf("%s/%s", datDir, subDir)
	}

	dirList, err := ioutil.ReadDir(datDir)

	if err != nil {
		response.Write([]byte(err.Error()))
		return
	}

	var html string = `<table border='1' cellSpacing='0' cellPadding='5' 
		style='border: 1px solid blue; font-size: 8pt;'>
	<tr>
	<td colspan='4' align='center' bgColor='#ededed'>Tabelas do MTSOL: ` + datDir + `</td>
	</tr><tr>`
	var cols int = 0

	for _, entry := range dirList {
		if entry.IsDir() {
			html += fmt.Sprintf("<td bgColor='#00ADAD'><a style='color: red' href='/?dir=%s'>%s</a></td>",
				entry.Name(), entry.Name())
			cols++
		} else {
			fileName := entry.Name()
			if strings.HasSuffix(fileName, ".tab") {
				html += fmt.Sprintf("<td><a href='/view?path=%s/%s'>%s</a></td>",
					datDir, entry.Name(), entry.Name())
				cols++
			}
		}

		if cols%4 == 0 {
			html += "</tr><tr>"
		}
	}

	for index := cols; index < 4; index++ {
		html += "<td bgColor='#efefef'>&nbsp;</td>"
	}

	html += "</tr></table></body></html>"

	response.Write([]byte(html))
}

func datViewHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Add("Content-Type", "text/html; charset=UTF-8")
	resp.WriteHeader(200)

	datFile := req.FormValue("path")
	if len(datFile) == 0 {
		resp.Write([]byte("Arquivo nao definido!"))
		return
	}

	_, err := os.Stat(datFile)
	if err != nil {
		resp.Write([]byte(err.Error()))
		return
	}

	// Vamos obter o nome do arquivo
	fileBase := filepath.Base(datFile)
	fileParts := strings.Split(fileBase, ".")
	fileExt := fileParts[1]

	if fileExt != "tab" {
		r := fmt.Sprintf("Arquivo do tipo '%s' ainda nao suportado.", fileExt)
		resp.Write([]byte(r))
		return
	}

	from := 0
	adList, err, _ := coredb.ReadFile(datFile, from)
	if err != nil {
		resp.Write([]byte(err.Error()))
		return
	}

	var html string = fmt.Sprintf(`
	<html>
	<head>
	  <title>DBSoler - %s</title>
	</head<
	<body>
	<div align=center style='font-family: verdana; font-size: 10pt; font-weight: bold'>%s</div>
	<br />
	<table border=1 cellPadding=5 cellSpacing=0 
		style='border: 1px solid black; font-family: verdana; font-size: 8pt'>`,
		fileBase, fileBase)

	hasHeader := false

	for _, dbItem := range adList {
		item := coredb.DatToMap(dbItem)

		keys := make([]string, 0, len(item))
		for key := range item {
			keys = append(keys, key)
		}

		sort.Slice(keys, func(a, b int) bool { return keys[a] < keys[b] })

		if !hasHeader {
			html += "<tr style='background: #dedede'>"
			for index := 0; index < len(keys); index++ {
				html += ("<td>" + keys[index] + "</td>")
			}
			html += "</tr>"

			hasHeader = true
		}

		html += "<tr>"
		for index := 0; index < len(keys); index++ {
			html += fmt.Sprintf("<td>%s</td>", item[keys[index]])
		}

		html += "</tr>"

	}

	html += "</table></body></html>"

	resp.Write([]byte(html))

}
