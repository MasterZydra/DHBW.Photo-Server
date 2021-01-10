/*
 * DHBW Mosbach project of subject "Programmieren 2" from:
 * 6439456
 * 8093702
 * 9752762
 */

package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"unicode"
)

func Greet() {
	asciiArt :=
		`
==================================================================================================================================
.______    __    __    ______   .___________.  ______             _______. _______ .______     ____    ____  _______ .______
|   _  \  |  |  |  |  /  __  \  |           | /  __  \           /       ||   ____||   _  \    \   \  /   / |   ____||   _  \
|  |_)  | |  |__|  | |  |  |  | '---|  |----'|  |  |  |  ______ |   (----'|  |__   |  |_)  |    \   \/   /  |  |__   |  |_)  |    
|   ___/  |   __   | |  |  |  |     |  |     |  |  |  | |______| \   \    |   __|  |      /      \      /   |   __|  |      /     
|  |      |  |  |  | |  '--'  |     |  |     |  '--'  |      .----)   |   |  |____ |  |\  \----.  \    /    |  |____ |  |\  \----.
| _|      |__|  |__|  \______/      |__|      \______/       |_______/    |_______|| _| '._____|   \__/     |_______|| _| '._____|

==================================================================================================================================
`
	fmt.Println(asciiArt)
}

func WaitForUserInput() UserInput {
	// create new reader to read user input
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Please enter parameters: ")

	// first parameter: username
	fmt.Print("Username: ")
	// user input is read until a new line is entered
	uname, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Error while reading username: %v", err)
	}
	// newline control character must be removed from the user input
	uname = strings.TrimSuffix(uname, "\n")
	//check that the username contains only numbers and letters
	err = validateUsername(uname)
	if err != nil {
		log.Fatal(err)
	}

	// second parameter: password
	// reading the remaining parameters works like reading the user name (up to the new line, remove new line)
	fmt.Print("Password: ")
	pswd, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Error while reading password: %v", err)
	}
	pswd = strings.TrimSuffix(pswd, "\n")

	// third parameter
	fmt.Print("Host of the server: ")
	host, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Error while reading password: %v", err)
	}
	host = strings.TrimSuffix(host, "\n")
	// check if host is valid
	err = validateHost(host)
	if err != nil {
		log.Fatal(err)
	}

	// fourth parameter
	fmt.Print("Path of the folder with the photos: ")
	path, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Error while reading password: %v", err)
	}
	path = strings.TrimSuffix(path, "\n")

	// create and return userInput variable filled with user input
	return UserInput{
		Username: uname,
		Password: pswd,
		Host:     host,
		Path:     path,
	}
}

func validateUsername(uname string) error {
	for _, c := range uname {
		// each character of the username should be a number or a letter
		if !unicode.IsDigit(c) && !unicode.IsLetter(c) {
			return errors.New("The username must only contain letters and numbers")
		}
	}
	return nil
}

func validateHost(host string) error {
	_, err := url.ParseRequestURI(host)
	if err != nil {
		return errors.New("Invalid Host" + err.Error())
	}
	return nil
}

func WriteMessage(msg string) {
	separator := "\n--------------\n"
	fmt.Print(separator, msg, separator)
}
