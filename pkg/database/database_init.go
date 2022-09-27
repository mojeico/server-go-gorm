package database

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/trucktrace/pkg/logger"

	"github.com/trucktrace/pkg/queue"

	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
	"github.com/olivere/elastic/v7"
	"github.com/trucktrace/internal/models"
)

func NewPostgresConnection() *gorm.DB {

	// dbUri := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
	// 	viper.GetString("postgres.host"),
	// 	viper.GetString("postgres.user"),
	// 	os.Getenv("postgres_pass"),
	// 	viper.GetString("postgres.name"),
	// 	viper.GetString("postgres.port"))

	dbUri := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		//viper.GetString("postgres.host"),
		"trucktrace.cimnhoapkcj7.us-east-2.rds.amazonaws.com",
		//viper.GetString("postgres.user"),
		"trucktrace",
		//os.Getenv("postgres_pass"),
		"BPxx~:)8)-+HxG8_",
		//viper.GetString("postgres.name"),
		"trucktrace",
		//viper.GetString("postgres.port"),
		"5432",
	)

	db, err := gorm.Open("postgres", dbUri)

	if err != nil {
		logger.SystemLoggerError("NewPostgresConnection", "Cant connect to Postgres").Error("Error - " + err.Error())
		panic(err.Error())
	}

	if migrationErr := initModels(db); migrationErr != nil {
		logger.SystemLoggerError("NewPostgresConnection", "Cant do migration").Error("Error - " + err.Error())
		panic(err.Error())
	}

	logger.SystemLoggerInfo("NewPostgresConnection").Info("Successfully connected to Postgres.")
	return db

}

func NewRedisConnection() *redis.Client {

	// redisDB := redis.NewClient(&redis.Options{
	// 	Addr: fmt.Sprintf("%s:%s",
	// 		viper.GetString("redis.host"),
	// 		viper.GetString("redis.port"),
	// 	),
	// 	Password: os.Getenv("redis_pass"),
	// 	DB:       viper.GetInt("redis.db"),
	// })

	redisDB := redis.NewClient(&redis.Options{
		Addr: "trucktrace-redis.6mjflf.ng.0001.use2.cache.amazonaws.com:6379",
	})

	ctx := context.Background()
	err := redisDB.Ping(ctx).Err()

	if err != nil {
		logger.SystemLoggerError("NewRedisConnection", "Cant connect to REDIS").Error("Error - " + err.Error())
		panic(err.Error())
	}

	logger.SystemLoggerInfo("NewRedisConnection").Info("Successfully connected to REDIS.")
	return redisDB

}

func NewElasticSearchConnection() *elastic.Client {

	// host := fmt.Sprintf("%s:%s",
	// 	viper.GetString("elastic.host"),
	// 	viper.GetString("elastic.port"),
	// )
	// client, err := elastic.NewClient(
	// 	elastic.SetURL(host),
	// 	elastic.SetSniff(false),
	// 	elastic.SetBasicAuth("elastic", "trucktrace"),
	// )

	client, err := elastic.NewClient(
		elastic.SetURL("https://trucktrace.es.us-central1.gcp.cloud.es.io:9243"),
		elastic.SetSniff(false),
		elastic.SetBasicAuth("elastic", "gMuOS1cKSTYJgMEGz55KlLaC"),
	)

	if err != nil {
		logger.SystemLoggerError("NewElasticSearchConnection", "Cant connect to Elastic Search").Error("Error - " + err.Error())
		panic(err.Error())
	}

	initModelsElastic(client)

	logger.SystemLoggerInfo("NewElasticSearchConnection").Info("Successfully connected to Elastic Search.")

	return client
}

func initModels(db *gorm.DB) error {
	var migrateOnce sync.Once

	migrateOnce.Do(func() {
		if err := db.AutoMigrate(
			&models.Groups{},
			&models.User{},
			&models.Driver{},
			&models.Trailer{},
			&models.File{},
			&models.Settlement{},
			&models.Truck{},
			&models.Customer{},
			&models.Order{},
			&models.Safety{},
			&models.Company{},
			&models.Charges{},
			&models.Invoicing{},
			&queue.UserNotification{},
			&queue.SystemNotification{},
			&models.ExtraPay{},
			&models.TrailerComment{},
		).Error; err != nil {
			log.Fatal("can't migrate to postgres")
		}
	})
	return nil
}

func initModelsElastic(db *elastic.Client) {

	var ctx = context.Background()

	db.CreateIndex("company").Do(ctx)
	db.CreateIndex("customer").Do(ctx)
	db.CreateIndex("driver").Do(ctx)
	db.CreateIndex("group").Do(ctx)
	db.CreateIndex("order").Do(ctx)
	db.CreateIndex("safety").Do(ctx)
	db.CreateIndex("settlement").Do(ctx)
	db.CreateIndex("trailer").Do(ctx)
	db.CreateIndex("truck").Do(ctx)
	db.CreateIndex("user").Do(ctx)
	db.CreateIndex("check").Do(ctx)
	db.CreateIndex("charges").Do(ctx)
	db.CreateIndex("invoicing").Do(ctx)
	db.CreateIndex("trailer_comments").Do(ctx)

	db.CreateIndex("user_notification").Do(ctx)
	db.CreateIndex("system_notification").Do(ctx)

	db.CreateIndex("sms").Do(ctx)
	db.CreateIndex("email").Do(ctx)

	db.CreateIndex("files").Do(ctx)

	db.CreateIndex("extrapay").Do(ctx)
}
