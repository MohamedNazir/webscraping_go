package controller

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/MohamedNazir/webscraper/webscraper/service"
)

const (
	QUERY       = "queryUrl"
	PARSING_ERR = "form parsing error : %v"
)

type HtmlParserController struct {
	Service service.HtmlParserService
}

var tpl = template.Must(template.ParseFiles("./asset/index.html"))

func (hpc *HtmlParserController) IndexHandler(w http.ResponseWriter, r *http.Request) {
	tpl.Execute(w, nil)
}

func (hpc *HtmlParserController) SearchHandler(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		fmt.Fprintln(w, fmt.Errorf(PARSING_ERR, err))
		return
	}

	u := r.FormValue(QUERY)
	result, err := hpc.Service.ParseHtml(u)
	if err != nil {
		fmt.Fprintln(w, fmt.Errorf(err.Error()))
		return
	}

	tpl.Execute(w, result)

}
