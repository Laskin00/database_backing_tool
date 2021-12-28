package main

import (
	"os"

	"github.com/go-pg/pg/v10"
)

func connect() (*pg.DB, error) {
	options := pg.Options{
		Addr:     currentFlags.Addr,
		User:     currentFlags.User,
		Password: currentFlags.Password,
		Database: currentFlags.Database,
	}
	db := pg.Connect(&options)

	var model []struct {
		First_name string
	}

	_, err := db.Query(&model, "Select first_name from Users")
	if err != nil {
		return nil, err
	}

	return db, nil
}

func seed(db *pg.DB) error {
	err := seedUsers(db)
	if err != nil {
		return err
	}

	err = seedOrganizations(db)
	if err != nil {
		return err
	}

	err = seedUserOrganizationConnections(db)
	if err != nil {
		return err
	}

	return nil
}

func recover(db *pg.DB) error {
	err := recoverUsers(db)
	if err != nil {
		return err
	}

	err = recoverOrganizations(db)
	if err != nil {
		return err
	}

	err = recoverUserOrganizationConnections(db)
	if err != nil {
		return err
	}

	return nil
}

func backup(db *pg.DB) error {
	os.RemoveAll(currentFlags.OutputFolderPath)

	err := os.Mkdir(currentFlags.OutputFolderPath, 0777)
	if err != nil {
		return err
	}

	err = backupUsers(db)
	if err != nil {
		return err
	}

	err = backupOrganizations(db)
	if err != nil {
		return err
	}

	err = backupUserOrganizationConnections(db)
	if err != nil {
		return err
	}

	return nil
}
