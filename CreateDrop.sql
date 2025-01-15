UPDATE Desks SET Occupied = false WHERE true;

SELECT * FROM Desks;

SELECT * FROM Reservations;

SELECT * FROM Users;


-- CREATE

-- Created by Vertabelo (http://vertabelo.com)
-- Last modification date: 2025-01-14 14:30:24.806

-- tables
-- Table: Desks
CREATE TABLE Desks (
    Id int NOT NULL AUTO_INCREMENT,
    Floor int NOT NULL,
    Occupied bool  NOT NULL,
    Body varchar(50)  NOT NULL,
    CONSTRAINT Desks_pk PRIMARY KEY (Id)
);

-- Table: Reservations
CREATE TABLE Reservations (
    User_Id int  NOT NULL,
    Desk_Id int  NOT NULL,
    `From` date  NOT NULL,
    `To` int  NOT NULL,
    CONSTRAINT Reservations_pk PRIMARY KEY (User_Id,Desk_Id)
);

-- Table: Users
CREATE TABLE Users (
    Id INT NOT NULL AUTO_INCREMENT,
    Name varchar(40)  NOT NULL,
    Surname varchar(40)  NOT NULL,
    CONSTRAINT Users_pk PRIMARY KEY (Id)
);

-- foreign keys
-- Reference: _Desks (table: Reservations)
ALTER TABLE Reservations ADD CONSTRAINT _Desks FOREIGN KEY _Desks (Desk_Id)
    REFERENCES Desks (Id);

-- Reference: _Users (table: Reservations)
ALTER TABLE Reservations ADD CONSTRAINT _Users FOREIGN KEY _Users (User_Id)
    REFERENCES Users (Id);



-- DROP

-- Created by Vertabelo (http://vertabelo.com)
-- Last modification date: 2025-01-14 14:30:24.806

-- foreign keys
ALTER TABLE Reservations
    DROP FOREIGN KEY _Desks;

ALTER TABLE Reservations
    DROP FOREIGN KEY _Users;

-- tables
DROP TABLE Desks;

DROP TABLE Reservations;

DROP TABLE Users;

-- End of file.

