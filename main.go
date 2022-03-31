package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/accurati-bi/csw-rpa-elt/elts"
	"github.com/gorilla/mux"
	"github.com/robfig/cron/v3"
)

var c = cron.New()

func main() {

	c.AddFunc("0 6 * * *", elts.EltGclese610)
	c.AddFunc("30 6 * * *", elts.EltCcesu665)
	c.AddFunc("45 6 * * *", elts.EltCcese320)
	c.AddFunc("0 7 * * *", elts.EltGcleft300)
	c.AddFunc("0 5 * * 1", elts.EltAslepme600) //Toda segunda-feira as 5:00
	c.AddFunc("0 6 * * *", elts.EltSigmaEstoque)

	c.Start()
	defer c.Stop()

	path, err := os.Executable()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(path)

	r := mux.NewRouter()
	r.HandleFunc("/elt/ccesu665", handlerCcesu665)
	r.HandleFunc("/elt/ccese320", handlerCcese320)
	r.HandleFunc("/elt/gcleese610", handlerGcleese610)
	r.HandleFunc("/elt/gcleft300", handlerGcleft300)
	r.HandleFunc("/elt/aslepme600", handlerAslepme600)

	r.HandleFunc("/elt/sigmaEstoque", handlerSigmaEstoque)

	r.HandleFunc("/elt/all", handlerAll)
	r.HandleFunc("/teste", testePost).Methods("POST")
	r.HandleFunc("/teste", testeGet).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", r))

}

func handlerCcese320(w http.ResponseWriter, r *http.Request) {
	elts.EltCcese320()
}

func handlerCcesu665(w http.ResponseWriter, r *http.Request) {
	elts.EltCcesu665()
}

func handlerGcleese610(w http.ResponseWriter, r *http.Request) {
	elts.EltGclese610()
}

func handlerGcleft300(w http.ResponseWriter, r *http.Request) {
	elts.EltGcleft300()
}

func handlerAslepme600(w http.ResponseWriter, r *http.Request) {
	elts.EltAslepme600()
}

func handlerSigmaEstoque(w http.ResponseWriter, r *http.Request) {
	elts.EltSigmaEstoque()
}

func handlerAll(w http.ResponseWriter, r *http.Request) {
	elts.EltGclese610()
	elts.EltCcesu665()
	elts.EltCcese320()
	elts.EltGcleft300()
	elts.EltAslepme600()
}
