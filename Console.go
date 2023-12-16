package main

import (
	"bufio"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
	"unicode"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
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

var input int
var currentSessionEmail string
var currentSessionPassword string
var currentSessionID string
var currentStatus string
var carPoolMap map[string]CarPool

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\nWelcome to Car-Pooling Platform\n-------------------------------")
	for {
		currentSessionEmail = ""
		currentSessionID = ""
		var checkLoginStatus = false
		var checkRegisterStatus = false
		currentStatus = "No Account"

		fmt.Print("\nLogin Menu\n----------\n1)Log into an account\n2)Register a new account\n0)Exit Platform\nChoose an option: ")
		fmt.Scanln(&input)
		if input == 0 {
			fmt.Println("Goodbye! See you again!")
			break
		} else if input == 1 || input == 2 { // checks if user chose input to login or to register account
			if input == 1 { // Runs if user chose to login
				fmt.Print(("Enter Email Address: "))
				fmt.Scanln(&currentSessionEmail)
				fmt.Print(("Enter password: "))
				fmt.Scanln(&currentSessionPassword)
				//Runs function to check entered credentials against database and returns boolean value
				var result = checkLoginCredentials(currentSessionEmail, currentSessionPassword)
				if result { //if true indicates login credentials matches database
					checkLoginStatus = true
					fmt.Println("Login successful!")
				} else { // if false indicates login credentials do not match any in database
					checkLoginStatus = false
				}
			} else { // Runs if user chose to register an account
				var firstName string
				fmt.Print("Enter your First Name: ")
				fmt.Scanln(&firstName)
				var lastName string
				fmt.Print("Enter your Last Name: ")
				fmt.Scanln(&lastName)
				var mobileNo int
				fmt.Print("Enter your Mobile Number: ")
				fmt.Scanln(&mobileNo)
				var email string
				fmt.Print("Enter your Email: ")
				fmt.Scanln(&email)
				var pass string
				fmt.Print("Enter your Password: ")
				fmt.Scanln(&pass)
				//Runs function to insert a new User record into database using the entered credentials
				checkRegisterStatus = RegisterAccount(firstName, lastName, mobileNo, email, pass)
				if checkRegisterStatus { //If true, indicates inserting of user record was successful
					fmt.Println("Registration successful!")
				}
			}
			if checkLoginStatus || checkRegisterStatus { // runs if login or register process were successful
				for {
					var loginInput int
					fmt.Print("\nMain page\n---------")
					//check current logged in user is passenger or car owner according to database
					if currentStatus == "Passenger" {
						fmt.Print("\n1)Update Account Information\n2)Register as a Car Owner\n3)Enrol for a trip\n4)View past trips taken\n5)Delete account\n0)Log out\nChoose an option: ")
						fmt.Scanln(&loginInput)
						if loginInput == 0 { //Logout option
							fmt.Println("Logged out successfully!")
							break
						} else if loginInput == 1 { //Update account information option
							UpdateAccount(loginInput)
						} else if loginInput == 2 { //Register as a car owner option (UPDATE FEATURE)
							UpdateAccount(loginInput)
						} else if loginInput == 3 { // Enrol for a trip option (INSERT)
							var enrolInput int
							fmt.Print("\nWould you like to:\n1)Browse all published trips\n2)Search for published trips\nChoose an option: ")
							fmt.Scanln(&enrolInput)
							if enrolInput == 1 { // Runs if user wants to browse all published trips
								var check = PrintAllCarPool() // Prints all records of car pool available to enrol in
								if check {
									//select and enrol
									var selectedCarPool int
									fmt.Print("\nPlease choose a published trip to enrol in: ")
									fmt.Scanln(&selectedCarPool)
									EnrolForATrip(selectedCarPool)
									//check and enrol
								} else {
									//Do nothing
								}
							} else if enrolInput == 2 { // Runs if user wants to search for a specific published trip
								var destination string
								fmt.Print("Enter your end destination to search for published trips: ")
								fmt.Scanln(&destination)
								destination, _ = reader.ReadString('\n')
								destination = trimRightSpace(destination)
								//Search using the string
								var che = PrintCarPoolUsingSubString(destination) // Prints all records of car pool available to enrol in if the destination value in database contains the string inputted by user
								if che {
									var selectedCarPool int
									fmt.Print("\nPlease choose a published trip to enrol in: ")
									fmt.Scanln(&selectedCarPool)
									EnrolForATrip(selectedCarPool)
								}
							} else {
								fmt.Println("Invalid input")
							}
						} else if loginInput == 4 { // View past trips taken option
							PrintTripsTaken()
						} else if loginInput == 5 { // Delete account option
							var checkk = DeleteAccount() // Checks if account is created over a year ago and deletes record from database if it is
							if checkk {                  // Reset current session ID and status if user deletes account
								currentSessionID = ""
								currentSessionEmail = ""
								currentSessionPassword = ""
								currentStatus = ""
								break
							}
						} else {
							fmt.Println("Invalid input. Please try again")
						}
					} else if currentStatus == "Car Owner" {
						fmt.Print("\n1)Update Account Information\n2)Publish a Car Pool Trip\n3)View published trips\n4)Delete account\n0)Log out\nChoose an option: ")
						fmt.Scanln(&loginInput)
						if loginInput == 0 { //Logout option
							fmt.Println("Logged out successfully!")
							break
						} else if loginInput == 1 { //Update account information option
							UpdateAccount(loginInput)
						} else if loginInput == 2 { //Publish a car pool trip option
							var resultOfPublish = PublishCarPoolTrip()
							if resultOfPublish == false { // Runs if function returns false indicating that there was an issue with inserting a new record into the database
								fmt.Println("Publication of car pool trip was unsuccessful. Please try again")
							}
						} else if loginInput == 3 { // View published trips option
							var checkForPublished = PrintCarPool() // Function prints all car pool published where the UserID matches the ID of the user currently logged in
							if checkForPublished {
								var carTrip int
								fmt.Print("\nSelect a car pool trip to start or cancel: ")
								fmt.Scanln(&carTrip)
								EditCarPool(carTrip) // Update car pool record in database to start or cancel if current time is 30 mins before start travelling time
							}
						} else if loginInput == 4 { // Delete account option
							var checkk = DeleteAccount()
							if checkk {
								currentSessionID = ""
								currentSessionEmail = ""
								currentSessionPassword = ""
								currentStatus = ""
								break
							}
						} else {
							fmt.Println("Invalid input. Please try again")
						}
					}
				}
			} else if checkLoginStatus == false && input == 1 {
				fmt.Println("Invalid email or password. Please try again.")
			} else if checkRegisterStatus == false && input == 2 {
				fmt.Print("Error in registration of account. Please try again.")
			} else {
				fmt.Println("An error occurred. Please try again.")
			}

		} else {
			fmt.Println("Invalid option. Please try again")
		}

	}
}

// Checks login credentials against all the user records in database and returns boolean value
func checkLoginCredentials(email string, password string) bool {
	db, _ := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/my_db") // Connect to database
	results, err := db.Query("select * from User")                        // Retrieve all user records
	var id string
	if err != nil {
		panic(err.Error())
	}
	for results.Next() {
		var u User
		err = results.Scan(&id, &u.FirstName, &u.LastName, &u.MobileNumber, &u.Email, &u.PW, &u.DateOfCreation, &u.AccountType, &u.DriverLicense, &u.CarPLateNumber, &u.AccountStatus)
		if err != nil {
			panic(err.Error())
		} else {
			if email == u.Email && password == u.PW { // Checks if the user record matches with the login credentials entered
				currentSessionID = id
				currentStatus = u.AccountType
				return true // Returns true if matches
			} else {
			}
		}
	}
	return false // Returns false if no record matches
}

// Generate unique id
func GenerateUUID() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "0", err
	}
	return strconv.Itoa(int(id.ID())), nil
}

// Checking database for any repeated user data and register new user into database
func RegisterAccount(fName string, lName string, mobile int, email string, pass string) bool {
	db, _ := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/my_db")
	var id string
	results, err := db.Query("select * from User")
	if err != nil {
		panic(err.Error())
	}
	for results.Next() {
		var u User
		err = results.Scan(&id, &u.FirstName, &u.LastName, &u.MobileNumber, &u.Email, &u.PW, &u.DateOfCreation, &u.AccountType, &u.DriverLicense, &u.CarPLateNumber, &u.AccountStatus)
		if err != nil {
			panic(err.Error())
		} else {
			if email == u.Email || mobile == u.MobileNumber { // Check if any exising user credentials is same as the new user
				return false
			} else {
			}
		}
	}
	//Create new User object with the new user credentials
	var userP User
	userP.FirstName = fName
	userP.LastName = lName
	userP.MobileNumber = mobile
	userP.Email = email
	userP.PW = pass
	userP.DateOfCreation = time.Now().Format("2006-01-02")
	userP.AccountType = "Passenger"
	userP.DriverLicense = ""
	userP.CarPLateNumber = ""
	userP.AccountStatus = true

	postBody, _ := json.Marshal(userP)
	resBody := bytes.NewBuffer(postBody)

	var userID, erro = GenerateUUID()
	if erro != nil {
		fmt.Println(erro)
		fmt.Println("Error in system please try again.")
	} else {
		client := &http.Client{}
		if req, err := http.NewRequest(http.MethodPost, "http://localhost:5000/api/v1/users/"+userID, resBody); err == nil { //Runs POST method using user api link from DBConnection Go File
			if res, err := client.Do(req); err == nil {
				if res.StatusCode == 202 { //Status code indicating success
					fmt.Println("Registration successfully")
					currentSessionID = id       // Update current session ID
					currentStatus = "Passenger" // Update current session status
					return true                 // Returns true if inserting record was successful
				} else if res.StatusCode == 409 { // Status code indicating failure
					fmt.Println("Error - userid", userID, "exists")
					return false // Returns false if inserting record was unsuccessful
				}
			} else {
				fmt.Println(2, err)
				return false
			}
		} else {
			fmt.Println(3, err)
			return false
		}
	}
	return false
}

// Updates a user record
func UpdateAccount(input int) {
	db, _ := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/my_db")
	var u User
	var id string
	result := db.QueryRow("select * from User where UserID=?", currentSessionID) // Retrieves user records where UserID matches current session UserID
	err := result.Scan(&id, &u.FirstName, &u.LastName, &u.MobileNumber, &u.Email, &u.PW, &u.DateOfCreation, &u.AccountType, &u.DriverLicense, &u.CarPLateNumber, &u.AccountStatus)
	if err == sql.ErrNoRows {
		fmt.Println("sql row error")
	}
	var userP User

	var updateInput int
	var updatedString string
	if input == 1 { //Runs if user chose to update account information
		// Request user to update a credential of theirs
		fmt.Print("\nChoose an option for the information you would like to update\n-------------------------------------------------------------\n1)First Name\n2)Last Name\n3)Mobile Number\n4)Email\nChoose an option: ")
		fmt.Scanln(&updateInput)
		if updateInput == 1 || updateInput == 2 || updateInput == 3 || updateInput == 4 {
			if updateInput == 1 {
				fmt.Print("Enter your new first name: ")
				fmt.Scanln(&updatedString)
				//Create new User object with updated credentials
				userP.FirstName = updatedString
				userP.LastName = u.LastName
				userP.MobileNumber = u.MobileNumber
				userP.Email = u.Email
				userP.PW = u.PW
				userP.DateOfCreation = u.DateOfCreation
				userP.AccountType = u.AccountType
				userP.DriverLicense = u.DriverLicense
				userP.CarPLateNumber = u.CarPLateNumber
				userP.AccountStatus = u.AccountStatus
			} else if updateInput == 2 {
				fmt.Print("Enter your new last name: ")
				fmt.Scanln(&updatedString)
				userP.FirstName = u.FirstName
				userP.LastName = updatedString
				userP.MobileNumber = u.MobileNumber
				userP.Email = u.Email
				userP.PW = u.PW
				userP.DateOfCreation = u.DateOfCreation
				userP.AccountType = u.AccountType
				userP.DriverLicense = u.DriverLicense
				userP.CarPLateNumber = u.CarPLateNumber
				userP.AccountStatus = u.AccountStatus
			} else if updateInput == 3 {
				var updatemobile int
				fmt.Print("Enter your new mobile number: ")
				fmt.Scanln(&updatemobile)
				userP.FirstName = u.FirstName
				userP.LastName = u.LastName
				userP.MobileNumber = updatemobile
				userP.Email = u.Email
				userP.PW = u.PW
				userP.DateOfCreation = u.DateOfCreation
				userP.AccountType = u.AccountType
				userP.DriverLicense = u.DriverLicense
				userP.CarPLateNumber = u.CarPLateNumber
				userP.AccountStatus = u.AccountStatus
			} else if updateInput == 4 {
				fmt.Print("Enter your new email: ")
				fmt.Scanln(&updatedString)
				userP.FirstName = u.FirstName
				userP.LastName = u.LastName
				userP.MobileNumber = u.MobileNumber
				userP.Email = updatedString
				userP.PW = u.PW
				userP.DateOfCreation = u.DateOfCreation
				userP.AccountType = u.AccountType
				userP.DriverLicense = u.DriverLicense
				userP.CarPLateNumber = u.CarPLateNumber
				userP.AccountStatus = u.AccountStatus
			}
			postBody, _ := json.Marshal(userP)
			client := &http.Client{}
			if req, err := http.NewRequest(http.MethodPut, "http://localhost:5000/api/v1/users/"+currentSessionID, bytes.NewBuffer(postBody)); err == nil { //Use PUT method of the user api link with the newly updated user object
				if res, err := client.Do(req); err == nil {
					if res.StatusCode == 202 {
						fmt.Println("Updated successfully")
						currentStatus = userP.AccountType // Updates current session status
					} else if res.StatusCode == 409 {
						fmt.Println("Error")
					}
				} else {
					fmt.Println(2, err)
				}
			} else {
				fmt.Println(3, err)
			}
		} else {
			fmt.Println("Invalid option. Please try again.")
		}
	} else if input == 2 { // Runs if user chose to register as a car owner
		var driverLicenseNo string
		var carPlateNo string
		var checkExist = false
		fmt.Print("\nEnter your driver license number: ")
		fmt.Scanln(&driverLicenseNo)
		fmt.Print("Enter your car plate number: ")
		fmt.Scanln(&carPlateNo)
		var id string
		results, err := db.Query("select * from User")
		if err != nil {
			panic(err.Error())
		}
		for results.Next() {
			var u User
			err = results.Scan(&id, &u.FirstName, &u.LastName, &u.MobileNumber, &u.Email, &u.PW, &u.DateOfCreation, &u.AccountType, &u.DriverLicense, &u.CarPLateNumber, &u.AccountStatus)
			if err != nil {
				panic(err.Error())
			} else {
				if driverLicenseNo == u.DriverLicense || carPlateNo == u.CarPLateNumber { //Checks for any overlapping driver license or car plate number
					fmt.Println("The inputted driver license or car plate number has already been registered with an exisiting account.")
					checkExist = true
					break
				} else {
				}
			}
		}
		if checkExist == false { // Runs if there are no overlapping driver license nor car plate number
			//Create new User object with driver license and car plate number entered
			userP.FirstName = u.FirstName
			userP.LastName = u.LastName
			userP.MobileNumber = u.MobileNumber
			userP.Email = u.Email
			userP.PW = u.PW
			userP.DateOfCreation = u.DateOfCreation
			userP.AccountType = "Car Owner"
			userP.DriverLicense = driverLicenseNo
			userP.CarPLateNumber = carPlateNo
			userP.AccountStatus = u.AccountStatus
			postBody, _ := json.Marshal(userP)
			client := &http.Client{}
			if req, err := http.NewRequest(http.MethodPut, "http://localhost:5000/api/v1/users/"+currentSessionID, bytes.NewBuffer(postBody)); err == nil { //Use PUT method of user api link with newly updated User object
				if res, err := client.Do(req); err == nil {
					if res.StatusCode == 202 {
						fmt.Println("Successfully registered as a car owner!")
						currentStatus = userP.AccountType // Update current session status
					} else if res.StatusCode == 409 {
						fmt.Println("Error")
					}
				} else {
					fmt.Println(2, err)
				}
			} else {
				fmt.Println(3, err)
			}
		}
	}
}

// Function to remove excess space on the right of the string when reading input
func trimRightSpace(s string) string {
	return strings.TrimRightFunc(s, unicode.IsSpace)
}

// Function to create new CarPool record in database
func PublishCarPoolTrip() bool {
	reader := bufio.NewReader(os.Stdin)
	var car CarPool
	var pickupL string
	var altpickupL string
	var startTravelTime string
	var addrOfDestination string
	var noOfPassengers int
	fmt.Print("\nEnter your pool's pick up location: ")
	pickupL, _ = reader.ReadString('\n')
	pickupL = trimRightSpace(pickupL)
	fmt.Print("\nEnter your pool's alternative pick up location: ")
	altpickupL, _ = reader.ReadString('\n')
	altpickupL = trimRightSpace(altpickupL)
	fmt.Print("\nEnter your pool's start travelling time (in the format of 'YYYY-MM-DD HH:mm'): ")
	startTravelTime, _ = reader.ReadString('\n')
	startTravelTime = trimRightSpace(startTravelTime)
	fmt.Print("\nEnter your pool's address of destination: ")
	addrOfDestination, _ = reader.ReadString('\n')
	addrOfDestination = trimRightSpace(addrOfDestination)
	fmt.Print("\nEnter your vehicle's capacity (number of passengers): ")
	fmt.Scanln(&noOfPassengers)
	date, error := time.ParseInLocation("2006-01-02 15:04", startTravelTime, time.Local)
	if error != nil {
		fmt.Println("Start travelling time inserted in the wrong format please try again")
		return false
	}
	if time.Now().After(date) { // Checks if start travelling time has already passed
		fmt.Println("Please input a start travelling time that has yet to passed.")
		return false
	} else {
		var carPoolId, erro = GenerateUUID() //Generate new car pool id
		if erro != nil {
			fmt.Println(erro)
			fmt.Println("Error in system please try again.")
			return false
		} else {
			//Create new CarPool object
			car.UserID = currentSessionID
			car.PickUpLocation = pickupL
			car.AlternatePickUp = altpickupL
			car.StartTravellingTime = startTravelTime
			car.AddressOfDestination = addrOfDestination
			car.NumberOfPassengers = noOfPassengers
			car.PoolStatus = "Awaiting"
			car.NumberOfVacancies = noOfPassengers
			postBody, _ := json.Marshal(car)
			resBody := bytes.NewBuffer(postBody)
			client := &http.Client{}
			if req, err := http.NewRequest(http.MethodPost, "http://localhost:5000/api/v1/carpool/"+carPoolId, resBody); err == nil { //Use POST method of the carpool api link with new CarPool object
				if res, err := client.Do(req); err == nil {
					if res.StatusCode == 202 { //Status code indicating success
						fmt.Println("Published successfully!")
						return true
					} else if res.StatusCode == 409 { //Status code indicating failure
						fmt.Println("Error in publishing car pool trip. PLease try again.")
						return false
					}
				} else {
					fmt.Println(2, err)
					return false
				}
			} else {
				fmt.Println(3, err)
				return false
			}
		}
	}
	return false
}

// Function to print all CarPool records from database if they were published by the current user
func PrintCarPool() bool {
	//Emptying map
	for k := range carPoolMap {
		delete(carPoolMap, k)
	} //Empties the CarPool map
	carPoolMap = make(map[string]CarPool)
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0) //Used to align the text when displaying records later
	db, _ := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/my_db")
	results, err := db.Query("select * from carpool where UserID=?", currentSessionID)
	var id string
	var listing = 1
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("\nPublished Car Pools\n-------------------")
	fmt.Fprintln(w, "   \tPick Up Location\tAlternative Pick Up Location\tStart Travelling Time\tAddress of Destination\tNumber Of Vacancies\tPool Status")
	for results.Next() {
		// fmt.Print(results)
		var c CarPool
		err = results.Scan(&id, &c.UserID, &c.PickUpLocation, &c.AlternatePickUp, &c.StartTravellingTime, &c.AddressOfDestination, &c.NumberOfPassengers, &c.PoolStatus, &c.NumberOfVacancies)
		date, error := time.ParseInLocation("2006-01-02 15:04", c.StartTravellingTime, time.Local)
		if error != nil {
			fmt.Println("Error")
		}
		if time.Now().After(date) { //Checks if the published trip has already passed
		} else {
			if c.PoolStatus == "Cancelled" || c.PoolStatus == "Started" { // Checks if the published trip has started or cancelled
			} else {
				fmt.Fprintln(w, (strconv.Itoa(listing) + ")\t" + c.AddressOfDestination + "\t" + c.AlternatePickUp + "\t" + c.StartTravellingTime + "\t" + c.AddressOfDestination + "\t" + strconv.Itoa(c.NumberOfVacancies) + "\t" + c.PoolStatus))
				listing++
				carPoolMap[id] = c //Updates the map of carpool
			}
		}
	}
	w.Flush()         // Print out the list of CarPool trips
	if listing == 1 { //Checks if any car pool were printed
		fmt.Println("You have not published any car pool trips yet")
		return false
	} else {
		return true
	}
}

// Function to print all car pool trips taken by the current user
func PrintTripsTaken() {
	var cid string
	var pid string
	var listing = 1
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	db, _ := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/my_db")
	//Database query retrieves all car pool records if for all PassengerTrip objects with a matching UserID and CarPoolID which is then displaying in descending order according to the start travelling time to display in reverse chronological order (newest to oldest)
	var qry = "SELECT * FROM CarPool INNER JOIN PassengerTrip ON PassengerTrip.CarPoolID = CarPool.CarPoolID WHERE PassengerTrip.UserID = " + currentSessionID + " AND CAST(CarPool.StartTravellingTime as datetime) < NOW() ORDER BY CAST(CarPool.StartTravellingTime as datetime) DESC"
	results, err := db.Query(qry)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("\nPast Trips Taken\n-------------------")
	fmt.Fprintln(w, "   \tPick Up Location\tAlternative Pick Up Location\tStart Travelling Time\tAddress of Destination")
	for results.Next() {
		var c CarPool
		var pt PassengerTrip
		err = results.Scan(&cid, &c.UserID, &c.PickUpLocation, &c.AlternatePickUp, &c.StartTravellingTime, &c.AddressOfDestination, &c.NumberOfPassengers, &c.PoolStatus, &c.NumberOfVacancies, &pid, pt.UserID, pt.CarPoolID, pt.TripStatus)
		if err != nil {
			fmt.Print(err)
		}
		fmt.Fprintln(w, (strconv.Itoa(listing) + ")\t" + c.AddressOfDestination + "\t" + c.AlternatePickUp + "\t" + c.StartTravellingTime + "\t" + c.AddressOfDestination))
		listing++
	}
	w.Flush()
	if listing == 1 {
		fmt.Println("You have not taken any car pool trips yet")
	}
}

// Prints all car pool trips in the database
func PrintAllCarPool() bool {
	//Emptying map of car pool
	for k := range carPoolMap {
		delete(carPoolMap, k)
	}
	carPoolMap = make(map[string]CarPool)
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	db, _ := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/my_db")
	results, err := db.Query("select * from carpool")
	var id string
	//var count int
	var listing = 1
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("\nPublished Car Pools\n-------------------")
	fmt.Fprintln(w, "   \tPick Up Location\tAlternative Pick Up Location\tStart Travelling Time\tAddress of Destination\tNumber Of Vacancies\tPool Status")
	for results.Next() {
		// fmt.Print(results)
		var c CarPool
		err = results.Scan(&id, &c.UserID, &c.PickUpLocation, &c.AlternatePickUp, &c.StartTravellingTime, &c.AddressOfDestination, &c.NumberOfPassengers, &c.PoolStatus, &c.NumberOfVacancies)
		date, error := time.ParseInLocation("2006-01-02 15:04", c.StartTravellingTime, time.Local)
		if error != nil {
			fmt.Println("Error")
		}
		if time.Now().After(date) { //Checks if trip has passed
		} else {
			if c.PoolStatus == "Cancelled" || c.PoolStatus == "Started" { //Checks if trip is cancelled or started
			} else {
				fmt.Fprintln(w, (strconv.Itoa(listing) + ")\t" + c.PickUpLocation + "\t" + c.AlternatePickUp + "\t" + c.StartTravellingTime + "\t" + c.AddressOfDestination + "\t" + strconv.Itoa(c.NumberOfVacancies) + "\t" + c.PoolStatus))
				listing++
				carPoolMap[id] = c
			}
		}
	}
	w.Flush()
	if listing == 1 { //Checks if any records were printed
		fmt.Println("There are currently no published car trips available to enrol in")
		return false
	} else {
		return true
	}
}

// Prints all car pool trips if the destination contains a specific substring
func PrintCarPoolUsingSubString(destination string) bool {
	//Emptying map of car pool
	for k := range carPoolMap {
		delete(carPoolMap, k)
	}
	carPoolMap = make(map[string]CarPool)
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	db, _ := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/my_db")
	var qry = "select * from carpool c WHERE c.AddressOfDestination LIKE '%" + destination + "%'" //Checks for car pool trips where the address of destination has contains the string entered by user
	results, err := db.Query(qry)
	var id string
	var listing = 1
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("\nPublished Car Pools\n-------------------")
	fmt.Fprintln(w, "   \tPick Up Location\tAlternative Pick Up Location\tStart Travelling Time\tAddress of Destination\tNumber Of Vacancies\tPool Status")
	for results.Next() {
		var c CarPool
		err = results.Scan(&id, &c.UserID, &c.PickUpLocation, &c.AlternatePickUp, &c.StartTravellingTime, &c.AddressOfDestination, &c.NumberOfPassengers, &c.PoolStatus, &c.NumberOfVacancies)
		date, error := time.ParseInLocation("2006-01-02 15:04", c.StartTravellingTime, time.Local)
		if error != nil {
			fmt.Println("Error")
		}
		if time.Now().After(date) {
		} else {
			if c.PoolStatus == "Cancelled" || c.PoolStatus == "Started" {
			} else {
				fmt.Fprintln(w, (strconv.Itoa(listing) + ")\t" + c.AddressOfDestination + "\t" + c.AlternatePickUp + "\t" + c.StartTravellingTime + "\t" + c.AddressOfDestination + "\t" + strconv.Itoa(c.NumberOfVacancies) + "\t" + c.PoolStatus))
				listing++
				carPoolMap[id] = c
			}
		}
	}
	w.Flush()
	if listing == 1 {
		fmt.Println("There are currently no published car trips available to enrol in")
		return false
	} else {
		return true
	}
}

// Function to update car pool trip object to start or cancel
func EditCarPool(index int) {
	var count = 1
	for k, v := range carPoolMap { //Map of carpool containing carpool objects that were displayed to user
		if index == count {
			var tripInput int
			fmt.Print("\n1)Start Trip\n2)Cancel Trip\nChoose an option: ")
			fmt.Scanln(&tripInput)
			date, error := time.ParseInLocation("2006-01-02 15:04", v.StartTravellingTime, time.Local)
			if error != nil {
				fmt.Println("Error")
			}
			if date.Sub(time.Now()).Minutes() < 30 { //Checks if current time is before 30 minutes from start travelling time
				if tripInput == 1 {
					//Start the trip
					v.PoolStatus = "Started"
					postBody, _ := json.Marshal(v)
					client := &http.Client{}
					if req, err := http.NewRequest(http.MethodPut, "http://localhost:5000/api/v1/carpool/"+k, bytes.NewBuffer(postBody)); err == nil { //Use PUT api link of carpool to update car pool trip
						if res, err := client.Do(req); err == nil {
							if res.StatusCode == 202 {
								fmt.Println("Successfully updated the trip status")
							} else if res.StatusCode == 409 {
								fmt.Println("Error in updating the trip status")
							}
						} else {
							fmt.Println(2, err)
						}
					} else {
						fmt.Println(3, err)
					}
					break
				} else if tripInput == 2 {
					//Cancel trip
					v.PoolStatus = "Cancelled"
					postBody, _ := json.Marshal(v)
					client := &http.Client{}
					if req, err := http.NewRequest(http.MethodPut, "http://localhost:5000/api/v1/carpool/"+k, bytes.NewBuffer(postBody)); err == nil {
						if res, err := client.Do(req); err == nil {
							if res.StatusCode == 202 {
								fmt.Println("Successfully updated the trip status")
							} else if res.StatusCode == 409 {
								fmt.Println("Error in updating the trip status")
							}
						} else {
							fmt.Println(2, err)
						}
					} else {
						fmt.Println(3, err)
					}
					break
				} else {
					fmt.Println("\nInvalid input. Please try again.")
				}
			} else {
				fmt.Println("\nYou are unable to cancel nor start the trip as it is more than 30 minutes before the scheduled travelling time")
				s := strconv.FormatFloat(date.Sub(time.Now()).Minutes(), 'g', 5, 64) //Displays time before the start travelling time of car pool object
				fmt.Printf("It is currently %v minutes from the trip's scheduled time\n", s)
			}
		} else {
			count++
		}
	}
}

// Function to enrol current user for a car pool trip
func EnrolForATrip(input int) {
	db, _ := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/my_db")
	var count = 1
	var id string
	var checkForConflicts = false
	for k, v := range carPoolMap { //Iterate through map of carpool objects displayed to user
		if count == input {
			results, err := db.Query("select * from PassengerTrip where UserID=?", currentSessionID)
			if err != nil {
				panic(err.Error())
			}
			for results.Next() {
				var p PassengerTrip
				err = results.Scan(&id, &p.UserID, &p.CarPoolID, &p.TripStatus)
				resultss, err := db.Query("select * from CarPool where CarPoolID=?", p.CarPoolID)
				if err != nil {
					panic(err.Error())
				}
				for resultss.Next() {
					var c CarPool
					err = results.Scan(&id, &c.UserID, &c.PickUpLocation, &c.AlternatePickUp, &c.StartTravellingTime, &c.AddressOfDestination, &c.NumberOfPassengers, &c.PoolStatus, &c.NumberOfVacancies)
					if c.StartTravellingTime == v.StartTravellingTime { //Check if the car pool trip user is enrolling in has time conflict with any other car pool trip they have enrolled in
						fmt.Println("This car pool conflicts with the timing of a pre-existing car pool trip you had already enrolled in")
						checkForConflicts = true
					} else {
					}
				}
			}
			if checkForConflicts == false {
				//create passengertrip data
				if v.NumberOfVacancies == 0 {
					fmt.Println("This car pool no longer has any vacancies. Please enrol in another.")
				} else {
					var ptId, erro = GenerateUUID()
					if erro != nil {
						fmt.Println(erro)
						fmt.Println("Error in system please try again.")
						// return false
					} else {

					}
					//Create new passengertrip object
					var pt PassengerTrip
					pt.UserID = currentSessionID
					pt.CarPoolID = k
					pt.TripStatus = v.PoolStatus
					postBody, _ := json.Marshal(pt)
					resBody := bytes.NewBuffer(postBody)
					client := &http.Client{}
					if req, err := http.NewRequest(http.MethodPost, "http://localhost:5000/api/v1/passengertrip/"+ptId, resBody); err == nil { //Use POST method of the passengertrip api link to insert new record
						if res, err := client.Do(req); err == nil {
							if res.StatusCode == 202 {
								fmt.Println("Published successfully!")
								//update car pool data
								v.NumberOfVacancies = v.NumberOfVacancies - 1
								postBody, _ := json.Marshal(v)
								client := &http.Client{}
								if req, err := http.NewRequest(http.MethodPut, "http://localhost:5000/api/v1/carpool/"+k, bytes.NewBuffer(postBody)); err == nil {
									if res, err := client.Do(req); err == nil {
										if res.StatusCode == 202 {
											fmt.Println("Successfully updated the number of vacancies")
										} else if res.StatusCode == 409 {
											fmt.Println("Error in updating the number of vacancies")
										}
									} else {
										fmt.Println(2, err)
									}
								} else {
									fmt.Println(3, err)
								}
							} else if res.StatusCode == 409 {
								fmt.Println("Error in enrolling in car pool trip. PLease try again.")
							}
						} else {
							fmt.Println(2, err)
						}
					} else {
						fmt.Println(3, err)
					}
				}
			} else {
			}
		} else {
			count++
		}
	}
}

// Function to delete user record
func DeleteAccount() bool {
	db, _ := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/my_db")
	var u User
	var id string
	result := db.QueryRow("select * from User where UserID=?", currentSessionID)
	err := result.Scan(&id, &u.FirstName, &u.LastName, &u.MobileNumber, &u.Email, &u.PW, &u.DateOfCreation, &u.AccountType, &u.DriverLicense, &u.CarPLateNumber, &u.AccountStatus)
	if err == sql.ErrNoRows {
		fmt.Println("sql row error")
		return false
	}
	date, error := time.ParseInLocation("2006-01-02", u.DateOfCreation, time.Local) //Convert string to date
	if error != nil {
		fmt.Println("Error")
		return false
	}
	diff := time.Now().Sub(date).Hours() / 24 / 365 //Checks if the account was created more than a year ago
	if diff < 1 {
		fmt.Println("Your account is not over 1 year old hence you are unable to delete just yet.")
		return false
	} else {
		//localhost to delete record
		client := &http.Client{}
		if req, err := http.NewRequest(http.MethodDelete, "http://localhost:5000/api/v1/users/"+currentSessionID, nil); err == nil { //Deletes user record
			if res, err := client.Do(req); err == nil {
				if res.StatusCode == 200 {
					fmt.Println("Account deleted successfully")
					return true
				} else if res.StatusCode == 404 {
					fmt.Println("Error in deleting Account")
					return false
				}
			} else {
				fmt.Println(2, err)
				return false
			}
		} else {
			fmt.Println(3, err)
			return false
		}
	}
	return false
}
