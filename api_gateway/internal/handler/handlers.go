package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func (s *HTTPServer) getProject(w http.ResponseWriter, r *http.Request) error {
	path := r.URL.Path
	parts := strings.Split(path, "/")
	id, err := strconv.ParseUint(parts[len(parts)-1], 10, 32)

	if err != nil {
		err = fmt.Errorf("failed to get id from URL: %w", err)
		return err
	}

	project, _ := s.client.GetProject(uint32(id))

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(project); err != nil {
		err = fmt.Errorf("failed to encode project to JSON: %w", err)
		return err
	}
	return err
}
