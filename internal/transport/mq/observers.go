package mq

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rfashwal/scs-room/internal"
	"github.com/rfashwal/scs-room/internal/dto"
	"github.com/rfashwal/scs-utilities/config"
	"github.com/rfashwal/scs-utilities/rabbit"
	"github.com/rfashwal/scs-utilities/rabbit/domain"
)

func TemperatureObserver(s internal.Service, manager rabbit.MQManager, conf config.Manager) error {
	observer, err := manager.InitObserver()
	if err != nil {
		return err
	}
	defer manager.CloseConnection()
	defer observer.Channel.Close()

	err = observer.DeclareTopicExchange(conf.TemperatureTopic())
	if err != nil {
		return err
	}
	err = observer.BindQueue(observer.Queue, conf.ReadingsRoutingKey()+".#", conf.TemperatureTopic())

	if err != nil {
		return err
	}

	deliveries := observer.Observe()

	for msg := range deliveries {
		measurementDTO := domain.TemperatureMeasurement{}
		err := json.Unmarshal(msg.Body, &measurementDTO)
		if err != nil {
			fmt.Printf("could not unmarshal expected measurement msg, %s\n", err.Error())
			continue
		}
		fmt.Printf("message is delivered %v\n", measurementDTO)

		room, err := s.FindRoomByName(measurementDTO.RoomId)
		if err != nil {
			fmt.Printf("could not find room, %s\n", err.Error())
			continue
		}

		if room != nil {
			if float64(room.TempratureRequired) != measurementDTO.Value {
				newActuatorValue := (room.TempratureRequired * float64(room.AculatorValue)) / measurementDTO.Value

				room.AculatorValue = uint(newActuatorValue)
				room.TempratureCurrent = measurementDTO.Value
				s.UpdateRoom(*room)

				err = s.PublishActuatorsData(dto.AculatorMessage{
					ProcessId:   uuid.New().String(),
					PublishedOn: time.Now(),
					SensorId:    measurementDTO.SensorId,
					Value:       newActuatorValue,
					Service:     conf.ServiceName(),
					RoomId:      measurementDTO.RoomId,
				})
				if err != nil {
					fmt.Printf("could not publish message, %s\n", err.Error())
					continue
				}
			}
		}
	}
	return nil
}
