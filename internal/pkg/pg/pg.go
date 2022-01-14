package pg

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/TheTeaParty/notnotes-platform/internal/config"
)

func New(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v", cfg.PGHost,
		cfg.PGUser, cfg.PGPass, cfg.PGName)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
}
