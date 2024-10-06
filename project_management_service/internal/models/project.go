package models

import "time"

type Project struct {
	Id             uint32    `db:"id"`
	Name           string    `db:"name"`
	Description    string    `db:"description"`
	StartDate      time.Time `db:"start_date"`
	PlannedEndDate time.Time `db:"planned_end_date"`
	ActualEndDate  time.Time `db:"actual_end_date"`
	Status         string    `db:"status"`
	Priority       uint32    `db:"priority"`
	ManagerId      uint32    `db:"manager_id"`
	Budget         float64   `db:"budget"`
}
