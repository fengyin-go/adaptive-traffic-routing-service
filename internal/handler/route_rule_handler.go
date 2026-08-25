package handler

import (
	"net/http"

	"loadbalancer/internal/model"
	"loadbalancer/pkg/httpx"
)

func (s *Server) registerRouteRuleRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/route-rules", s.createRouteRule)
	mux.HandleFunc("GET /api/route-rules", s.listRouteRules)
	mux.HandleFunc("GET /api/route-rules/{id}", s.getRouteRule)
	mux.HandleFunc("PUT /api/route-rules/{id}", s.updateRouteRule)
	mux.HandleFunc("DELETE /api/route-rules/{id}", s.deleteRouteRule)
}

type createRouteRuleRequest struct {
	PathPrefix string   `json:"path_prefix"`
	StrategyID string   `json:"strategy_id"`
	BackendIDs []string `json:"backend_ids"`
	Enabled    bool     `json:"enabled"`
}

func (s *Server) createRouteRule(w http.ResponseWriter, r *http.Request) {
	var req createRouteRuleRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	rr, err := s.svc.CreateRouteRule(model.RouteRule{PathPrefix: req.PathPrefix, StrategyID: req.StrategyID, BackendIDs: req.BackendIDs, Enabled: req.Enabled})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, rr)
}

func (s *Server) listRouteRules(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.RouteRuleFilter{
		PathPrefix: r.URL.Query().Get("path_prefix"),
		Keyword:    r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListRouteRules(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getRouteRule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	rr, err := s.svc.GetRouteRule(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, rr)
}

type updateRouteRuleRequest struct {
	PathPrefix string   `json:"path_prefix"`
	StrategyID string   `json:"strategy_id"`
	BackendIDs []string `json:"backend_ids"`
	Enabled    bool     `json:"enabled"`
}

func (s *Server) updateRouteRule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateRouteRuleRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	rr, err := s.svc.UpdateRouteRule(id, model.RouteRule{PathPrefix: req.PathPrefix, StrategyID: req.StrategyID, BackendIDs: req.BackendIDs, Enabled: req.Enabled})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, rr)
}

func (s *Server) deleteRouteRule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteRouteRule(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
