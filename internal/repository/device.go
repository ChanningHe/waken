package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/channinghe/waken/internal/model"
)

type DeviceRepository struct {
	db *sql.DB
}

func NewDeviceRepository(db *sql.DB) *DeviceRepository {
	return &DeviceRepository{db: db}
}

func (r *DeviceRepository) List() ([]model.Device, error) {
	rows, err := r.db.Query(`
		SELECT id, name, mac, broadcast_addr, port, created_at, updated_at
		FROM devices ORDER BY name ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("query devices: %w", err)
	}
	defer rows.Close()

	var devices []model.Device
	for rows.Next() {
		var d model.Device
		if err := rows.Scan(&d.ID, &d.Name, &d.MAC, &d.BroadcastAddr, &d.Port, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan device: %w", err)
		}
		devices = append(devices, d)
	}
	return devices, rows.Err()
}

func (r *DeviceRepository) GetByID(id string) (*model.Device, error) {
	var d model.Device
	err := r.db.QueryRow(`
		SELECT id, name, mac, broadcast_addr, port, created_at, updated_at
		FROM devices WHERE id = ?
	`, id).Scan(&d.ID, &d.Name, &d.MAC, &d.BroadcastAddr, &d.Port, &d.CreatedAt, &d.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get device: %w", err)
	}
	return &d, nil
}

func (r *DeviceRepository) GetByName(name string) (*model.Device, error) {
	var d model.Device
	err := r.db.QueryRow(`
		SELECT id, name, mac, broadcast_addr, port, created_at, updated_at
		FROM devices WHERE name = ?
	`, name).Scan(&d.ID, &d.Name, &d.MAC, &d.BroadcastAddr, &d.Port, &d.CreatedAt, &d.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get device by name: %w", err)
	}
	return &d, nil
}

func (r *DeviceRepository) Create(req model.CreateDeviceRequest) (*model.Device, error) {
	now := time.Now().UTC()
	id := model.DeviceID(req.MAC)

	_, err := r.db.Exec(`
		INSERT INTO devices (id, name, mac, broadcast_addr, port, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, id, req.Name, req.MAC, req.BroadcastAddr, req.Port, now, now)
	if err != nil {
		return nil, fmt.Errorf("insert device: %w", err)
	}

	return &model.Device{
		ID:            id,
		Name:          req.Name,
		MAC:           req.MAC,
		BroadcastAddr: req.BroadcastAddr,
		Port:          req.Port,
		CreatedAt:     now,
		UpdatedAt:     now,
	}, nil
}

func (r *DeviceRepository) Update(id string, req model.CreateDeviceRequest) (*model.Device, error) {
	now := time.Now().UTC()
	newID := model.DeviceID(req.MAC)

	// MAC changed → ID changes, need delete + insert
	if newID != id {
		result, err := r.db.Exec("DELETE FROM devices WHERE id = ?", id)
		if err != nil {
			return nil, fmt.Errorf("delete old device: %w", err)
		}
		rows, _ := result.RowsAffected()
		if rows == 0 {
			return nil, nil
		}

		_, err = r.db.Exec(`
			INSERT INTO devices (id, name, mac, broadcast_addr, port, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?)
		`, newID, req.Name, req.MAC, req.BroadcastAddr, req.Port, now, now)
		if err != nil {
			return nil, fmt.Errorf("insert device with new id: %w", err)
		}

		return r.GetByID(newID)
	}

	// MAC unchanged → simple update
	result, err := r.db.Exec(`
		UPDATE devices SET name = ?, mac = ?, broadcast_addr = ?, port = ?, updated_at = ?
		WHERE id = ?
	`, req.Name, req.MAC, req.BroadcastAddr, req.Port, now, id)
	if err != nil {
		return nil, fmt.Errorf("update device: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return nil, nil
	}

	return r.GetByID(id)
}

func (r *DeviceRepository) Delete(id string) (bool, error) {
	result, err := r.db.Exec("DELETE FROM devices WHERE id = ?", id)
	if err != nil {
		return false, fmt.Errorf("delete device: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("rows affected: %w", err)
	}
	return rows > 0, nil
}
