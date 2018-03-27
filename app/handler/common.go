package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/fadhlika/aviana/app/config"
	"github.com/gorilla/websocket"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func RespondJson(doc interface{}, status int, w http.ResponseWriter, r *http.Request) {
	res, err := json.Marshal(doc)
	if err != nil {
		log.Println("Respond ", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(res)
}

func RespondExcel(xlsx *excelize.File, deviceID string, status int, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.xlsx", deviceID))
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	w.WriteHeader(status)
	xlsx.Write(w)
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

	err := json.Unmarshal(data, doc)
	if err != nil {
		log.Println("Unmarshal ", err)
	}

	if col.Name == "data" {
		jakarta, _ := time.LoadLocation("Asia/Jakarta")
		date := time.Now().In(jakarta)
		(*doc)["date"] = date
	}

	err = col.Insert(doc)
	if err != nil {
		log.Println("Insert document ", err)
	}

	var result bson.M
	err = col.Find(doc).One(&result)
	if err != nil {
		log.Println("Find ", err)
	}

	return result
}
