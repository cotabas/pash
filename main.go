package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type entry struct {

  Service string
  Username string
  Password string
}

func addEntry() (string, string, string) {

  var service string
  var username string
  var password string
  
  fmt.Print("Service: ")
  fmt.Scan(&service)

  fmt.Print("Username: ")
  fmt.Scan(&username)

  fmt.Print("Password: ")
  fmt.Scan(&password)

  return service, username, password
}

func main() {
  
  test, _ := os.Create("tester.json")

  t := &entry{
    Service: "google",
    Username: "cotabas",
    Password: "qwerty"}

    mt := map[int]*entry{}

    mt[0] = t

    mt[1] = t

    jt, _ := json.Marshal(mt)

  fmt.Println(string(jt))

  test.Write(jt)

  fmt.Println("now read it back")

  reed, _ := os.ReadFile("tester.json")
  tj := map[int]entry{}

  fmt.Println(string(reed))

  json.Unmarshal(reed, &tj)

  fmt.Println(tj)

  iconv := tj[0]
  fmt.Println(iconv.Service)
}
