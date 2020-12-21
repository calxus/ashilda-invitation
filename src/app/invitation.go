package main

import (
	"database/sql"
	"log"
)

// Invitation type represents the table by the same name in the database
type Invitation struct {
	UserID  int `json:"user_id,omitempty"`
	EventID int `json:"event_id,omitempty"`
}

func (i *Invitation) populate(rows *sql.Rows) {
	err := rows.Scan(&i.UserID, &i.EventID)
	if err != nil {
		log.Print(err)
	}
}
