package main

import (
	"encoding/json"
	"net/http"
	"sync"
	"fmt"
	"io/ioutil"
	"time"
)

// data types to return data
type User struct {
	ID           string `json:"id"`
	FirstName    string `json:"firstname"`
	LastName     string `json:"lastname"`
	Dateofbirth  string `json:"dateofbirth"`
	Email        string `json:"email"`
	Phonenumber  int    `json:"phonenumber"`
}

// database
type userHandlers struct {
	sync.Mutex  // to handle concurrent request 
	store map[string]User
}

// function to switch between post and get method
func (h *userHandlers) users(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// calling get method 
		h.get(w, r)
		return
	case "POST":
		// calling post method
		h.post(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
		return
	}
}

// get method function to serve teh data
func (h *userHandlers) get(w http.ResponseWriter, r *http.Request){
	// getting data of map in a list
	users := make([]User, len(h.store))

	h.Lock()
	i := 0
	// list of user
	for _, user := range h.store {
		users[i] = user
		i++
	}
	h.Unlock()

	// converting list in json format and checking error
	jsonBytes, err := json.Marshal(users)

	if err != nil {
		// TODO
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

//post method function to add data into api
func (h *userHandlers) post(w http.ResponseWriter, r *http.Request){
	// reading the body of json file sent
	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// chechking for content type
	ct := r.Header.Get("content-type")
	if ct != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("need content-type 'application/json', but got '%s'", ct)))
		return
	}

	var user User
	// converting from json format to list
	err = json.Unmarshal(bodyBytes, &user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	// generating user id from timestamps as ano third parties packages are involved
	user.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	// 
	h.Lock()
	h.store[user.ID] = user
	defer h.Unlock()
}

// a constructor function (don't take any argument and return on of our user object)
func newUserHandlers() *userHandlers {
	// point to a new user handlers
	return &userHandlers{
		store: map[string]User{
			// hardcoded data for testing initial get request
			// "id1": User{
			// 	ID: "someid",
			// 	FirstName: "sparsh",
			// 	LastName: "saxena",
			// 	Dateofbirth: "17-12-2000",
			// 	Email: "sparsh0987654321@gmail.com",
			// 	Phonenumber: +918439803019,
			// },
		},
	}
}

// main function
func main() {
	userHandlers := newUserHandlers()  // calling constructor function 
	http.HandleFunc("/users", userHandlers.users)
	// building a http server
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		panic(err)
	}
}