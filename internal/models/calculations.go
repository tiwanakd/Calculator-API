package models

import (
	"database/sql"
	"errors"
	"time"
)

type Calculation struct {
	ID        int
	Operation string
	NumberA   int
	NumberB   int
	Result    float64
	Created   time.Time
}

type CalculationModel struct {
	DB *sql.DB
}

func (m *CalculationModel) Insert(operation string, numberA, numberB int, result float64) error {
	stmt := `INSERT INTO calculations (operation, number_a, number_b, result, created)
	VALUES ($1, $2, $3, $4, NOW());`

	_, err := m.DB.Exec(stmt, operation, numberA, numberB, result)
	if err != nil {
		return err
	}

	return nil
}

func (m *CalculationModel) Get(id int) (Calculation, error) {
	stmt := "SELECT * FROM calculations WHERE id = $1"
	var c Calculation
	err := m.DB.QueryRow(stmt, id).Scan(&c.ID, &c.Operation, &c.NumberA, &c.NumberB, &c.Result, &c.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Calculation{}, ErrNotFound
		} else {
			return Calculation{}, err
		}
	}

	return c, nil
}

func (m *CalculationModel) GetAll() ([]Calculation, error) {
	calculations, err := m.getMultiple("SELECT * FROM calculations")
	if err != nil {
		return nil, err
	}

	return calculations, nil
}

func (m *CalculationModel) GetCalculations(operation string) ([]Calculation, error) {
	stmt := `SELECT * FROM calculations WHERE LOWER(operation)=LOWER($1)`
	calculations, err := m.getMultiple(stmt, operation)
	if err != nil {
		return nil, err
	}
	return calculations, nil
}

func (m *CalculationModel) GetLatestCalculations() ([]Calculation, error) {
	stmt := `SELECT * FROM calculations ORDER BY created DESC LIMIT 5`

	calculations, err := m.getMultiple(stmt)
	if err != nil {
		return nil, nil
	}
	return calculations, nil
}

func (m *CalculationModel) getMultiple(stmt string, args ...any) ([]Calculation, error) {
	rows, err := m.DB.Query(stmt, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var calculations []Calculation

	for rows.Next() {
		var c Calculation
		err = rows.Scan(&c.ID, &c.Operation, &c.NumberA, &c.NumberB, &c.Result, &c.Created)
		if err != nil {
			return nil, err
		}
		calculations = append(calculations, c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return calculations, nil
}
