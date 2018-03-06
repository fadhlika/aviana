package handler

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/fadhlika/aviana/app/config"
	"github.com/gorilla/mux"

	"github.com/rs/xid"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//CreateDevice handler
func CreateDevice(cfg *config.Config, DB *mgo.Database, w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Panicln("Read body ", err)
	}
	defer r.Body.Close()

	doc := bson.M{
		"_id": xid.New().String(),
	}

	col := DB.C("devices")
	res := InsertDocument(col, b, &doc)

	SendWs(cfg, res)
	RespondJson(res, 200, w, r)
}

//GetAllDevice handler
func GetAllDevice(DB *mgo.Database, w http.ResponseWriter, r *http.Request) {
	var result []interface{}
	err := DB.C("devices").Find(nil).All(&result)
	if err != nil {
		log.Println(err)
	}
	RespondJson(result, 200, w, r)
}

//GetDevice handler
func GetDevice(DB *mgo.Database, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]

	var result interface{}
	err := DB.C("devices").FindId(id).One(&result)
	if err != nil {
		log.Println(err)
	}
	RespondJson(result, 200, w, r)
}

//UpdateDevice handler
func UpdateDevice(cfg *config.Config, DB *mgo.Database, w http.ResponseWriter, r *http.Request) {

}

//DeleteDevice handler
func DeleteDevice(cfg *config.Config, DB *mgo.Database, w http.ResponseWriter, r *http.Request) {

}
