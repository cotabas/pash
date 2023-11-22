package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type entry struct {
  Service string
  Username string
  Password string
}

func check(e error) {
  if e != nil {
    panic(e)
  }
}

func addEntry(scanner *bufio.Scanner, ent *entry) {
  fmt.Print("Service: ")
  if scanner.Scan() {
    ent.Service = scanner.Text()
  }

  fmt.Print("Username: ")
  if scanner.Scan() {
    ent.Username = scanner.Text()
  }

  fmt.Print("Password: ")
  if scanner.Scan() {
    ent.Password = scanner.Text()
  }
}

func addFile(fileName string) {
  var out strings.Builder

  _, err := os.Create(fileName + ".json")
  check(err)
  
  keyFile, err := os.Create(fileName + ".key")
  check(err)

  createKey := exec.Command("openssl", "genrsa", "4096")
  createKey.Stdin = strings.NewReader("")
  createKey.Stdout = &out

  err = createKey.Run()
  check(err)

  keyFile.Write([]byte(out.String()))
}

func showEntries(store map[int]entry) {
  fmt.Println(".  Service             Username            Password")
  for i := 0; i < len(store); i++ {
    spacer := ""
    for j := len(store[i].Service); j < 20; j++ { spacer += " " }
    fmt.Print(fmt.Sprint(i + 1) + ". " + store[i].Service + spacer)

    spacer = ""
    for j := len(store[i].Username); j < 20; j++ { spacer += " " }
    fmt.Print(store[i].Username + spacer)

    fmt.Println(store[i].Password)
  }
} 

func readFile(fileName string) (map[int]entry) {
  file, err := os.ReadFile(fileName + ".json")
  check(err)
  store := map[int]entry{}
  json.Unmarshal(file, &store)

  return store
}

func writeFile(fileName string, addition entry) {
  store := readFile(fileName)

  file, err := os.Create(fileName + ".json")
  check(err)
  
  store[len(store)] = addition

  writeByte, err := json.Marshal(store)
  check(err)

  file.Write(writeByte)
  
}

func loginMenu(fileName string) {
  
  var choice string
  scanner := bufio.NewScanner(os.Stdin)
  fileMap := readFile(fileName)

  fmt.Println(" .:" + fileName + ":. ")
  showEntries(fileMap)

  fmt.Println("\n\n\n1. Add new")
  fmt.Println("2. Change")
  fmt.Println("3. Remove")
  fmt.Print("Select: ")
  
  if scanner.Scan() {
    choice = scanner.Text()
  }

  switch choice {
  case "1":
    var newEntry entry
    addEntry(scanner, &newEntry)
    writeFile(fileName, newEntry)
    loginMenu(fileName)
  }

}

func main() {
  
  scanner := bufio.NewScanner(os.Stdin)
  var choice string

  fmt.Println(".:Pash:.")
  fmt.Println("1. Login")
  fmt.Println("2. New table")
  fmt.Println("3. Exit")

  fmt.Print("Select: ")
  if scanner.Scan() {
    choice = scanner.Text()
  }

  switch choice {
  case "1":
    fmt.Println("Login")
    fmt.Print("Table name: ")
    if scanner.Scan() {
      loginMenu(scanner.Text())
    }
  case "2":
    fmt.Println("Create new table")
    fmt.Print("Table name: ")
    if scanner.Scan() {
      addFile(scanner.Text())
      loginMenu(scanner.Text())
    }
  }

  fileName := "blah"

  createKey := exec.Command("openssl", "rsa", "-text", "-in", fileName + ".key", "-noout")
  //("openssl rsa -text file.key")

  createKey.Stdin = strings.NewReader("")

  var out strings.Builder

  createKey.Stdout = &out

  err := createKey.Run()
  check(err)

  fmt.Println(out.String())
}
