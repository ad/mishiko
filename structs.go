package main

type auth struct {
	Email           string `json:"email"`
	ID              int    `json:"id"`
	LastAppActivity string `json:"lastAppActivity"`
	LastKnownArea   string `json:"lastKnownArea"`
	Name            string `json:"name"`
}

type activity struct {
	PetActivityAim   int    `json:"petActivityAim"`
	PetAvatarPath    string `json:"petAvatarPath"`
	BatteryCharge    int    `json:"batteryCharge"`
	SubscriptionDate string `json:"subscriptionDate"`
	CurrentEnergy    int    `json:"currentEnergy"`
	CurrentActivity  int    `json:"currentActivity"`
	CurrentDistance  int    `json:"currentDistance"`
	UpdateDate       string `json:"updateDate"`
	DeviceStatus     string `json:"deviceStatus"`
	LightStatus      bool   `json:"lightStatus"`
	PetID            int    `json:"petId"`
	SignalLevel      string `json:"signalLevel"`
}

type pet struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Birthday    int     `json:"birthday"`
	Weight      float32 `json:"weight"`
	Height      float32 `json:"height"`
	ActivityAim int     `json:"activityAim"`
}

type locations struct {
	BoundDataID int            `json:"boundDataId"`
	SosModeTime int            `json:"sosModeTime"`
	Pets        []locationsPet `json:"pets"`
}

type locationsPet struct {
	ID            int     `json:"id"`
	DeviceID      int     `json:"deviceId"`
	Accuracy      float64 `json:"accuracy"`
	BatteryCharge int     `json:"batteryCharge"`
	Alt           float64 `json:"alt"`
	Lat           float64 `json:"lat"`
	Lon           float64 `json:"lon"`
	Date          int64   `json:"date"`
	DeviceStatus  string  `json:"deviceStatus"`
	SosModeTime   int     `json:"sosModeTime"`
}

type telegramResponse struct {
	body string
}
