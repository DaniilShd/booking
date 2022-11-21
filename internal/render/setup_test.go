package render

import (
	"encoding/gob"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/DaniilShd/booking/internal/config"
	"github.com/DaniilShd/booking/internal/models"
	"github.com/alexedwards/scs/v2"
)

var session *scs.SessionManager
var testApp config.AppConfig

func TestMain(m *testing.M) {

	// what am i going to put in the session
	gob.Register(models.Reservation{})

	// change this to true when in production
	testApp.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = testApp.InProduction

	testApp.Session = session

	app = &testApp

	os.Exit(m.Run())
}

type MyWriter struct{}

func (tw *MyWriter) Header() http.Header {
	var h http.Header
	return h
}

func (tw *MyWriter) WriteHeader(statusCode int) {

}

func (tw *MyWriter) Write(b []byte) (int, error) {
	length := len(b)
	return length, nil
}
