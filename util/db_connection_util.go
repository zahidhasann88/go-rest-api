package util

import (
	"fmt"
	"log"
	"os"

	"strconv"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)


// GetConfigString returns values of string variable from local-config.json
func GetConfigString(key string) string {
	viper.AddConfigPath("./config")
	viper.SetConfigName("local-config")
	viper.SetConfigType("json")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}
	value, ok := viper.Get(key).(string)
	if !ok {
		log.Fatalf("Invalid type assertion: " + key)
		return ""
	}
	return value
}

// GetConfigInt returns values of int variable from local-config.json
func GetConfigInt(key string) int {
	viper.AddConfigPath("./config")
	viper.SetConfigName("local-config")
	viper.SetConfigType("json")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}
	value := viper.GetInt(key)
	return value
}

// ViperReturnStringConfigVariableFromLocalConfigJSON returns values of string variable from local-config.json
func ViperReturnStringConfigVariableFromLocalConfigJSON(key string) string {
	// viper.SetConfigFile("local-config.json")
	viper.SetConfigName("local-config") // name of config file (without extension)
	viper.SetConfigType("json")         // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("./config")     // path to look for the config file in
	// viper.AddConfigPath("$HOME/.appname") // call multiple times to add many search paths
	// viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}
	value, ok := viper.Get(key).(string)
	if !ok {
		fmt.Println(key)
		fmt.Println(value)
		log.Fatalf("Invalid type assertion")
		return ""
	}
	return value
}

// ViperReturnIntegerConfigVariableFromLocalConfigJSON returns values of int variable from local-config.json
func ViperReturnIntegerConfigVariableFromLocalConfigJSON(key string) int {
	// viper.SetConfigFile("local-config.json")
	viper.SetConfigName("local-config") // name of config file (without extension)
	viper.SetConfigType("json")         // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("./config")     // path to look for the config file in
	// viper.AddConfigPath("$HOME/.appname") // call multiple times to add many search paths
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}
	value := viper.GetInt(key)
	return value
}

// CreateConnectionUsingGormToProcurementSchema creates database connection using gorm to procurement schema
func CreateConnectionUsingGormToCommonSchema() *gorm.DB {
	fmt.Println("Connecting....")
	dbHost := ViperReturnStringConfigVariableFromLocalConfigJSON("db_host")
	dbPort := ViperReturnIntegerConfigVariableFromLocalConfigJSON("db_port")
	dbName := ViperReturnStringConfigVariableFromLocalConfigJSON("db_name")
	dbUsername := ViperReturnStringConfigVariableFromLocalConfigJSON("db_username")
	dbPassword := ViperReturnStringConfigVariableFromLocalConfigJSON("db_password")

	dataSourceName := "host=" + dbHost + " user=" + dbUsername + " password=" + dbPassword + " dbname=" + dbName + " port=" + strconv.Itoa(dbPort) + " sslmode=disable"
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			LogLevel: logger.Info, // Log level
		},
	)
	db, err := gorm.Open(postgres.Open(dataSourceName), &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   ViperReturnStringConfigVariableFromLocalConfigJSON("my_schema_name") + ".",
			SingularTable: true,
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		// fmt.Println("failed to connect database")
		panic(err)
	} else {
		return db
	}
}
