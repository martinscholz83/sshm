package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

type Connection struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	IP       string `json:"ip"`
	Port     *int   `json:"port"`
}

type Connections []Connection

type ConnectionManager struct {
	Connections []Connection
	Path        string
}

const folder = ".sshm"
const fileName = "connections.json"

func newConnectionManager() *ConnectionManager {
	usr, _ := user.Current()
	configFile := filepath.Join(usr.HomeDir, folder, fileName)
	_, err := os.Stat(configFile)
	if os.IsNotExist(err) {
		os.MkdirAll(filepath.Join(usr.HomeDir, folder), 0700)
		os.Create(configFile)
	}

	return &ConnectionManager{
		Path: configFile,
	}
}

func (c ConnectionManager) LoadConnections() ([]Connection, error) {
	file, err := os.Open(c.Path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	bytes, err := ioutil.ReadAll(file)

	if err != nil {
		log.Println(err)
	}

	var connections Connections
	json.Unmarshal([]byte(bytes), &connections)

	return connections, nil
}

func (c ConnectionManager) AddConnection(con Connection) error {
	fileIn, err := os.OpenFile(c.Path, os.O_RDONLY, 0755)
	if err != nil {
		return err
	}

	bytes, err := ioutil.ReadAll(fileIn)
	if err != nil {
		log.Fatal(err)
	}

	if err := fileIn.Close(); err != nil {
		log.Fatal(err)
	}

	fileOut, err := os.OpenFile(c.Path, os.O_RDWR, 0755)
	if err != nil {
		return err
	}

	var connections Connections
	json.Unmarshal([]byte(bytes), &connections)

	connections = append(connections, con)

	js, err := json.MarshalIndent(connections, "", "    ")
	if err != nil {
		return err
	}

	_, err = fileOut.Write(js)

	if err != nil {
		return err
	}

	if err := fileOut.Close(); err != nil {
		log.Fatal(err)
	}

	return nil
}

func (c ConnectionManager) DeletConnection(id int) error {
	fileIn, err := os.OpenFile(c.Path, os.O_RDONLY, 0755)
	if err != nil {
		return err
	}

	bytes, err := ioutil.ReadAll(fileIn)
	if err != nil {
		log.Fatal(err)
	}

	if err := fileIn.Close(); err != nil {
		log.Fatal(err)
	}

	// delete all content because new content is smaller then existing
	if err := os.Truncate(c.Path, 0); err != nil {
		return err
	}

	fileOut, err := os.OpenFile(c.Path, os.O_RDWR, 0644)
	if err != nil {
		return err
	}

	var conns Connections
	json.Unmarshal([]byte(bytes), &conns)

	conns = append(conns[:id], conns[id+1:]...)

	js, err := json.MarshalIndent(conns, "", "    ")
	if err != nil {
		return err
	}
	log.Println(conns)

	_, err = fileOut.Write(js)

	if err != nil {
		return err
	}

	if err := fileOut.Close(); err != nil {
		log.Fatal(err)
	}

	return nil
}
