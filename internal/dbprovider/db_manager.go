package dbprovider

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/rfashwal/scs-room/internal/config"
	"github.com/rfashwal/scs-room/internal/dto"
	"github.com/rfashwal/scs-room/internal/model"
)

var Mgr DBManager

type manager struct {
	db *gorm.DB
}
type DBManager interface {
	GetDB() *gorm.DB

	MapToRoomEntity(roomDto dto.RoomDTO) *model.Room
	MapToRoomDTO(room *model.Room) dto.RoomDTO

	SaveRoom(roomDto dto.RoomDTO) (*model.Room, error)
	GetAllRooms() ([]model.Room, error)
	GetRoomByName(name string) (*model.Room, error)
	UpdateRoom(roomDto dto.RoomDTO) (*model.Room, error)
}

func NewDBManager() (DBManager, error) {
	dbPath := config.Config().DBPath()
	if path, exists := adjustDBPath(dbPath); !exists {
		dbPath = path
	}

	db2, err := gorm.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to init db[%s]", dbPath)
	}

	db2.SingularTable(true)
	db2.AutoMigrate(&model.Room{})

	return &manager{db: db2}, nil
}

func adjustDBPath(dbPath string) (string, bool) {
	var exists = true

	if _, err := os.Stat(dbPath); err != nil {
		exists = false
	}

	if !exists {
		dbPath = dbPath + "?mode=rwc"
	}

	return dbPath, exists
}

func (m *manager) SaveRoom(roomDto dto.RoomDTO) (*model.Room, error) {
	exists, err := m.GetRoomByName(roomDto.Name)
	if err != nil && err.Error() != "record not found" {
		return nil, err
	}

	if exists != nil {
		return nil, errors.New("room already exists")
	}

	room := m.MapToRoomEntity(roomDto)

	err = m.db.Create(room).Error

	if err != nil {
		return nil, err
	}

	return room, err
}

func (m *manager) GetAllRooms() ([]model.Room, error) {
	var rooms []model.Room
	err := m.db.Find(&rooms).Error

	return rooms, err
}

func (m *manager) GetRoomByName(name string) (*model.Room, error) {
	room := &model.Room{}
	err := m.db.Where("name=?", name).First(room).Error
	if err != nil {
		return nil, err
	}
	return room, nil
}

func (m *manager) UpdateRoom(roomDTO dto.RoomDTO) (*model.Room, error) {

	room := &model.Room{}
	err := m.db.Where("id=?", roomDTO.ID).First(room).Error
	if err != nil {
		return nil, err
	}

	*room.Name = roomDTO.Name
	*room.AculatorValue = roomDTO.AculatorValue
	*room.TempratureRequired = roomDTO.TempratureRequired
	err = m.db.Save(room).Error
	if err != nil {
		return nil, err
	}
	return room, nil
}

func (m *manager) GetDB() *gorm.DB {
	return m.db
}

func (mgr *manager) MapToRoomEntity(roomDto dto.RoomDTO) *model.Room {

	now := time.Now()

	room := &model.Room{
		Name:               &roomDto.Name,
		AculatorValue:      &roomDto.AculatorValue,
		TempratureRequired: &roomDto.TempratureRequired,
		TempratureCurrent:  &roomDto.TempratureCurrent,
		CreatedAt:          &now,
		ModifiedAt:         &now,
	}

	return room
}

func (mgr *manager) MapToRoomDTO(room *model.Room) dto.RoomDTO {

	if room == nil {
		return dto.RoomDTO{}
	}

	dto := dto.RoomDTO{
		ID:                 *room.ID,
		AculatorValue:      *room.AculatorValue,
		Name:               *room.Name,
		TempratureRequired: *room.TempratureRequired,
		TempratureCurrent:  *room.TempratureCurrent,
	}

	return dto
}
