package repo

import (
	"context"
	"time"

	"sos/internal/model"
)

// RouteRepository определяет методы для работы с маршрутами.
type RouteRepository interface {
	GetRouteByID(ctx context.Context, routeID int) (*model.Route, error)
	GetAllRoutes(ctx context.Context, routeList []int) ([]model.Route, error)
}

// StopRepository определяет методы для работы с остановками.
type StopRepository interface {
	GetStopsByRoute(ctx context.Context, routeID int) ([]model.Stop, error)
}

// TrafficRepository определяет методы для работы с пассажиропотоком.
type TrafficRepository interface {
	GetTraffic(ctx context.Context, stopID int, startTime time.Time, endTime time.Time) ([]model.Traffic, error)
}
