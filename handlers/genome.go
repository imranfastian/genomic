package handlers

import (
	"encoding/json"
	"net/http"

	"genomic/config"
)

// Genome represents a genome record from the DB
type Genome struct {
	ID    int    `json:"id"`
	Label string `json:"label"`
}

// GenomesHandler returns all genomes (protected by JWT middleware)
func GenomesHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := config.DB.Query("SELECT id, label FROM genomes")
	if err != nil {
		http.Error(w, "Database query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var results []Genome
	for rows.Next() {
		var g Genome
		if err := rows.Scan(&g.ID, &g.Label); err != nil {
			http.Error(w, "Data scan failed", http.StatusInternalServerError)
			return
		}
		results = append(results, g)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
