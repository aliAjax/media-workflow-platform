package httpapi

import (
	"encoding/json"
	"errors"
	"example.com/media-workflow/internal/application"
	"example.com/media-workflow/internal/domain"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	svc       *application.Service
	assets    domain.AssetRepository
	pipelines domain.PipelineRepository
	jobs      domain.JobRepository
}

func New(s *application.Service, a domain.AssetRepository, p domain.PipelineRepository, j domain.JobRepository) *Handler {
	return &Handler{svc: s, assets: a, pipelines: p, jobs: j}
}
func (h *Handler) Routes() *http.ServeMux {
	m := http.NewServeMux()
	m.HandleFunc("/healthz", h.health)
	m.HandleFunc("/readyz", h.health)
	m.HandleFunc("/api/v1/assets", h.assetsRoute)
	m.HandleFunc("/api/v1/pipelines", h.pipelineRoute)
	m.HandleFunc("/api/v1/jobs", h.jobsRoute)
	m.HandleFunc("/api/v1/jobs/", h.jobDetail)
	return m
}
func (h *Handler) health(w http.ResponseWriter, _ *http.Request) {
	write(w, 200, map[string]any{"status": "ok", "service": "media-workflow"})
}
func (h *Handler) assetsRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		l, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		v, e := h.assets.ListAssets(r.Context(), r.URL.Query().Get("tenant_id"), l)
		respond(w, v, e)
		return
	}
	if r.Method != http.MethodPost {
		w.WriteHeader(405)
		return
	}
	var a domain.Asset
	if !decode(r, &a) {
		writeError(w, domain.ErrInvalidInput)
		return
	}
	created, e := h.svc.CreateAsset(r.Context(), a)
	if e != nil {
		writeError(w, e)
		return
	}
	write(w, 201, created)
}
func (h *Handler) pipelineRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		v, e := h.pipelines.ListPipelines(r.Context(), r.URL.Query().Get("tenant_id"))
		respond(w, v, e)
		return
	}
	if r.Method != http.MethodPost {
		w.WriteHeader(405)
		return
	}
	var p domain.Pipeline
	if !decode(r, &p) {
		writeError(w, domain.ErrInvalidInput)
		return
	}
	created, e := h.svc.CreatePipeline(r.Context(), p)
	if e != nil {
		writeError(w, e)
		return
	}
	write(w, 201, created)
}
func (h *Handler) jobsRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		v, e := h.jobs.ListJobs(r.Context(), r.URL.Query().Get("tenant_id"))
		respond(w, v, e)
		return
	}
	if r.Method != http.MethodPost {
		w.WriteHeader(405)
		return
	}
	var j domain.Job
	if !decode(r, &j) {
		writeError(w, domain.ErrInvalidInput)
		return
	}
	created, e := h.svc.CreateJob(r.Context(), j)
	if e != nil {
		writeError(w, e)
		return
	}
	write(w, 202, created)
}
func (h *Handler) jobDetail(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/jobs/")
	j, e := h.jobs.GetJob(r.Context(), id)
	if e != nil {
		writeError(w, e)
		return
	}
	write(w, 200, j)
}
func decode(r *http.Request, v any) bool {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v) == nil
}
func write(w http.ResponseWriter, s int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(s)
	_ = json.NewEncoder(w).Encode(v)
}
func respond(w http.ResponseWriter, v any, e error) {
	if e != nil {
		writeError(w, e)
		return
	}
	write(w, 200, v)
}
func writeError(w http.ResponseWriter, e error) {
	s := 500
	if errors.Is(e, domain.ErrInvalidInput) {
		s = 400
	}
	if errors.Is(e, domain.ErrNotFound) {
		s = 404
	}
	if errors.Is(e, domain.ErrConflict) {
		s = 409
	}
	write(w, s, map[string]string{"error": e.Error()})
}
