package handler

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/fadhlika/aviana/app/config"
	"github.com/gorilla/mux"
	"github.com/rs/xid"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//CreateData handler
func CreateData(cfg *config.Config, DB *mgo.Database, w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	var result bson.M
	err := DB.C("devices").FindId(id).One(&result)
	if err != nil {
		w.WriteHeader(204)
		log.Panicln(err)
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Panicln(err)
	}
	defer r.Body.Close()

	doc := bson.M{
		"_id":         xid.New().String(),
		"device_id":   result["_id"],
		"device_name": result["name"],
		"type":        result["type"],
	}

	col := DB.C("data")
	res := InsertDocument(col, b, &doc)
	SendWs(cfg, res)
	RespondJson(res, 200, w, r)
}

//GetAllData handler
func GetAllData(DB *mgo.Database, w http.ResponseWriter, r *http.Request) {
	var result []bson.M
	DB.C("data").Find(nil).Sort("-date").All(&result)

	RespondJson(result, 200, w, r)
}

//GetData handler
func GetData(DB *mgo.Database, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	deviceID := vars["id"]
	limit, err := strconv.Atoi(vars["limit"])
	if err != nil {
		log.Println(err)
	}

	var result []bson.M
	DB.C("data").Find(bson.M{"device_id": deviceID}).Sort("-date").Limit(limit).All(&result)

	RespondJson(result, 200, w, r)
}

//UpdateData handler
func UpdateData(cfg *config.Config, DB *mgo.Database, w http.ResponseWriter, r *http.Request) {

}

//DeleteData handler
func DeleteData(cfg *config.Config, DB *mgo.Database, w http.ResponseWriter, r *http.Request) {

}
