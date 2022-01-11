package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/go-pg/pg/v10"
)

type user struct {
	First_name    string
	Last_name     string
	User_uuid     string
	Pwd           string
	Session_token string
	Email         string
	Is_admin      bool
}

func getUsers(db *pg.DB) ([]user, error) {
	var users []user
	_, err := db.Query(&users, "SELECT * FROM users")
	if err != nil {
		return nil, err
	}

	return users, nil
}

func backupUsers(db *pg.DB) error {
	users, err := getUsers(db)
	if err != nil {
		return err
	}

	err = backupEntity(db, "users", users)
	if err != nil {
		return err
	}

	fmt.Println("Successfully backed up users.")

	return nil
}

func recoverUsers(db *pg.DB) error {

	users, err := getUsersFromJson(currentFlags.OutputFolderPath + "/users.json")
	if err != nil {
		return err
	}

	for _, v := range users {
		err := insertIntoUsers(db, v)
		if err != nil {
			return err
		}
		fmt.Println("Succesfully recovered user: ", v.User_uuid, v.First_name)
	}

	return nil
}

func seedUsers(db *pg.DB) error {
	users, err := getUsersFromJson("./seed/users.json")
	if err != nil {
		return err
	}

	for _, v := range users {
		err := insertIntoUsers(db, v)
		if err != nil {
			return err
		}
		fmt.Println("Succesfully seeded user: ", v.User_uuid, v.First_name)
	}

	return nil
}

func insertIntoUsers(db *pg.DB, u user) error {
	st, err := db.Prepare(`INSERT INTO users(user_uuid,first_name,last_name,pwd,is_admin,session_token,email) VALUES($1,$2,$3,$4,$5,$6,$7)`)
	if err != nil {
		return err
	}
	_, err = st.Exec(u.User_uuid, u.First_name, u.Last_name, u.Pwd, u.Is_admin, u.Session_token, u.Email)
	return err
}

func getUsersFromJson(filePath string) ([]user, error) {
	usersJson, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	fmt.Println("Successfully Opened: " + filePath)
	defer usersJson.Close()

	byteUsers, err := ioutil.ReadAll(usersJson)
	if err != nil {
		return nil, err
	}

	var users []user
	err = json.Unmarshal(byteUsers, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}
