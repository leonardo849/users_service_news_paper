package helper

import "gorm.io/gorm"

func SetTx(tx **gorm.DB, db *gorm.DB) {
	if *tx == nil {
		*tx = db
	}
}