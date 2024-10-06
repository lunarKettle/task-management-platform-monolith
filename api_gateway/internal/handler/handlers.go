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
	if err != nil {
		err = fmt.Errorf("failed to get parameter from request: %w", err)
		log.Print(err)
		return err
	}
	plannedEndDate, err := time.Parse(timeLayout, urlQuery.Get("planned_end_date"))
	if err != nil {
		err = fmt.Errorf("failed to get parameter from request: %w", err)
		log.Print(err)
		return err
	}
	//actualEndDate, err := time.Parse(timeLayout, urlQuery.Get("actual_end_date"))

	status := urlQuery.Get("status")

	priority, err := strconv.ParseUint(urlQuery.Get("priority"), 10, 32)
	if err != nil {
		err = fmt.Errorf("failed to get parameter from request: %w", err)
		log.Print(err)
		return err
	}
	managerId, err := strconv.ParseUint(urlQuery.Get("manager_id"), 10, 32)
	if err != nil {
		err = fmt.Errorf("failed to get parameter from request: %w", err)
		log.Print(err)
		return err
	}

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
		Priority:       uint32(priority),
		ManagerId:      uint32(managerId),
		Budget:         budget,
	}

	id, err := s.client.CreateProject(project)

	//TODO Add response type
	if err := json.NewEncoder(w).Encode(id); err != nil {
		err = fmt.Errorf("failed to encode project to JSON: %w", err)
		log.Print(err)
		return err
	}
	return err
}

func (s *HTTPServer) updateProject(w http.ResponseWriter, r *http.Request) error {
	var project models.Project
	err := json.NewDecoder(r.Body).Decode(&project)
	if err != nil {
		err = fmt.Errorf("error while decoding request body: %w", err)
		return err
	}
	if err := json.NewEncoder(w).Encode(project); err != nil {
		err = fmt.Errorf("failed to encode project to JSON: %w", err)
		log.Print(err)
		return err
	}
	return err
}
