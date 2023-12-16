CREATE database my_db;
USE my_db;

CREATE TABLE User (
UserID varchar(255) NOT NULL PRIMARY KEY,
FirstName varchar(255),
LastName varchar(255),
MobileNumber bigint,
Email varchar(255),
PW varchar(255),
DateOfCreation varchar(255),
AccountType varchar(255),
DriverLicense varchar(255),
CarPlateNumber varchar(255),
AccountStatus boolean
);

CREATE TABLE CarPool (
CarPoolID varchar(255) NOT NULL PRIMARY KEY,
UserID varchar(255),
PickUpLocation varchar(255),
AlternatePickUp varchar(255),
StartTravellingTime varchar(255),
AddressOfDestination varchar(255),
NumberOfPassengers int,
PoolStatus varchar(255),
NumberOfVacancies int
);

CREATE TABLE PassengerTrip (
TripID varchar(255) NOT NULL PRIMARY KEY,
UserID varchar(255),
CarPoolID varchar(255),
TripStatus varchar(255)
);