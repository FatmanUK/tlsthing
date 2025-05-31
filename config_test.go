package main

import (
	"testing"
)

func initTestData() ConfigRepo {
	args := Args{
		Path: "/tmp/testrepo",
		Repo: "https://github.com/FatmanUK/tlsthing",
		Refresh: 30}
	re := ConfigRepo{args, internal_repo_path}
	return re
}

func TestValidate(t *testing.T) {
	//re := initTestData()
	//err := re.validate()
	// Test that validation works
	// supply faulty data and see if it's rejected
}

func TestHead(t *testing.T) {
	//re := initTestData()
	// Test that we get a 200 or 204 status
}

func TestIsCloned(t *testing.T) {
	//re := initTestData()
	// Test that the repo is cloned
}

func TestIsOkRepo(t *testing.T) {
	//re := initTestData()
	// Test that the test repo is ok
}

func TestClone(t *testing.T) {
	//re := initTestData()
	// Test that we can clone the test repo
}

func TestPull(t *testing.T) {
	//re := initTestData()
	// Test that we can pull on the test repo
}

func TestArchive(t *testing.T) {
	//re := initTestData()
	// Test that we can tar under random name
}

func TestErase(t *testing.T) {
	//re := initTestData()
	// Test that we can erase te test data
}

