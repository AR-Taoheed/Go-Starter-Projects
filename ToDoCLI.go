package main

import (
	"encoding/json"
	"fmt"
	"os"
	"flag"
)

type Task struct{
	ID        int     `json:"id"`
	Title     string   `json:"title"`
	Completed bool     `json: "completed"`
}

var tasks []Task
var nextID=1

func addTask(title string) {
	task := Task{
		ID: nextID,
		Title: title,
		Completed: false,
	}
	tasks = append(tasks, task)
	nextID++
	fmt.Printf("Added task: %s\n", title)
}

func listTasks() {
	if len(tasks) == 0{
		fmt.Println("No tasks found.")
		return
	}
	for _, task := range tasks {
		status := "Incomplete"
		if task.Completed{
			status = "Completed"
		}
		fmt.Printf("%d. %s [%s]\n", task.ID, task.Title, status)
	}
}

func completeTask(id int) {
	for i, task := range tasks{
		if task.ID == id {
			tasks[i].Completed = true
			fmt.Printf("Task %d marked as completed.\n", id)
			return
		}
	}
	fmt.Printf("Task with ID %d not found.\n", id)
}

func deleteTask(id int) {
	for i, task := range tasks{
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			fmt.Printf("Task %d deleted.\n", id)
			return
		}
	}
	fmt.Printf("Task with ID %d not found.\n", id)
}

func saveTasks() error{
	file, err := os.Create("tasks.json")
	if err != nil {
		return err
	}
	defer file.Close()

	encode := json.NewEncoder(file)
	return encode.Encode(tasks)
}

func loadTasks() error{
	file, err := os.Open("tasks.json")
	if err != nil {
		if os.IsNotExist(err){
			return nil
		}
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	return decoder.Decode(&tasks)
}

func help(){
	fmt.Println(`To-Do List CLI Tool
	A simple command-line tool to manage your to-do list

	USAGE:
	todo-cli [command]

	COMMANDS:
	  -add <task>      Add a new task.
	  -list            List all tasks.
	  -complete <id>   Mark a task as completed.
	  -delete <id>     Delete a task.
	  -help            Display this help message
	fmt.Println`)

}

func main() {

	flag.Usage = func(){
		help()
	}
	err := loadTasks()
	if err != nil{
		fmt.Println("Error loading task", err)
	}
	add := flag.String("add", "", "Add a new task")
	list := flag.Bool("list", false, "List all tasks")
	complete := flag.Int("complete", 0, "Mark a task as completed")
	delete := flag.Int("delete", 0, "Delete a task")

	flag.Parse()

	switch{
	case *add != "":addTask(*add)
	case *list: listTasks()
	case *complete > 0: completeTask(*complete)
	case *delete > 0: deleteTask(*delete)
	default: flag.Usage(); os.Exit(1)
	}

	err = saveTasks()
	if err != nil {
		fmt.Println("Error saving tasks: ", err)
	}
}