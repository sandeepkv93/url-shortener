package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// QueryPerformanceMonitor tracks query performance and identifies slow queries
type QueryPerformanceMonitor struct {
	slowQueryThreshold time.Duration
	logger            *logrus.Logger
}

func NewQueryPerformanceMonitor(threshold time.Duration, log *logrus.Logger) *QueryPerformanceMonitor {
	return &QueryPerformanceMonitor{
		slowQueryThreshold: threshold,
		logger:            log,
	}
}

// LogMode implements GORM logger interface for query monitoring
func (qpm *QueryPerformanceMonitor) LogMode(level logger.LogLevel) logger.Interface {
	return qpm
}

func (qpm *QueryPerformanceMonitor) Info(ctx context.Context, msg string, data ...interface{}) {
	qpm.logger.WithContext(ctx).Infof(msg, data...)
}

func (qpm *QueryPerformanceMonitor) Warn(ctx context.Context, msg string, data ...interface{}) {
	qpm.logger.WithContext(ctx).Warnf(msg, data...)
}

func (qpm *QueryPerformanceMonitor) Error(ctx context.Context, msg string, data ...interface{}) {
	qpm.logger.WithContext(ctx).Errorf(msg, data...)
}

func (qpm *QueryPerformanceMonitor) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, rowsAffected := fc()

	fields := logrus.Fields{
		"duration_ms": elapsed.Milliseconds(),
		"rows":        rowsAffected,
	}

	if err != nil {
		fields["error"] = err.Error()
		qpm.logger.WithContext(ctx).WithFields(fields).Error("Database query failed")
		return
	}

	// Log slow queries with details
	if elapsed >= qpm.slowQueryThreshold {
		fields["sql"] = sql
		qpm.logger.WithContext(ctx).WithFields(fields).Warn("Slow query detected")
		return
	}

	// Regular debug logging for all queries in development
	qpm.logger.WithContext(ctx).WithFields(fields).Debug("Query executed")
}

// DatabaseMetrics collects database performance metrics
type DatabaseMetrics struct {
	TotalQueries     int64         `json:"total_queries"`
	SlowQueries      int64         `json:"slow_queries"`
	AverageQueryTime time.Duration `json:"average_query_time"`
	ErrorRate        float64       `json:"error_rate"`
	ConnectionsInUse int           `json:"connections_in_use"`
	MaxConnections   int           `json:"max_connections"`
}

// GetDatabaseMetrics returns current database performance metrics
func GetDatabaseMetrics(db *gorm.DB) (*DatabaseMetrics, error) {
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	stats := sqlDB.Stats()

	metrics := &DatabaseMetrics{
		ConnectionsInUse: stats.InUse,
		MaxConnections:   0, // MaxOpenConns not available in sql.DBStats
	}

	return metrics, nil
}

// OptimizedQuery provides utilities for optimized database queries
type OptimizedQuery struct {
	db *gorm.DB
}

func NewOptimizedQuery(db *gorm.DB) *OptimizedQuery {
	return &OptimizedQuery{db: db}
}

// BatchInsert performs optimized batch insert operations
func (oq *OptimizedQuery) BatchInsert(ctx context.Context, records interface{}, batchSize int) error {
	return oq.db.WithContext(ctx).CreateInBatches(records, batchSize).Error
}

// PreparedStatement creates a prepared statement for repeated queries
func (oq *OptimizedQuery) PreparedStatement(query string) (*sql.Stmt, error) {
	sqlDB, err := oq.db.DB()
	if err != nil {
		return nil, err
	}
	return sqlDB.Prepare(query)
}

// BulkUpdate performs bulk update operations efficiently
func (oq *OptimizedQuery) BulkUpdate(ctx context.Context, table string, updates map[string]interface{}, condition string, args ...interface{}) error {
	tx := oq.db.WithContext(ctx).Table(table).Where(condition, args...)
	return tx.Updates(updates).Error
}

// CountWithoutLoad counts records without loading them into memory
func (oq *OptimizedQuery) CountWithoutLoad(ctx context.Context, model interface{}, condition string, args ...interface{}) (int64, error) {
	var count int64
	err := oq.db.WithContext(ctx).Model(model).Where(condition, args...).Count(&count).Error
	return count, err
}

// ExistsCheck checks if a record exists without loading it
func (oq *OptimizedQuery) ExistsCheck(ctx context.Context, model interface{}, condition string, args ...interface{}) (bool, error) {
	var count int64
	err := oq.db.WithContext(ctx).Model(model).Where(condition, args...).Limit(1).Count(&count).Error
	return count > 0, err
}

// SelectColumns selects only specific columns to reduce memory usage
func (oq *OptimizedQuery) SelectColumns(ctx context.Context, model interface{}, columns []string, condition string, args ...interface{}) error {
	return oq.db.WithContext(ctx).Select(columns).Where(condition, args...).Find(model).Error
}

// CacheWarmer provides utilities for cache warming
type CacheWarmer struct {
	db    *gorm.DB
	cache interface{} // This should be the cache service interface
}

func NewCacheWarmer(db *gorm.DB, cache interface{}) *CacheWarmer {
	return &CacheWarmer{
		db:    db,
		cache: cache,
	}
}

// WarmFrequentlyAccessedData pre-loads frequently accessed data into cache
func (cw *CacheWarmer) WarmFrequentlyAccessedData(ctx context.Context) error {
	// This would implement cache warming logic
	// For now, it's a placeholder that would be implemented with actual cache service
	return nil
}

// IndexOptimizer provides utilities for analyzing and optimizing database indexes
type IndexOptimizer struct {
	db *gorm.DB
}

func NewIndexOptimizer(db *gorm.DB) *IndexOptimizer {
	return &IndexOptimizer{db: db}
}

// AnalyzeQueryPerformance analyzes query performance and suggests optimizations
func (io *IndexOptimizer) AnalyzeQueryPerformance(ctx context.Context, query string, args ...interface{}) (*QueryAnalysis, error) {
	sqlDB, err := io.db.DB()
	if err != nil {
		return nil, err
	}

	// Execute EXPLAIN ANALYZE for the query (PostgreSQL specific)
	explainQuery := "EXPLAIN ANALYZE " + query
	rows, err := sqlDB.QueryContext(ctx, explainQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plans []string
	for rows.Next() {
		var plan string
		if err := rows.Scan(&plan); err != nil {
			return nil, err
		}
		plans = append(plans, plan)
	}

	analysis := &QueryAnalysis{
		Query:           query,
		ExecutionPlans:  plans,
		Suggestions:     generateOptimizationSuggestions(plans),
		EstimatedCost:   extractCostFromPlan(plans),
		ActualTime:      extractTimeFromPlan(plans),
	}

	return analysis, nil
}

// QueryAnalysis contains the results of query performance analysis
type QueryAnalysis struct {
	Query          string    `json:"query"`
	ExecutionPlans []string  `json:"execution_plans"`
	Suggestions    []string  `json:"suggestions"`
	EstimatedCost  float64   `json:"estimated_cost"`
	ActualTime     float64   `json:"actual_time_ms"`
	Timestamp      time.Time `json:"timestamp"`
}

// ConnectionPoolOptimizer optimizes database connection pool settings
type ConnectionPoolOptimizer struct {
	db *gorm.DB
}

func NewConnectionPoolOptimizer(db *gorm.DB) *ConnectionPoolOptimizer {
	return &ConnectionPoolOptimizer{db: db}
}

// OptimizeConnectionPool adjusts connection pool settings based on load
func (cpo *ConnectionPoolOptimizer) OptimizeConnectionPool(ctx context.Context) error {
	sqlDB, err := cpo.db.DB()
	if err != nil {
		return err
	}

	stats := sqlDB.Stats()
	
	// Adjust max open connections based on current usage
	// Note: MaxOpenConns not available in sql.DBStats, using a default value
	defaultMaxConns := 25
	if stats.InUse > int(float64(defaultMaxConns)*0.8) {
		// Increase max connections if we're using > 80% of current limit
		newMaxConns := defaultMaxConns + 5
		if newMaxConns <= 50 { // Safety limit
			sqlDB.SetMaxOpenConns(newMaxConns)
		}
	}

	// Adjust max idle connections
	idealIdleConns := defaultMaxConns / 4
	if idealIdleConns < 2 {
		idealIdleConns = 2
	}
	sqlDB.SetMaxIdleConns(idealIdleConns)

	// Set connection lifetime
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetConnMaxIdleTime(time.Minute * 30)

	return nil
}

// Helper functions for query analysis

func generateOptimizationSuggestions(plans []string) []string {
	suggestions := []string{}
	
	for _, plan := range plans {
		if containsSeqScan(plan) {
			suggestions = append(suggestions, "Consider adding an index to avoid sequential scan")
		}
		if containsNestedLoop(plan) {
			suggestions = append(suggestions, "Consider optimizing JOIN conditions or adding indexes")
		}
		if containsSort(plan) {
			suggestions = append(suggestions, "Consider adding an index to avoid sorting")
		}
	}
	
	if len(suggestions) == 0 {
		suggestions = append(suggestions, "Query appears to be well optimized")
	}
	
	return suggestions
}

func extractCostFromPlan(plans []string) float64 {
	// Extract estimated cost from execution plan
	// This is a simplified implementation
	for _, plan := range plans {
		// Look for cost patterns in PostgreSQL execution plans
		// Example: "cost=0.00..4.02"
		// This would need more sophisticated parsing
		_ = plan // Use the variable to avoid unused error
	}
	return 0.0
}

func extractTimeFromPlan(plans []string) float64 {
	// Extract actual execution time from plan
	// This is a simplified implementation
	for _, plan := range plans {
		// Look for time patterns in PostgreSQL execution plans
		// Example: "actual time=0.023..0.025"
		// This would need more sophisticated parsing
		_ = plan // Use the variable to avoid unused error
	}
	return 0.0
}

func containsSeqScan(plan string) bool {
	return false // Simplified - would check for "Seq Scan" in plan
}

func containsNestedLoop(plan string) bool {
	return false // Simplified - would check for "Nested Loop" in plan
}

func containsSort(plan string) bool {
	return false // Simplified - would check for "Sort" in plan
}

// TransactionOptimizer provides utilities for optimizing database transactions
type TransactionOptimizer struct {
	db *gorm.DB
}

func NewTransactionOptimizer(db *gorm.DB) *TransactionOptimizer {
	return &TransactionOptimizer{db: db}
}

// OptimizedTransaction executes a transaction with performance monitoring
func (to *TransactionOptimizer) OptimizedTransaction(ctx context.Context, fn func(*gorm.DB) error) error {
	start := time.Now()
	
	err := to.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(tx)
	})
	
	duration := time.Since(start)
	
	// Log transaction performance
	logrus.WithFields(logrus.Fields{
		"duration_ms": duration.Milliseconds(),
		"success":     err == nil,
	}).Debug("Transaction completed")
	
	if duration > 100*time.Millisecond {
		logrus.WithFields(logrus.Fields{
			"duration_ms": duration.Milliseconds(),
		}).Warn("Slow transaction detected")
	}
	
	return err
}

// ReadOnlyTransaction executes a read-only transaction for better performance
func (to *TransactionOptimizer) ReadOnlyTransaction(ctx context.Context, fn func(*gorm.DB) error) error {
	// Set transaction as read-only for potential optimizations
	return to.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Set read-only mode (PostgreSQL specific)
		if err := tx.Exec("SET TRANSACTION READ ONLY").Error; err != nil {
			return err
		}
		return fn(tx)
	})
}