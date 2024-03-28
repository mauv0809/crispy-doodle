package db

import (
	"database/sql"

	"github.com/charmbracelet/log"
	_ "github.com/mattn/go-sqlite3"
)

type Task struct {
	Id          int
	Name        string
	Description string
	Status      string
}

func Init() {
	db, err := sql.Open("sqlite3", "./tasks.db")
	if err != nil {
		log.Fatal(err)
	} else {
		log.Info("Database initialised")
	}
	//Setup task table
	tx, err := db.Begin()
	if err != nil {
		log.Error("Transaction failed: ", err)
	} else {
		_, err = tx.Exec("CREATE TABLE IF NOT EXISTS tasks (id INTEGER PRIMARY KEY, name TEXT, description TEXT, status TEXT)")
		if err != nil {
			tx.Rollback()
			log.Error("Table creation failed: ", err)
		} else {
			log.Debug("Table created")
			tx.Commit()
		}
	}
}
func AddTask(name, description, status string) (Task, error) {
	db, err := sql.Open("sqlite3", "./tasks.db")
	if err != nil {
		log.Fatal(err)
	}
	task := Task{Name: name, Description: description, Status: status}
	tx, err := db.Begin()
	if err != nil {
		log.Error("Transaction failed: ", err)
	} else {
		_, err = tx.Exec("INSERT INTO tasks (name, description, status) VALUES (?, ?, ?)", name, description, "created")
		if err != nil {
			tx.Rollback()
			log.Error("Insert failed: ", err)
		} else {
			log.Info("Task added")
			tx.Commit()
		}
	}
	return task, err
}
func GetTasks() ([]Task, error) {
	db, err := sql.Open("sqlite3", "./tasks.db")
	if err != nil {
		log.Fatal(err)
	}
	tasks := []Task{}
	rows, err := db.Query("SELECT id, name, description, status FROM tasks")
	if err != nil {
		log.Error("Query failed: ", err)
	} else {
		for rows.Next() {
			var id int
			var name string
			var description string
			var status string
			err = rows.Scan(&id, &name, &description, &status)
			if err != nil {
				log.Error("Row scan failed: ", err)
			} else {
				log.Debug("Task: ", id, name, description, status)
				tasks = append(tasks, Task{id, name, description, status})
			}
		}
	}
	log.Debug("Tasks fetched")
	return tasks, err
}
func DeleteTask(id int) {
	db, err := sql.Open("sqlite3", "./tasks.db")
	if err != nil {
		log.Fatal(err)
	}
	tx, err := db.Begin()
	if err != nil {
		log.Error("Transaction failed: ", err)
	} else {
		_, err = tx.Exec("DELETE FROM tasks WHERE id = ?", id)
		if err != nil {
			tx.Rollback()
			log.Error("Delete failed: ", err)
		} else {
			log.Debug("Task deleted")
			tx.Commit()
		}
	}
}
