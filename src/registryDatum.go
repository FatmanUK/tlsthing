package main

import (
	"time"
	"gorm.io/gorm"
)

type RegistryDatum struct {
	Name string   `
	json:"Name"
	validate:"required"`
	Host string   `
	json:"Host"
	validate:"required"`
	Type string   `
	json:"Type"
	validate:"required"`
	CertPath string  `
	json:"CertPath"
	validate:"required"`
	KeyPath string  `
	json:"KeyPath"
	validate:"required"`
	RegDate time.Time  `
	json:"RegDate"
	sql:"DEFAULT:CURRENT_TIMESTAMP"
	validate:"required"`
	NotBefore time.Time  `
	json:"NotBefore"
	sql:"DEFAULT:NULL"
	validate:"required"`
	NotAfter time.Time  `
	json:"NotAfter"
	sql:"DEFAULT:NULL"
	validate:"required"`
}

func (re RegistryDatum) Duration() time.Duration {
	return re.NotAfter.Sub(re.NotBefore)
}

func (re RegistryDatum) Create(db interface{}) error {
	db.(*gorm.DB).Create(&re)
	return nil
}

func (re RegistryDatum) Update(db interface{}) error {
	db.(*gorm.DB).Where(
		"name = ? AND host = ?",
		re.Name, re.Host).Model(&re).Updates(re)
	return nil
}

func (re RegistryDatum) Delete(db interface{}) error {
	db.(*gorm.DB).Where(
		"name = ? AND host = ?",
		re.Name, re.Host).Delete(&re)
	return nil
}
