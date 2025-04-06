package clickhouse

import (
	"context"
	"log"
	"sos/internal/model"

	"github.com/ClickHouse/clickhouse-go/v2"
	"golang.org/x/crypto/bcrypt"
)

type Repository struct {
	conn   clickhouse.Conn
	log    *log.Logger
	dbName string
}

func NewRepository(conn clickhouse.Conn, log *log.Logger, dbName string) *Repository {
	return &Repository{
		conn:   conn,
		log:    log,
		dbName: dbName,
	}
}

func (r *Repository) Close(_ context.Context) error {
	return r.conn.Close()
}

func (r *Repository) GetTransportsByRoute(routeID int) (vehicles []model.Vehicle, err error) {
	err = r.conn.Select(context.Background(), &vehicles, "SELECT id, route_id, capacity FROM transport WHERE route_id = ?", routeID)
	if err != nil {
		return nil, err
	}

	return vehicles, nil
}

func (r *Repository) GetRouteLoad(routeID int) (float64, error) {
	var load float64
	err := r.conn.QueryRow(context.Background(), "SELECT load FROM routes WHERE id = ?", routeID).Scan(&load)
	return load, err
}

//func (r *Repository) RecordMovement(movement model.Movement) error {
//	query := `
//        INSERT INTO movement (transport_id, stop_time, passengers_in, passengers_out)
//        VALUES (?, ?, ?, ?)`
//	err := r.conn.Exec(context.Background(), query, movement.TransportID, movement.StopTime, movement.PassengersIn, movement.PassengersOut)
//	return err
//}

type User struct {
	ID           int
	Username     string
	PasswordHash string
	Role         int
}

var users []User

func GetUsers(name, login string) ([]User, error) {
	query := "SELECT username, role FROM users WHERE 1=1"
	args := []interface{}{}

	if name != "" {
		query += " AND name LIKE ?"
		args = append(args, "%"+name+"%")
	}

	if login != "" {
		query += " AND username LIKE ?"
		args = append(args, "%"+login+"%")
	}

	rows := users
	return rows, nil
}

func CreateUser(username, password, role string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte("admin"), 14)
	if err != nil {
		log.Fatal(err)
	}

	var roleInt int
	switch role {
	case "admin":
		roleInt = 0
	case "moderator":
		roleInt = 2
	case "user":
		roleInt = 3
	}

	users = append(users, User{ID: len(users), Username: username, PasswordHash: string(hash), Role: roleInt})
	return err
}
