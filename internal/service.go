package internal

import (
	"encoding/json"

	"github.com/rfashwal/scs-room/internal/dbprovider"
	"github.com/rfashwal/scs-room/internal/dto"
	"github.com/rfashwal/scs-utilities/config"
	"github.com/rfashwal/scs-utilities/rabbit/publishing"
)

type Service interface {
	PublishActuatorsData(req dto.AculatorMessage) error

	SaveRoom(registerDTO dto.RoomDTO) (*dto.RoomDTO, error)
	UpdateRoom(registerDTO dto.RoomDTO) (*dto.RoomDTO, error)
	GetAllRooms() ([]dto.RoomDTO, error)
	FindRoomByName(room string) (*dto.RoomDTO, error)
}

func NewService(publisher *publishing.Publisher, conf config.Manager, m dbprovider.DBManager) (Service, error) {
	return service{
		publisher:       publisher,
		conf:            conf,
		databaseManager: m}, nil
}

type service struct {
	publisher       *publishing.Publisher
	conf            config.Manager
	databaseManager dbprovider.DBManager
}

func (s service) PublishActuatorsData(msg dto.AculatorMessage) error {

	encodedMsg, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return s.publisher.Publish(s.conf.ActuatorTopic()+"."+msg.RoomId, s.conf.ReadingsRoutingKey(), string(encodedMsg))
}

func (s service) SaveRoom(roomDto dto.RoomDTO) (*dto.RoomDTO, error) {
	register, err := s.databaseManager.SaveRoom(roomDto)
	dto := s.databaseManager.MapToRoomDTO(register)
	return &dto, err
}
func (s service) UpdateRoom(roomDto dto.RoomDTO) (*dto.RoomDTO, error) {
	updated, err := s.databaseManager.UpdateRoom(roomDto)
	dto := s.databaseManager.MapToRoomDTO(updated)
	return &dto, err
}
func (s service) GetAllRooms() ([]dto.RoomDTO, error) {
	rooms, err := s.databaseManager.GetAllRooms()
	if err != nil {
		return nil, err
	}

	var dtos []dto.RoomDTO
	for _, item := range rooms {
		dtos = append(dtos, s.databaseManager.MapToRoomDTO(&item))
	}
	return dtos, nil
}
func (s service) FindRoomByName(name string) (*dto.RoomDTO, error) {
	room, err := s.databaseManager.GetRoomByName(name)
	if err != nil {
		return nil, err
	}

	dto := s.databaseManager.MapToRoomDTO(room)
	return &dto, nil
}
