package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/flosch/pongo2"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
)

type Movie struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Thumbnail string `json:"thumbnail"`
	Path      string `json:"path"`
}

type Comment struct {
	ID        int       `json:"id" db:"id"`
	MovieID   int       `json:"movie_id" db:"movie_id"`
	Text      string    `json:"text" db:"text"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

var (
	db *sqlx.DB
)

func main() {
	rand.Seed(time.Now().UnixNano())
	_db, err := sqlx.Open("mysql", "root:mysql@tcp(127.0.0.1:3306)/donuts?parseTime=true")

	if err != nil {
		panic(err)
	}

	db = _db
	e := echo.New()
	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	e.GET("/movies/register", func(c echo.Context) error {
		temp, err := pongo2.FromFile("./template/index.html")
		if err != nil {
			return err
		}

		return temp.ExecuteWriter(pongo2.Context{}, c.Response())
	})

	e.POST("/movies", postMovies)
	e.GET("/movies", getMovies)
	e.GET("/", getMovies)
	e.GET("/movies/:id", getMovie)
	e.POST("/movies/:id", deleteMovie)

	e.GET("/movies/:id/edit", getMovieEdit)
	e.POST("/movies/:id/update", updateMovie)
	e.POST("/movies/:id/comments", postComment)

	e.Static("/static/thumbnails", "./thumbnails")
	e.Static("/static/movies", "./movies")

	e.Start(":1323")
}

func postComment(c echo.Context) error {
	id := c.Param("id")
	movie := Movie{}
	if err := db.Get(&movie, "SELECT * FROM movie WHERE id = ?", id); err != nil {
		return err
	}

	if movie.ID == 0 {
		temp, err := pongo2.FromFile("./template/404.html")
		if err != nil {
			return err
		}
		c.Response().Status = http.StatusNotFound
		return temp.ExecuteWriter(pongo2.Context{}, c.Response())
	}

	text := c.FormValue("comment")

	db.Exec("INSERT INTO comment (text, movie_id) VALUES (?, ?)", text, id)
	return c.Redirect(http.StatusSeeOther, "/movies/"+id)

}

func updateMovie(c echo.Context) error {
	id := c.Param("id")
	movie := Movie{}
	if err := db.Get(&movie, "SELECT * FROM movie WHERE id = ?", id); err != nil {
		return err
	}

	if movie.ID == 0 {
		temp, err := pongo2.FromFile("./template/404.html")
		if err != nil {
			return err
		}
		c.Response().Status = http.StatusNotFound
		return temp.ExecuteWriter(pongo2.Context{}, c.Response())
	}

	title := c.FormValue("title")
	filename := strings.Split(movie.Path, "/")[2]

	file, err := c.FormFile("thumbnail")
	if err != nil {
		return err
	}
	if file.Size > 0 {
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		// Destination
		dst, err := os.Create("./thumbnails/" + filename)
		if err != nil {
			return err
		}
		defer dst.Close()

		// Copy
		if _, err = io.Copy(dst, src); err != nil {
			return err
		}

	}

	movieFile, err := c.FormFile("movie")
	if err != nil {
		return err
	}
	if movieFile.Size > 0 {
		movieSrc, err := movieFile.Open()
		if err != nil {
			return err
		}
		defer movieSrc.Close()

		// Destination
		movieDst, err := os.Create("./movies/" + filename)
		if err != nil {
			return err
		}
		defer movieDst.Close()

		// Copy
		if _, err = io.Copy(movieDst, movieSrc); err != nil {
			return err
		}

	}

	_, err = db.Exec("UPDATE movie SET title = ? WHERE id = ?", title, movie.ID)
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, "/")
}

func getMovieEdit(c echo.Context) error {
	id := c.Param("id")
	movie := Movie{}
	if err := db.Get(&movie, "SELECT * FROM movie WHERE id = ?", id); err != nil {
		return err
	}

	if movie.ID == 0 {
		temp, err := pongo2.FromFile("./template/404.html")
		if err != nil {
			return err
		}
		c.Response().Status = http.StatusNotFound
		return temp.ExecuteWriter(pongo2.Context{}, c.Response())
	}

	temp, err := pongo2.FromFile("./template/edit.html")
	if err != nil {
		return err
	}

	return temp.ExecuteWriter(pongo2.Context{"movie": movie}, c.Response())
}

func deleteMovie(c echo.Context) error {
	id := c.Param("id")
	movie := Movie{}
	if err := db.Get(&movie, "SELECT * FROM movie WHERE id = ?", id); err != nil {
		return err
	}

	if movie.ID == 0 {
		temp, err := pongo2.FromFile("./template/404.html")
		if err != nil {
			return err
		}
		c.Response().Status = http.StatusNotFound

		return temp.ExecuteWriter(pongo2.Context{}, c.Response())
	}

	if err := os.Remove("." + movie.Path); err != nil {
		return err
	}
	if err := os.Remove("." + movie.Thumbnail); err != nil {
		return err
	}

	if _, err := db.Exec("DELETE FROM movie WHERE id = ?", movie.ID); err != nil {
		return err
	}

	return c.Redirect(http.StatusMovedPermanently, "/")

}

func getMovie(c echo.Context) error {
	id := c.Param("id")
	movie := Movie{}
	if err := db.Get(&movie, "SELECT * FROM movie WHERE id = ?", id); err != nil {
		return err
	}

	if movie.ID == 0 {
		temp, err := pongo2.FromFile("./template/404.html")
		if err != nil {
			return err
		}
		c.Response().Status = http.StatusNotFound

		return temp.ExecuteWriter(pongo2.Context{}, c.Response())
	}

	comments := []Comment{}
	if err := db.Select(&comments, "SELECT comment.id, comment.movie_id, comment.text, comment.created_at FROM comment LEFT JOIN movie ON comment.movie_id = movie.id WHERE movie.id = ?", movie.ID); err != nil {
		fmt.Println(err)
	}
	fmt.Println(movie.ID)
	fmt.Println(comments)

	temp, err := pongo2.FromFile("./template/movie.html")
	if err != nil {
		return err
	}

	return temp.ExecuteWriter(pongo2.Context{"movie": movie, "comments": comments}, c.Response())
}

func getMovies(c echo.Context) error {
	movies := []Movie{}
	db.Select(&movies, "SELECT * FROM movie")

	temp, err := pongo2.FromFile("./template/movies.html")
	if err != nil {
		return err
	}

	return temp.ExecuteWriter(pongo2.Context{"movies": movies}, c.Response())

}

func postMovies(c echo.Context) error {
	title := c.FormValue("title")
	filename := randomString(5)

	file, err := c.FormFile("thumbnail")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Destination
	dst, err := os.Create("./thumbnails/" + filename)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	movieFile, err := c.FormFile("movie")
	if err != nil {
		return err
	}
	movieSrc, err := movieFile.Open()
	if err != nil {
		return err
	}
	defer movieSrc.Close()

	// Destination
	movieDst, err := os.Create("./movies/" + filename)
	if err != nil {
		return err
	}
	defer movieDst.Close()

	// Copy
	if _, err = io.Copy(movieDst, movieSrc); err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO movie (title, path, thumbnail) VALUES (?, ?, ?)", title, "/movies/"+filename, "/thumbnails/"+filename)
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, "/")
}

func randomString(count int) string {
	gen := "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM"
	res := ""
	for i := 0; i < count; i++ {
		res += string(gen[rand.Intn(len(gen))])
	}
	return res
}
