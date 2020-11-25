package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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
	fmt.Print("Password: ")
	pswd, _ := reader.ReadString('\n')
	fmt.Print("Host of the server: ")
	host, _ := reader.ReadString('\n')
	fmt.Print("Path of the folder with the photos: ")
	path, _ := reader.ReadString('\n')
	path = strings.TrimSuffix(path, "\n")
	var userInput UserInput
	userInput.Username = uname
	userInput.Password = pswd
	userInput.Host = host
	userInput.Path = path
	// TODO: validate Input (so everyone can enter anything he wants -> not good)
	return userInput
}

func WriteMessage(msg string) {
	separator := "\n--------------\n"
	fmt.Print(separator, msg, separator)
}