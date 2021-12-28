package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/go-pg/pg/v10"
)

type organization struct {
	Organization_uuid string
}

func getOrganizations(db *pg.DB) ([]organization, error) {
	var organizations []organization
	_, err := db.Query(&organizations, "SELECT * FROM organizations")
	if err != nil {
		return nil, err
	}

	return organizations, nil
}

func backupOrganizations(db *pg.DB) error {
	organizations, err := getOrganizations(db)
	if err != nil {
		return err
	}

	err = backupEntity(db, "organizations", organizations)
	if err != nil {
		return err
	}

	fmt.Println("Successfully backed up organizations.")

	return nil
}

func recoverOrganizations(db *pg.DB) error {
	organizations, err := getOrganizationsFromJson(currentFlags.OutputFolderPath + "/organizations.json")
	if err != nil {
		return err
	}

	for _, v := range organizations {
		err := insertIntoOrganizations(db, v)
		if err != nil {
			return err
		}
		fmt.Println("Succesfully recovered organization: ", v.Organization_uuid)
	}

	return nil
}

func seedOrganizations(db *pg.DB) error {
	organizations, err := getOrganizationsFromJson("./seed/organizations.json")
	if err != nil {
		return err
	}

	for _, v := range organizations {
		err := insertIntoOrganizations(db, v)
		if err != nil {
			return err
		}
		fmt.Println("Succesfully seeded organization: ", v.Organization_uuid)
	}

	return nil
}

func insertIntoOrganizations(db *pg.DB, o organization) error {
	st, err := db.Prepare(`INSERT INTO organizations(organization_uuid) VALUES($1)`)
	if err != nil {
		return err
	}
	_, err = st.Exec(o.Organization_uuid)
	return err
}

func getOrganizationsFromJson(filePath string) ([]organization, error) {
	organizationsJson, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	fmt.Println("Successfully Opened: " + filePath)
	defer organizationsJson.Close()

	byteOrganizations, err := ioutil.ReadAll(organizationsJson)
	if err != nil {
		return nil, err
	}

	var organizations []organization
	err = json.Unmarshal(byteOrganizations, &organizations)
	if err != nil {
		return nil, err
	}

	return organizations, nil
}
