package service

import (
	"context"
	"log"
	"math"
	"sync"

	"sos/internal/model"
	clickhouseRepo "sos/internal/repo/clickhouse"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/pkg/errors"
)

const (
	targetLoad  = 60.0 // Целевая загруженность (в %)
	maxLoad     = 80.0 // Максимальная допустимая загруженность
	busCapacity = 100  // Вместимость одного автобуса
)

type OptimizationService struct {
	Repository *clickhouseRepo.Repository
}

func NewOptimizationService(log *log.Logger, conn clickhouse.Conn, dbName string) *OptimizationService {
	return &OptimizationService{
		Repository: clickhouseRepo.NewRepository(conn, log, dbName),
	}
}

func (s *OptimizationService) OptimizeRoute(ctx context.Context, targetRouteID int) (model.OptimizationResult, error) {
	// 1. Получаем данные о загруженности всех остановок
	stopsWorkload, err := s.Repository.GetWorkloads(ctx, targetRouteID)
	if err != nil {
		return model.OptimizationResult{}, err
	}

	// 2. Формируем данные по маршрутам
	routes := s.buildRouteData(stopsWorkload)

	// 3. Определяем целевой маршрут
	targetRoute := routes[targetRouteID]
	if targetRoute.GetID() == 0 {
		return model.OptimizationResult{}, errors.New("target route not found")
	}

	// 4. Определяем необходимое количество ТС
	requiredBuses := calculateRequiredBuses(targetRoute)

	// 5. Находим "доноров"
	donors := findDonorRoutes(routes)

	return model.OptimizationResult{
		RequiredBuses: requiredBuses,
		DonorRoutes:   donors,
	}, nil
}

func calculateRequiredBuses(route *model.Route) int {
	totalRequiredPassengers := 0

	for _, stop := range route.GetStops() {
		if stop.Workload > targetLoad {
			// Рассчитываем количество пассажиров, которых нужно перевести
			currentPassengers := stop.Workload / 100 * float64(route.GetVehicleCapacity())
			targetPassengers := targetLoad / 100 * float64(route.GetVehicleCapacity())
			requiredPassengers := int(currentPassengers - targetPassengers)

			if requiredPassengers > 0 {
				totalRequiredPassengers += requiredPassengers
			}
		}
	}

	// Рассчитываем необходимое количество автобусов
	requiredBuses := int(math.Ceil(float64(totalRequiredPassengers) / float64(busCapacity)))

	// Ограничиваем количество автобусов, которые можно перенаправить
	maxBusesToRedirect := 3 // Максимальное количество автобусов, которое можно перенаправить
	if requiredBuses > maxBusesToRedirect {
		return maxBusesToRedirect
	}

	return requiredBuses
}

func findDonorRoutes(routes map[int]*model.Route) []model.DonorRoute {
	var donorRoutes []model.DonorRoute

	for _, route := range routes {
		// Получаем процент загруженности
		workloadPercentage := route.GetWorkload() // Предполагается, что это значение от 0 до 100
		vehicleCapacity := route.GetVehicleCapacity()
		vehicleNumber := route.GetVehicleNumber()

		// Проверяем, меньше ли загруженность 60%
		if workloadPercentage < 60 {
			// Вычисляем фактическую загруженность
			actualWorkload := int(float64(vehicleCapacity) * (workloadPercentage / 100))
			availableBuses := vehicleNumber - actualWorkload // Расчет доступных автобусов

			if availableBuses > 0 {
				donorRoutes = append(donorRoutes, model.DonorRoute{
					RouteID: route.GetID(),
					Buses:   availableBuses,
				})
			}
		}
	}

	return donorRoutes
}

//func findDonorRoutes(routes map[int]*model.Route, targetRoute *model.Route, requiredBuses int) []model.DonorRoute {
//	var donors []model.DonorRoute
//
//	for id, route := range routes {
//		// Пропускаем целевой маршрут
//		if id == targetRoute.GetID() {
//			continue
//		}
//
//		// Проверяем наличие остановок
//		stops := route.GetStops()
//		if len(stops) == 0 {
//			continue // Пропускаем маршруты без остановок
//		}
//
//		// Ищем остановки с низкой загруженностью
//		suitableStops := stops[:0] // Инициализируем срез для подходящих остановок
//		for _, stop := range stops {
//			if stop.Workload < route.GetWorkload()*0.6 {
//				suitableStops = append(suitableStops, stop)
//			}
//		}
//
//		// Рассчитываем, сколько ТС можно забрать
//		availableBuses := route.CalculateAvailableBuses()
//		if availableBuses > 0 {
//			donors = append(donors, model.DonorRoute{
//				RouteID: id,
//				Stops:   suitableStops,
//				Buses:   availableBuses,
//			})
//		}
//	}
//
//	// Сортируем доноров по количеству доступных ТС (и по ID для стабильности)
//	sort.Slice(donors, func(i, j int) bool {
//		if donors[i].Buses == donors[j].Buses {
//			return donors[i].RouteID < donors[j].RouteID // Сортировка по ID в случае равенства
//		}
//		return donors[i].Buses > donors[j].Buses
//	})
//
//	// Возвращаем только необходимое количество ТС
//	totalBuses := 0
//	result := []model.DonorRoute{}
//	for _, d := range donors {
//		if totalBuses >= requiredBuses {
//			break
//		}
//		busesToTake := min(d.Buses, requiredBuses-totalBuses)
//		result = append(result, model.DonorRoute{
//			RouteID: d.RouteID,
//			Stops:   d.Stops,
//			Buses:   busesToTake,
//		})
//		totalBuses += busesToTake
//	}
//
//	return result
//}

//func calculateAvailableBuses(stops []model.Stop) int {
//	// Пример расчета: 1 ТС на каждые 100 пассажиров
//	totalLoad := 0.0
//	for _, stop := range stops {
//		totalLoad += stop.Workload
//	}
//	return int(totalLoad / 100)
//}

// Параллельное группирование данных по маршрутам
func (s *OptimizationService) buildRouteData(routes []*model.Route) map[int]*model.Route {
	routeData := make(map[int]*model.Route)
	var wg sync.WaitGroup

	for _, sl := range routes {
		wg.Add(1)
		go func(sl *model.Route) {
			defer wg.Done()
			// Обновление данных для маршрута
			rd := routeData[sl.GetID()]
			stops := append(rd.GetStops(), sl.GetStops()...)
			rd.SetStops(stops)
			routeData[sl.GetID()] = rd
		}(sl)
	}
	wg.Wait()
	return routeData
}
