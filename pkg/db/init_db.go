package db

import (
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

const (
	databaseName            string = "rubixcore"
	queuesCollectionName    string = "queues"
	customersCollectionName string = "customers"
	usersCollectionName     string = "users"
)

func createUsersTable(db *sqlx.DB, logger *zap.Logger) error {
	sql := `
	CREATE TABLE IF NOT EXISTS system_users(
		id 			SERIAL PRIMARY KEY,
		username 	VARCHAR(255) NOT NULL,
		fullname	VARCHAR(255) NOT NULL,
		password 	VARCHAR(255) NOT NULL,
		is_admin 	BOOLEAN NOT NULL,
		is_active 	BOOLEAN NOT NULL,
		created_at 	TIMESTAMPTZ DEFAULT Now(),
		updated_at 	TIMESTAMPTZ DEFAULT Now()
	);
	`
	logger.Info("creating system_users table")
	_, err := db.Exec(sql)
	if err != nil {
		return err
	}
	logger.Info("successfully created system_users table")

	return nil
}

func createQueuesTable(db *sqlx.DB, logger *zap.Logger) error {
	sql := `
	CREATE TABLE IF NOT EXISTS queues(
		id 			 	SERIAL PRIMARY KEY,
		name 	 		VARCHAR(255) NOT NULL,
		description 	VARCHAR(255) NOT NULL,
		title 			BOOLEAN NOT NULL,
		is_active 		BOOLEAN NOT NULL,
		created_by  	INTEGER REFERENCES system_users (id),
		created_at 		TIMESTAMPTZ DEFAULT Now(),
		updated_at 		TIMESTAMPTZ DEFAULT Now()
	);
	`
	logger.Info("creating queues table")
	_, err := db.Exec(sql)
	if err != nil {
		return err
	}
	logger.Info("successfully created queues table")

	return nil
}

func createCustomersTable(db *sqlx.DB, logger *zap.Logger) error {
	sql := `
	CREATE TABLE IF NOT EXISTS customers(
		id 			 	SERIAL PRIMARY KEY,
		msisdn 	 		VARCHAR(255) NOT NULL,
		queue_id		INTEGER REFERENCES queues (id),
		ticket_number 	VARCHAR(255) NOT NULL,
		served_at 		TIMESTAMPTZ,
		served_by		INTEGER REFERENCES system_users (id),
		created_at 		TIMESTAMPTZ DEFAULT Now(),
		updated_at 		TIMESTAMPTZ DEFAULT Now()
	);
	`
	logger.Info("creating queues table")
	_, err := db.Exec(sql)
	if err != nil {
		return err
	}
	logger.Info("successfully created queues table")

	return nil
}

func createSettingsTable(db *sqlx.DB, logger *zap.Logger) error {
	sql := `
	CREATE TABLE IF NOT EXISTS settings(
		id 			SERIAL PRIMARY KEY,
		key 	 	VARCHAR(255) NOT NULL,
		value 		VARCHAR(255) NOT NULL,
		created_by  INTEGER REFERENCES system_users (id),
		created_at 	TIMESTAMPTZ DEFAULT Now(),
		updated_at 	TIMESTAMPTZ DEFAULT Now()
	);
	`
	logger.Info("creating settings table")
	_, err := db.Exec(sql)
	if err != nil {
		return err
	}
	logger.Info("successfully created settings table")

	return nil
}

func createCustomersMsisdnIndex(db *sqlx.DB, logger *zap.Logger) error {
	sql := `
	CREATE INDEX IF NOT EXISTS customers_msisdn_index
	ON customers (msisdn);
	`
	logger.Info("creating customers msisdn index")
	_, err := db.Exec(sql)
	if err != nil {
		return err
	}
	logger.Info("successfully created customers msisdn index")

	return nil
}

func createSystemAdminUser(db *sqlx.DB, logger *zap.Logger, fullname, username, password string) error {
	sql := "INSERT INTO system_users (fullname, username, password, is_admin, is_active, created_at) VALUES($1, $2, $3, $4, $5, $6);"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	logger.Info("creating default admin user", zap.String("fullname", fullname))
	_, err = db.Exec(sql, fullname, username, string(hashedPassword), true, true, time.Now())
	if err != nil {
		return err
	}
	logger.Info("successfully created default admin user")

	return nil
}

// InitDB creates tables, sets constraints, sets indexes, and seeds appropriate tables
func InitDB(db *sqlx.DB, logger *zap.Logger, adminFullname, adminUsername, adminPass string) error {
	err := createUsersTable(db, logger)
	if err != nil {
		return err
	}

	err = createQueuesTable(db, logger)
	if err != nil {
		return err
	}

	err = createCustomersTable(db, logger)
	if err != nil {
		return err
	}

	err = createSettingsTable(db, logger)
	if err != nil {
		return err
	}

	err = createCustomersMsisdnIndex(db, logger)
	if err != nil {
		return err
	}

	err = createSystemAdminUser(db, logger, adminFullname, adminUsername, adminPass)
	if err != nil {
		return err
	}

	return nil
}
