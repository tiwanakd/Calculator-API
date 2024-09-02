package models

import (
	"database/sql"
	"time"
)

type Calculation struct {
	ID        int
	Operation string
	NumberA   int
	NUmberB   int
	Result    float64
	Created   time.Time
}

type CalculationModel struct {
	DB *sql.DB
}

func (m *CalculationModel) Insert(operation string, numberA, numberB int, result float64) error {
	stmt := `INSERT INTO calculations (operation, number_a, number_b, result, created)
	VALUES (?, ?, ?, ?, NOW());`

	_, err := m.DB.Exec(stmt, operation, numberA, numberB, result)
	if err != nil {
		return err
	}

	return nil
}

func (m *CalculationModel) GetAll() ([]Calculation, error) {
	return nil, nil
}
