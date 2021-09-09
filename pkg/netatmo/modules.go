package netatmo

type DeviceListResponse struct {
	Body   DeviceListResponseBody `json:"body"`
	Status string                 `json:"status"`
}

type DeviceListResponseBody struct {
	Devices []Device `json:"devices"`
	// User
}

type Device struct {
	Id       string                 `json:"_id"`
	Name     string                 `json:"station_name"`
	DataType []string               `json:"data_type"`
	Data     map[string]interface{} `json:"dashboard_data"`
	Modules  []Module               `json:"modules"`
}

type Module struct {
	Id       string                 `json:"_id"`
	Name     string                 `json:"module_name"`
	DataType []string               `json:"data_type"`
	Data     map[string]interface{} `json:"dashboard_data"`
}

type DashboardData struct {
	Temperature      float64 `json:"Temperature"`
	CO2              int     `json:"CO2"`
	Humidity         int     `json:"Humidity"`
	Noise            int     `json:"Noise"`
	Pressure         float64 `json:"Pressure"`
	AbsolutePressure float64 `json:"AbsolutePressure"`
	Min_temp         float64 `json:"min_temp"`
	Max_temp         float64 `json:"max_temp"`
	Date_max_temp    float64 `json:"date_max_temp"`
	Date_min_temp    float64 `json:"date_min_temp"`
	Temp_trend       string  `json:"temp_trend"`
	Pressure_trend   string  `json:"pressure_trend"`
}
