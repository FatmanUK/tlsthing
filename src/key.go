package main

import (
	//"os"
	//"log"
	"time"
	//"strconv"
	"gorm.io/gorm"
	//"path/filepath"
	//"github.com/glebarez/sqlite" // pure Go?
	//"github.com/FatmanUK/fatgo/doh"
)

// private and public keys
// not sure what the fields will be yet
// status 0=current, 1=expired, 2=error state etc.
type Key struct {
	gorm.Model
	Name string
	Type int
	CreationDate time.Time `sql:"DEFAULT:CURRENT_TIMESTAMP"`
	ExpiryDate time.Time `sql:"DEFAULT:NULL"`
	Ident []byte
	Status int
}

// What's stored in the database?

// Not Vault details; those can be file config/env config. Obviously, database details also in config.

// In database:
/*
  type      format   expiry
  key/cert  pem/pem  2025-12-01 06:43:02
  jks       pkcs1/pem
  jks       pkcs8/pem

pkcs1 = rsa only privkey, der, can be enc/unenc
pkcs8 = any algo privkey, der, can be enc/unenc
pkcs12/pfx = can use pkcs8?

keep it simple.

for now just have a "typeId" field for future foreign keys.
*/

/*
1   package main
  1 
  2 import (
  3         "gorm.io/gorm"
  4         "github.com/glebarez/sqlite" // pure Go?
  5 )
  6 
  7 type Page struct {
  8         gorm.Model
  9         Title   string
 10         Body    []byte
 11 }
 12 
 13 func (re *Page) Debug() string {
 14         output := `
14 ## Pages
 13 ___`
 12         return output
 11 }
 10 
  9 var db *gorm.DB
  8 
  7 func (re *Page) OpenDatabase(dbfile string) {
  6         d, err := gorm.Open(sqlite.Open(dbfile), &gorm.Config{})
  5         if err != nil {
  4                 panic(err.Error())
  3         }
  2         db = d
  1         db.AutoMigrate(&Page{})
  9 }
  8 
  7 func (re *Page) LoadPage(title string) error {
  6         result := db.Where("title = ?", title).Last(re)
  5         return result.Error
  4 }
  3 
  2 func (re *Page) Save() {
  1         db.Create(re)
39  }


*/
