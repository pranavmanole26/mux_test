package main

import (
	"encoding/json"
	"fmt"
	"mux_test/person"
	"net/http"
	"strconv"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func testRedis() {
	fmt.Println("Go Redis Tutorial")

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := client.Ping().Result()
	name := client.Get("name")
	fmt.Println(name, err)
}

func main() {

	router := mux.NewRouter()

	// router.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
	// 	fmt.Println("Hello world")
	// 	resData := "Hello world"

	// 	qp := req.URL.Query().Get("id")

	// 	if qp != "" {
	// 		resData += qp
	// 	}
	// 	res.Write([]byte(resData))
	// })

	// router.HandleFunc("/{id}", func(res http.ResponseWriter, req *http.Request) {
	// 	vars := mux.Vars(req)
	// 	res.Write([]byte(vars["id"]))
	// }).Methods(http.MethodGet)

	// router.HandleFunc("/{id}", func(res http.ResponseWriter, req *http.Request) {
	// 	vars := mux.Vars(req)
	// 	res.Write([]byte(vars["id"]))
	// })

	http.HandleFunc("/hello", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("Hello world"))
	})

	router.Path("/home").Methods(http.MethodPost).HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("Welcome to home"))
	})

	router.HandleFunc("/helloM", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("Hello Mux"))
	}).Methods(http.MethodPost, http.MethodGet)

	psRouter := router.PathPrefix("/person").Subrouter()

	psRouter.Path("").Methods(http.MethodGet).HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		ps := person.GetPersons()
		resBody, err := json.Marshal(ps)
		if err != nil {
			fmt.Printf("Error while marshaling records: %v", err)
			res.Write([]byte("Error while fecthing records"))
			return
		}
		fmt.Printf("Records fetched")
		res.Write([]byte(resBody))
	})

	psRouter.Path("").Methods(http.MethodPost).HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		var p person.Person
		if err := json.NewDecoder(req.Body).Decode(&p); err != nil {
			fmt.Printf("Error while reading body: %v", err)
		}
		lastId, err := person.AddPerson(p)
		if err != nil {
			fmt.Printf("Error occured while adding new person record: %v", err)
			res.Write([]byte("Error while adding new person"))
			return
		}
		res.Write([]byte(fmt.Sprintf("New person added with id: %d", lastId)))
	})

	psRouter.Path("/{id}").Methods(http.MethodGet).HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		pId, err := strconv.Atoi(mux.Vars(req)["id"])
		if err != nil {
			fmt.Println("Error while parsing id path param")
			res.WriteHeader(http.StatusBadRequest)
			res.Write([]byte("Improper path param for person id."))
			return
		}
		p, err := person.GetPerson(pId)
		if err != nil || p.Id <= 0 {
			fmt.Printf("Error while fecthing person details: %v", err)
			res.WriteHeader(http.StatusNotFound)
			res.Write([]byte(fmt.Sprintf("Person details with %d does not found.", pId)))
			return
		}
		res.WriteHeader(http.StatusOK)
		resBody, _ := json.Marshal(p)
		res.Write([]byte(resBody))
	})

	http.ListenAndServe(":82", router)
}
