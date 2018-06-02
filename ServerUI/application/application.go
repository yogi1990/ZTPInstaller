package application

import (
	"github.com/carbocation/interpose"
	gorilla_mux "github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/spf13/viper"
	"net/http"

	"github.com/ZTPInstaller/ServerUI/handlers"
	"github.com/ZTPInstaller/ServerUI/middlewares"
)

// New is the constructor for Application struct.
func New(config *viper.Viper) (*Application, error) {
	cookieStoreSecret := config.Get("cookie_secret").(string)

	app := &Application{}
	app.config = config
	app.sessionStore = sessions.NewCookieStore([]byte(cookieStoreSecret))

	return app, nil
}

// Application is the application object that runs HTTP server.
type Application struct {
	config      *viper.Viper
	sessionStore sessions.Store
}

func (app *Application) MiddlewareStruct() (*interpose.Middleware, error) {
	middle := interpose.New()
	middle.Use(middlewares.SetSessionStore(app.sessionStore))

	middle.UseHandler(app.mux())

	return middle, nil
}

func (app *Application) mux() *gorilla_mux.Router {
	router := gorilla_mux.NewRouter()

	router.Handle("/web-ui-dashboard", http.HandlerFunc(handlers.GetHome)).Methods("GET")
	router.Handle("/create", http.HandlerFunc(handlers.GetCreate)).Methods("GET")
	router.Handle("/edit", http.HandlerFunc(handlers.GetEdit)).Methods("GET")
	router.Handle("/submit",http.HandlerFunc(handlers.ProcessCreate)).Methods("POST")
	router.Handle("/submit-update",http.HandlerFunc(handlers.ProcessEdit)).Methods("POST")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("static")))

	return router
}
