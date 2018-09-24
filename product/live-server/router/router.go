package router

import "github.com/jinzhu/gorm"

type H map[string]string

type router struct {
	db *gorm.DB
}

func New(db *gorm.DB) *router {
	return &router{db: db}
}
