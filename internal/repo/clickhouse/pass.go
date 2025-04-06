package clickhouse

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"time"

	"sos/internal/model"
)

func (r *Repository) GetRouteByID(
	ctx context.Context,
	routeID int) (route *model.Route, err error) {
	err = r.conn.Select(ctx, &route, fmt.Sprintf(`SELECT * FROM %[1]s WHERE RouteID = %[2]s`, r.dbName, routeID))
	if err != nil {
		return nil, errors.Wrap(err, "failed to conn.Select")
	}

	return route, nil
}

func (r *Repository) GetAllRoutes(
	ctx context.Context,
	routeList []int,
) (routes []model.Route, err error) {
	var args []interface{}
	where, whereArgs := Delimiter("WHERE", "AND", []Statement{
		{Statement: "RouteID IN (@RouteID)", Arg: routeList},
	})
	args = append(args, whereArgs...)

	err = r.conn.Select(ctx, routes, fmt.Sprintf(`
			SELECT * FROM %[1]s 
		         %[2]s`, r.dbName, where),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to conn.Select")
	}

	return routes, nil
}

func (r *Repository) GetStopsByRoute(
	ctx context.Context,
	routeID int) ([]model.Stop, error) {
	return nil, nil
}

func (r *Repository) GetTraffic(
	ctx context.Context,
	stopID int, startTime time.Time, endTime time.Time) ([]model.Traffic, error) {
	return nil, nil
}

func (r *Repository) GetWorkloads(ctx context.Context, targetRouteID int) ([]*model.Route, error) {
	// Формируем SQL-запрос
	query := `
        SELECT 
            r.id,
            r.name,
            s.id AS stop_id,
            s.name AS stop_name,
            s.location,
            s.workload,
            s.updated_at
        FROM routes r
        JOIN stops s ON arrayContains(r.stops, s.name)
        WHERE r.id != ? AND s.workload IS NOT NULL
    `

	// Подготавливаем параметры
	params := []interface{}{targetRouteID}

	// Выполняем запрос
	rows, err := r.conn.Query(ctx, query, params...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	// Парсим результаты
	routes := make(map[int]*model.Route)
	for rows.Next() {
		var (
			routeID      int
			routeName    string
			stopID       int
			stopName     string
			stopLocation model.Point
			stopWorkload float64
			updatedAt    time.Time
		)

		err = rows.Scan(&routeID, &routeName, &stopID, &stopName, &stopLocation, &stopWorkload, &updatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Формируем структуру Stop
		stop := model.Stop{
			ID:        stopID,
			Name:      stopName,
			Location:  stopLocation,
			Workload:  stopWorkload,
			UpdatedAt: updatedAt,
		}

		// Добавляем в соответствующий маршрут
		if _, ok := routes[routeID]; !ok {
			newRoute := model.NewRoute(routeID, routeName, []model.Stop{stop}, 0, 70, 10)
			routes[routeID] = newRoute
		} else {
			//s := append(routes[routeID].Stops, stop)
			//routes[routeID].Stops = s
		}
	}

	// Преобразуем карту в слайс и рассчитываем Workload
	result := make([]*model.Route, 0, len(routes))
	for _, r := range routes {
		// Рассчитываем среднюю загруженность
		var totalWorkload float64
		for _, stop := range r.GetStops() {
			totalWorkload += stop.Workload
		}
		r.SetWorkload(totalWorkload / float64(len(r.GetStops())))
		result = append(result, r)
	}

	return result, nil
}

func (r *Repository) GetStops(ctx context.Context, routeIDs []string, timeFrom, timeTo time.Time) ([]model.Stop, error) {
	where, whereArgs := Delimiter("WHERE", "AND", []Statement{
		{Statement: "RouteID in (@RouteID)", Arg: routeIDs},
		{Statement: "DeviceTime <= @DateTo", Arg: timeTo},
		{Statement: "DeviceTime >= @DateFrom", Arg: timeFrom},
	})

	query := fmt.Sprintf(`
		SELECT 
			stop_id,
			AVG(passengers_in - passengers_out) AS avg_load,
			MAX(stop_time) AS last_update
		FROM telematics
		%[1]s
		GROUP BY stop_id
	`, where)

	rows, err := r.conn.Query(ctx, query, whereArgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var loads []model.Stop
	for rows.Next() {
		var stop model.Stop
		var workload float64
		var lastUpdate time.Time
		err := rows.Scan(&stop.ID, &workload, &lastUpdate)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		stop.Workload = workload
		stop.UpdatedAt = lastUpdate
		loads = append(loads, stop)
	}

	return loads, nil
}
