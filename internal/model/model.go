package model

import "time"

// Route представляет маршрут общественного транспорта.
type Route struct {
	id              int     // Идентификатор маршрута
	name            string  // Название маршрута
	stops           []Stop  // Список остановок на маршруте
	workload        float64 // Загруженность маршрута
	vehicleCapacity int     // Вместимость ТС на маршруте
	vehicleNumber   int     // Количество ТС на маршруте
}

func NewRoute(id int, name string, stops []Stop, workload float64, vehicleCapacity int, vehicleNumber int) *Route {
	return &Route{
		id:              id,
		name:            name,
		stops:           stops,
		workload:        workload,
		vehicleCapacity: vehicleCapacity,
		vehicleNumber:   vehicleNumber,
	}
}

func (r *Route) GetID() int {
	return r.id
}

func (r *Route) GetName() string {
	return r.name
}

func (r *Route) CalculateAvailableBuses() int {
	totalLoad := 0.0
	for _, stop := range r.stops {
		totalLoad += stop.Workload
	}
	return int(totalLoad / float64(r.vehicleCapacity))
}

func (r *Route) GetStops() []Stop {
	return r.stops
}

func (r *Route) GetWorkload() float64 {
	return r.workload
}

func (r *Route) GetVehicleCapacity() int {
	return r.vehicleCapacity
}

func (r *Route) SetStops(stops []Stop) {
	r.stops = stops
}

func (r *Route) GetVehicleNumber() int {
	return r.vehicleNumber
}

func (r *Route) SetWorkload(workload float64) {
	r.workload = workload
}

type DonorRoute struct {
	RouteID int
	Stops   []Stop
	Buses   int
}

type OptimizationResult struct {
	RequiredBuses int          // Количество ТС, которое нужно добавить
	DonorRoutes   []DonorRoute // Маршруты, с которых можно забрать ТС
	IdealBuses    int          // Идеальное количество ТС для оптимальной загруженности
	ActualBuses   int          // Фактическое количество ТС, которое можно получить
	Errors        []string     // Ошибки оптимизации
}

//func (r *Route) CalculateWorkload(in, out int) {
//	r.Workload = r.Workload + in - out
//}

// Stop представляет остановку общественного транспорта.
type Stop struct {
	ID        int     // Идентификатор остановки.
	Name      string  // Название остановки.
	Location  Point   // Географические координаты остановки.
	Workload  float64 // Загруженность остановки.
	UpdatedAt time.Time
}

// Point представляет географическую точку с координатами.
type Point struct {
	Latitude  float64 // Широта.
	Longitude float64 // Долгота.
}

// Traffic представляет данные о пассажиропотоке на остановке.
type Traffic struct {
	StopID    int       // Идентификатор остановки
	Timestamp time.Time // Время регистрации
	In        int       // Количество вошедших пассажиров
	Out       int       // Количество вышедших пассажиров
}

// AnalysisRequest представляет запрос на анализ маршрута.
type AnalysisRequest struct {
	RouteID   int       // Идентификатор маршрута
	StartTime time.Time // Начало анализа
	EndTime   time.Time // Конец анализа
}

// OptimizationRecommendation представляет рекомендации по оптимизации.
type OptimizationRecommendation struct {
	RouteID          int      // Идентификатор маршрута
	SuggestedChanges []string // Рекомендации по изменениям
}

type Vehicle struct {
	ID       int
	RouteID  int
	Capacity int
}

type Movement struct {
	ID              int
	Items           []Item
	RouteID         int
	VehicleID       int
	MovingStopTime  time.Time
	MovingStartTime time.Time
}

type Item struct {
	ID              int
	StopName        string
	MovingStopTime  time.Time
	MovingStartTime time.Time
	In              int
	Out             int
}
