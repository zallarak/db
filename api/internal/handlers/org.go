package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/zallarak/db/api/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OrgHandler struct {
	db *sql.DB
}

func NewOrgHandler(db *sql.DB) *OrgHandler {
	return &OrgHandler{db: db}
}

type CreateOrgRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdateOrgRequest struct {
	Name string `json:"name"`
}

func (h *OrgHandler) ListOrgs(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	query := `
		SELECT o.id, o.name, o.created_at, o.updated_at, m.role
		FROM orgs o
		JOIN memberships m ON o.id = m.org_id
		WHERE m.user_id = $1
		ORDER BY o.created_at DESC
	`

	rows, err := h.db.Query(query, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get organizations"})
		return
	}
	defer rows.Close()

	var orgs []gin.H
	for rows.Next() {
		var org models.Org
		var role models.UserRole
		
		err := rows.Scan(&org.ID, &org.Name, &org.CreatedAt, &org.UpdatedAt, &role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan organization"})
			return
		}

		orgs = append(orgs, gin.H{
			"id":         org.ID,
			"name":       org.Name,
			"created_at": org.CreatedAt,
			"updated_at": org.UpdatedAt,
			"role":       role,
		})
	}

	c.JSON(http.StatusOK, gin.H{"orgs": orgs})
}

func (h *OrgHandler) CreateOrg(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req CreateOrgRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Start transaction
	tx, err := h.db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer tx.Rollback()

	// Create org
	org := models.Org{
		ID:        uuid.New().String(),
		Name:      req.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	orgQuery := `
		INSERT INTO orgs (id, name, created_at, updated_at) 
		VALUES ($1, $2, $3, $4)
	`
	_, err = tx.Exec(orgQuery, org.ID, org.Name, org.CreatedAt, org.UpdatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create organization"})
		return
	}

	// Add user as owner
	memberQuery := `
		INSERT INTO memberships (user_id, org_id, role) 
		VALUES ($1, $2, $3)
	`
	_, err = tx.Exec(memberQuery, userID, org.ID, models.RoleOwner)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create membership"})
		return
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"org": org})
}

func (h *OrgHandler) GetOrg(c *gin.Context) {
	orgID := c.Param("orgId")
	userID := c.GetString("user_id")

	// Check if user has access to this org
	var role models.UserRole
	roleQuery := "SELECT role FROM memberships WHERE user_id = $1 AND org_id = $2"
	err := h.db.QueryRow(roleQuery, userID, orgID).Scan(&role)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check access"})
		return
	}

	// Get org
	var org models.Org
	orgQuery := "SELECT id, name, created_at, updated_at FROM orgs WHERE id = $1"
	err = h.db.QueryRow(orgQuery, orgID).Scan(&org.ID, &org.Name, &org.CreatedAt, &org.UpdatedAt)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get organization"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"org":  org,
		"role": role,
	})
}

func (h *OrgHandler) UpdateOrg(c *gin.Context) {
	orgID := c.Param("orgId")
	userID := c.GetString("user_id")

	// Check if user is admin or owner
	var role models.UserRole
	roleQuery := "SELECT role FROM memberships WHERE user_id = $1 AND org_id = $2"
	err := h.db.QueryRow(roleQuery, userID, orgID).Scan(&role)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check access"})
		return
	}

	if role != models.RoleOwner && role != models.RoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}

	var req UpdateOrgRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := "UPDATE orgs SET name = $1, updated_at = $2 WHERE id = $3"
	_, err = h.db.Exec(query, req.Name, time.Now(), orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update organization"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Organization updated successfully"})
}

func (h *OrgHandler) DeleteOrg(c *gin.Context) {
	orgID := c.Param("orgId")
	userID := c.GetString("user_id")

	// Check if user is owner
	var role models.UserRole
	roleQuery := "SELECT role FROM memberships WHERE user_id = $1 AND org_id = $2"
	err := h.db.QueryRow(roleQuery, userID, orgID).Scan(&role)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check access"})
		return
	}

	if role != models.RoleOwner {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only owners can delete organizations"})
		return
	}

	// TODO: Check for dependent resources (projects, instances)

	query := "DELETE FROM orgs WHERE id = $1"
	_, err = h.db.Exec(query, orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete organization"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Organization deleted successfully"})
}