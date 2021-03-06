//
// Author: Vinhthuy Phan, 2015 - 2017
//
package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

//-----------------------------------------------------------------
func informIPAddress() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err.Error() + "\n")
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && ipnet.IP.IsGlobalUnicast() {
			return ipnet.IP.String() + ":" + PORT
		}
	}
	return ""
}

//-----------------------------------------------------------------

func prepareCleanup() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, syscall.SIGTERM)

	go func() {
		for {
			select {
			case <-quit:
				fmt.Println("Preparing to stop server...")
				writeDB()
				os.Exit(1)
			}
		}
	}()
}

//-----------------------------------------------------------------
// Make sure teacher runs server on his own laptop.
//-----------------------------------------------------------------
func Authorize(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Host != "localhost:4030" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Unauthorized")
		} else {
			fn(w, r)
		}
	}
}

//-----------------------------------------------------------------
// Register automatically if a student is not yet registered.
//-----------------------------------------------------------------
func AutoRegister(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, ok := Boards[r.FormValue("uid")]; !ok {
			RegisterStudent(r.FormValue("uid"))
		}
		fn(w, r)
	}
}

//-----------------------------------------------------------------
func main() {
	SERVER = informIPAddress()
	fmt.Println("*********************************************")
	fmt.Printf("*   Code4Brownies (%s)\n", VERSION)
	fmt.Printf("*   Server address: %s\n", SERVER)
	fmt.Println("*********************************************\n")
	rand.Seed(time.Now().UnixNano())
	USER_DB = filepath.Join(".", "C4B_DB.csv")
	flag.StringVar(&USER_DB, "db", USER_DB, "user database in csv format, which consists of UID,POINTS.")
	flag.Parse()
	prepareCleanup()

	// student handlers
	// http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/share", AutoRegister(shareHandler))
	http.HandleFunc("/my_points", AutoRegister(my_pointsHandler))
	http.HandleFunc("/receive_broadcast", AutoRegister(receive_broadcastHandler))
	http.HandleFunc("/check_broadcast", AutoRegister(check_broadcastHandler))

	// teacher handlers
	http.HandleFunc("/query_poll", Authorize(query_pollHandler))
	http.HandleFunc("/view_poll", Authorize(view_pollHandler))
	http.HandleFunc("/answer_poll", Authorize(answer_pollHandler))
	// http.HandleFunc("/start_poll", Authorize(start_pollHandler))
	http.HandleFunc("/points", Authorize(pointsHandler))
	http.HandleFunc("/give_points", Authorize(give_pointsHandler))
	http.HandleFunc("/peek", Authorize(peekHandler))
	http.HandleFunc("/broadcast", Authorize(broadcastHandler))
	http.HandleFunc("/get_post", Authorize(get_postHandler))
	http.HandleFunc("/get_posts", Authorize(get_postsHandler))

	empty_file, _ := loadDB()
	if empty_file {
		initDB()
	}
	err := http.ListenAndServe("0.0.0.0:"+PORT, nil)
	if err != nil {
		panic(err.Error() + "\n")
	}
}
