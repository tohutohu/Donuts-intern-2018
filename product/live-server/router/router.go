package router

import "github.com/jinzhu/gorm"

type Live struct {
	gorm.Model
	Name string
	E    string
	St   string
	Done bool
}

type router struct {
	db *gorm.DB
}

func New(db *gorm.DB) *router {
	return &router{db: db}
}
