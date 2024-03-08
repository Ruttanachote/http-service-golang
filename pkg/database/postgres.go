package database

import (
	"log"
	"time"

	"github.com/Stream-I-T-Consulting/stream-http-service-go/config"
	"github.com/Stream-I-T-Consulting/stream-http-service-go/pkg/utils/color"
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

func Initialize(config *config.Config) *gorm.DB {
	if !fiber.IsChild() {
		log.Println("Database connecting...")
	}

	dbConn, err := gorm.Open(
		postgres.New(
			postgres.Config{
				DSN: config.Database.DatabaseDSN,
			},
		),
		&gorm.Config{},
	)
	dbConn.Use(
		dbresolver.Register(dbresolver.Config{
			Sources:           []gorm.Dialector{},
			Replicas:          []gorm.Dialector{},
			Policy:            nil,
			TraceResolverMode: false,
		}).
			SetConnMaxIdleTime(time.Hour).
			SetConnMaxLifetime(24 * time.Hour).
			SetMaxIdleConns(config.Database.DatabaseMaxIdleConns).
			SetMaxOpenConns(config.Database.DatabaseMaxOpenConns),
	)

	if err != nil {
		log.Printf("Cannot connect to database")
		log.Fatal("DatabaseError:", err)
	}

	if !fiber.IsChild() {
		log.Println("Database connected", color.Format(color.GREEN, "successfully!"))
	}

	return dbConn
}
