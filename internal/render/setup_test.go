package render

import (
	"encoding/gob"
	"github.com/alexedwards/scs/v2"
	"github.com/satya-kr/bookings/internal/config"
	"github.com/satya-kr/bookings/internal/models"
	"net/http"
	"os"
	"testing"
	"time"
)

var session *scs.SessionManager
var testApp config.AppConfig

func TestMain(m *testing.M) {

	//put in the session
	gob.Register(models.Reservation{})

	//change this to true when we are in production
	testApp.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false //set false because we use http insted of https
	testApp.Session = session

	app = &testApp

	os.Exit(m.Run())
}
