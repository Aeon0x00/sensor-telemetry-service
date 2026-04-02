package repository

import (
	"context"
	"errors"
	"fmt"
	"machine-telemetry-service/internal/domain"
	"time"

	"github.com/jackc/pgx/v5"
)

// func (r *PostgresRepo) CreateIncidentSimple(ctx context.Context, machineId string, temp float64) error {
// 	_, err := r.pool.Exec(ctx, "INSERT INTO machine_incidents (machine_id, temp_c) VALUES ($1, $2);", machineId, temp)
// 	if err != nil {
// 		// logger будет обрабатывать
// 		return err
// 	}
// 	return nil
// }

func (r *PostgresRepo) CreateIncident(ctx context.Context, inc domain.Incident) error {
	fmt.Println("create incident for machine:", inc.MachineId)
	sqlQuery := `
	INSERT INTO machine_incidents (id, machine_id, event_id, temp_c, vibro, reason, is_resolved, created_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8);
	`
	
	_, err := r.pool.Exec(
		ctx,
		sqlQuery,
		inc.Id,
		inc.MachineId,
		inc.EventId,
		inc.TempC,
		inc.Vibro,
		inc.Reason,
		inc.IsResolved,
		inc.CreatedAt,
	)
	if err != nil {
		r.logger.Errorf("create incident error %v",err)
		return err
	}
	
	return  nil
}
func (r *PostgresRepo) ListActiveIncidents(ctx context.Context, machineId string) ([]domain.Incident, error) {
	sqlQuery := `
	SELECT id, machine_id, event_id, temp_c, vibro, reason, is_resolved, created_at, resolved_at
	FROM machine_incidents
	WHERE machine_id = $1 AND is_resolved = FALSE
	ORDER BY created_at DESC
	`
	rows, err := r.pool.Query(ctx, sqlQuery, machineId)
	if err != nil {
		return nil, err
	}
	
	defer rows.Close()

	list := []domain.Incident{}
	for rows.Next(){
		var inc domain.Incident
		err := rows.Scan(
			&inc.Id,
			&inc.MachineId,
			&inc.EventId,
			&inc.TempC,
			&inc.Vibro,
			&inc.Reason,
			&inc.IsResolved,
			&inc.CreatedAt,
			&inc.ResolvedAt,
		)
		if err != nil {
			return nil, err
		}


		list = append(list, inc)
	}
	return  list, nil
}


func (r *PostgresRepo) ResolveIncident(ctx context.Context, incidentId string) (string, error) {
	sqlQuery := `
	UPDATE machine_incidents
	SET is_resolved = TRUE, resolved_at = $1
	WHERE id = $2
	RETURNING machine_id;
	`
	var machineId string
	
	err := r.pool.QueryRow(ctx, sqlQuery, time.Now(), incidentId).Scan(&machineId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", domain.ErrIncidentNotFound
		}
		return "", err
	}
	
	return machineId, nil
}
