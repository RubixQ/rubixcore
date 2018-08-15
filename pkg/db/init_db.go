package db

import (
	"database/sql"

	"go.uber.org/zap"
)

const (
	databaseName            string = "rubixcore"
	queuesCollectionName    string = "queues"
	customersCollectionName string = "customers"
	usersCollectionName     string = "users"
)

func createUsersTable(db *sql.DB, logger *zap.Logger) error {
	sql := `
	CREATE TABLE IF NOT EXISTS system_users(
		id 			SERIAL PRIMARY KEY,
		username 	VARCHAR(255) NOT NULL,
		fullname	VARCHAR(255) NOT NULL,
		password 	VARCHAR(255) NOT NULL,
		is_admin 	BOOLEAN NOT NULL,
		is_active 	BOOLEAN NOT NULL,
		created_by  INTEGER REFERENCES system_users (id),
		created_at 	TIMESTAMPTZ DEFAULT Now(),
		updated_at 	TIMESTAMPTZ
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

func createQueuesTable(db *sql.DB, logger *zap.Logger) error {
	sql := `
	CREATE TABLE IF NOT EXISTS queues(
		id 			 	SERIAL PRIMARY KEY,
		name 	 		VARCHAR(255) NOT NULL,
		description 	VARCHAR(255) NOT NULL,
		title 			BOOLEAN NOT NULL,
		is_active 		BOOLEAN NOT NULL,
		created_by  	INTEGER REFERENCES system_users (id),
		created_at 		TIMESTAMPTZ DEFAULT Now(),
		updated_at 		TIMESTAMPTZ
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

func createCustomersTable(db *sql.DB, logger *zap.Logger) error {
	sql := `
	CREATE TABLE IF NOT EXISTS customers(
		id 			 	SERIAL PRIMARY KEY,
		msisdn 	 		VARCHAR(255) NOT NULL,
		queue_id		INTEGER REFERENCES queues (id),
		ticket_number 	VARCHAR(255) NOT NULL,
		served_at 		TIMESTAMPTZ,
		served_by		INTEGER REFERENCES system_users (id),
		created_at 		TIMESTAMPTZ DEFAULT Now(),
		updated_at 		TIMESTAMPTZ
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

func createSettingsTable(db *sql.DB, logger *zap.Logger) error {
	sql := `
	CREATE TABLE IF NOT EXISTS settings(
		id 			SERIAL PRIMARY KEY,
		key 	 	VARCHAR(255) NOT NULL,
		value 		VARCHAR(255) NOT NULL,
		created_by  INTEGER REFERENCES system_users (id),
		created_at 	TIMESTAMPTZ DEFAULT Now(),
		updated_at 	TIMESTAMPTZ
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

func createCustomersMsisdnIndex(db *sql.DB, logger *zap.Logger) error {
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

// InitDB sets db contraints and indexes
func InitDB(db *sql.DB, logger *zap.Logger) error {
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

	return nil
}
