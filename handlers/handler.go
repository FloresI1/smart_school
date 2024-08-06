package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"	
	"net/http"
	"strconv"
	"time"

	"github.com/FloresI1/smart_school/db"
	"github.com/FloresI1/smart_school/model"
	"github.com/google/uuid"
)

func handleError(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func CreateMaterialHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		handleError(w, http.StatusMethodNotAllowed, "Invalid request method")
		return
	}

	var material model.Material
	err := json.NewDecoder(r.Body).Decode(&material)
	if err != nil {
		handleError(w, http.StatusBadRequest, err.Error())
		return
	}

	materialUUID := uuid.New()
	material.UUID = materialUUID.String()
	material.CreatedAt = time.Now()
	material.UpdatedAt = time.Now()

	db, err := db.DBConn()
	if err != nil {
		handleError(w, http.StatusInternalServerError, "Failed to connect to database")
		return
	}
	defer db.Close()

	sqlStatement := `
        INSERT INTO materials (uuid, material_type, status, title, content, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `
	_, err = db.Exec(sqlStatement, material.UUID, material.MaterialType, material.Status, material.Title, material.Content, material.CreatedAt, material.UpdatedAt)
	if err != nil {
		handleError(w, http.StatusInternalServerError, "Failed to insert material")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(material)
}

func GetMaterialHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		handleError(w, http.StatusMethodNotAllowed, "Invalid request method")
		return
	}

	uuid := r.URL.Query().Get("uuid")
	if uuid == "" {
		handleError(w, http.StatusBadRequest, "UUID is required")
		return
	}

	db, err := db.DBConn()
	if err != nil {
		handleError(w, http.StatusInternalServerError, "Failed to connect to database")
		return
	}
	defer db.Close()

	var material model.Material
	sqlStatement := `SELECT uuid, material_type, status, title, content, created_at, updated_at FROM materials WHERE uuid=$1`
	row := db.QueryRow(sqlStatement, uuid)
	err = row.Scan(&material.UUID, &material.MaterialType, &material.Status, &material.Title, &material.Content, &material.CreatedAt, &material.UpdatedAt)
	if err == sql.ErrNoRows {
		handleError(w, http.StatusNotFound, "Material not found")
		return
	} else if err != nil {
		handleError(w, http.StatusInternalServerError, "Failed to query material")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(material)
}

func UpdateMaterialHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		handleError(w, http.StatusMethodNotAllowed, "Invalid request method")
		return
	}

	var material model.Material
	err := json.NewDecoder(r.Body).Decode(&material)
	if err != nil {
		handleError(w, http.StatusBadRequest, err.Error())
		return
	}

	if material.UUID == "" {
		handleError(w, http.StatusBadRequest, "UUID is required")
		return
	}

	material.UpdatedAt = time.Now()

	db, err := db.DBConn()
	if err != nil {
		handleError(w, http.StatusInternalServerError, "Failed to connect to database")
		return
	}
	defer db.Close()

	sqlStatement := `
        UPDATE materials 
        SET status=$1, title=$2, content=$3, updated_at=$4 
        WHERE uuid=$5 AND material_type=$6
    `
	result, err := db.Exec(sqlStatement, material.Status, material.Title, material.Content, material.UpdatedAt, material.UUID, material.MaterialType)
	if err != nil {
		handleError(w, http.StatusInternalServerError, "Failed to update material")
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		handleError(w, http.StatusInternalServerError, "Failed to retrieve rows affected")
		return
	}

	if rowsAffected == 0 {
		handleError(w, http.StatusNotFound, "Material not found or no changes made")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(material)
}

func GetAllMaterialsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		handleError(w, http.StatusMethodNotAllowed, "Invalid request method")
		return
	}

	db, err := db.DBConn()
	if err != nil {
		handleError(w, http.StatusInternalServerError, "Failed to connect to database")
		return
	}
	defer db.Close()

	query := `SELECT uuid, material_type, title, created_at, updated_at FROM materials WHERE status = 'активный'`
	var filters []interface{}
	var filterCount int

	materialType := r.URL.Query().Get("material_type")
	if materialType != "" {
		filterCount++
		query += fmt.Sprintf(" AND material_type = $%d", filterCount)
		filters = append(filters, materialType)
	}

	dateFrom := r.URL.Query().Get("date_from")
	if dateFrom != "" {
		filterCount++
		query += fmt.Sprintf(" AND created_at >= $%d", filterCount)
		filters = append(filters, dateFrom)
	}

	dateTo := r.URL.Query().Get("date_to")
	if dateTo != "" {
		filterCount++
		query += fmt.Sprintf(" AND created_at <= $%d", filterCount)
		filters = append(filters, dateTo)
	}

	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")
	if page != "" && limit != "" {
		pageInt, err := strconv.Atoi(page)
		if err == nil && pageInt > 0 {
			limitInt, err := strconv.Atoi(limit)
			if err == nil && limitInt > 0 {
				offset := (pageInt - 1) * limitInt
				query += fmt.Sprintf(" LIMIT %d OFFSET %d", limitInt, offset)
			}
		}
	}

	rows, err := db.Query(query, filters...)
	if err != nil {
		handleError(w, http.StatusInternalServerError, "Failed to query materials")
		return
	}
	defer rows.Close()

	var materials []model.Material
	for rows.Next() {
		var material model.Material
		err := rows.Scan(&material.UUID, &material.MaterialType, &material.Title, &material.CreatedAt, &material.UpdatedAt)
		if err != nil {
			handleError(w, http.StatusInternalServerError, "Failed to scan material")
			return
		}
		materials = append(materials, material)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(materials)
}
