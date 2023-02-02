package database

import (
	"database/sql"
	"fmt"
	"os"

	"bitecodelabs.com/librarian/config"
	"bitecodelabs.com/librarian/logger"
	"github.com/go-sql-driver/mysql"
	"github.com/jamf/go-mysqldump"
)

func DumpDatabase(config config.DatabaseConfig) {
	// Open connection to database
	db_config := mysql.NewConfig()
	db_config.User = config.User
	db_config.Passwd = config.Password
	db_config.DBName = config.DBName
	db_config.Net = "tcp"
	db_config.Addr = config.Host

	dumpDir := config.Output_Directory                                    // you should create this directory
	dumpFilenameFormat := fmt.Sprintf("%s_02Jan2006_1504", config.DBName) // accepts time layout string and add .sql at the end of file

	err := os.MkdirAll(config.Output_Directory, 0700)
	if err != nil {
		logger.ErrorLog.Println("Error making database dump dir: ", err)
		return
	}

	db, err := sql.Open("mysql", db_config.FormatDSN()+"?parseTime=true")
	if err != nil {
		logger.ErrorLog.Println("Error opening database: ", err)
		return
	}

	// Register database with mysqldump
	dumper, err := mysqldump.Register(db, dumpDir, dumpFilenameFormat)
	if err != nil {
		logger.ErrorLog.Println("Error registering database:", err)
		return
	}

	// Dump database to file
	err = dumper.Dump()
	if err != nil {
		logger.ErrorLog.Println("Error dumping:", err)
		return
	}

	// Close dumper, connected database and file stream.
	dumper.Close()
}
