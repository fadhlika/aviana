package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"

	"github.com/fadhlika/aviana/app/config"
	"github.com/fadhlika/aviana/app/handler"

	"github.com/rs/cors"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
)

//App Application struct
type App struct {
	AppCfg   *config.Config
	Session  *mgo.Session
	Router   *mux.Router
	Upgrader *websocket.Upgrader
	Clients  *handler.Clients
}

//Init application
func (a *App) Init() {
	dat, err := ioutil.ReadFile("app.conf")
	if err != nil {
		log.Panicln(err)
	}

	cfg := &config.Config{}
	err = json.Unmarshal(dat, cfg)
	if err != nil {
		fmt.Println(string(dat))
		log.Panicln(err)
	}

	a.AppCfg = cfg

	session, err := mgo.Dial(a.AppCfg.DbURL)
	if err != nil {
		log.Panicln(err)
	}
	a.Session = session

	a.Router = mux.NewRouter()
	a.SetupRouter()

	a.Upgrader = &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	a.Clients = &handler.Clients{
		make(map[*websocket.Conn]bool),
		sync.Mutex{},
		make(chan interface{}),
	}

}

//CreateDevice register new device wrapper
func (a *App) CreateDevice(w http.ResponseWriter, r *http.Request) {
	session := a.GetSession()
	defer session.Close()

	db := a.GetDB(session)
	handler.CreateDevice(a.AppCfg, db, w, r)
}

//GetAllDevice register new device wrapper
func (a *App) GetAllDevice(w http.ResponseWriter, r *http.Request) {
	session := a.GetSession()
	defer session.Close()

	db := a.GetDB(session)
	handler.GetAllDevice(db, w, r)
}

//GetDevice register new device wrapper
func (a *App) GetDevice(w http.ResponseWriter, r *http.Request) {
	session := a.GetSession()
	defer session.Close()

	db := a.GetDB(session)
	handler.GetDevice(db, w, r)
}

//UpdateDevice register new device wrapper
func (a *App) UpdateDevice(w http.ResponseWriter, r *http.Request) {
	session := a.GetSession()
	defer session.Close()

	db := a.GetDB(session)
	handler.UpdateDevice(a.AppCfg, db, w, r)
}

//DeleteDevice Delete selected device
func (a *App) DeleteDevice(w http.ResponseWriter, r *http.Request) {
	session := a.GetSession()
	defer session.Close()

	db := a.GetDB(session)
	handler.DeleteDevice(a.AppCfg, db, w, r)
}

//CreateData register new device wrapper
func (a *App) CreateData(w http.ResponseWriter, r *http.Request) {
	session := a.GetSession()
	defer session.Close()

	db := a.GetDB(session)
	handler.CreateData(a.AppCfg, db, w, r)
}

//GetAllData register new device wrapper
func (a *App) GetAllData(w http.ResponseWriter, r *http.Request) {
	session := a.GetSession()
	defer session.Close()

	db := a.GetDB(session)
	handler.GetAllData(db, w, r)
}

//GetData register new device wrapper
func (a *App) GetData(w http.ResponseWriter, r *http.Request) {
	session := a.GetSession()
	defer session.Close()

	db := a.GetDB(session)
	handler.GetData(db, w, r)
}

//UpdateData register new device wrapper
func (a *App) UpdateData(w http.ResponseWriter, r *http.Request) {
	session := a.GetSession()
	defer session.Close()

	db := a.GetDB(session)
	handler.UpdateData(a.AppCfg, db, w, r)
}

//DeleteData register new device wrapper
func (a *App) DeleteData(w http.ResponseWriter, r *http.Request) {
	session := a.GetSession()
	defer session.Close()

	db := a.GetDB(session)
	handler.DeleteData(a.AppCfg, db, w, r)
}

func (a *App) Echo(w http.ResponseWriter, r *http.Request) {
	handler.HandleConnection(a.Upgrader, a.Clients, w, r)
}

//GetSession Copy mongodb session
func (a *App) GetSession() *mgo.Session {
	return a.Session.Copy()
}

//GetDB get database
func (a *App) GetDB(session *mgo.Session) *mgo.Database {
	return session.DB(a.AppCfg.DbName)
}

//SetupRouter setup router
func (a *App) SetupRouter() {
	a.Router.HandleFunc("/device", a.CreateDevice).Methods("POST")
	a.Router.HandleFunc("/device", a.GetAllDevice).Methods("GET")
	a.Router.HandleFunc("/device/{id}", a.GetDevice).Methods("GET")
	a.Router.HandleFunc("/device", a.UpdateDevice).Methods("PUT")
	a.Router.HandleFunc("/device", a.DeleteDevice).Methods("DELETE")

	a.Router.HandleFunc("/data", a.CreateData).Methods("POST")
	a.Router.HandleFunc("/data", a.GetAllData).Methods("GET")
	a.Router.HandleFunc("/data/{id}/{limit}", a.GetData).Methods("GET")
	a.Router.HandleFunc("/data", a.UpdateData).Methods("PUT")
	a.Router.HandleFunc("/data", a.DeleteData).Methods("DELETE")

	a.Router.HandleFunc("/websocket", a.Echo)
}

//Run application
func (a *App) Run() {
	log.Printf("Running on http://%s:%d", a.AppCfg.URL, a.AppCfg.Port)
	go handler.HandleMessage(a.Clients)

	err := http.ListenAndServe(
		fmt.Sprintf("%s:%d", a.AppCfg.URL, a.AppCfg.Port),
		cors.Default().Handler(a.Router),
	)
	if err != nil {
		log.Println("HTTP ", err)
	}
}
