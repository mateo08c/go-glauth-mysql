package main

import (
	"errors"
	"github.com/joho/godotenv"
	"github.com/mateo08c/go-glauth-mysql/glauth"
	"github.com/mateo08c/go-glauth-mysql/glauth/ressources"
	"gorm.io/gorm"
	"log"
	"os"
	"testing"
)

var context *glauth.Context

// start on init
func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	log.Println("Loaded .env file")

	context = &glauth.Context{
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Hostname: os.Getenv("DB_HOSTNAME"),
		Port:     os.Getenv("DB_PORT"),
		Database: os.Getenv("DB_NAME"),
	}
}

func TestNew(t *testing.T) {
	client, err := glauth.New(context)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(client)

	//try to get the group users
	group, err := client.GetGroupByName("users")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(group)
}

func TestGroup(t *testing.T) {
	client, err := glauth.New(context)
	if err != nil {
		t.Fatal(err)
	}

	err = client.CreateGroup(&ressources.CreateGroup{
		Name: "test",
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Group created")

	//try to get the group
	group, err := client.GetGroupByName("test")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(group)

	//delete the group
	err = client.DeleteGroup(group.GIDNumber)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Group deleted")

	//try to get the group
	group, err = client.GetGroupByName("test")
	if err == nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			t.Log("Group not found")
		} else {
			t.Fatal(err)
		}
	}

	t.Log(group)

	if group == nil {
		t.Log("Group not found")
	}
}

func TestUser(t *testing.T) {
	client, err := glauth.New(context)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Helper function for error checking
	checkError := func(t *testing.T, err error, msg string) {
		t.Helper()
		if err != nil {
			t.Fatalf("%s: %v", msg, err)
		}
	}

	t.Run("CreateUser", func(t *testing.T) {
		err = client.CreateUser(&ressources.CreateUser{Name: "test"})
		checkError(t, err, "Failed to create user")
		t.Log("User created")
	})

	t.Run("GetUser", func(t *testing.T) {
		user, err := client.GetUserByName("test")
		checkError(t, err, "Failed to get user")
		t.Logf("User retrieved: %v", user)
	})

	t.Run("AddCapability", func(t *testing.T) {
		user, err := client.GetUserByName("test")
		checkError(t, err, "Failed to get user")

		err = client.CreateCapability(&ressources.Capability{
			UserID: user.UIDNumber,
			Action: "search",
			Object: "*",
		})
		checkError(t, err, "Failed to create capability")
		t.Log("Capability created")
	})

	t.Run("CheckCapabilities", func(t *testing.T) {
		user, err := client.GetUserByName("test")
		checkError(t, err, "Failed to get user")
		t.Logf("User capabilities: %v", user.Capabilities)
	})

	t.Run("UpdateUser", func(t *testing.T) {
		gn := "test2"
		err := client.UpdateUser("test", &ressources.UpdateUser{GivenName: &gn})
		checkError(t, err, "Failed to update user")
		t.Log("User updated")

		user, err := client.GetUserByName("test")
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			t.Fatalf("Unexpected error when getting updated user: %v", err)
		}

		if user.GivenName != "test2" {
			t.Fatalf("Expected user given name to be test2, got %s", user.GivenName)
		}
		t.Logf("User updated: %v", user)
	})

	t.Run("DeleteUser", func(t *testing.T) {
		user, err := client.GetUserByName("test")
		checkError(t, err, "Failed to get user before deletion")

		err = client.DeleteUser(user.UIDNumber)
		checkError(t, err, "Failed to delete user")
		t.Log("User deleted")

		user, err = client.GetUserByName("test")
		if err == nil || !errors.Is(err, gorm.ErrRecordNotFound) {
			t.Fatalf("Expected user to be not found after deletion, got %v", err)
		}
		t.Log("User not found after deletion")
	})
}
