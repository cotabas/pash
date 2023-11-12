package main

import "fmt"

func main() {

  var service string
  var username string
  var password string


  fmt.Print("Service: ")
  fmt.Scan(&service)

  fmt.Print("Username: ")
  fmt.Scan(&username)

  fmt.Print("Password: ")
  fmt.Scan(&password)

  fmt.Println("Service   Username   Password")
  fmt.Printf("%s %s %s", service, username, password)
}
