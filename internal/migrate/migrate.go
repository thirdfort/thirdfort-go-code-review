package migrate

import (
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/thirdfort/thirdfort-go-code-review/internal/models"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Actor{},
		&models.Transaction{},
	)
}

func DropTables(db *gorm.DB) error {
	return db.Migrator().DropTable(
		&models.Actor{},
		&models.Transaction{},
	)
}
