package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type User struct {
	FirstName      string
	LastName       string
	MobileNumber   int
	Email          string
	PW             string
	DateOfCreation string
	AccountType    string
	DriverLicense  string
	CarPLateNumber string
	AccountStatus  bool
}

type CarPool struct {
	UserID               string
	PickUpLocation       string
	AlternatePickUp      string
	StartTravellingTime  string
	AddressOfDestination string
	NumberOfPassengers   int
	PoolStatus           string
	NumberOfVacancies    int
}

type PassengerTrip struct {
	UserID     string
	CarPoolID  string
	TripStatus string
}

var (
	db  *sql.DB
	err error
)

func main() {
	db, err = sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/my_db")

	if err != nil {
		panic(err.Error())
	}

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/users/{userid}", user).Methods("DELETE", "POST", "PUT")
	router.HandleFunc("/api/v1/carpool/{carpoolid}", carpool).Methods("POST", "PUT")
	router.HandleFunc("/api/v1/passengertrip/{passengertripid}", passengertrip).Methods("POST", "PUT")
	fmt.Println("Listening at port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))
}

func user(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	if r.Method == "POST" {
		if body, err := ioutil.ReadAll(r.Body); err == nil {
			var data User
			fmt.Println(string(body))
			if err := json.Unmarshal(body, &data); err == nil {
				fmt.Println(data)
				insertUser(params["userid"], data)
				w.WriteHeader(http.StatusAccepted)
			} else {
				fmt.Println(err)
			}
		}
	} else if r.Method == "PUT" {
		if body, err := ioutil.ReadAll(r.Body); err == nil {
			var data User

			if err := json.Unmarshal(body, &data); err == nil {
				fmt.Println(data)
				updateUser(params["userid"], data)
				w.WriteHeader(http.StatusAccepted)
			} else {
				fmt.Println(err)
			}
		}
	} else if r.Method == "DELETE" {
		deluser(params["userid"])
		fmt.Fprintf(w, params["userid"]+" Deleted")
	} else {
	}
}

func carpool(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	if r.Method == "POST" {
		if body, err := ioutil.ReadAll(r.Body); err == nil {
			var data CarPool
			fmt.Println(string(body))
			if err := json.Unmarshal(body, &data); err == nil {
				fmt.Println(data)
				insertCarPool(params["carpoolid"], data)
				w.WriteHeader(http.StatusAccepted)
			} else {
				fmt.Println(err)
			}
		}
	} else if r.Method == "PUT" {
		if body, err := ioutil.ReadAll(r.Body); err == nil {
			var data CarPool

			if err := json.Unmarshal(body, &data); err == nil {
				fmt.Println(data)
				updateCarPool(params["carpoolid"], data)
				w.WriteHeader(http.StatusAccepted)
			} else {
				fmt.Println(err)
			}
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func passengertrip(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	if r.Method == "POST" {
		if body, err := ioutil.ReadAll(r.Body); err == nil {
			var data PassengerTrip
			fmt.Println(string(body))
			if err := json.Unmarshal(body, &data); err == nil {
				fmt.Println(data)
				insertPassengerTrip(params["passengertripid"], data)
				w.WriteHeader(http.StatusAccepted)
			} else {
				fmt.Println(err)
			}
		}
	} else if r.Method == "PUT" {
		if body, err := ioutil.ReadAll(r.Body); err == nil {
			var data PassengerTrip

			if err := json.Unmarshal(body, &data); err == nil {
				fmt.Println(data)
				w.WriteHeader(http.StatusAccepted)
			} else {
				fmt.Println(err)
			}
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func deluser(id string) (int64, error) {
	result, err := db.Exec("delete from User where UserID=?", id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func insertUser(id string, u User) {
	_, err := db.Exec("insert into user (UserID, FirstName, LastName, MobileNumber, Email, PW, DateOfCreation, AccountType, DriverLicense, CarPlateNumber, AccountStatus) values(?,?,?,?,?,?,?,?,?,?,?)", id, u.FirstName, u.LastName, u.MobileNumber, u.Email, u.PW, u.DateOfCreation, u.AccountType, u.DriverLicense, u.CarPLateNumber, u.AccountStatus)
	if err != nil {
		panic(err.Error())
	}
}

func updateUser(id string, u User) {
	_, err := db.Exec("update user set FirstName=?, LastName=?, MobileNumber=?, Email=?, PW=?, DateOfCreation=?, AccountType=?, DriverLicense=?, CarPlateNumber=?, AccountStatus=? where UserID=?", u.FirstName, u.LastName, u.MobileNumber, u.Email, u.PW, u.DateOfCreation, u.AccountType, u.DriverLicense, u.CarPLateNumber, u.AccountStatus, id)
	if err != nil {
		panic(err.Error())
	}
}

func insertCarPool(id string, c CarPool) {
	_, err := db.Exec("insert into Carpool (CarPoolID, UserID, PickUpLocation, AlternatePickUp, StartTravellingTime, AddressOfDestination, NumberOfPassengers, PoolStatus, NumberOfVacancies) values(?,?,?,?,?,?,?,?,?)", id, c.UserID, c.PickUpLocation, c.AlternatePickUp, c.StartTravellingTime, c.AddressOfDestination, c.NumberOfPassengers, c.PoolStatus, c.NumberOfVacancies)
	if err != nil {
		panic(err.Error())
	}
}

func updateCarPool(id string, c CarPool) {
	_, err := db.Exec("update Carpool set CarPoolID=?, UserID=?, PickUpLocation=?, AlternatePickUp=?, StartTravellingTime=?, AddressOfDestination=?, NumberOfPassengers=?, PoolStatus=?, NumberOfVacancies=? where CarPoolID=?", id, c.UserID, c.PickUpLocation, c.AlternatePickUp, c.StartTravellingTime, c.AddressOfDestination, c.NumberOfPassengers, c.PoolStatus, c.NumberOfVacancies, id)
	if err != nil {
		panic(err.Error())
	}
}

func insertPassengerTrip(id string, p PassengerTrip) {
	_, err := db.Exec("insert into PassengerTrip (TripID, UserID, CarPoolID, TripStatus) values(?,?,?,?)", id, p.UserID, p.CarPoolID, p.TripStatus)
	if err != nil {
		panic(err.Error())
	}
}
