package database

import (
	"context"
	"database/sql"

	_ "share/config"

	// register mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/sunliang711/goutils/mongodb"
	umysql "github.com/sunliang711/goutils/mysql"

	//"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/spf13/viper"
)

var (
	MysqlConn    *sql.DB
	MysqlORMConn *gorm.DB
	MongoConn    *mongo.Client
)

func init() {
	logrus.Infof("database init()...")
	if viper.GetBool("mysql.enable") {
		dsn := viper.GetString("mysql.dsn")
		initMysql(dsn)
	} else {
		logrus.Info("Mysql is disabled.")
	}

	if viper.GetBool("mongodb.enable") {
		dsn := viper.GetString("mongodb.url")
		initMongo(dsn)
	} else {
		logrus.Info("Mongodb is disabled.")
	}
}

//initMysql open mysql with dsn
func initMysql(dsn string) {
	logrus.Infof("Try to connect to mysql: '%v'", dsn)
	var err error
	if viper.GetBool("mysql.orm") {
		MysqlORMConn, err = gorm.Open("mysql", dsn)
		logrus.Infof("Use gorm driver...")
		if err != nil {
			panic(err)
		}

	} else {
		MysqlConn, err = umysql.New(dsn, viper.GetInt("mysql.maxIdleConns"), viper.GetInt("mysql.maxOpenConns"))
		if err != nil {
			panic(err)
		}
	}
	logrus.Infof("Connected to mysql")
}

//CloseMysql close mysql connection
func closeMysql() {
	if viper.GetBool("mysql.enable") {
		logrus.Infoln("Close mysql.")
		if viper.GetBool("mysql.orm") {
			MysqlORMConn.Close()
		} else {
			MysqlConn.Close()
		}
	}
}

// InitMongo opens a mongodb connection
func initMongo(url string) {
	logrus.Infof("Try to connect to mongodb: '%v'", url)
	var err error
	MongoConn, err = mongodb.New(url, 5)
	if err != nil {
		panic(err)
	}
	logrus.Infof("Connected to mongodb")
}

func closeMongo() {
	if viper.GetBool("mongodb.enable") {
		logrus.Infof("Close mongodb")
		MongoConn.Disconnect(context.Background())
	}
}

func Release() {
	closeMongo()
	closeMysql()
	closeRedis()
}
