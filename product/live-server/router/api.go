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
	lives := []Live{}
	r.db.Find(&lives, "done = false")
	return c.JSON(http.StatusOK, lives)
}
