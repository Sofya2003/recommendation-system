package service

import (
	"github.com/magiconair/properties/assert"
	"testing"

	"sos/internal/model"
)

func Test_calculateRequiredBuses(t *testing.T) {
	tests := []struct {
		name  string
		route *model.Route
		want  int
	}{
		{
			name: "test1",
			route: model.NewRoute(
				1, "test", []model.Stop{
					{ID: 1, Name: "stop1", Workload: 37.0},
					{ID: 1, Name: "stop2", Workload: 55.0},
					{ID: 1, Name: "stop3", Workload: 87.0},
					{ID: 1, Name: "stop4", Workload: 84.0},
					{ID: 1, Name: "stop5", Workload: 48.0},
				},
				83.0,
				70,
				10,
			),
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculateRequiredBuses(tt.route); got != tt.want {
				t.Errorf("calculateRequiredBuses() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFindDonorRoutes(t *testing.T) {
	tests := []struct {
		name     string
		routes   map[int]*model.Route
		expected []model.DonorRoute
	}{
		{
			name: "Route with mixed workloads",
			routes: map[int]*model.Route{
				1: model.NewRoute(1, "Route 1", []model.Stop{
					{ID: 1, Name: "Stop 1", Workload: 50.0},
					{ID: 2, Name: "Stop 2", Workload: 60.0},
				}, 40.0, 70, 5), // Вместимость одного автобуса 70, количество ТС 5
			},
			expected: []model.DonorRoute{
				{RouteID: 1, Buses: 0}, // 70 - (50% + 60%) = 70 - 42 = 0 доступных автобусов
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := findDonorRoutes(tt.routes)
			if len(got) != len(tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, got)
				return
			}
			for i, donorRoute := range got {
				assert.Equal(t, tt.expected[i].RouteID, donorRoute.RouteID)
			}
		})
	}
}
