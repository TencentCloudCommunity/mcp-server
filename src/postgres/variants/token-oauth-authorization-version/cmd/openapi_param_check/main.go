package main

import (
	"fmt"
	"os"

	"postgres_server/tools"
)

func main() {
	fmt.Println("========== OPENAPI PARAM ALIGNMENT CHECK ==========")
	results, err := tools.ValidateAllOpenAPIArgumentAlignment()
	passed := 0
	failed := 0
	for _, result := range results {
		if result.OK {
			fmt.Printf("✓ %s\n", result.Name)
			passed++
			continue
		}
		fmt.Printf("✗ %s: %s\n", result.Name, result.Error)
		failed++
	}
	fmt.Println()
	fmt.Printf("Summary: passed=%d failed=%d total=%d\n", passed, failed, len(results))
	if err != nil {
		fmt.Println("Result:", err)
		os.Exit(1)
	}
}
