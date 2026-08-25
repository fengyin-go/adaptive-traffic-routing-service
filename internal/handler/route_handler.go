package handler

import (
	"net/http"

	"loadbalancer/pkg/httpx"
)

func (s *Server) registerRouteRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/route/select", s.selectRoute)
	mux.HandleFunc("GET /api/route/available", s.listAvailableBackends)
}

type selectRouteRequest struct {
	Path string `json:"path"`
}

func (s *Server) selectRoute(w http.ResponseWriter, r *http.Request) {
	var req selectRouteRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	res, err := s.svc.SelectRoute(req.Path)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, res)
}

func (s *Server) listAvailableBackends(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		httpx.BadRequest(w, "path 参数不能为空")
		return
	}
	backends, err := s.svc.ListAvailableBackendsForPath(path)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, backends)
}
