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
)

type entry struct {
  Service string
  Username string
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
    ent.Username = scanner.Text()
  }

  fmt.Print("Password: ")
  if scanner.Scan() {
    ent.Password = []byte(scanner.Text())
  }
}

func encryptPass(ent *entry, publicKey rsa.PublicKey) {
  encryptedBytes, err := rsa.EncryptOAEP(
    sha256.New(),
    rand.Reader,
    &publicKey,
    []byte(ent.Password),
    nil)
  check(err)
  ent.Password = encryptedBytes
}


func addFile(fileName string) {

  _, err := os.Create(fileName + ".json")
  check(err)
  
  keyFile, err := os.Create(fileName + ".pem")
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
//  var out strings.Builder
//  keyFile, err := os.Create(fileName + ".key")
//  check(err)
//
//  createKey := exec.Command("openssl", "genrsa", "4096")
//  createKey.Stdin = strings.NewReader("")
//  createKey.Stdout = &out
//
//  err = createKey.Run()
//  check(err)
//
//  keyFile.Write([]byte(out.String()))
}

func showEntries(store map[int]entry, privateKey *rsa.PrivateKey) {
  fmt.Println(".  Service             Username            Password")
  for i := 0; i < len(store); i++ {
    spacer := ""
    for j := len(store[i].Service); j < 20; j++ { spacer += " " }
    fmt.Print(fmt.Sprint(i + 1) + ". " + store[i].Service + spacer)

    spacer = ""
    for j := len(store[i].Username); j < 20; j++ { spacer += " " }
    fmt.Print(store[i].Username + spacer)

    decryptedBytes, err := privateKey.Decrypt(nil, store[i].Password, &rsa.OAEPOptions{Hash: crypto.SHA256})
    check(err)
    fmt.Println(string(decryptedBytes))
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

  pkPem, err := os.ReadFile(fileName + ".pem")
  check(err)

  block, _ := pem.Decode(pkPem)

  privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
  check(err)

  return privateKey
}

func loginMenu(fileName string) {
  
  var choice string
  scanner := bufio.NewScanner(os.Stdin)
  fileMap := readFile(fileName)

  privateKey := getPrivateKey(fileName)
  publicKey := privateKey.PublicKey

  fmt.Println(" .:" + fileName + ":. ")
  showEntries(fileMap, privateKey)

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
    encryptPass(&newEntry, publicKey)
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

  //fileName := "t"

  privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
  check(err)

  publicKey := privateKey.PublicKey

  encryptedBytes, err := rsa.EncryptOAEP(
    sha256.New(),
    rand.Reader,
    &publicKey,
    []byte("super secret message"),
    nil)
    check(err)


	fmt.Println("encrypted bytes: ", encryptedBytes)

  decryptedBytes, err := privateKey.Decrypt(nil, encryptedBytes, &rsa.OAEPOptions{Hash: crypto.SHA256})

  check(err)
  fmt.Println("decrypted message: ", string(decryptedBytes))

  pkPem := pem.EncodeToMemory(
    &pem.Block{
      Type:  "RSA PRIVATE KEY",
      Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
    },
  )

  fmt.Println(pkPem)

  kf, _ := os.Create("t.pem")

  kf.Write(pkPem)

      block, _ := pem.Decode([]byte(pkPem))

    priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)

    fmt.Println(priv)

  decryptedBytes, err = priv.Decrypt(nil, encryptedBytes, &rsa.OAEPOptions{Hash: crypto.SHA256})

  check(err)
  fmt.Println("decrypted message: ", string(decryptedBytes))
//
//  createKey := exec.Command("openssl", "rsa", "-text", "-in", fileName + ".key", "-noout")
//
//  createKey.Stdin = strings.NewReader("")
//
//  var out strings.Builder
//
//  createKey.Stdout = &out
//
//  err := createKey.Run()
//  check(err)
//
//  s := strings.Fields(out.String())
//  fmt.Println(s)
//
//  hex := strings.Join(s[6:11], "")
//  hex = strings.ReplaceAll(hex, ":", "")
//  fmt.Println(hex)
//  modulus, _ := new(big.Int).SetString(hex, 16)
//
//  hex = strings.Join(s[15:20], "")
//  hex = strings.ReplaceAll(hex, ":", "")
//  fmt.Println(hex)
//  private, _ := new(big.Int).SetString(hex, 16)
//  fmt.Println(modulus)
//  fmt.Println(private)
//  var public float64 = 65537
//
//  pow := new(big.Int)
//  pow = big.NewInt(int64(math.Pow(65, public)))
//  _, answer := new(big.Int).DivMod(pow, modulus, new(big.Int)) 
//
//  fmt.Println(answer)
  
  //so slow..
  //pexp := new(big.Int).Exp(answer, private, nil)


  //_, answer = new(big.Int).DivMod(pexp, modulus, new(big.Int))



}
