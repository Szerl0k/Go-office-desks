UPDATE Desk SET Occupied = false WHERE true;

INSERT INTO User (Email, Name, Surname, Password, Role) VALUES ('test@gmail.com', 'imie', 'nazwisko', '$2a$10$Alh89NdVkHJfV1l/0DzTxOJa4604BL3CSO43lZrcR1SJOlYheHqmu', 'admin');

SELECT * FROM Desk;

SELECT * FROM Reservation;

SELECT * FROM User;


-- CREATE

-- Created by Vertabelo (http://vertabelo.com)
-- Last modification date: 2025-01-27 14:57:53.06

-- tables
-- Table: Desk
CREATE TABLE Desk (
    Id int  NOT NULL,
    Floor int  NOT NULL,
    Occupied bool  NOT NULL,
    Body varchar(50)  NOT NULL,
    CONSTRAINT Desk_pk PRIMARY KEY (Id)
);

-- Table: Reservation
CREATE TABLE Reservation (
    Id int  NOT NULL,
    `From` date  NOT NULL,
    `To` date  NOT NULL,
    Desk int  NOT NULL,
    User_Email varchar(40)  NOT NULL,
    CONSTRAINT Reservation_pk PRIMARY KEY (Id)
);

-- Table: User
CREATE TABLE User (
    Email varchar(40)  NOT NULL,
    Name varchar(40)  NOT NULL,
    Surname varchar(40)  NOT NULL,
    Password char(60)  NOT NULL,
    Role varchar(40) NOT NULL,
    CONSTRAINT User_pk PRIMARY KEY (Email)
);

-- foreign keys
-- Reference: Reservation_Desk (table: Reservation)
ALTER TABLE Reservation ADD CONSTRAINT Reservation_Desk FOREIGN KEY Reservation_Desk (Desk)
    REFERENCES Desk (Id);

-- Reference: Reservation_User (table: Reservation)
ALTER TABLE Reservation ADD CONSTRAINT Reservation_User FOREIGN KEY Reservation_User (User_Email)
    REFERENCES User (Email);

-- End of file.


-- DROP

-- Created by Vertabelo (http://vertabelo.com)
-- Last modification date: 2025-01-27 14:57:53.06

-- foreign keys
ALTER TABLE Reservation
    DROP FOREIGN KEY Reservation_Desk;

ALTER TABLE Reservation
    DROP FOREIGN KEY Reservation_User;

-- tables
DROP TABLE Desk;

DROP TABLE Reservation;

DROP TABLE User;

-- End of file.

