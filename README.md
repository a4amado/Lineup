# Lineup

A lightweight, thread-safe queue management library for Go that allows you to control concurrent processing with automatic position tracking.

## Features

- **Thread-safe queue management** - Built with mutex locks for safe concurrent access
- **Concurrent processing control** - Limit the number of items processed simultaneously
- **Automatic position tracking (O(N))** - Items automatically receive their position in the queue
- **Queue item purging** - Remove items from the queue at any time
- **Simple API** - Easy to use with minimal configuration

## Installation

```bash
go get github.com/a4amado/Lineup
```

## Usage

### Basic Example

```go
package main

import (
    "fmt"
    "time"
    "github.com/a4amado/Lineup"
)

func main() {
    // Create a new queue that allows 3 items to be processed concurrently
    queue := Lineup.New(Lineup.QueueOptions{
        MaxProcessing: 3,
    })

    // Place items in the queue
    for i := 0; i < 10; i++ {
        item := queue.Place()
        
        go func(id int, item *Lineup.QueueItem) {
            defer item.Purge() // Ensure cleanup when done
            
            // Wait for your turn (position 1 means you can proceed)
            position := <-item.Halt
            
            if position == 1 {
                fmt.Printf("Item %d is now processing\n", id)
                // Do your work here
                time.Sleep(2 * time.Second)
                fmt.Printf("Item %d finished processing\n", id)
            } else {
                fmt.Printf("Item %d is waiting at position %d\n", id, position)
            }
        }(i, item)
    }

    // Keep the program running
    time.Sleep(10 * time.Second)
}
```


### Purging Items

You can remove items from the queue before they are processed. It's recommended to use `defer` to ensure cleanup:

```go
queue := Lineup.New(Lineup.QueueOptions{
    MaxProcessing: 2,
})

item1 := queue.Place()
defer item1.Purge() // Clean up when done

item2 := queue.Place()
defer item2.Purge()

item3 := queue.Place()
defer item3.Purge()

// Or manually purge an item before it's processed
item2.Purge() // item1 and item3 will proceed, item2 will be skipped
```

## API Reference

All types and functions are available from the `Lineup` package: `github.com/a4amado/Lineup`

### `New(opts QueueOptions) *Queue`

Creates a new queue instance.

**Parameters:**
- `opts QueueOptions` - Configuration options for the queue
  - `MaxProcessing int` - Maximum number of items that can be processed concurrently

**Returns:** A pointer to a new `Queue` instance

**Example:**
```go
queue := Lineup.New(Lineup.QueueOptions{
    MaxProcessing: 3,
})
```

### `Queue.Place() *QueueItem`

Adds a new item to the queue and returns a `QueueItem`.

**Returns:** A pointer to a `QueueItem` that can be used to track position and purge the item

**Example:**
```go
item := queue.Place()
```

### `QueueItem.Purge()`

Removes the item from the queue. The item will be skipped during processing.

**Example:**
```go
item.Purge()
```

### `QueueItem.Halt chan int`

A channel that receives the item's position in the queue:
- `1` - The item can proceed (within the processing limit)
- `> 1` - The item is waiting (position in line)

The queue automatically recalibrates positions every 300ms, so items will receive updated positions as the queue progresses.

**Example:**
```go
item := queue.Place()
defer item.Purge() // Clean up when done

position := <-item.Halt
if position == 1 {
    // Process the item
}
```

## How It Works

1. When you call `Place()`, a new item is added to the queue
2. A background goroutine continuously recalibrates item positions every 300ms
3. Items within the `MaxProcessing` limit receive position `1` (can proceed)
4. Items beyond the limit receive their actual position in line
5. Items can be purged at any time and will be removed from the queue

## License

See [LICENSE.txt](LICENSE.txt) for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

