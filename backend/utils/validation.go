package utils

import (
	"fmt"
)

// ValidateDependencies validates dependencies
func ValidateDependencies(depends []string, currentTaskUUID string) error {
	if len(depends) == 0 {
		return nil
	}

	// check for self-dependency
	for _, dep := range depends {
		if dep == currentTaskUUID {
			return fmt.Errorf("task cannot depend on itself: %s", dep)
		}
	}

	return nil
}

// Task represents a minimal task structure for dependency validation
type TaskDependency struct {
	UUID    string   `json:"uuid"`
	Depends []string `json:"depends"`
	Status  string   `json:"status"`
}

func ValidateCircularDependencies(depends []string, currentTaskUUID string, existingTasks []TaskDependency) error {
	if len(depends) == 0 {
		return nil
	}

	for _, dep := range depends {
		if dep == currentTaskUUID {
			return fmt.Errorf("task cannot depend on itself: %s", dep)
		}
	}

	dependencyGraph := make(map[string][]string)
	for _, task := range existingTasks {
		if task.Status == "pending" { // Only consider pending tasks
			dependencyGraph[task.UUID] = task.Depends
		}
	}

	dependencyGraph[currentTaskUUID] = depends

	if hasCycle := detectCycle(dependencyGraph, currentTaskUUID); hasCycle {
		return fmt.Errorf("circular dependency detected: adding these dependencies would create a cycle")
	}

	return nil
}

// White (0): unvisited, Gray (1): visiting, Black (2): visited
func detectCycle(graph map[string][]string, startNode string) bool {
	color := make(map[string]int)

	for node := range graph {
		color[node] = 0
	}

	return dfsHasCycle(graph, startNode, color)
}

func dfsHasCycle(graph map[string][]string, node string, color map[string]int) bool {
	color[node] = 1

	for _, dep := range graph[node] {
		if color[dep] == 1 {
			return true
		}
		if color[dep] == 0 && dfsHasCycle(graph, dep, color) {
			return true
		}
	}

	color[node] = 2
	return false
}
