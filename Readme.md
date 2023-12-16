# ETI Assignment 1

## Design Considerations of the microservices
When designing the microservices for the assignment, i decided to utilise 2 microservices. They are the console.go and the dbConnection.go files. 

### Database Microservice
As the file name suggests, the dbConnection.go file microservice's main purpose is to connect to the MySQL database and run commands to POST, PUT, and DELETE data rows according to what was needed in the assignment. After which, the dbConnection microservice creates API links which parses data to and from the database when calling the API, which returns  the results and/or queries retrieved from the database to wherever is calling the API.

### Console Microservice
As for the Console.go file microservice, it's main purpose was to act as the console application that users can utilise to access the Car Pool platform. It would provide the services stated in the assignment writeup displayed in the console by accessing the MySQL database when required by calling the API links generated in the dbConnection microservice through HTTP Requests.

Through the usage of both microservice, a console application that connects to the MySQL databse is created with the ability to perform all the stated requirements in the assignment writeup.

### Database (MySQL)
**User Table**
* UserID 
* FirstName 
* LastName
* MobileNumber
* Email
* PW (Password)
* DateOfCreation (Date when account was created)
* AccountType (Identifies if user is a passenger or car owner)
* DriverLicense
* CarPlateNumber
* AccountStatus (Indicates if the account is active)
  
**CarPool Table (Tracks published car pool trips)**
* CarPoolID
* UserID
* PickUpLocation
* AlternatePickUp
* StartTravellingTime
* AddressOfDestination
* NumberOfPassengers
* PoolStatus (Indicates if trip is awaiting, started, or cancelled)
* NumberOfVacancies

**PassengerTrip Table (Tracks the car pool trips passengers enrol in)**
* TripID
* UserID
* CarPoolID
* TripStatus

## Architecture diagram
![Architecture Diagram](https://github.com/ChunKaii/ETI_Assignment_1/blob/main/Architecture%20Diagram.png?raw=true))

## Instructions for setting up and running the microservices
1. Connect or create a connection in your MySQL Workbench
2. In the connection, create a new database with the relevant tables by running the SQL script in your connection
3. Open the dbConnection.go file
4. In line 51, update the username and password to what is configured in your MySQL connection, the port number the connection utilises, as well as the database name to "my_db"
5. Open the Console.go file
6. Repeat step 4 for lines 223, 257, 322, 557, 599, 632, 677, 780, and 869.
7. Save the changes made for both files.
8. Run the dbConnection.go file.
9. Allow pop-up notification to allow public and private networks to access the app.
10. Run the Console.go file to begin using the Car Pool platform.
