package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/channinghe/waken/internal/config"
	"github.com/channinghe/waken/internal/model"
	"github.com/channinghe/waken/internal/repository"
	"github.com/channinghe/waken/internal/wol"
	"github.com/go-chi/chi/v5"
)

type WakeHandler struct {
	repo *repository.DeviceRepository
	cfg  config.Config
}

func NewWakeHandler(repo *repository.DeviceRepository, cfg config.Config) *WakeHandler {
	return &WakeHandler{repo: repo, cfg: cfg}
}

// WakeByID sends a magic packet to a stored device by its hash ID.
func (h *WakeHandler) WakeByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "device id is required")
		return
	}

	device, err := h.repo.GetByID(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get device")
		return
	}
	if device == nil {
		writeError(w, http.StatusNotFound, "device not found")
		return
	}

	h.wakeDevice(w, device)
}

// WakeByName sends a magic packet to a stored device by its unique name.
func (h *WakeHandler) WakeByName(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if name == "" {
		writeError(w, http.StatusBadRequest, "device name is required")
		return
	}

	device, err := h.repo.GetByName(name)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get device")
		return
	}
	if device == nil {
		writeError(w, http.StatusNotFound, "device not found")
		return
	}

	h.wakeDevice(w, device)
}

// WakeByMAC sends a magic packet to an arbitrary MAC address.
func (h *WakeHandler) WakeByMAC(w http.ResponseWriter, r *http.Request) {
	var req model.WakeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := req.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	broadcastAddr := req.BroadcastAddr
	if broadcastAddr == "" {
		broadcastAddr = h.cfg.BroadcastAddr
	}
	port := req.Port
	if port == 0 {
		port = h.cfg.WOLPort
	}

	if err := wol.Send(req.MAC, broadcastAddr, port); err != nil {
		log.Printf("failed to send magic packet to %s: %v", req.MAC, err)
		writeError(w, http.StatusInternalServerError, "failed to send magic packet")
		return
	}

	log.Printf("magic packet sent to %s", req.MAC)
	writeJSON(w, http.StatusOK, map[string]string{
		"message": "magic packet sent",
		"mac":     req.MAC,
	})
}

func (h *WakeHandler) wakeDevice(w http.ResponseWriter, device *model.Device) {
	if err := wol.Send(device.MAC, device.BroadcastAddr, device.Port); err != nil {
		log.Printf("failed to send magic packet to %s (%s): %v", device.Name, device.MAC, err)
		writeError(w, http.StatusInternalServerError, "failed to send magic packet")
		return
	}

	log.Printf("magic packet sent to %s (%s)", device.Name, device.MAC)
	writeJSON(w, http.StatusOK, map[string]string{
		"message": "magic packet sent",
		"device":  device.Name,
		"mac":     device.MAC,
	})
}
