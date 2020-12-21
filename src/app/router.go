package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	models "github.com/gtadam/ashilda-common"
)

type router struct {
	basePath string
	database models.Database
	mux      *mux.Router
}

func newRouter(bp string) *router {
	return &router{
		basePath: bp,
		database: *models.NewDatabase(),
		mux:      mux.NewRouter().StrictSlash(true),
	}
}

func (rt *router) getInvitation(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["user_id"]
	eventID := mux.Vars(r)["event_id"]
	statement := models.NewDatabaseSelect(Table)
	statement.AddColumn(UserIDField)
	statement.AddColumn(EventIDField)
	if userID != "0" {
		statement.AddCondition(UserIDField, models.EQUALS, userID)
	}
	if eventID != "0" {
		statement.AddCondition(EventIDField, models.EQUALS, eventID)
	}
	rows, _ := rt.database.ExecuteSelect(statement)
	invitations := []Invitation{}
	for rows.Next() {
		invitation := Invitation{}
		invitation.populate(rows)
		invitations = append(invitations, invitation)
	}
	rows.Close()
	if len(invitations) == 0 {
		w.WriteHeader(404)
		return
	}
	json, _ := json.Marshal(invitations)
	fmt.Fprintf(w, string(json))
}

func (rt *router) postInvitation(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	invitation := Invitation{}
	json.Unmarshal(body, &invitation)
	statement := models.NewDatabaseInsert(Table)
	statement.AddEntry(UserIDField, strconv.Itoa(invitation.UserID))
	statement.AddEntry(EventIDField, strconv.Itoa(invitation.EventID))
	rt.database.ExecuteInsert(statement)
}

func (rt *router) deleteInvitation(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	invitation := Invitation{}
	json.Unmarshal(body, &invitation)
	statement := models.NewDatabaseDelete(Table)
	statement.AddCondition(UserIDField, models.EQUALS, strconv.Itoa(invitation.UserID))
	statement.AddCondition(EventIDField, models.EQUALS, strconv.Itoa(invitation.EventID))
	rt.database.ExecuteDelete(statement)
}

func (rt *router) populateRoutes() {
	rt.database.Connect()
	rt.mux.HandleFunc(rt.basePath+"/invitation/{user_id:[0-9]+}/{event_id:[0-9]+}", rt.getInvitation).Methods("GET")
	rt.mux.HandleFunc(rt.basePath+"/invitation", rt.postInvitation).Methods("POST")
	rt.mux.HandleFunc(rt.basePath+"/invitation", rt.deleteInvitation).Methods("DELETE")
}
