package parking

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/raynine/parking-go/parking/handlers"
)

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Init() {
	r := gin.Default()
	r.Use(cors.Default())

	handler := handlers.NewParkingHandler()
	go handler.ParkCar()

	parkingGroup := r.Group("/parking")
	parkingGroup.GET("/available", handler.GetAvailableParkingLosts)
	parkingGroup.POST("/park", handler.ParkInAvailableSpace)
	parkingGroup.DELETE("/leave", handler.LeaveParkingLot)

	panic(http.ListenAndServe(":8080", r))
}
