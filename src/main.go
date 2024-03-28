/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

// import "crispy-doodle/cmd"
import (
	"crispy-doodle/src/internal/db"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	choices  []db.Task       // items on the to-do list
	cursor   int             // which to-do list item our cursor is pointing at
	selected map[int]db.Task // which to-do items are selected
}

func initialModel() model {
	return model{
		// Our to-do list is a grocery list
		choices: []db.Task{},

		// A map which indicates which choices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		selected: make(map[int]db.Task),
	}
}
func getTasksFromDB() tea.Msg {
	tasks, err := db.GetTasks()
	if err != nil {
		return getTasksFromDBMsgError{err}
	}
	return getTasksFromDBMsg(tasks)

}

type getTasksFromDBMsgError struct{ err error }
type getTasksFromDBMsg []db.Task

func (m model) Init() tea.Cmd {
	return getTasksFromDB
}
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case getTasksFromDBMsg:
		m.choices = msg

	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {
		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter":
			_, ok := m.selected[m.cursor]
			if ok {
				//Remove task from list and db
				delete(m.selected, m.cursor)
				db.DeleteTask(m.choices[m.cursor].Id)
				//Update the UI
				return m, getTasksFromDB
			}
		case " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = db.Task{}
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}
func (m model) View() string {
	// The header
	s := "What should we buy at the market?\n\n"

	// Iterate over our choices
	for i, choice := range m.choices {
		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if _, ok := m.selected[i]; ok {
			checked = "x" // selected!
		}
		// Render the row
		s += fmt.Sprintf("%s [%s] Id: %d | Name: %s | Description: %s | Status: %s\n", cursor, checked, choice.Id, choice.Name, choice.Description, choice.Status)
	}
	s += "\n"

	s += "Press A to add a task\n"
	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}
func main() {
	db.Init()
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
