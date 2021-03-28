package repository

import (
	pkg "amazing_talker/pkg"

	"gitlab.com/howmay/gopher/db"
	gorm "gorm.io/gorm"
)

type repository struct {
	writeDB *gorm.DB
	readDB  *gorm.DB
}

// NewRepository 依賴注入
func NewRepository(conn *db.Connections) pkg.IRepository {
	return &repository{
		readDB:  conn.ReadDB,
		writeDB: conn.WriteDB,
	}
}

// WriteDB IRepository的實作
func (repo *repository) WriteDB() *gorm.DB {
	return repo.writeDB
}
