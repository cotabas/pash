package main

import "fmt"

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

  service, username, password := addEntry()

  fmt.Println("Service   Username   Password")
  fmt.Printf("%s %s %s", service, username, password)

}
