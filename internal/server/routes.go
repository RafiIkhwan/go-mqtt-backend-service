package server

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/", s.HelloWorldHandler)

	mux.HandleFunc("/health", s.healthHandler)
	
	mux.HandleFunc("/api/data/latest", s.GetLatestData)
	mux.HandleFunc("/api/data/history", s.GetHistoryData)
	mux.HandleFunc("/api/data/average", s.GetAverageData)

	// Wrap the mux with CORS middleware
	return s.corsMiddleware(mux)
}

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Replace "*" with specific origins if needed
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "false") // Set to "true" if credentials are required

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Proceed with the next handler
		next.ServeHTTP(w, r)
	})
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{"message": "Hello World!"}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(jsonResp); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := json.Marshal(s.db.Health())
	if err != nil {
		http.Error(w, "Failed to marshal health check response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(resp); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}


func (s *Server) GetLatestData(w http.ResponseWriter, r *http.Request) {
	data, err := s.db.GetLatestData()
	if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
	}
	json.NewEncoder(w).Encode(data)
}

func (s *Server) GetHistoryData(w http.ResponseWriter, r *http.Request) {
	deviceID := r.URL.Query().Get("device_id")
	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")

	startDate, err := time.Parse(time.RFC3339, start)
	if err != nil {
			http.Error(w, "Invalid start date format", http.StatusBadRequest)
			return
	}

	endDate, err := time.Parse(time.RFC3339, end)
	if err != nil {
			http.Error(w, "Invalid end date format", http.StatusBadRequest)
			return
	}

	data, err := s.db.GetHistoryData(deviceID, startDate, endDate)
	if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
	}
	json.NewEncoder(w).Encode(data)
}

func (s *Server) GetAverageData(w http.ResponseWriter, r *http.Request) {
	deviceID := r.URL.Query().Get("device_id")
	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")

	startDate, err := time.Parse(time.RFC3339, start)
	if err != nil {
			http.Error(w, "Invalid start date format", http.StatusBadRequest)
			return
	}

	endDate, err := time.Parse(time.RFC3339, end)
	if err != nil {
			http.Error(w, "Invalid end date format", http.StatusBadRequest)
			return
	}

	data, err := s.db.GetAverageData(deviceID, startDate, endDate)
	if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
	}
	json.NewEncoder(w).Encode(data)
}