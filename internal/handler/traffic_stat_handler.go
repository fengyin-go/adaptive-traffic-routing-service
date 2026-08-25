package handler

import (
	"net/http"
	"time"

	"loadbalancer/internal/model"
	"loadbalancer/pkg/httpx"
)

func (s *Server) registerTrafficStatRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/traffic-stats", s.createTrafficStat)
	mux.HandleFunc("GET /api/traffic-stats", s.listTrafficStats)
	mux.HandleFunc("GET /api/traffic-stats/{id}", s.getTrafficStat)
	mux.HandleFunc("PUT /api/traffic-stats/{id}", s.updateTrafficStat)
	mux.HandleFunc("DELETE /api/traffic-stats/{id}", s.deleteTrafficStat)
	mux.HandleFunc("GET /api/traffic-stats/summary/global", s.globalSummary)
	mux.HandleFunc("GET /api/traffic-stats/summary/by-strategy", s.summaryByStrategy)
	mux.HandleFunc("GET /api/traffic-stats/summary/backend/{backend_id}", s.backendSummary)
}

type createTrafficStatRequest struct {
	BackendID       string    `json:"backend_id"`
	RequestCount    int64     `json:"request_count"`
	ErrorCount      int64     `json:"error_count"`
	TotalDurationMs int64     `json:"total_duration_ms"`
	LastRequestAt   time.Time `json:"last_request_at"`
}

func (s *Server) createTrafficStat(w http.ResponseWriter, r *http.Request) {
	var req createTrafficStatRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	t, err := s.svc.CreateTrafficStat(model.TrafficStat{BackendID: req.BackendID, RequestCount: req.RequestCount, ErrorCount: req.ErrorCount, TotalDurationMs: req.TotalDurationMs, LastRequestAt: req.LastRequestAt})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, t)
}

func (s *Server) listTrafficStats(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.TrafficStatFilter{
		BackendID: r.URL.Query().Get("backend_id"),
	}
	items, total, err := s.svc.ListTrafficStats(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getTrafficStat(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	t, err := s.svc.GetTrafficStat(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, t)
}

type updateTrafficStatRequest struct {
	RequestCount    int64     `json:"request_count"`
	ErrorCount      int64     `json:"error_count"`
	TotalDurationMs int64     `json:"total_duration_ms"`
	LastRequestAt   time.Time `json:"last_request_at"`
}

func (s *Server) updateTrafficStat(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateTrafficStatRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	t, err := s.svc.UpdateTrafficStat(id, model.TrafficStat{RequestCount: req.RequestCount, ErrorCount: req.ErrorCount, TotalDurationMs: req.TotalDurationMs, LastRequestAt: req.LastRequestAt})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, t)
}

func (s *Server) deleteTrafficStat(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteTrafficStat(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) globalSummary(w http.ResponseWriter, r *http.Request) {
	summary, err := s.svc.GlobalSummary()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, summary)
}

func (s *Server) summaryByStrategy(w http.ResponseWriter, r *http.Request) {
	summary, err := s.svc.SummaryByStrategy()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, summary)
}

func (s *Server) backendSummary(w http.ResponseWriter, r *http.Request) {
	backendID := r.PathValue("backend_id")
	summary, err := s.svc.BackendSummary(backendID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, summary)
}
