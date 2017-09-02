package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/itscaro/go-tools/upnp"
)

func main() {
	f, err := os.OpenFile("log.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Sprintf("error opening file: %v", err))
	}
	defer f.Close()

	log.SetOutput(f)

	u, err := upnp.NewUPNP("Go Webhook", []string{"192.168.0.0/16"})
	if err != nil {
		log.Panic(err)
	}
	u.LogEnabled = true

	err = u.AddPortMapping(8080, 9500, "TCP")
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("UPNP: added mapping wan: %v => %v %v\n", 8080, 9500, "TCP")

	ip, _ := u.ExternalIPAddress()
	if ip == nil {
		log.Panic("Missing external IP")
	}
	fmt.Printf("IP: %+v\n", ip)

	// Disable Console Color
	// gin.DisableConsoleColor()

	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()

	webhook := router.Group("/webhook")
	{
		webhook.GET("/", webhookFunc)
		webhook.POST("/", webhookFunc)
		webhook.PUT("/", webhookFunc)
		webhook.DELETE("/", webhookFunc)
		webhook.PATCH("/", webhookFunc)
		webhook.HEAD("/", webhookFunc)
		webhook.OPTIONS("/", webhookFunc)
	}

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	router.Run(":8080")
	// router.Run(":3000") for a hard coded port
}

func webhookFunc(c *gin.Context) {
	data, _ := c.GetRawData()

	fmt.Printf("Data:\n%+v\n", string(data))
}
