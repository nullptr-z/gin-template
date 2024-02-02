package dao

import (
	"github.com/jmoiron/sqlx"
	"github.com/nullptr-z/gin-template/settings"
)

var sq *sqlx.DB

func InitializeDao() {
	sq = settings.GetDB()
}
