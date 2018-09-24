package router

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/tohutohu/Donuts/product/live-server/db"
)

func (r *router) StartLive(c echo.Context) error {
	name := c.QueryParam("name")
	e := c.QueryParam("e")
	st := c.QueryParam("st")

	live := &db.Live{
		Name: name,
		E:    e,
		St:   st,
		Done: false,
	}

	r.db.Create(live)

	return c.String(http.StatusOK, "OK")
}

func (r *router) EndLive(c echo.Context) error {
	name := c.QueryParam("name")
	e := c.QueryParam("e")
	st := c.QueryParam("st")

	live := &db.Live{
		Name: name,
		E:    e,
		St:   st,
	}

	r.db.First(live)
	if live.ID == 0 {
		return c.NoContent(http.StatusNotFound)
	}

	live.Done = true
	r.db.Save(live)

	return c.String(http.StatusOK, "OK")
}
