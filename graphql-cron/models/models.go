// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package models

type Attend struct {
	ID            string  `db:"id"`
	StaffID       string  `db:"staff_id"`
	InTime        string  `db:"in_time"`
	OutTime       string  `db:"out_time"`
	IsAttend      bool    `db:"is_attend"`
	DateStartTime string  `db:"date_start_time"`
	CreatedAt     *string `db:"created_at"`
	UpdatedAt     string  `db:"updated_at"`
}

type Staff struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	Age       int       `db:"age"`
	Attends   []*Attend `db:"attends"`
	CreatedAt *string   `db:"created_at"`
	UpdatedAt string    `db:"updated_at"`
}
