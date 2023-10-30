package models

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func readDBConfig() string {
	var config struct {
		DB struct {
			Username string `json:"username"`
			Password string `json:"password"`
			DBName   string `json:"dbname"`
			Host     string `json:"host"`
			Port     string `json:"port"`
			SSLMode  string `json:"sslmode"`
		} `json:"db"`
	}

	file, err := os.Open("../.dsn")
	if err != nil {
		log.Fatalf("Failed to open secrets file: %s", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		log.Fatalf("Failed to decode secrets: %s", err)
	}

	dsn := "user=" + config.DB.Username +
		" password=" + config.DB.Password +
		" dbname=" + config.DB.DBName +
		" host=" + config.DB.Host +
		" port=" + config.DB.Port +
		" sslmode=" + config.DB.SSLMode

	return dsn
}

func setupTestDB() *gorm.DB {
	dsn := readDBConfig()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("Failed to connect to the database: %s", err)
	}

	// Run the migrations.
	db.AutoMigrate(&ShoppingList{}, &Item{}, &Store{})

	return db
}

/*
Test to check if database connection can be set up successfully.
*/
func TestMain(m *testing.M) {
	db = setupTestDB()

	exitCode := m.Run()
	os.Exit(exitCode)
}

/* ****** TESTS FOR SHOPPINGLIST OBJECT */
func TestCreateShoppingList(t *testing.T) {
	tx := db.Begin()
	defer tx.Rollback()

	// Create and commit to database.
	shoppingList := ShoppingList{Name: "Giant", Creator: "Soumya Sengupta"}
	result := tx.Create(&shoppingList)
	//tx.Commit()

	// Check if there were errors.
	assert.NoErrorf(t, result.Error, "Error creating shopping list: %s", result.Error)

	// Check if the shopping list id was set.
	assert.NotEqual(t, uint(0), shoppingList.ID, "Shopping list id was not set")

	var expected string
	var received string
	var msg string

	// Check the name of the newly inserted shopping list.
	expected = "Giant"
	received = shoppingList.Name
	msg = "Shopping list name was not set; expected: Giant, received: %v"
	assert.Equalf(t, expected, received, msg, expected, received)

	// Check the creator of the newly inserted shopping list.
	expected = "Soumya Sengupta"
	received = shoppingList.Creator
	msg = "Shopping list creator was not set; expected: Soumya Sengupta, received: %v"
	assert.Equalf(t, expected, received, msg, expected, received)
}

func TestCreateShoppingListNameNotNullConstraint(t *testing.T) {
	tx := db.Begin()
	defer tx.Rollback()

	shoppingList := ShoppingList{Name: "", Creator: "Soumya Sengupta"}
	result := tx.Create(&shoppingList)
	err := result.Error

	// Check if there were errors.
	if assert.Error(t, result.Error, "Name being null should not have been accepted") {
		// Check if the right error message was thrown.
		var expected string
		var received string
		var msg string

		expected = "chk_shopping_lists_name"
		received = err.Error()
		msg = "Unexpected violation flagged: expected %s, received %s"
		assert.Containsf(t, received, expected, msg, expected, received)
	}
}

func TestCreateShoppingListCreatorNotNullConstraint(t *testing.T) {
	tx := db.Begin()
	defer tx.Rollback()

	shoppingList := ShoppingList{Name: "Giant", Creator: ""}
	result := tx.Create(&shoppingList)
	err := result.Error

	// Check if there were errors.
	if assert.Error(t, result.Error, "Creator being null should not have been accepted") {
		// Check if the right error message was thrown.
		var expected string
		var received string
		var msg string

		expected = "chk_shopping_lists_creator"
		received = err.Error()
		msg = "Unexpected violation flagged: expected %s, received %s"
		assert.Containsf(t, received, expected, msg, expected, received)
	}
}

func TestCreateShoppingListAndOneItem(t *testing.T) {
	tx := db.Begin()
	defer tx.Rollback()

	// Create shopping list with 1 item and commit.
	shoppingList := ShoppingList{
		Name:    "Giant",
		Creator: "Soumya Sengupta",
		Items: []Item{
			{
				Name:     "Milk",
				Quantity: 2.0,
				Uom:      "L",
			},
		},
	}
	result := tx.Create(&shoppingList)
	tx.Commit()

	// Check if there were errors.
	assert.NoErrorf(t, result.Error, "Error creating shopping list: %s", result.Error)

	// Check if the shopping list id was set.
	assert.NotEqual(t, uint(0), shoppingList.ID, "Shopping list id was not set")

	var expected string
	var received string
	var msg string

	// Check if only 1 shopping list got created.
	var count int64
	db.Model(&ShoppingList{}).Where("name = ?", "Giant").Count(&count)
	expected = "1"
	received = strconv.Itoa(int(count))
	msg = "Wrong number of shopping lists got created; expected: %v, received: %v"
	assert.Equalf(t, expected, received, msg, expected, received)

	var items = shoppingList.Items

	// Check if only one item was created.
	expected = "1"
	received = strconv.Itoa(len(items))
	msg = "Wrong number of items got created; expected: %v, received: %v"
	assert.Equal(t, expected, received, msg, expected, received)

	var item = items[0]

	// Check if the item id was set.
	assert.NotEqual(t, uint(0), item.ID, "Item id was not set")

	// Check if the name of the item was set.
	expected = "Milk"
	received = item.Name
	msg = "Item name was not set; expected: %v, received: %v"
	assert.Equal(t, expected, received, msg, expected, received)

	// Check if the quantity of the item was set.
	expected = strconv.FormatFloat(2.0, 'f', -1, 64)
	received = strconv.FormatFloat(float64(item.Quantity), 'f', -1, 64)
	msg = "Item quantity was not set; expected: %v, received: %v"
	assert.Equal(t, expected, received, msg, expected, received)

	// Check if the uom of the item was set.
	expected = "L"
	received = item.Uom
	msg = "Item uom was not set; expected: %v, received: %v"
	assert.Equal(t, expected, received, msg, expected, received)
}

func TestCreateItemNameNotNullConstraint(t *testing.T) {

}

func TestCreateItemQuantityNotNegativeConstraint(t *testing.T) {

}

func TestCreateItemWithoutShoppingList(t *testing.T) {

}
