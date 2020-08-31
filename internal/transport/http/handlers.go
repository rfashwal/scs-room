package http

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rfashwal/scs-room/internal"
	"github.com/rfashwal/scs-room/internal/dto"
)

type SensorReadingRequest struct {
	SensorID string
	Type     string
	Value    interface{}
}

func NewRouter(s internal.Service) *gin.Engine {
	r := gin.Default()

	r.POST("/rooms", addRoomHandler(s))
	r.PATCH("/rooms/:name", updateHandler(s))
	r.GET("/rooms", listHandler(s))
	r.GET("/rooms/:name", getHandler(s))

	return r
}

func publishSensorDataHandler(s internal.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
		var sensorReadingRequest dto.AculatorMessage

		if err := c.Bind(&sensorReadingRequest); err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		err := s.PublishActuatorsData(sensorReadingRequest)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprint(err)})
			return
		}

		c.JSON(http.StatusOK, "sensor data published")
	}
}

func addRoomHandler(s internal.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
		roomDto := &dto.RoomDTO{}
		if err := c.Bind(roomDto); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprint(err)})
			return
		}
		created, err := s.SaveRoom(*roomDto)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprint(err)})
			return
		}

		c.JSON(http.StatusCreated, created)
	}
}

func updateHandler(s internal.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))

		roomDto := &dto.RoomDTO{}
		if err := c.Bind(roomDto); err != nil {
			if strings.Contains(err.Error(), "415") {
				c.JSON(http.StatusUnsupportedMediaType, gin.H{"error": fmt.Sprint(err)})
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprint(err)})
			return
		}

		updated, err := s.UpdateRoom(*roomDto)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprint(err)})
			return
		}

		c.JSON(http.StatusAccepted, updated)
	}
}

func listHandler(s internal.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
		dtos, err := s.GetAllRooms()
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		if len(dtos) == 0 {
			c.JSON(http.StatusNoContent, dtos)
			return
		}
		c.JSON(http.StatusOK, dtos)
	}
}

func getHandler(s internal.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
		room := c.Param("room")

		sensor, err := s.FindRoomByName(room)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		if sensor == nil {
			c.JSON(http.StatusNoContent, errors.New("unknown room given"))
			return
		}

		c.JSON(http.StatusOK, sensor)
	}
}
