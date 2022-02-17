package main

import (
	"encoding/json"
	"net/http"
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
	store map[string]User
}

// get method function to serve teh data
func (h *userHandlers) get(w http.ResponseWriter, r *http.Request){
	// getting data of map in a list
	users := make([]User, len(h.store))

	i := 0
	// list of user
	for _, user := range h.store {
		users[i] = user
		i++
	}

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

// a constructor function (don't take any argument and return on of our user object)
func newUserHandlers() *userHandlers {
	// point to a new user handlers
	return &userHandlers{
		store: map[string]User{
			// hardcoded data for testing
			"id1": User{
				ID: "someid",
				FirstName: "sparsh",
				LastName: "saxena",
				Dateofbirth: "17-12-2000",
				Email: "sparsh0987654321@gmail.com",
				Phonenumber: +918439803019,
			},
		},
	}
}

// main function
func main() {
	userHandlers := newUserHandlers()  // calling constructor function 
	http.HandleFunc("/users", userHandlers.get)
	// building a http server
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		panic(err)
	}
}