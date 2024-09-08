package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	addr := flag.String("addr", ":4000", "HHTP network address")
	flag.Parse() //en caso que se escriba algo en la linea de commando lo lee y asgina a addr

	//parametros: el destino donde escribira el log, un prefijo que va a tener el msj y flags para indicar que informacion extra va a incluir
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	// lo hacemos para usar nuestra configuracion de logs en errores http, y no el defaul del http.Server
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	//log.Printf("Starting server on %s", *addr)
	infoLog.Printf("Starting server on %s", *addr)
	err := srv.ListenAndServe()
	//log.Fatal(err)
	errorLog.Fatal(err)
}
