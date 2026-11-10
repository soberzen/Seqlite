package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"qingmo-server/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PostgresService interface {
	Health() map[string]string
	Close() error
	DB() *gorm.DB
}

type postgresService struct {
	db *gorm.DB
}

var postgresInstance *postgresService

func NewPostgres(cfg config.PostgresConfig) PostgresService {
	if postgresInstance != nil {
		return postgresInstance
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s search_path=%s",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.Name,
		cfg.Port,
		cfg.SSLMode,
		cfg.Schema,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:         logger.Default.LogMode(logger.Warn),
		TranslateError: true,
	})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
		return nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get database handle: %v", err)
		return nil
	}

	sqlDB.SetMaxOpenConns(50)                 // 设置最大打开连接数
	sqlDB.SetMaxIdleConns(25)                 // 设置最大空闲连接数
	sqlDB.SetConnMaxIdleTime(5 * time.Minute) // 设置最大空闲时间

	postgresInstance = &postgresService{db: db}

	return postgresInstance
}

func (p *postgresService) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)

	defer cancel()

	stats := make(map[string]string)

	sqlDB, err := p.db.DB()
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("failed to get database handle %v", err)
		return stats
	}

	err = sqlDB.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("failed to ping database %v", err)
		return stats
	}

	stats["status"] = "up"
	stats["message"] = "database is healthy"

	dbStats := sqlDB.Stats()
	stats["open_connections"] = strconv.Itoa(
		dbStats.OpenConnections,
	) // 已打开连接数
	stats["in_use"] = strconv.Itoa(
		dbStats.InUse,
	) // 正在使用的连接数
	stats["idle"] = strconv.Itoa(
		dbStats.Idle,
	) // 空闲连接数
	stats["wait_count"] = strconv.FormatInt(
		dbStats.WaitCount,
		10,
	) // 等待连接数
	stats["wait_duration"] = dbStats.WaitDuration.String() // 等待连接时间
	stats["max_idle_closed"] = strconv.FormatInt(
		dbStats.MaxIdleClosed,
		10,
	) // 最大空闲连接数
	stats["max_lifetime_closed"] = strconv.FormatInt(
		dbStats.MaxLifetimeClosed,
		10,
	) // 最大生命周期连接数

	// 已打开连接数超过40，说明数据库负载较高，考虑调整最大打开连接数
	if dbStats.OpenConnections > 40 {
		stats["message"] = "The database is experiencing heavy load."
	}
	// 等待连接数超过1000，说明数据库负载较高，考虑调整最大打开连接数
	if dbStats.WaitCount > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}
	// 最大空闲连接数超过总连接数的一半，说明空闲连接数过多，考虑调整最大空闲连接数
	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}
	// 最大生命周期连接数超过总连接数的一半，说明连接生命周期过长，考虑增加最大生命周期或调整连接使用模式
	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats
}

func (p *postgresService) Close() error {
	sqlDB, err := p.db.DB()
	if err != nil {
		return err
	}
	log.Printf("Disconnected from database: %s", os.Getenv("POSTGRES_NAME"))
	return sqlDB.Close()
}

func (p *postgresService) DB() *gorm.DB {
	return p.db
}
