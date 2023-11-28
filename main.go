package main

import (
	"bufio"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var SYS_COPY string = "wl-copy"

type entry struct {
  Service string
  Username []byte 
  Password []byte
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
    ent.Username = []byte(scanner.Text())
  }

  fmt.Print("Password: ")
  if scanner.Scan() {
    ent.Password = []byte(scanner.Text())
  }
}

func removeEntry(fileMap map[int]entry, enteredNumber string, fileName string) {
  remove, err := strconv.Atoi(enteredNumber)
  check(err)
  file, err := os.Create(fileName + ".json")
  check(err)

  newMap := make(map[int]entry)

  delete(fileMap, remove - 1)

  for i := 0; i < len(fileMap); i++ {
    if i >= remove - 1 { 
      newMap[i] = fileMap[i + 1]
    } else {
      newMap[i] = fileMap[i]
    }
  }
  writeBytes, err := json.Marshal(newMap)
  check(err)

  file.Write(writeBytes)
}

func encryptPass(ent *entry, publicKey rsa.PublicKey) {
  encryptedBytes, err := rsa.EncryptOAEP(
    sha256.New(),
    rand.Reader,
    &publicKey,
    ent.Password,
    nil)
  check(err)
  ent.Password = encryptedBytes
}

func encryptUser(ent *entry, publicKey rsa.PublicKey) {
  encryptedBytes, err := rsa.EncryptOAEP(
    sha256.New(),
    rand.Reader,
    &publicKey,
    ent.Username,
    nil)
  check(err)
  ent.Username = encryptedBytes
}

func addFile(fileName string) {

  _, err := os.Create(fileName + ".json")
  check(err)
  
  keyFile, err := os.Create(fileName)
  check(err)
  
  privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
  check(err)

  pkPem := pem.EncodeToMemory(
    &pem.Block{
      Type:  "RSA PRIVATE KEY",
      Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
    },
  )

  keyFile.Write(pkPem)
}

func showEntries(store map[int]entry, privateKey *rsa.PrivateKey) {
  fmt.Println(".  Service             Username            Password")
  for i := 0; i < len(store); i++ {
    spacer := ""
    for j := len(store[i].Service); j < 20; j++ { spacer += " " }
    fmt.Print(fmt.Sprint(i + 1) + ". " + store[i].Service + spacer)

    spacer = ""
    userBytes, err := privateKey.Decrypt(nil, store[i].Username, &rsa.OAEPOptions{Hash: crypto.SHA256})
    check(err)
    for j := len(userBytes); j < 20; j++ { spacer += " " }
    fmt.Print(string(userBytes) + spacer)

    passBytes, err := privateKey.Decrypt(nil, store[i].Password, &rsa.OAEPOptions{Hash: crypto.SHA256})
    check(err)
    fmt.Println(string(passBytes))
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

func getPrivateKey(fileName string) (*rsa.PrivateKey) {

  pkPem, err := os.ReadFile(fileName)
  check(err)

  block, _ := pem.Decode(pkPem)

  privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
  check(err)

  return privateKey
}

func copyPass(entryNumber string, fileMap map[int]entry, privateKey *rsa.PrivateKey) {
  num, err := strconv.Atoi(entryNumber)
  check(err)
  pass := fileMap[num - 1].Password
  passBytes, err := privateKey.Decrypt(nil, pass, &rsa.OAEPOptions{Hash: crypto.SHA256})
  check(err)
  copyCmd := exec.Command(SYS_COPY, string(passBytes))
  err = copyCmd.Run()
  check(err)
}

func loginMenu(fileName string, privateKey *rsa.PrivateKey) {
  
  var choice string
  scanner := bufio.NewScanner(os.Stdin)
  fileMap := readFile(fileName)

  publicKey := privateKey.PublicKey

  fmt.Println(" .:" + fileName + ":. ")
  showEntries(fileMap, privateKey)

  fmt.Println("\n\n\n1. Add new")
  fmt.Println("2. Copy to Clipboard")
  fmt.Println("3. Remove")
  fmt.Println("4. Change")
  fmt.Println("5. Exit")
  fmt.Print("Select: ")
  
  if scanner.Scan() {
    choice = scanner.Text()
  }

  switch choice {
  case "1":
    var newEntry entry
    addEntry(scanner, &newEntry)
    encryptPass(&newEntry, publicKey)
    encryptUser(&newEntry, publicKey)
    writeFile(fileName, newEntry)
    loginMenu(fileName, privateKey)
  case "2":
    fmt.Print("Number: ")
    if scanner.Scan() {
      copyPass(scanner.Text(), fileMap, privateKey)
      loginMenu(fileName, privateKey)
    }
  case "3":
    fmt.Print("Number: ")
    if scanner.Scan() {
      removeEntry(fileMap, scanner.Text(), fileName)
      loginMenu(fileName, privateKey)
    }
  }

}

func main() {

  pemFile := os.Args[1:]

  if len(pemFile) == 0 {
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
        privateKey := getPrivateKey(scanner.Text())
        loginMenu(scanner.Text(), privateKey)
      }
    case "2":
      fmt.Println("Create new table")
      fmt.Print("Table name: ")
      if scanner.Scan() {
        addFile(scanner.Text())
        privateKey := getPrivateKey(scanner.Text())
        loginMenu(scanner.Text(), privateKey)
      }
    }
  }
  if len(pemFile) != 0 {
    privateKey := getPrivateKey(pemFile[0])

    par := strings.FieldsFunc(pemFile[0], func(r rune) bool {
      if r == '/' { return true }
      return false
    })

    fileName := par[len(par) - 1]
    loginMenu(fileName, privateKey)

  }

}
