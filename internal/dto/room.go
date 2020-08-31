package dto

type RoomDTO struct {
	ID                 uint    `json:"id"`
	Name               string  `json:"name"`
	AculatorValue      uint    `json:"aculator_value"`
	TempratureCurrent  float64 `json:"temprature_current"`
	TempratureRequired float64 `json:"temprature_required"`
}
