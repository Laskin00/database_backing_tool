package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/go-pg/pg/v10"
)

type userOrganizationConnection struct {
	Connection_id     int
	Is_manager        bool
	User_uuid         string
	Organization_uuid string
	Position          string
}

func getUserOrganizationConnections(db *pg.DB) ([]userOrganizationConnection, error) {
	var uoc []userOrganizationConnection
	_, err := db.Query(&uoc, "SELECT * FROM userorganizationconnections")
	if err != nil {
		return nil, err
	}

	return uoc, nil
}

func backupUserOrganizationConnections(db *pg.DB) error {
	uoc, err := getUserOrganizationConnections(db)
	if err != nil {
		return err
	}

	err = backupEntity(db, "userOrganizationConnections", uoc)
	if err != nil {
		return err
	}

	fmt.Println("Successfully backed up userOrganizationConnections.")

	return nil
}

func recoverUserOrganizationConnections(db *pg.DB) error {
	userOrganizationConnections, err := getUserOrganizationConnectionsFromJson(currentFlags.OutputFolderPath + "/userOrganizationConnections.json")
	if err != nil {
		return err
	}

	for _, v := range userOrganizationConnections {
		err := insertIntoUserOrganizationConnections(db, v)
		if err != nil {
			return err
		}
		fmt.Println("Succesfully recovered userOrganizationConnection: ", v.Connection_id)
	}

	return nil
}

func seedUserOrganizationConnections(db *pg.DB) error {
	userOrganizationConnections, err := getUserOrganizationConnectionsFromJson("./seed/userOrganizationConnections.json")
	if err != nil {
		return err
	}

	for _, v := range userOrganizationConnections {
		err := insertIntoUserOrganizationConnections(db, v)
		if err != nil {
			return err
		}
		fmt.Println("Succesfully seeded userOrganizationConnection: ", v.Connection_id)
	}

	return nil
}

func insertIntoUserOrganizationConnections(db *pg.DB, uoc userOrganizationConnection) error {
	st, err := db.Prepare(`INSERT INTO userOrganizationConnections(connection_id,user_uuid,organization_uuid,is_manager) VALUES($1,$2,$3,$4,$5)`)
	if err != nil {
		return err
	}
	_, err = st.Exec(uoc.Connection_id, uoc.User_uuid, uoc.Organization_uuid, uoc.Is_manager, uoc.Position)
	return err
}

func getUserOrganizationConnectionsFromJson(filePath string) ([]userOrganizationConnection, error) {
	userOrganizationConnectionsJson, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	fmt.Println("Successfully Opened: " + filePath)
	defer userOrganizationConnectionsJson.Close()

	byteuserOrganizationConnections, err := ioutil.ReadAll(userOrganizationConnectionsJson)
	if err != nil {
		return nil, err
	}

	var userOrganizationConnections []userOrganizationConnection
	err = json.Unmarshal(byteuserOrganizationConnections, &userOrganizationConnections)
	if err != nil {
		return nil, err
	}

	return userOrganizationConnections, nil
}
