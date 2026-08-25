package handler

import (
	"net/http"

	"loadbalancer/internal/model"
	"loadbalancer/pkg/httpx"
)

func (s *Server) registerBackendRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/backends", s.createBackend)
	mux.HandleFunc("GET /api/backends", s.listBackends)
	mux.HandleFunc("GET /api/backends/{id}", s.getBackend)
	mux.HandleFunc("PUT /api/backends/{id}", s.updateBackend)
	mux.HandleFunc("DELETE /api/backends/{id}", s.deleteBackend)
	mux.HandleFunc("POST /api/backends/{id}/status", s.transitionBackendStatus)
	mux.HandleFunc("POST /api/backends/{id}/health-check", s.recordHealthCheck)
}

type createBackendRequest struct {
	Name   string `json:"name"`
	Host   string `json:"host"`
	Port   int    `json:"port"`
	Weight int    `json:"weight"`
	Status string `json:"status"`
}

func (s *Server) createBackend(w http.ResponseWriter, r *http.Request) {
	var req createBackendRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	b, err := s.svc.CreateBackend(model.Backend{Name: req.Name, Host: req.Host, Port: req.Port, Weight: req.Weight, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, b)
}

func (s *Server) listBackends(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.BackendFilter{
		Status:  r.URL.Query().Get("status"),
		Keyword: r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListBackends(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getBackend(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	b, err := s.svc.GetBackend(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, b)
}

type updateBackendRequest struct {
	Name   string `json:"name"`
	Host   string `json:"host"`
	Port   int    `json:"port"`
	Weight int    `json:"weight"`
}

func (s *Server) updateBackend(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateBackendRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	b, err := s.svc.UpdateBackend(id, model.Backend{Name: req.Name, Host: req.Host, Port: req.Port, Weight: req.Weight})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, b)
}

func (s *Server) deleteBackend(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteBackend(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

type transitionStatusRequest struct {
	Status string `json:"status"`
}

func (s *Server) transitionBackendStatus(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req transitionStatusRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	b, err := s.svc.TransitionBackendStatus(id, req.Status)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, b)
}

type healthCheckRequest struct {
	Success bool `json:"success"`
}

func (s *Server) recordHealthCheck(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req healthCheckRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	b, err := s.svc.RecordHealthCheck(id, req.Success)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, b)
}
