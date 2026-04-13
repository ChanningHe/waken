package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/channinghe/waken/internal/config"
	"github.com/channinghe/waken/internal/model"
	"github.com/channinghe/waken/internal/repository"
	"github.com/go-chi/chi/v5"
)

type DeviceHandler struct {
	repo *repository.DeviceRepository
	cfg  config.Config
}

func NewDeviceHandler(repo *repository.DeviceRepository, cfg config.Config) *DeviceHandler {
	return &DeviceHandler{repo: repo, cfg: cfg}
}

func (h *DeviceHandler) List(w http.ResponseWriter, r *http.Request) {
	devices, err := h.repo.List()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list devices")
		return
	}
	if devices == nil {
		devices = []model.Device{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"devices": devices})
}

func (h *DeviceHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.BroadcastAddr == "" {
		req.BroadcastAddr = h.cfg.BroadcastAddr
	}
	if req.Port == 0 {
		req.Port = h.cfg.WOLPort
	}

	if err := req.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	device, err := h.repo.Create(req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create device")
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"device": device})
}

func (h *DeviceHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid device id")
		return
	}

	var req model.CreateDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.BroadcastAddr == "" {
		req.BroadcastAddr = h.cfg.BroadcastAddr
	}
	if req.Port == 0 {
		req.Port = h.cfg.WOLPort
	}

	if err := req.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	device, err := h.repo.Update(id, req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update device")
		return
	}
	if device == nil {
		writeError(w, http.StatusNotFound, "device not found")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"device": device})
}

func (h *DeviceHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid device id")
		return
	}

	found, err := h.repo.Delete(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete device")
		return
	}
	if !found {
		writeError(w, http.StatusNotFound, "device not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
