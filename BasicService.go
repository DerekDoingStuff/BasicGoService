// BasicService
package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"
	"io/ioutil"
	"strings"
	"github.com/gorilla/mux"
	_ "github.com/denisenkom/go-mssqldb"
)

var (
	db *sql.DB
)

func simpleHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, www112!")
}

func derekHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Sup ")
	fmt.Fprintf(w, time.Now().Local().Format("15:04:05"))
	
}

func ioHandler(w http.ResponseWriter, r *http.Request) {
	thingie, err := ioutil.ReadFile("BasicService.go")
	if err != nil{
		panic(err)
	}
	
	fmt.Fprintf(w, string(thingie))
}

func ioWriteHandler(w http.ResponseWriter, r *http.Request) {
	thingie, err := ioutil.ReadFile("BasicService.go")
	if err != nil{
		panic(err)
	}
	
	mystuff := strings.Replace(string(thingie), "www", "www1", 1)
	
	err = ioutil.WriteFile("BasicService.go", []byte(mystuff), 0644)
	if err != nil{
		panic(err)
	}
	
	fmt.Fprintf(w, string(mystuff))
}

func sessionHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	sessionID := params["SessionID"]
	uid := params["uid"]
		
	fmt.Fprintln(w, sessionID)
	fmt.Fprintf(w, string(uid))
}

func dbHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["uid"]
	
	rows, err := db.Query("SELECT AppUserID, Handle FROM AppUser where AppUserID = ? order by AppUserID desc", id)
	
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	
	for rows.Next(){
		var Handle string
		var AppUserID string
		if err := rows.Scan(&AppUserID, &Handle); err != nil {
			panic(err)
		}
		fmt.Fprintln(w, AppUserID, Handle)
	}
	if err := rows.Err(); err != nil {
            panic(err)
    }
}

func main() {
	
	database, err := sql.Open("mssql", "server=IP_ADDR;database=DB_NAME;user id=USER_ID;password=SUPER_PASSWORD")
	if err != nil {
		panic(err)
	}
	
	db = database
	
	r:= mux.NewRouter()
	r.HandleFunc("/", simpleHandler)
	r.HandleFunc("/derek", derekHandler)
	r.HandleFunc("/file", ioHandler)
	r.HandleFunc("/write", ioWriteHandler)
	r.HandleFunc("/db/User/{uid:[0-9]+}", dbHandler).Name("User")
	r.HandleFunc("/Session/{SessionID}/{uid:[0-9]+}", sessionHandler).Name("Session")
	
	http.ListenAndServe(":8082", r)
}
