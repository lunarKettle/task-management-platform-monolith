package handler

import (
	"api_gateway/internal/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
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

func (s *HTTPServer) createProject(w http.ResponseWriter, r *http.Request) error {
	urlQuery := r.URL.Query()

	name := urlQuery.Get("name")
	description := urlQuery.Get("description")

	timeLayout := "02-01-2006"
	startDate, err := time.Parse(timeLayout, urlQuery.Get("start_date"))
	plannedEndDate, err := time.Parse(timeLayout, urlQuery.Get("planned_end_date"))
	//actualEndDate, err := time.Parse(timeLayout, urlQuery.Get("actual_end_date"))

	status := urlQuery.Get("status")
	priority := urlQuery.Get("priority")
	managerId, err := strconv.ParseUint(urlQuery.Get("manager_id"), 10, 32)
	budget, err := strconv.ParseFloat(strings.TrimSpace(urlQuery.Get("budget")), 64)

	if err != nil {
		err = fmt.Errorf("failed to get parameter from request: %w", err)
		log.Print(err)
		return err
	}

	project := models.Project{
		Name:           name,
		Description:    description,
		StartDate:      startDate,
		PlannedEndDate: plannedEndDate,
		Status:         status,
		Priority:       priority,
		ManagerId:      uint32(managerId),
		Budget:         budget,
	}

	id, err := s.client.CreateProject(project)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(id); err != nil {
		err = fmt.Errorf("failed to encode project to JSON: %w", err)
		log.Print(err)
		return err
	}
	return err
}
