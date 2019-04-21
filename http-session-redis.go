package main

import (
	"fmt"
	"github.com/gorilla/sessions"
	"gopkg.in/boj/redistore.v1"
	"log"
	"net/http"
)

const
(
	CONN_HOST = "localhost"
	CONN_PORT = "8080"
)

//we declared a private store to store session using secure cookies
var store *redistore.RediStore

var err error

func init(){
	store, err = redistore.NewRediStore(10, "tcp", ":6379", "",
		[]byte("secret-key"))
	if err != nil {
		log.Fatal("error getting redis store : ", err)
	}
}

//function for home
func home (w http.ResponseWriter, r *http.Request){
	session,_ := store.Get(r,"sesion-name")
	var authenticate = session.Values["authenticated"]
	if authenticate != nil{
		isAuthenticated := session.Values["authenticated"].(bool)
		if(isAuthenticated){
			http.Error(w,"You are unauthorized to view the page",
				http.StatusForbidden)
			return
		}
		fmt.Fprintln(w,"Home page")
	}else
	{
		http.Error(w, "You are unauthorized to view the page",
			http.StatusForbidden)
		return
	}
}

//for logging user to session
func login(w http.ResponseWriter, r *http.Request){
	//Get value from cookie store with same name
	session, _ := store.Get(r,"session-name")
	//Set authenticated to true
	session.Values["authenticated"]=true
	//if error occured show message dialog error else allright
	if err = sessions.Save(r, w); err != nil{
		log.Fatalf("Error saving session: %v", err)
	}
	fmt.Fprintln(w, "You have successfully logged in.")
}

func logout(w http.ResponseWriter, r *http.Request)  {
	//Get value from session
	session, _ := store.Get(r,"sesion_name")
	//set user to unauthorized by setting false
	session.Values["authenticated"]=false
	//Save session
	session.Save(r,w)
	//Show error function
	fmt.Println(w, "You are logout")
}

func main() {
	http.HandleFunc("/home", home)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	err := http.ListenAndServe(CONN_HOST+":"+CONN_PORT, nil)
	defer store.Close()
	if err != nil{
		log.Fatal("error starting http server : ", err)
		return
	}
}