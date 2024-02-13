package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/raynine/parking-go/models"
)

type ParkingHandler struct {
	ParkingLots chan models.Parking
	queuedCars  chan models.Parking
	leavingCars chan models.Parking
}

func NewParkingHandler() *ParkingHandler {
	return &ParkingHandler{
		ParkingLots: make(chan models.Parking, 2),
		queuedCars:  make(chan models.Parking, 4),
		leavingCars: make(chan models.Parking, 1),
	}
}

func (handler *ParkingHandler) GetAvailableParkingLosts(c *gin.Context) {
	availableParkings := handler.getAvailableParkingLosts()
	c.JSON(http.StatusOK, gin.H{
		"available_parkings": availableParkings,
	})
}

func (handler *ParkingHandler) ParkInAvailableSpace(c *gin.Context) {
	var dto models.ParkingDTO

	err := c.BindJSON(&dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
		})
		return
	}

	availableSpace := handler.getAvailableParkingLosts()

	parking := models.Parking{
		Id:  len(handler.ParkingLots) + 1,
		Car: dto.Car,
	}

	fmt.Println("Available space: ", availableSpace)
	fmt.Println("Cap: ", cap(handler.ParkingLots))

	if availableSpace > 0 {
		fmt.Println("There are available parkings, you are assigned to it.")
		handler.ParkingLots <- parking

	} else {
		fmt.Println("There are no available parkings, you will be queued.")
		handler.queuedCars <- parking
	}

	c.JSON(http.StatusAccepted, nil)
}

func (handler *ParkingHandler) LeaveParkingLot(c *gin.Context) {
	if handler.getAvailableParkingLosts() == cap(handler.ParkingLots) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Parking lot is empty, no car can leave",
		})
		return
	}

	car := <-handler.ParkingLots
	handler.leavingCars <- car
}

func (handler *ParkingHandler) getAvailableParkingLosts() int {
	availableParkings := cap(handler.ParkingLots) - len(handler.ParkingLots)
	return availableParkings
}

func (handler *ParkingHandler) ParkCar() {
	fmt.Println("Running go routine")
	for leaving := range handler.leavingCars {
		fmt.Printf("Car: %v is leaving, moving next car in the queue to the parking lot\n", leaving.Car)
		nextCar := <-handler.queuedCars
		handler.ParkingLots <- nextCar
		fmt.Printf("Car: %v has been parked\n", nextCar.Car)
	}
}
