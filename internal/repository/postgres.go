package repository

import (
	"context"
	"errors"
	"machine-telemetry-service/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type PostgresRepo struct {
	pool   		*pgxpool.Pool
	logger 		*zap.SugaredLogger
}

func NewPostgresRepo (pool *pgxpool.Pool,logger *zap.SugaredLogger) *PostgresRepo {
	return &PostgresRepo{
		pool:   	pool,
		logger: 	logger,
	}
}

// func (r *PostgresRepo) GetByCode(ctx context.Context, code string) (*domain.Machine, error) {
// 	var m domain.Machine
// 	m.Code = code
// 	err := r.pool.QueryRow(ctx, "SELECT id, workshop, max_temp_c, max_vibro, created_at FROM machines WHERE code = $1;", code).Scan(&m.Id, &m.Workshop, &m.MaxTempC, &m.MaxVibro, &m.CreatedAt)
// 	if err != nil {
// 		// logger будет обрабатывать
// 		return nil, err
// 	}
// 	return &m, nil
// }

func (r *PostgresRepo) CreateMachine(ctx context.Context, m domain.Machine) error {
	sqlQuery := `
	INSERT INTO machines (id, code, workshop, eq_type, max_temp_c, max_vibro, created_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7);
	`
	
	_, err := r.pool.Exec(ctx, sqlQuery, m.Id, m.Code, m.Workshop, m.EqType, m.MaxTempC, m.MaxVibro, m.CreatedAt)
	if err != nil {
		r.logger.Errorf("create machine error %v",err)
		
		return err
	}
	
	return  nil
}
func (r *PostgresRepo) UpdateLimits(ctx context.Context, machineId string, maxTempC, maxVibro float64) error {
	sqlQuery := `
	UPDATE machines
	SET max_temp_c = $1, max_vibro = $2
	WHERE id = $3;
	`
	tag, err := r.pool.Exec(ctx, sqlQuery, maxTempC, maxVibro, machineId)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrMachineNotFound
	}
	return nil
}

func (r *PostgresRepo) GetMachineById(ctx context.Context, machineId string) (*domain.Machine, error) {
	sqlQuery := `
	SELECT id, code, workshop, eq_type, max_temp_c, max_vibro, created_at
	FROM machines
	WHERE id = $1
	`
	var m domain.Machine
	
	err := r.pool.QueryRow(ctx, sqlQuery, machineId).Scan(
		&m.Id,
		&m.Code,
		&m.Workshop,
		&m.EqType,
		&m.MaxTempC,
		&m.MaxVibro,
		&m.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrMachineNotFound
		}
		
		return  nil, err
	}
	
	return &m, nil
}

func (r *PostgresRepo) SaveTelemetryEvent(ctx context.Context, ev domain.TelemetryEvent) error {
	sqlQuery := `
	INSERT INTO telemetry_events (event_id, machine_id, temp_c, vibro, rpm, power_kw, is_overload, created_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	ON CONFLICT (event_id) DO NOTHING;
	`
	
	_, err := r.pool.Exec(
		ctx,
		sqlQuery,
		ev.EventId,
		ev.MachineId,
		ev.TempC,
		ev.Vibro,
		ev.Rpm,
		ev.PowerKw,
		ev.IsOverload,
		ev.CreatedAt,
	)
	if err != nil {
		r.logger.Errorf("save telemetry event error %v",err)
		return err
		
	}


	return nil
}
func (r *PostgresRepo) GetRecentEvents(ctx context.Context, machineId string, limit int) ([]domain.TelemetryEvent, error) {
	sqlQuery := `
	SELECT event_id, machine_id, temp_c, vibro, rpm, power_kw, is_overload, created_at
	FROM telemetry_events
	WHERE machine_id = $1
	ORDER BY created_at DESC
	LIMIT $2
	`
	rows, err := r.pool.Query(ctx, sqlQuery, machineId, limit)
	if err != nil {
		return nil, err
	}
	
	defer rows.Close()

	events := []domain.TelemetryEvent{}
	for rows.Next(){
		var ev domain.TelemetryEvent
		err := rows.Scan(
			&ev.EventId,
			&ev.MachineId,
			&ev.TempC,
			&ev.Vibro,
			&ev.Rpm,
			&ev.PowerKw,
			&ev.IsOverload,
			&ev.CreatedAt,
		)
		if err != nil {
			return nil, err
		}


		events = append(events, ev)
	}
	return  events, nil
}


func (r *PostgresRepo) CreateShiftReport(ctx context.Context, rep domain.ShiftReport) error {
	sqlQuery := `
	INSERT INTO shift_reports (id, machine_id, status, avg_temp_c, max_temp_c, avg_vibro, max_vibro, overload_count, risk_level, created_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);
	`
	_, err := r.pool.Exec(ctx, sqlQuery, rep.Id, rep.MachineId, rep.Status, rep.AvgTempC, rep.MaxTempC, rep.AvgVibro, rep.MaxVibro, rep.OverloadCount, rep.RiskLevel, rep.CreatedAt)
	return err
}

func (r *PostgresRepo) UpdateShiftReport(ctx context.Context, rep domain.ShiftReport) error {
	sqlQuery := `
	UPDATE shift_reports
	SET status = $1, avg_temp_c = $2, max_temp_c = $3, avg_vibro = $4, max_vibro = $5, overload_count = $6, risk_level = $7
	WHERE id = $8;
	`
	
	_, err := r.pool.Exec(ctx, sqlQuery, rep.Status, rep.AvgTempC, rep.MaxTempC, rep.AvgVibro, rep.MaxVibro, rep.OverloadCount, rep.RiskLevel, rep.Id)
	return err
}
