package tests

import (
	"log"

	"github.com/gin-gonic/gin"
)

// StaticServer : the static file server used for testing the scrapper
type StaticServer struct {
	Router *gin.Engine
}

// Initialize : configures the static server for testing
func (s *StaticServer) Initialize() {
	gin.SetMode(gin.ReleaseMode)
	s.Router = gin.Default()
	s.initializeRoutes()
}

// Run : runs the static server for testing
func (s *StaticServer) Run(appPort string) {
	go func() {
		if err := s.Router.Run(":" + appPort); err != nil {
			log.Fatalf("Error while serving tests: %v", err)
		}
	}()
}

func (s *StaticServer) initializeRoutes() {
	s.Router.Static("/telemetry-data", "./assets/telemetry-data")
	s.Router.Static("/inverter", "./assets/inverter")
}
