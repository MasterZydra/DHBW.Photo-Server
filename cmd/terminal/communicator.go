package main

import (
	"bufio"
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
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Please enter parameters: ")

	fmt.Print("Username: ")
	uname, _ := reader.ReadString('\n')
	uname = strings.TrimSuffix(uname, "\n")
	//check that the username contains only numbers and letters
	validateUsername(uname)

	fmt.Print("Password: ")
	pswd, _ := reader.ReadString('\n')
	pswd = strings.TrimSuffix(pswd, "\n")

	fmt.Print("Host of the server: ")
	host, _ := reader.ReadString('\n')
	host = strings.TrimSuffix(host, "\n")
	// check if host is valid
	validateHost(host)

	fmt.Print("Path of the folder with the photos: ")
	path, _ := reader.ReadString('\n')
	path = strings.TrimSuffix(path, "\n")
	var userInput UserInput
	userInput.Username = uname
	userInput.Password = pswd
	userInput.Host = host
	userInput.Path = path
	
	return userInput
}

func validateUsername(uname string) {
	for _, c := range uname {
		if !unicode.IsDigit(c) && !unicode.IsLetter(c) {
			log.Fatal("The username must only contain letters and numbers.")
		}
	}
}

func validateHost(host string) {
	_, err := url.ParseRequestURI(host)
	if err != nil {
		log.Fatal("Invalid Host" + err.Error())
	}
}

func WriteMessage(msg string) {
	separator := "\n--------------\n"
	fmt.Print(separator, msg, separator)
}