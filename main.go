package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"

	c "github.com/drizion/wabot-go/client"
	"github.com/drizion/wabot-go/command"
	"github.com/drizion/wabot-go/handler"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal/v3"
)

func main() {
	dbLog := waLog.Stdout("Database", "ERROR", true)
	// Make sure you add appropriate DB connector imports, e.g. github.com/mattn/go-sqlite3 for SQLite
	container, err := sqlstore.New("sqlite3", "file:database.db?_foreign_keys=on", dbLog)
	if err != nil {
		panic(err)
	}
	// If you want multiple sessions, remember their JIDs and use .GetDevice(jid) or .GetAllDevices() instead.
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		panic(err)
	}
	clientLog := waLog.Stdout("client", "ERROR", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)
	c.Wabot = client

	client.AddEventHandler(handler.EventHandler)

	if client.Store.ID == nil {
		// No ID stored, new login
		qrChan, _ := client.GetQRChannel(context.Background())
		err = client.Connect()
		if err != nil {
			panic(err)
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				// Render the QR code here
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
			} else {
				fmt.Println("Login event:", evt.Event)
			}

		}
	} else {
		// Already logged in, just connect
		fmt.Println("Connecting...")
		command.SetupCommands()
		err = client.Connect()
		if err != nil {
			panic(err)
		}
		fmt.Println("Connected!")
	}

	// Listen to Ctrl+C (you can also do something else that prevents the program from exiting)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	client.Disconnect()
}
