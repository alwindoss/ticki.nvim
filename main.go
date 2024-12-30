package main

import (
	"fmt"
	"log"
	"os"

	"github.com/neovim/go-client/nvim"
	"github.com/neovim/go-client/nvim/plugin"
)

// TodoManager manages the TODO list.
type TodoManager struct {
	todos []string
}

// AddTodo adds a new TODO item.
func (tm *TodoManager) AddTodo(args []string) string {
	if len(args) == 0 {
		return "Error: TODO item cannot be empty!"
	}
	tm.todos = append(tm.todos, args[0])
	return fmt.Sprintf("Added TODO: %s", args[0])
}

// ListTodos lists all TODO items.
func (tm *TodoManager) ListTodos() []string {
	if len(tm.todos) == 0 {
		return []string{"No TODO items found."}
	}
	return tm.todos
}

// RemoveTodo removes a TODO item by index.
func (tm *TodoManager) RemoveTodo(index int) string {
	if index < 0 || index >= len(tm.todos) {
		return "Error: Invalid index!"
	}
	removed := tm.todos[index]
	tm.todos = append(tm.todos[:index], tm.todos[index+1:]...)
	return fmt.Sprintf("Removed TODO: %s", removed)
}

func main() {
	// Log to a file for debugging
	logFile, err := os.OpenFile("/tmp/ticki_plugin.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	log.SetOutput(logFile)
	log.Println("Plugin started")
	todoManager := &TodoManager{}

	// Initialize the plugin.
	plugin.Main(func(p *plugin.Plugin) error {
		log.Println("Registering commands")
		// Register AddTodo command.
		p.HandleCommand(&plugin.CommandOptions{Name: "TickiTodoAdd"}, func(v *nvim.Nvim, args []string) (string, error) {
			log.Println("TodoAdd called with args:", args)
			return todoManager.AddTodo(args), nil
		})

		// Register ListTodos command.
		p.HandleFunction(&plugin.FunctionOptions{Name: "TickiTodoList"}, func(v *nvim.Nvim, args []string) ([]string, error) {
			log.Println("TodoList called")
			return todoManager.ListTodos(), nil
		})

		// Register RemoveTodo command.
		p.HandleCommand(&plugin.CommandOptions{Name: "TickiTodoRemove"}, func(v *nvim.Nvim, args []string) (string, error) {
			log.Println("TodoRemove called with args:", args)
			if len(args) == 0 {
				return "Error: Index required to remove TODO!", nil
			}
			return todoManager.RemoveTodo(parseIndex(args[0])), nil
		})

		return nil
	})
}

// Helper function to parse index safely.
func parseIndex(arg string) int {
	var index int
	fmt.Sscanf(arg, "%d", &index)
	return index
}
