package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/DaniilShd/booking/internal/config"
	"github.com/DaniilShd/booking/internal/driver"
	"github.com/DaniilShd/booking/internal/handlers"
	"github.com/DaniilShd/booking/internal/helpers"
	"github.com/DaniilShd/booking/internal/models"
	"github.com/DaniilShd/booking/internal/render"
	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {

	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() (*driver.DB, error) {
	// what am i going to put in the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})

	// change this to true when in production
	app.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	//connect to database
	log.Println("Connecting to database...")
	dsn := fmt.Sprintf("host=localhost port=5432 dbname=booking user=testuser password=1234")
	db, err := driver.ConnectSQL(dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()
	log.Println("Connected to database!")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	app.TemplateCache = tc
	app.UseCache = false
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	handlers.NewHandlers(handlers.NewRepository(&app, db))

	return db, nil
}
