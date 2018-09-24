package router

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/tohutohu/Donuts/product/live-server/db"
	"golang.org/x/crypto/bcrypt"
)

func (r *router) GetLiveEndpoint(c echo.Context) error {
	rp := strings.NewReplacer("+", "-", "/", "_", "=", "")
	unix := time.Now().Add(time.Hour).Unix()
	h := md5.New()
	req := &struct {
		Name string `json:"name"`
	}{}
	c.Bind(req)
	if req.Name == "" {
		return c.NoContent(http.StatusBadRequest)
	}
	h.Write([]byte("CocoroIsGodlive/" + req.Name + strconv.FormatInt(unix, 10)))
	fmt.Println(unix)
	key := rp.Replace(base64.StdEncoding.EncodeToString(h.Sum(nil)))

	return c.JSON(http.StatusOK, map[string]string{"url": fmt.Sprintf("http://localhost:1935/live/%s?e=%d&st=%s", req.Name, unix, key)})
}

func (r *router) GetLives(c echo.Context) error {
	lives := []db.Live{}
	r.db.Find(&lives, "done = false")
	return c.JSON(http.StatusOK, lives)
}

func (r *router) PostUsers(c echo.Context) error {
	req := &struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}

	c.Bind(req)
	if req.Username == "" || req.Password == "" {
		return c.NoContent(http.StatusBadRequest)
	}

	u := &db.User{}
	r.db.Where("username = ?", req.Username).First(u)
	if !r.db.NewRecord(u) {
		return c.JSON(http.StatusConflict, H{"message": "username conflict"})
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	user := &db.User{
		Username: req.Username,
		Password: string(hashedPass),
	}

	r.db.Create(user)

	sess, _ := session.Get("sessions", c)
	sess.Values["user_id"] = user.ID
	sess.Values["username"] = user.Username
	sess.Save(c.Request(), c.Response())

	return c.NoContent(http.StatusCreated)
}

func (r *router) PostLogin(c echo.Context) error {
	req := &struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}

	c.Bind(req)
	if req.Username == "" || req.Password == "" {
		return c.NoContent(http.StatusBadRequest)
	}

	u := &db.User{}
	r.db.Where("username = ?", req.Username).First(u)

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		return c.NoContent(http.StatusForbidden)
	}

	sess, _ := session.Get("sessions", c)
	sess.Values["user_id"] = u.ID
	sess.Values["username"] = u.Username
	sess.Save(c.Request(), c.Response())
	return c.NoContent(http.StatusOK)
}

func (r *router) GetWhoAmI(c echo.Context) error {
	sess, _ := session.Get("sessions", c)
	if sess.Values["user_id"] == nil {
		return c.JSON(http.StatusOK, H{"user_id": "0", "username": "anonymous"})
	}
	id := strconv.FormatUint(uint64(sess.Values["user_id"].(uint)), 10)
	username := sess.Values["username"].(string)
	return c.JSON(http.StatusOK, H{"user_id": id, "username": username})
}
