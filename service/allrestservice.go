package service

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var _logger = logrus.New()

//APIResponse returns the service response
type APIResponse struct {
	Message   string      `json:"serviceMessage"`
	Payload   interface{} `json:"payload,omitempty"`
	ServiceTS string      `json:"ts"`
	IsSuccess bool        `json:"isSuccess"`
	Token     *string     `json:"token,omitempty"`
}

//TbsCommonRestService implements the rest API for HRM transactions
type AllRestService struct {
	colorService   *ColorService
}

//NewSobarRestService retuens a new initialized version of the service
func NewJOTRestService(config []byte, verbose bool) *AllRestService {
	service := new(AllRestService)
	if err := service.Init(config, verbose); err != nil {
		_logger.Errorf("Unable to intialize service instance %v", err)
		return nil
	}
	return service
}

//Init initializes the service
func (allRestService *AllRestService) Init(config []byte, verbose bool) error {
	if verbose {
		_logger.SetLevel(logrus.DebugLevel)
	}

	_logger.Info("Service instance intialized. Waiting to lauch the service")

	return nil
}

//Serve runs the infinite service method.
func (allRestService *AllRestService) Serve(address string, port int, stopSignal chan bool) {
	//gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	//TODO: Following to be changed for production
	router.MaxMultipartMemory = 8 << 21 //16 MB Max file size
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AddAllowMethods("GET", "POST", "DELETE", "PUT", "OPTIONS", "HEAD")
	corsConfig.AddAllowHeaders("Authorization", "WWW-Authenticate", "Content-Type", "Accept", "X-Requested-With")
	corsConfig.AddExposeHeaders("Authorization", "WWW-Authenticate", "Content-Type", "Accept", "X-Requested-With")
	router.Use(cors.New(corsConfig))
	// JWT Auth handler
	
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, buildResponse(true, "Service Available", nil))
	})

	allRestService.colorService.AddRouters(router)

	//ReportGeneration service
	httpServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", address, port),
		Handler: router,
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		// service connections
		_logger.Infof("Started listening@%v", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			_logger.Errorf("Unable to listen: %v\n", err)
			return
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-stopSignal
		_logger.Infof("Received stop signal")
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		_logger.Infof("Sending shutdown to http server")
		httpServer.Shutdown(ctx)
	}()
	time.Sleep(1 * time.Second)
	//Wait indefinitely
	wg.Wait()

}

func buildResponse(isOk bool, msg string, payload interface{}) APIResponse {
	return APIResponse{
		IsSuccess: isOk,
		Message:   msg,
		Payload:   payload,
		ServiceTS: time.Now().Format("2006-01-02-15:04:05.000"),
	}
}
