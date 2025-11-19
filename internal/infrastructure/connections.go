// connections.go placeholder
package infrastructure

import (
	"fmt"
	"log"

	"app/internal/config"
	"app/internal/connections"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// ✅ HIGH-LEVEL: Holds all connections
type Connections struct {
	DB    *gorm.DB
	Redis *redis.Client
	// S3 client can be added
}

// ✅ HIGH-LEVEL: Initialize ALL connections
func InitializeConnections(cfg *config.Config) (*Connections, error) {
	log.Println("Initializing connections...")

	// 1. Database - Use LOW-LEVEL function
	var db *gorm.DB
	var err error

	switch cfg.Database.Type {
	case "postgres":
		db, err = connections.NewPostgresConnection(
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.Name,
			cfg.Database.SSLMode,
		)
	case "mysql":
		db, err = connections.NewMySQLConnection(
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.Name,
		)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", cfg.Database.Type)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}
	log.Printf("✅ Database connected (%s)", cfg.Database.Type)

	// 2. Redis - Use LOW-LEVEL function
	var redisClient *redis.Client
	if cfg.Redis.Enabled {
		redisClient, err = connections.NewRedisConnection(
			&cfg.Redis,
		)
		if err != nil {
			log.Printf("⚠️  Redis not available: %v", err)
			redisClient = nil
		} else {
			log.Println("✅ Redis connected")
		}
	}

	return &Connections{
		DB:    db,
		Redis: redisClient,
	}, nil
}

// Cleanup closes all connections
// func (c *Connections) Cleanup() {
// 	log.Println("Cleaning up connections...")

// 	if c.DB != nil {
// 		c.DB.Close()
// 		log.Println("✅ Database closed")
// 	}

// 	if c.Redis != nil {
// 		c.Redis.Close()
// 		log.Println("✅ Redis closed")
// 	}
// }
