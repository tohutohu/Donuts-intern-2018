package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	db *sqlx.DB
)

type Movie struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

func main() {
	_db, err := sqlx.Open("mysql", "root:mysql@tcp(127.0.0.1:3306)/minich_local?parseTime=true")
	if err != nil {
		panic(err)
	}
	db = _db

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("please input command > ")
		if !scanner.Scan() {
			break
		}
		command := strings.TrimRight(scanner.Text(), "\r\n")
		switch command {
		case "exit":
			break
		case "read":
			movies, err := readMovies()
			if err != nil {
				panic(err)
			}
			for _, movie := range *movies {
				fmt.Printf("ID: %d Title: %s\n", movie.ID, movie.Title)
			}

		case "add":
			fmt.Print("please input title > ")
			if !scanner.Scan() {
				break
			}
			title := strings.TrimRight(scanner.Text(), "\r\n")
			res, err := insertMovie(title)
			if err != nil {
				panic(err)
			}
			id, _ := res.LastInsertId()
			fmt.Printf("%s is inserted. ID is %d\n", title, id)

		case "delete":
			fmt.Print("please input title > ")
			if !scanner.Scan() {
				break
			}
			title := strings.TrimRight(scanner.Text(), "\r\n")
			res, err := deleteMovieByTitle(title)
			if err != nil {
				panic(err)
			}
			if rows, _ := res.RowsAffected(); rows > 0 {
				fmt.Printf("%s is deleted\n", title)
			} else {
				fmt.Printf("%s is not deleted\n", title)
			}

		case "update":
			fmt.Print("please input ID > ")
			if !scanner.Scan() {
				break
			}
			id, err := strconv.Atoi(strings.TrimRight(scanner.Text(), "\r\n"))
			if err != nil {
				panic(err)
			}
			fmt.Print("please input title > ")
			if !scanner.Scan() {
				break
			}
			title := strings.TrimRight(scanner.Text(), "\r\n")
			res, err := updateTitle(id, title)
			if err != nil {
				panic(err)
			}
			if rows, _ := res.RowsAffected(); rows > 0 {
				fmt.Printf("%s is updated\n", title)
			} else {
				fmt.Printf("%s is not deleted\n", title)
			}

		}
	}
}

func insertMovie(title string) (sql.Result, error) {
	res, err := db.Exec("INSERT INTO `movie` (`title`) VALUES (?);", title)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func readMovies() (*[]Movie, error) {
	movies := &[]Movie{}
	err := db.Select(movies, "SELECT * FROM movie")
	if err != nil {
		return nil, err
	}
	return movies, nil
}

func deleteMovieByID(id int) (sql.Result, error) {
	res, err := db.Exec("DELETE FROM movie WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func deleteMovieByTitle(title string) (sql.Result, error) {
	res, err := db.Exec("DELETE FROM movie WHERE title = ?", title)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func updateTitle(id int, title string) (sql.Result, error) {
	res, err := db.Exec("UPDATE movie SET title = ? WHERE id = ?", title, id)
	if err != nil {
		return nil, err
	}
	return res, nil
}
