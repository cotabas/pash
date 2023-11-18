package main

import (
	//	"crypto/rsa"
	//	"encoding/asn1"
	//	"errors"
	"encoding/json"
	"fmt"
	"os"
	// "math/big"
)

// type pkcs1PublicKey struct {
// 	N *big.Int
// 	E int
// }
//
//
// func ParsePKCS1PublicKey(der []byte) (*rsa.PublicKey, error) {
// 	var pub pkcs1PublicKey
// 	rest, _ := asn1.Unmarshal(der, &pub)
// 	if len(rest) > 0 {
// 		return nil, asn1.SyntaxError{Msg: "trailing data"}
// 	}
//
// 	if pub.N.Sign() <= 0 || pub.E <= 0 {
// 		return nil, errors.New("x509: public key contains zero or negative value")
// 	}
// 	if pub.E > 1<<31-1 {
// 		return nil, errors.New("x509: public key contains large public exponent")
// 	}
//
// 	return &rsa.PublicKey{
// 		E: pub.E,
// 		N: pub.N,
// 	}, nil
// }

type entry struct {

  Service string
  Username string
  Password string
}

type database struct {

  Name string
  Entries entry

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
  //test.Write([]byte("hello"))

  t := &entry{
    Service: "google",
    Username: "cotabas",
    Password: "qwerty"}

    mt := map[int]*entry{}

    mt[0] = t

    mt[1] = t

    jt, _ := json.Marshal(mt)

  fmt.Print(string(jt))
  test.Write(jt)
  test.Write(jt)

  reed, _ := os.ReadFile("tester.json")

  //fmt.Print(reed)

  var tj = entry{}

  fmt.Print(string(reed))

  json.Unmarshal(reed, &tj)

  fmt.Print(tj)
  //service, username, password := addEntry()

  //fmt.Println("Service   Username   Password")
  //fmt.Printf("%s %s %s", service, username, password)
  //key, _ := os.ReadFile("/home/cptmo/test.pub")
  
  //fmt.Print(key)
  // e, n := ParsePKCS1PublicKey([]byte("AAAAB3NzaC1yc2EAAAADAQABAAABgQC3P8rM8f+YmtHHLNAEmHJ0XL5IRW5LfzK8jtOvfHIxkViubetC9utH0w339KJCi3tTzp/Eu8YJybfMtBHjz99EBcI1DkfVoOJeXYf10/A4DvAYyHxsuVfi5czNUTECVaBkOhJykC5/oDsdpy1GsmZuJmqCESaTnvyQoLApXQHW01wguW2AIM0JZ1ujqAu+S7Iy1yoi5TbxVbh0/SoYkkFLeDAq6QKnRlAN/r8voShv6l/aM9/H/kGxsyEJgGmJfSPe+2jZqtiviKMCQy9kXzFXPL5vWhrFWv16C31vIObxgDGERDF2pebbFPf5Tmy85ijvQwmiKCz5CT5xfN3ePZfnzI64fdQ05PHHQv7gORogN5x8PXAh5tfIUHdZE2w8OfJvRqSNEfqKv6HlGhhWAjebUe0eEopuinC0TLiM3FJNZNlmp+5CvgP8OvzH0GGEuq8pksHVx4lzxOHbGNMU7sZYpJ45IPqvpXgCTj0RPtAMc6ZGgGvyooi1tgEGBMVfWas="))
//  fmt.Printf("%x %x", e, n)
}
