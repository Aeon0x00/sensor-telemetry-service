package domain

type KafkaTelemetryDTO struct {
	EventId   	string  	`json:"event_id"`
	MachineId 	string  	`json:"machine_id"`
	TempC     	float64 	`json:"temp_c"`
	Vibro     	float64 	`json:"vibro"`
	Rpm       	int     	`json:"rpm"`
	PowerKw   	float64 	`json:"power_kw"`
	
}

type KafkaAlertDTO struct {
	IncidentId 	string  	`json:"incident_id"`
	MachineId  	string  	`json:"machine_id"`
	Code       	string  	`json:"code"`
	Workshop   	string  	`json:"workshop"`
	TempC      	float64 	`json:"temp_c"`
	Vibro      	float64 	`json:"vibro"`
	Reason     	string  	`json:"reason"`
}


type MachineStateDTO struct {
	MachineId  	string  	`json:"machine_id"`
	Code       	string  	`json:"code"`
	Workshop   	string  	`json:"workshop"`
	LastTempC  	float64 	`json:"last_temp_c"`
	LastVibro  	float64 	`json:"last_vibro"`
	LastRpm    	int     	`json:"last_rpm"`
	LastPower  	float64 	`json:"last_power"`
	IsOverload 	bool    	`json:"is_overload"`
	
}
