package domain

import "time"

type Machine struct {
	Id         	string    	`json:"id"`
	Code       	string    	`json:"code"`
	Workshop   	string    	`json:"workshop"`
	EqType     	string    	`json:"eq_type"`
	MaxTempC   	float64   	`json:"max_temp_c"`
	MaxVibro   	float64   	`json:"max_vibro"`
	
	CreatedAt  	time.Time 	`json:"created_at"`
}

type TelemetryEvent struct {
	EventId    	string    	`json:"event_id"`
	MachineId  	string    	`json:"machine_id"`
	TempC      	float64   	`json:"temp_c"`
	Vibro      	float64   	`json:"vibro"`
	
	Rpm        	int       	`json:"rpm"`
	PowerKw    	float64   	`json:"power_kw"`
	IsOverload 	bool      	`json:"is_overload"`
	CreatedAt  	time.Time 	`json:"created_at"`
}


type ShiftReport struct {
	Id            	string    	`json:"id"`
	MachineId     	string    	`json:"machine_id"`
	Status        	string    	`json:"status"`
	AvgTempC      	float64   	`json:"avg_temp_c"`
	MaxTempC      	float64   	`json:"max_temp_c"`
	AvgVibro      	float64   	`json:"avg_vibro"`
	MaxVibro      	float64   	`json:"max_vibro"`
	OverloadCount 	int       	`json:"overload_count"`
	RiskLevel     	string    	`json:"risk_level"`
	
	CreatedAt     	time.Time 	`json:"created_at"`
}

type Incident struct {
	Id         	string     	`json:"id"`
	MachineId  	string     	`json:"machine_id"`
	EventId    	string     	`json:"event_id"`
	TempC      	float64    	`json:"temp_c"`
	Vibro      	float64    	`json:"vibro"`
	Reason     	string     	`json:"reason"`
	IsResolved 	bool       	`json:"is_resolved"`
	
	CreatedAt  	time.Time  	`json:"created_at"`
	ResolvedAt 	*time.Time 	`json:"resolved_at"`
}
