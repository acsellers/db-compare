package main

import (
	"database/sql"
	"strconv"
	"time"
)

func toIntPtr(v sql.NullInt64) *int64 {
	if v.Valid {
		return &v.Int64
	}
	return nil
}
func toTimePtr(v sql.NullTime) *time.Time {
	if v.Valid {
		return &v.Time
	}
	return nil
}
func toStringPtr(v sql.NullString) *string {
	if v.Valid {
		return &v.String
	}
	return nil
}
func toStringOrZero(v sql.NullString) string {
	if v.Valid {
		return v.String
	}
	return ""
}
func toBoolPtr(v sql.NullBool) *bool {
	if v.Valid {
		return &v.Bool
	}
	return nil
}
func parseStringFloat(v string) float64 {
	fv, _ := strconv.ParseFloat(v, 64)
	return fv
}
func parseNullStringFloat(v sql.NullString) float64 {
	if v.Valid {
		return parseStringFloat(v.String)
	}
	return 0
}
