package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"plugin"

	"strings"

	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/itscaro/go-tools/upnp"
	"github.com/joho/godotenv"
)

var GitCommit string
var BuildDatetime string

var loadedHooks = make(map[string]Hook)

type Hook interface {
	Exec(data []byte) interface{}
}

type JsonResponse struct {
	Message string
	Code    int
}

func init() {
	fmt.Printf("Version %s - Built at %s\n", GitCommit, BuildDatetime)

	err := godotenv.Load()
	if err == nil {
		fmt.Println(".env file loadded")

		if len(os.Getenv(gin.ENV_GIN_MODE)) > 0 {
			gin.SetMode(os.Getenv(gin.ENV_GIN_MODE))
		}
	}

	if runtime.GOOS == "linux" {
		_ = os.Mkdir("./hook", os.FileMode(0755))
	}
}

func main() {
	fmt.Println("Starting...")
	f, err := os.OpenFile("log.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Sprintf("error opening file: %v", err))
	}
	defer f.Close()

	log.SetOutput(f)

	if os.Getenv("UPNP_ENABLED") == "true" {
		log.Println("Start UPNP service")

		var u *upnp.UPNP
		var err error

		if len(os.Getenv("UPNP_LOCAL_IP_RANGE")) > 0 {
			u, err = upnp.NewUPNP("Go Webhook", []string{os.Getenv("UPNP_LOCAL_IP_RANGE")})
		} else {
			u, err = upnp.NewUPNP("Go Webhook", nil)
		}

		if err != nil {
			log.Println(err)
		} else {
			u.LogEnabled = true

			err = u.AddPortMapping(8080, 9500, "TCP")
			if err != nil {
				log.Println(err)
			}
			fmt.Printf("UPNP: added mapping wan: %v => %v %v\n", 8080, 9500, "TCP")

			ip, _ := u.ExternalIPAddress()
			if ip == nil {
				log.Println("Missing external IP")
			}
			fmt.Printf("IP: %+v\n", ip)
		}
	}

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

		webhook.GET("/:name", webhookByPluginFunc)
		webhook.POST("/:name", webhookByPluginFunc)
		webhook.PUT("/:name", webhookByPluginFunc)
		webhook.DELETE("/:name", webhookByPluginFunc)
		webhook.PATCH("/:name", webhookByPluginFunc)
		webhook.HEAD("/:name", webhookByPluginFunc)
		webhook.OPTIONS("/:name", webhookByPluginFunc)
	}

	admin := router.Group("/admin")
	{
		admin.GET("/hooks", getHooksFunc)
		admin.DELETE("/hooks", clearHooksFunc)
	}

	fmt.Println("Ready for incoming requests")
	router.Run()
}

// getHooksFunc get all loaded hook plugins
func getHooksFunc(c *gin.Context) {
	loadedHooksList := []string{}

	for hookName := range loadedHooks {
		loadedHooksList = append(loadedHooksList, hookName)
	}

	c.JSON(http.StatusOK, loadedHooksList)
}

// clearHookFunc clears all loaded hook plugins
func clearHooksFunc(c *gin.Context) {
	loadedHooks = make(map[string]Hook)

	c.JSON(http.StatusOK, JsonResponse{
		Message: "Done",
	})
}

// webhookFunc is a static hook
func webhookFunc(c *gin.Context) {
	// Always send back response
	defer func() {
		if r := recover(); r != nil {
			c.JSON(http.StatusInternalServerError, JsonResponse{
				Message: fmt.Sprintf("Internal Error: %+v", r),
			})
		}
	}()

	if rawData, err := c.GetRawData(); err == nil {
		c.JSON(http.StatusOK, rawData)
	} else {
		c.JSON(http.StatusInternalServerError, JsonResponse{
			Message: err.Error(),
		})
	}
}

// webhookByPluginFunc uses dynamic plugin system
func webhookByPluginFunc(c *gin.Context) {
	// Always send back response
	defer func() {
		if r := recover(); r != nil {
			c.JSON(http.StatusInternalServerError, JsonResponse{
				Message: fmt.Sprintf("Internal Error: %+v", r),
			})
		}
	}()

	hookName := c.Param("name")
	log.Printf("Hook name: %s\n", hookName)
	if rawData, err := c.GetRawData(); err == nil {
		if h, err := getHook(hookName); err == nil {
			if h == nil {
				c.JSON(http.StatusNotFound, JsonResponse{
					Message: fmt.Sprintf("The hook '%s' does not exist", hookName),
				})
			} else {
				log.Printf("Hook to use: %v", &h)
				c.JSON(http.StatusOK, h.Exec(rawData))
			}
		} else {
			c.JSON(http.StatusInternalServerError, JsonResponse{
				Message: err.Error(),
			})
		}
	} else {
		c.JSON(http.StatusInternalServerError, JsonResponse{
			Message: err.Error(),
		})
	}
}

func getHook(name string) (h Hook, err error) {
	if _, ok := loadedHooks[name]; ok {
		return loadedHooks[name], nil
	}

	if p, err := plugin.Open("hook/" + name + ".so"); err == nil {
		symbol, err := p.Lookup(strings.Title(name))
		if err != nil {
			log.Println(err)
		}

		if h, ok := symbol.(Hook); ok {
			loadedHooks[name] = h
			return loadedHooks[name], nil
		} else {
			log.Println("Hook declaration does not match Hook interface")
			return nil, fmt.Errorf("Error while executing the hook '%s'", name)
		}
	} else {
		log.Println(err)
	}

	return
}
