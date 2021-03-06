package db

import (
	"strings"
	"testing"

	"github.com/syncore/a2sapi/src/constants"
	"github.com/syncore/a2sapi/src/models"

	_ "github.com/mattn/go-sqlite3"
)

var testData map[string]string

func init() {
	testData = make(map[string]string, 2)
	testData["10.0.0.10"] = "Reflex"
	testData["172.16.0.1"] = "QuakeLive"
}

func TestCreateServerDBtable(t *testing.T) {
	err := createServerDBtable(constants.TestServerDbFilePath)
	if err != nil {
		t.Fatalf("Unable to create test DB file: %s", err)
	}
}

func TestAddServersToDB(t *testing.T) {
	db, err := OpenServerDB()
	if err != nil {
		t.Fatalf("Unable to open test database: %s", err)
	}
	defer db.Close()
	db.AddServersToDB(testData)
}

func TestGetIDsForServerList(t *testing.T) {
	c := make(chan map[string]int64, 2)
	db, err := OpenServerDB()
	if err != nil {
		t.Fatalf("Unable to open test database: %s", err)
	}
	defer db.Close()
	db.GetIDsForServerList(c, testData)
	result := <-c
	if len(result) != 2 {
		t.Fatalf("Expected 2 results, got: %d", len(result))
	}
	if _, ok := result["10.0.0.10"]; !ok {
		t.Fatalf("Expected value 10.0.0.10 to exist.")
	}
	if _, ok := result["172.16.0.1"]; !ok {
		t.Fatalf("Expected value 172.16.0.1 to exist.")
	}
}

func TestGetIDsAPIQuery(t *testing.T) {
	c1 := make(chan *models.DbServerID, 1)
	c2 := make(chan *models.DbServerID, 1)
	db, err := OpenServerDB()
	if err != nil {
		t.Fatalf("Unable to open test database: %s", err)
	}
	defer db.Close()
	h1 := []string{"10.0.0.10"}
	h2 := []string{"172.16.0.1"}
	db.GetIDsAPIQuery(c1, h1)
	r1 := <-c1
	if len(r1.Servers) != 1 {
		t.Fatalf("Expected 1 server, got: %d", len(r1.Servers))
	}
	if !strings.EqualFold(r1.Servers[0].Game, "Reflex") {
		t.Fatalf("Expected result 1 to be Reflex, got: %v", r1.Servers[0].Game)
	}
	db.GetIDsAPIQuery(c2, h2)
	r2 := <-c2
	if len(r2.Servers) != 1 {
		t.Fatalf("Expected 1 server, got: %d", len(r2.Servers))
	}
	if !strings.EqualFold(r2.Servers[0].Game, "QuakeLive") {
		t.Fatalf("Expected result 2 to be QuakeLive, got: %v", r2.Servers[0].Game)
	}
}

func TestGetHostsAndGameFromIDAPIQuery(t *testing.T) {
	c := make(chan map[string]string, 2)
	db, err := OpenServerDB()
	if err != nil {
		t.Fatalf("Unable to open test database: %s", err)
	}
	defer db.Close()
	ids := []string{"1", "2"}
	db.GetHostsAndGameFromIDAPIQuery(c, ids)
	result := <-c
	if len(result) != 2 {
		t.Fatalf("Expected 2 results, got: %d", len(result))
	}
	if !strings.EqualFold(result["10.0.0.10"], "Reflex") {
		t.Fatalf("Expected result Reflex, got: %v", result["10.0.0.10"])
	}
	if !strings.EqualFold(result["172.16.0.1"], "QuakeLive") {
		t.Fatalf("Expected result QuakeLive, got: %v", result["1172.16.0.1"])
	}
}
