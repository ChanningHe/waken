package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

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

// WakeByID sends a magic packet to a stored device.
func (h *WakeHandler) WakeByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid device id")
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
