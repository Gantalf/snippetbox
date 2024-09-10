package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"snippetbox.gantalf.net/internal/models"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	snippets *models.SnippetModel
}

func main() {
	addr := flag.String("addr", ":4000", "HHTP network address")
	flag.Parse() //en caso que se escriba algo en la linea de commando lo lee y asgina a addr

	// Define a new command-line flag for the MySQL DSN string.
	dsn := flag.String("dsn", "web:webUser1.@tcp(127.0.0.1:3306)/snippetbox?parseTime=true", "MySQL data source name")

	//parametros: el destino donde escribira el log, un prefijo que va a tener el msj y flags para indicar que informacion extra va a incluir
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		snippets: &models.SnippetModel{DB: db},
	}

	// lo hacemos para usar nuestra configuracion de logs en errores http, y no el defaul del http.Server
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	//log.Printf("Starting server on %s", *addr)
	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	//log.Fatal(err)
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
