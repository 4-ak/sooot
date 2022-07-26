package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func NewDB() {
	db, openErr := sql.Open("sqlite3", "db.db")
	if openErr != nil {
		fmt.Print(openErr)
		panic(openErr)
	}
	DB = db
	if _, err := os.Stat("/path/to/whatever"); err != nil {
		CreateTable()
	}
}

const (
	createLectureTable = `
	CREATE TABLE IF NOT EXISTS "lecture" (
		"uid"	INTEGER NOT NULL UNIQUE,
		"name"	TEXT NOT NULL,
		"professor_name"	TEXT NOT NULL,
		"season"	INTEGER NOT NULL,
		"grade"	INTEGER NOT NULL,
		"credit"	INTEGER NOT NULL,
		"category"	NUMERIC NOT NULL,
		PRIMARY KEY("uid" AUTOINCREMENT)
	);
	`
	createReviewTable = `
	CREATE TABLE IF NOT EXISTS "review" (
		"uid"	INTEGER NOT NULL UNIQUE,
		"lecture_id"	INTEGER NOT NULL,
		"beneficial_point"	INTEGER NOT NULL CHECK("beneficial_point" BETWEEN 1 AND 5),
		"honey_point"	INTEGER NOT NULL CHECK("honey_point" BETWEEN 1 AND 5),
		"professor_point"	INTEGER NOT NULL CHECK("professor_point" BETWEEN 1 AND 5),
		"is_team"	INTEGER NOT NULL DEFAULT 0 CHECK("is_team" BETWEEN 0 AND 1),
		"is_presentation"	INTEGER NOT NULL DEFAULT 0 CHECK("is_presentation" BETWEEN 0 AND 1),
		PRIMARY KEY("uid" AUTOINCREMENT),
		FOREIGN KEY("lecture_id") REFERENCES "lecture"("uid") ON DELETE CASCADE
	);
	`
	createUserTable = `
	CREATE TABLE IF NOT EXISTS "user" (
		"uid"	INTEGER NOT NULL UNIQUE,
		"id"	TEXT NOT NULL UNIQUE,
		"pass"	TEXT NOT NULL,
		"is_cert"	INTEGER,
		PRIMARY KEY("uid" AUTOINCREMENT)
	);
	`
)

func CreateTable() {
	_, reviewErr := DB.Exec(createReviewTable)
	_, lectureErr := DB.Exec(createLectureTable)
	_, userErr := DB.Exec(createUserTable)
	if reviewErr != nil {
		fmt.Print(reviewErr)
	}
	if lectureErr != nil {
		fmt.Print(lectureErr)
	}
	if userErr != nil {
		fmt.Print(userErr)
	}
}
