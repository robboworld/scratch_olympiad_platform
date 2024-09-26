package db

import (
	"github.com/robboworld/scratch_olympiad_platform/internal/consts"
	"github.com/robboworld/scratch_olympiad_platform/internal/models"
	"github.com/robboworld/scratch_olympiad_platform/pkg/logger"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

type PostgresClient struct {
	Db         *gorm.DB
	InfoLogger *log.Logger
}

func InitPostgresClient(m consts.Mode, loggers logger.Loggers) (postgresClient PostgresClient, err error) {
	// set stdout gorm logger depends on app mode
	var dbLogger gormLogger.Interface
	switch m {
	case consts.Production:
		gormF, err := os.OpenFile(viper.GetString("logger.gorm"), os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			loggers.Err.Fatalf("%s", err.Error())
		}
		defer func(gormF *os.File) {
			err := gormF.Close()
			if err != nil {
				loggers.Err.Fatalf("%s", err.Error())
			}
		}(gormF)
		dbLogger = gormLogger.New(
			log.New(gormF, "[GORM]\t", log.LstdFlags),
			gormLogger.Config{
				SlowThreshold:             time.Second,     // Slow SQL threshold
				LogLevel:                  gormLogger.Info, // Log level
				IgnoreRecordNotFoundError: false,           // Ignore ErrRecordNotFound error for logger
				Colorful:                  true,            // Disable color
			},
		)
	case consts.Development:
		dbLogger = gormLogger.New(
			log.New(os.Stdout, "[GORM]\t", log.LstdFlags),
			gormLogger.Config{
				SlowThreshold:             time.Second,     // Slow SQL threshold
				LogLevel:                  gormLogger.Info, // Log level
				IgnoreRecordNotFoundError: false,           // Ignore ErrRecordNotFound error for logger
				Colorful:                  true,            // Disable color
			},
		)
	}

	db, err := gorm.Open(postgres.Open(viper.GetString("postgres_dsn")), &gorm.Config{Logger: dbLogger})
	if err != nil {
		loggers.Err.Fatalf("Failed to initialize postgres client: %s", err.Error())
		return
	}
	postgresClient = PostgresClient{
		Db:         db,
		InfoLogger: loggers.Info,
	}
	if migrateErr := postgresClient.Migrate(); migrateErr != nil {
		loggers.Err.Fatalf("Failed to migrate: %s", migrateErr.Error())
	}
	return postgresClient, err
}

func (c *PostgresClient) Migrate() (err error) {
	err = c.Db.AutoMigrate(
		&models.UserCore{},
		&models.ProjectPageCore{},
		&models.ProjectCore{},
		&models.ParentRelCore{},
		&models.SettingsCore{},
		&models.ApplicationCore{},
		&models.NominationCore{},
		&models.CountryCore{},
		&models.SolutionCore{},
	)
	if err != nil {
		return err
	}
	var count int64
	if err := c.Db.First(&models.SettingsCore{ID: 1}).Count(&count).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}
	}
	if count > 0 {
		return nil
	} else {
		return c.Db.Create(&models.SettingsCore{ID: 1, ActivationByLink: true}).Error
	}
}
