package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/fadhlika/aviana/app/config"
	"github.com/gorilla/websocket"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func RespondJson(doc interface{}, status int, w http.ResponseWriter, r *http.Request) {
	res, err := json.Marshal(doc)
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(res)
}

func SendWs(cfg *config.Config, data interface{}) {
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%d/websocket", cfg.URL, cfg.Port), nil)
	if err != nil {
		log.Printf("Dial websocket: %s", err)
	}
	defer conn.Close()

	conn.WriteJSON(data)
}

func InsertDocument(col *mgo.Collection, data []byte, doc *bson.M) bson.M {

	err := json.Unmarshal(data, &doc)
	if err != nil {
		log.Println(err)
	}

	err = col.Insert(doc)
	if err != nil {
		log.Println(err)
	}

	var result bson.M
	err = col.Find(doc).One(&result)
	if err != nil {
		log.Println(err)
	}

	return result
}
