package postgres

import (
	"database/sql"
	"fmt"
	model "project/internal/storage/postgres/models"
)

func CreateProject(db *sql.DB, title, description string) error {
	_, err := db.Exec(`
 INSERT INTO projects (title, description, created_at)
 VALUES ($1, $2, NOW())`, title, description,
	)

	if err != nil {
		fmt.Println(" Ошибка при сохранении :", err)
	}
	return err

}

func GetProjects(db *sql.DB) ([]model.Project, error) {
	rows, err := db.Query("SELECT id, title, description, created_at FROM projects")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []model.Project

	for rows.Next() {
		var p model.Project
		err := rows.Scan(&p.ID, &p.Title, &p.Description, &p.Created_at)
		if err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}

	return projects, nil
}

func UpdateProject(db *sql.DB, id int, title, description string) error {
	_, err := db.Exec(`
  UPDATE projects 
  SET title =$1, description = $2
  WHERE ID =$3`, title, description, id)
	return err

}

func DeleteProject(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM projects WHERE id =$1", id)
	return err
}
