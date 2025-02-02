package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"

	c "github.com/drizion/wabot-go/client"
	"github.com/drizion/wabot-go/command"
	"github.com/drizion/wabot-go/config"
	"github.com/drizion/wabot-go/handler"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal/v3"
)

func main() {

	config.SetupConfig()

	fmt.Println("Starting wabot...")
	fmt.Printf("Prefix: %+v\n", config.Bot.Prefix)
	fmt.Printf("OwnerNumbers: %+v\n", config.Bot.OwnerNumbers)

	// g := gen.NewGenerator(gen.Config{
	// 	OutPath: "./database/models",
	// 	// Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	// })

	// db, err := gorm.Open(postgres.New(postgres.Config{
	// 	DriverName: "postgres",
	// 	DSN:        "host=147.135.9.120 user=drizion password=--- dbname=wabot.net port=32768 sslmode=disable TimeZone=America/Sao_Paulo",
	// }), &gorm.Config{})

	// if err != nil {
	// 	panic(err)
	// }

	// // gormdb, _ := gorm.Open(mysql.Open("root:@(127.0.0.1:3306)/demo?charset=utf8mb4&parseTime=True&loc=Local"))
	// g.UseDB(db) // reuse your gorm db

	// // Generate basic type-safe DAO API for struct `model.User` following conventions

	// g.ApplyBasic(
	// 	// Generate struct `User` based on table `users`
	// 	g.GenerateModel("Ads"),
	// 	g.GenerateModel("BotUser"),
	// 	g.GenerateModel("Chatgpt"),
	// 	g.GenerateModel("DailyBonus"),
	// 	g.GenerateModel("EuNunca"),
	// 	g.GenerateModel("Groups"),
	// 	g.GenerateModel("SiteUser"),
	// 	g.GenerateModel("VerificationRequest"),
	// 	g.GenerateModel("CasinoSession"),
	// 	g.GenerateModel("CasinoBet"),
	// 	// g.GenerateAllTable(),

	// 	// Generate struct `Employee` based on table `users`
	// 	// g.GenerateModelAs("users", "Employee"),

	// 	// Generate struct `User` based on table `users` and generating options
	// 	// g.GenerateModel("users", gen.FieldIgnore("address"), gen.FieldType("id", "int64")),

	// 	// Generate struct `Customer` based on table `customer` and generating options
	// 	// customer table may have a tags column, it can be JSON type, gorm/gen tool can generate for your JSON data type
	// 	// g.GenerateModel("customer", gen.FieldType("tags", "datatypes.JSON")),
	// )
	// // g.ApplyBasic(
	// // Generate structs from all tables of current database
	// // g.GenerateAllTable()...,
	// // )
	// // Generate the code
	// g.Execute()

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
