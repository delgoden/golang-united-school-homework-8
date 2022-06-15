package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

var (
	userNotFound      = "Item with id %s not found"
	itemAlreadyExists = "Item with id %s already exists"
)

type user struct {
	Id    string `json:"id"`
	Email string `json:"email,omitempty"`
	Age   int    `json:"age,omitempty"`
}

func List(fileName string) (string, error) {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return "", err
	}

	defer func() {
		cerr := file.Close()
		if cerr != nil {
			err = cerr
		}
	}()
	data, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	if len(data) == 0 {
		return "", errors.New("file is empty")
	}

	return string(data), nil

}

func Add(item string, fileName string, writer io.Writer) error {
	var userItem user
	err := json.Unmarshal([]byte(item), &userItem)
	if err != nil {
		return err
	}
	var users []user
	itemsInFile, err := List(fileName)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(itemsInFile), &users)
	if err != nil {
		return err
	}

	for _, user := range users {
		if user.Id == userItem.Id {
			_, err = fmt.Fprintf(writer, itemAlreadyExists, userItem.Id)
			if err != nil {
				return err
			}
			return nil
		}
	}

	users = append(users, userItem)

	data, err := json.Marshal(users)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = os.WriteFile(fileName, []byte(data), 0644)
	if err != nil {
		return err
	}

	_, err = fmt.Fprint(writer, string(item))
	if err != nil {
		return err
	}
	return nil
}

func FindById(id string, fileName string, writer io.Writer) (err error) {
	var itemFind user
	itemsInFile, err := List(fileName)
	if err != nil {
		return err
	}
	var users []user

	err = json.Unmarshal([]byte(itemsInFile), &users)
	if err != nil {
		return err
	}

	for i, user := range users {
		if user.Id == id {
			itemFind = user

			break
		}

		if i == len(users)-1 {
			_, err = fmt.Fprint(writer, "")
			if err != nil {
				return
			}
			return
		}
	}

	data, err := json.Marshal(itemFind)
	if err != nil {
		return
	}
	_, err = fmt.Fprint(writer, string(data))
	if err != nil {
		return
	}
	return
}

func Remove(id string, fileName string, writer io.Writer) (err error) {
	itemsInFile, err := List(fileName)
	if err != nil {
		return err
	}
	var users []user

	err = json.Unmarshal([]byte(itemsInFile), &users)
	if err != nil {
		return err
	}

	for i, user := range users {
		if user.Id == id {
			users = append(users[:i], users[i+1:]...)
			break
		}

		if i == len(users)-1 {
			_, err = fmt.Fprintf(writer, userNotFound, id)
			if err != nil {
				return
			}
		}
	}

	data, err := json.Marshal(users)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = os.WriteFile(fileName, []byte(data), 0644)
	if err != nil {
		return err
	}

	return
}
