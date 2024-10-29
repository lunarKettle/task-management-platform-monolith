package handler

import (
	"api_gateway/internal/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func (s *HTTPServer) getProject(w http.ResponseWriter, r *http.Request) error {
	path := r.URL.Path
	parts := strings.Split(path, "/")
	id, err := strconv.ParseUint(parts[len(parts)-1], 10, 32)

	if err != nil {
		http.Error(w, "invalid project id", http.StatusBadRequest)
		return err
	}

	project, err := s.grpcClient.GetProject(uint32(id))

	if err != nil {
		return err
	}

	if err := json.NewEncoder(w).Encode(project); err != nil {
		err = fmt.Errorf("failed to encode project to JSON: %w", err)
		return err
	}
	return err
}

func (s *HTTPServer) createProject(w http.ResponseWriter, r *http.Request) error {
	var project models.Project
	err := json.NewDecoder(r.Body).Decode(&project)
	if err != nil {
		err = fmt.Errorf("error while decoding request body: %w", err)
		return err
	}
	defer r.Body.Close()

	id, err := s.grpcClient.CreateProject(project)
	//TODO Add response type
	if err := json.NewEncoder(w).Encode(id); err != nil {
		err = fmt.Errorf("failed to encode project to JSON: %w", err)
		log.Print(err)
		return err
	}
	return err
}

func (s *HTTPServer) updateProject(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPut {
		return fmt.Errorf("invalid request method: %s", r.Method)
	}

	var project models.Project
	err := json.NewDecoder(r.Body).Decode(&project)
	if err != nil {
		err = fmt.Errorf("error while decoding request body: %w", err)
		return err
	}
	defer r.Body.Close()

	err = s.grpcClient.UpdateProject(project)

	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)

	return nil
}

func (s *HTTPServer) deleteProject(w http.ResponseWriter, r *http.Request) error {
	path := r.URL.Path
	parts := strings.Split(path, "/")
	id, err := strconv.ParseUint(parts[len(parts)-1], 10, 32)

	if err != nil {
		http.Error(w, "invalid project id", http.StatusBadRequest)
		return err
	}

	err = s.grpcClient.DeleteProject(uint32(id))

	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)

	return nil
}
