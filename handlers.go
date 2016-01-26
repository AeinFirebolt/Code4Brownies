//
// Author: Vinhthuy Phan, 2015
//
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

//-----------------------------------------------------------------

func verifyPasscode(w http.ResponseWriter, r *http.Request) error {
	if r.FormValue("passcode") != PASSCODE {
		w.WriteHeader(http.StatusUnauthorized)
		return errors.New("Unauthorized access")
	}
	return nil
}

//-----------------------------------------------------------------
// users query to know their current points
//-----------------------------------------------------------------
func my_pointsHandler(w http.ResponseWriter, r *http.Request) {
	user := r.FormValue("uid")
	_, ok := AllUsers[user]
	if !ok {
		AllUsers[user] = &User{0}
	}
	record := AllUsers[user]
	cur_points := Points.get(user)
	mesg := fmt.Sprintf("Points for %s\nCurrent points: %d\nTotal points: %d\n", user, cur_points, record.points)
	fmt.Fprintf(w, mesg)
}

//-----------------------------------------------------------------
// users submit their codes
//-----------------------------------------------------------------
func submit_postHandler(w http.ResponseWriter, r *http.Request) {
	user, body := r.FormValue("uid"), r.FormValue("body")
	if _, ok := AllUsers[user]; !ok {
		AllUsers[user] = &User{0}
	}
	Posts.Add(user, body)
	fmt.Println(user, "submitted.")
	fmt.Fprintf(w, user+" submitted succesfully.")
}

//-----------------------------------------------------------------
// return points of currently awarded users
//-----------------------------------------------------------------
func pointsHandler(w http.ResponseWriter, r *http.Request) {
	if verifyPasscode(w, r) == nil {
		js, err := json.Marshal(Points.data)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
		}
	}
}

//-----------------------------------------------------------------
// give one brownie point to a user
//-----------------------------------------------------------------
func give_pointHandler(w http.ResponseWriter, r *http.Request) {
	if verifyPasscode(w, r) == nil {
		Points.addOne(r.FormValue("uid"))
		fmt.Println("+1", r.FormValue("uid"))
		fmt.Fprintf(w, "Point awarded to "+r.FormValue("uid"))
	}
}

//-----------------------------------------------------------------
// return all current posts
//-----------------------------------------------------------------
func peekHandler(w http.ResponseWriter, r *http.Request) {
	if verifyPasscode(w, r) == nil {
		js, err := json.Marshal(Posts.queue)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
		}
	}
}

//-----------------------------------------------------------------
// Instructor retrieves code
//-----------------------------------------------------------------
func get_postHandler(w http.ResponseWriter, r *http.Request) {
	if verifyPasscode(w, r) == nil {
		e, err := strconv.Atoi(r.FormValue("post"))
		if err != nil {
			fmt.Println(err.Error)
		} else {
			js, err := json.Marshal(Posts.Remove(e))
			if err != nil {
				fmt.Println(err.Error())
			} else {
				w.Header().Set("Content-Type", "application/json")
				w.Write(js)
			}
		}
	}
}

//-----------------------------------------------------------------
// Instructor retrieves all codes
//-----------------------------------------------------------------
func get_postsHandler(w http.ResponseWriter, r *http.Request) {
	if verifyPasscode(w, r) == nil {
		js, err := json.Marshal(Posts.queue)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			Posts.Clear()
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
		}
	}
}
