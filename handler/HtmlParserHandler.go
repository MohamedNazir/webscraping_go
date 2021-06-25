package handler

import (
	"html/template"
	"log"
	"net/http"

	"github.com/MohamedNazir/webscraper/service"
)

const (
	QUERY         = "queryUrl"
	DATA          = "data"
	ERR           = "err"
	PARSING_ERR   = "form parsing error : %v"
	ERR_MSG_ERROR = "Sorry, something went wrong"
	RESULT        = "result"
	SUCCESS_RES   = "sending response back"
)

type HtmlParserHandler struct {
	Service service.HtmlParserService
}

var (
	tmpl = template.Must(template.ParseFiles("../asset/index.html"))
)

func (hpc *HtmlParserHandler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Request received for /")
	render(w, nil)
}

func (hpc *HtmlParserHandler) SearchHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Request received for /search")
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		return
	}

	u := r.FormValue(QUERY)
	result, err := hpc.Service.ParseHtml(u)
	if err != nil {
		m := map[string]interface{}{DATA: string(err.Error())}
		err := map[string]interface{}{ERR: m}
		render(w, err)
		return
	}
	log.Println(SUCCESS_RES)
	render(w, map[string]interface{}{RESULT: result})
}

func render(w http.ResponseWriter, data interface{}) {

	if err := tmpl.Execute(w, data); err != nil {
		log.Println(err)
		http.Error(w, ERR_MSG_ERROR, http.StatusInternalServerError)
	}
}
