package main

import (
    "fmt"
    "sync"
)

func worker(id int, wg *sync.WaitGroup) {
    defer wg.Done() // Mark this goroutine as done
    fmt.Printf("Worker %d starting\n", id)
    // Simulate work
    fmt.Printf("Worker %d done\n", id)
}

func main() {
    var wg sync.WaitGroup

    for i := 1; i <= 3; i++ {
        wg.Add(1)
        go worker(i, &wg)
    }

    wg.Wait() // Wait for all goroutines to finish
    fmt.Println("All workers finished")
}
