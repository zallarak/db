package models

import (
	"time"
)

type User struct {
	ID        string    `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	PwHash    string    `json:"-" db:"pw_hash"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type UserRole string

const (
	RoleOwner  UserRole = "owner"
	RoleAdmin  UserRole = "admin"
	RoleMember UserRole = "member"
	RoleViewer UserRole = "viewer"
)

type Membership struct {
	UserID string   `json:"user_id" db:"user_id"`
	OrgID  string   `json:"org_id" db:"org_id"`
	Role   UserRole `json:"role" db:"role"`
}

type Org struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Project struct {
	ID     string `json:"id" db:"id"`
	OrgID  string `json:"org_id" db:"org_id"`
	Name   string `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Instance struct {
	ID        string    `json:"id" db:"id"`
	ProjectID string    `json:"project_id" db:"project_id"`
	Name      string    `json:"name" db:"name"`
	Plan      string    `json:"plan" db:"plan"`
	PgVersion int       `json:"pg_version" db:"pg_version"`
	Node      string    `json:"node" db:"node"`
	CTID      int       `json:"ctid" db:"ctid"`
	FQDN      string    `json:"fqdn" db:"fqdn"`
	Status    string    `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type ApiKey struct {
	ID         string    `json:"id" db:"id"`
	UserID     string    `json:"user_id" db:"user_id"`
	OrgID      string    `json:"org_id" db:"org_id"`
	Hash       string    `json:"-" db:"hash"`
	Prefix     string    `json:"prefix" db:"prefix"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	LastUsedAt *time.Time `json:"last_used_at" db:"last_used_at"`
}

type AuditLog struct {
	ID          string    `json:"id" db:"id"`
	ActorUserID string    `json:"actor_user_id" db:"actor_user"`
	OrgID       string    `json:"org_id" db:"org_id"`
	ResourceURN string    `json:"resource_urn" db:"resource_urn"`
	Action      string    `json:"action" db:"action"`
	DiffJSON    string    `json:"diff_json" db:"diff_json"`
	IP          string    `json:"ip" db:"ip"`
	Timestamp   time.Time `json:"timestamp" db:"ts"`
}

type Job struct {
	ID          string    `json:"id" db:"id"`
	Type        string    `json:"type" db:"type"`
	PayloadJSON string    `json:"payload_json" db:"payload_json"`
	Status      string    `json:"status" db:"status"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}