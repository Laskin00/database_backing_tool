package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/go-pg/pg/v10"
)

func backupEntity(db *pg.DB, entityName string, entities interface{}) error {
	entitiesJson, err := json.MarshalIndent(entities, "", " ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(currentFlags.OutputFolderPath+"/"+entityName+".json", entitiesJson, 0777)
	if err != nil {
		return err
	}

	err = sshConnection.sendFile(currentFlags.OutputFolderPath + "/" + entityName + ".json")
	if err != nil {
		return err
	}

	return nil
}
