package handler

import (
	"encoding/json"
	"net/http"
	"ws-socket/models"
)

// User represents Temporary User on Joining
var User map[string]bool

func init() {
	User = make(map[string]bool)
}

// Join HandleFunc
func Join(w http.ResponseWriter, r *http.Request) {
	// create decoder to read the body of request (readcloser)
	decoder := json.NewDecoder(r.Body)
	user := models.User{}
	// bind to user struct
	if err := decoder.Decode(&user); err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	// check if user can join (already exist username)
	if !canJoin(user.Username) {
		w.Write([]byte("Duplicate"))
		return
	}
	// set the username to make it exist
	User[user.Username] = true
	w.Write([]byte("ok"))
	return
}

func canJoin(user string) bool {
	if isExist, ok := User[user]; isExist && ok {
		return false
	}
	return true
}
