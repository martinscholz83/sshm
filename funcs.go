package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/urfave/cli/v2"
)

var cm = newConnectionManager()

func ListEntries(c *cli.Context) error {
	conn, err := cm.LoadConnections()
	if err != nil {
		return err
	}

	for i := 0; i < len(conn); i++ {
		fmt.Printf("++++++++Connection %d++++++++\n", i)
		printConnection(&conn[i])
	}
	return nil
}

func NewConnection(c *cli.Context) error {
	conn := Connection{}

	fmt.Println("Add new connection")
	fmt.Println("Name:")
	input := bufio.NewScanner(os.Stdin)

	input.Scan()

	if len(input.Text()) <= 0 {
		return errors.New("please enter a valid Name")
	}

	conn.Name = input.Text()

	fmt.Println("Username:")
	input.Scan()

	if len(input.Text()) <= 0 {
		return errors.New("please enter a valid Username")
	}

	conn.Username = input.Text()

	fmt.Println("IP:")
	input.Scan()

	if len(input.Text()) <= 0 {
		return errors.New("please enter a valid IP address")
	}

	conn.IP = input.Text()

	err := cm.AddConnection(conn)

	if err != nil {
		return err
	}

	return nil
}

func DeletConnection(c *cli.Context) error {
	conns, err := cm.LoadConnections()
	if err != nil {
		return err
	}

	fmt.Println("Please select a connection to delete")
	fmt.Println("")
	for i := 0; i < len(conns); i++ {
		fmt.Printf("++++++++Connection %d++++++++\n", i)
		printConnection(&conns[i])
	}
	fmt.Println("")
	fmt.Println("Connection:")
	input := bufio.NewScanner(os.Stdin)

	input.Scan()

	if len(input.Text()) <= 0 {
		return errors.New("please enter a valid id")
	}

	id, err := strconv.Atoi(input.Text())
	if err != nil {
		return err
	}

	if err = cm.DeletConnection(id); err != nil {
		return err
	}

	fmt.Println("successfully deleted connection")

	return nil
}

func Connect(c *cli.Context) error {
	conns, err := cm.LoadConnections()
	if err != nil {
		return err
	}

	fmt.Println("please select a connection")
	fmt.Println("")
	for i := 0; i < len(conns); i++ {
		fmt.Printf("++++++++Connection %d++++++++\n", i)
		printConnection(&conns[i])
	}
	fmt.Println("")
	fmt.Println("connection number:")
	input := bufio.NewScanner(os.Stdin)

	input.Scan()

	if len(input.Text()) <= 0 {
		return errors.New("please enter a valid Name")
	}

	id, err := strconv.Atoi(input.Text())
	if err != nil {
		return err
	}

	con := conns[id]

	com := fmt.Sprintf("ssh -l %s %s", con.Username, con.IP)
	log.Println(com)
	cmd := exec.Command("bash", "-c", com)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	err = cmd.Wait()
	if err != nil {
		return err
	}

	log.Println("closed connection")
	return nil
}

func printConnection(con *Connection) {
	fmt.Printf("Name: %s\n", con.Name)
	fmt.Printf("Username: %s\n", con.Username)
	fmt.Printf("IP: %s\n", con.IP)
	fmt.Println("++++++++++++++++++++++++++++")
	fmt.Println("")
}
