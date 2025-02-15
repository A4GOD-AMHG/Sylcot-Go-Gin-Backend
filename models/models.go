package models

import "gorm.io/gorm"

func MigrateAll(db *gorm.DB) error {
	err := db.AutoMigrate(
		&User{},
		&Task{},
		&Category{},
	)

	if db.Dialector.Name() == "mysql" {
		db.Exec("ALTER TABLE tasks MODIFY priority ENUM('high', 'medium', 'low')")
	}

	return err
}
