package postgres

import (
	"context"
	"fmt"

	"uvico-api/internal/domain/shared"
	"uvico-api/internal/platform/database"
)

type AnalyticsRepo struct {
	pool *database.Pool
}

func NewAnalyticsRepo(pool *database.Pool) *AnalyticsRepo {
	return &AnalyticsRepo{pool: pool}
}

type VisitInput struct {
	TenantID  shared.TenantID
	SessionID string
	Path      string
	Device    string
	Browser   string
	OS        string
	Referrer  string
	UTMSource string
}

func (r *AnalyticsRepo) TrackVisit(ctx context.Context, v VisitInput) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO store_visits (tenant_id, session_id, path, device, browser, os, referrer, utm_source)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, v.TenantID.String(), v.SessionID, v.Path, v.Device, v.Browser, v.OS, v.Referrer, v.UTMSource)
	
	if err != nil {
		return fmt.Errorf("analytics_repo: track visit: %w", err)
	}
	return nil
}

type VisitPoint struct {
	Time   string `json:"time"`
	Visits int64  `json:"visits"`
}

func (r *AnalyticsRepo) GetVisitsTimeSeries(ctx context.Context, tenantID shared.TenantID, timeRange string) ([]VisitPoint, error) {
	var query string

	// Agrupamos según el rango de tiempo solicitado
	switch timeRange {
	case "24h":
		query = `
			SELECT to_char(date_trunc('hour', created_at), 'HH24:00') as time_bucket, COUNT(*)::bigint as visits
			FROM store_visits
			WHERE tenant_id = $1 AND created_at >= NOW() - INTERVAL '24 hours'
			GROUP BY time_bucket ORDER BY time_bucket ASC
		`
	case "7d":
		query = `
			SELECT to_char(date_trunc('day', created_at), 'DD Mon') as time_bucket, COUNT(*)::bigint as visits
			FROM store_visits
			WHERE tenant_id = $1 AND created_at >= NOW() - INTERVAL '7 days'
			GROUP BY time_bucket ORDER BY MIN(created_at) ASC
		`
	case "30d":
		query = `
			SELECT to_char(date_trunc('day', created_at), 'DD Mon') as time_bucket, COUNT(*)::bigint as visits
			FROM store_visits
			WHERE tenant_id = $1 AND created_at >= NOW() - INTERVAL '30 days'
			GROUP BY time_bucket ORDER BY MIN(created_at) ASC
		`
	default: // "1y"
		query = `
			SELECT to_char(date_trunc('month', created_at), 'Mon YYYY') as time_bucket, COUNT(*)::bigint as visits
			FROM store_visits
			WHERE tenant_id = $1 AND created_at >= NOW() - INTERVAL '1 year'
			GROUP BY time_bucket ORDER BY MIN(created_at) ASC
		`
	}

	rows, err := r.pool.Query(ctx, query, tenantID.String())
	if err != nil {
		return nil, fmt.Errorf("analytics_repo: get visits: %w", err)
	}
	defer rows.Close()

	var points []VisitPoint
	for rows.Next() {
		var p VisitPoint
		if err := rows.Scan(&p.Time, &p.Visits); err != nil {
			return nil, err
		}
		points = append(points, p)
	}
	return points, rows.Err()
}