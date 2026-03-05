package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

// Package main implements a console-based color sliding puzzle game.
// The game displays a solved puzzle on the left and a scrambled puzzle on the right.
// The player moves the blank tile using w/a/s/d to match the solved puzzle.

// Tile represents a single tile in the sliding puzzle.
type Tile struct {
	Color      string // ANSI color code for tile background
	CorrectPos int    // Correct position index in solved puzzle
	IsBlank    bool   // Indicates if the tile is the blank tile
}

// colorCodes contains ANSI background color codes for tiles.
var colorCodes = []string{
	"\033[41m",  // Red
	"\033[42m",  // Green
	"\033[44m",  // Blue
	"\033[45m",  // Magenta
	"\033[46m",  // Cyan
	"\033[43m",  // Yellow
	"\033[47m",  // White
	"\033[101m", // Bright Red
}

const resetCode = "\033[0m" // ANSI reset code

// createTiles creates a slice of 3x3 Tile pointers including one blank tile.
// Returns a 1D slice of 9 tiles.
func createTiles() []*Tile {
	tiles := []*Tile{}

	// Create 8 colored tiles
	for i := 0; i < 8; i++ {
		tiles = append(tiles, &Tile{
			Color:      colorCodes[i%len(colorCodes)],
			CorrectPos: i,
			IsBlank:    false,
		})
	}

	// Add blank tile as the 9th tile
	tiles = append(tiles, &Tile{
		Color:      "",
		CorrectPos: 8,
		IsBlank:    true,
	})
	return tiles
}

// createBoard converts a 1D slice of tiles into a 2D 3x3 board.
// tiles: 1D slice of 9 tiles.
// Returns a 2D slice [3][3] of Tile pointers.
func createBoard(tiles []*Tile) [][]*Tile {
	board := make([][]*Tile, 3)
	index := 0

	// Fill the 2D board row by row
	for i := 0; i < 3; i++ {
		board[i] = make([]*Tile, 3)
		for j := 0; j < 3; j++ {
			board[i][j] = tiles[index]
			index++
		}
	}
	return board
}

// countInversions counts the number of inversions in a 1D slice of tiles (ignores blank).
// tiles: 1D slice of tiles.
// Returns the total number of inversions.
func countInversions(tiles []*Tile) int {
	count := 0
	for i := 0; i < len(tiles); i++ {
		if tiles[i].IsBlank {
			continue // skip blank tile
		}
		for j := i + 1; j < len(tiles); j++ {
			if tiles[j].IsBlank {
				continue // skip blank tile
			}
			// Count inversion if earlier tile has higher correct position than later tile
			if tiles[i].CorrectPos > tiles[j].CorrectPos {
				count++
			}
		}
	}
	return count
}

// shuffleSolvable shuffles tiles randomly until a solvable 3x3 puzzle is created.
// tiles: slice of Tile pointers (modified in place).
func shuffleSolvable(tiles []*Tile) {
	rand.Seed(time.Now().UnixNano())

	for {
		// Fisher-Yates shuffle: iterate backwards and swap each element with a random one before it
		for i := len(tiles) - 1; i > 0; i-- {
			j := rand.Intn(i + 1)
			tiles[i], tiles[j] = tiles[j], tiles[i]
		}

		// Puzzle is solvable if inversion count is even
		if countInversions(tiles)%2 == 0 {
			break
		}
	}

	// Update CorrectPos to match shuffled order
	for i, t := range tiles {
		t.CorrectPos = i
	}
}

// shuffleBoard performs random legal moves on the blank tile to simulate scrambling.
// board: 2D slice of tiles (modified in place)
// moves: number of random moves to perform
func shuffleBoard(board [][]*Tile, moves int) {
	dirs := []string{"w", "a", "s", "d"} // possible move directions
	for i := 0; i < moves; i++ {
		move(board, dirs[rand.Intn(4)]) // perform random move
	}
}

// printBoard prints the goal board (left) and current board (right) side by side.
// current: the player's current board
// goal: the solved puzzle board
func printBoard(current, goal [][]*Tile) {
	fmt.Println()
	for i := 0; i < 3; i++ {
		// Print goal board row
		for j := 0; j < 3; j++ {
			tile := goal[i][j]
			if tile.IsBlank {
				fmt.Print("   ") // truly blank tile
			} else {
				fmt.Print(tile.Color + "   " + resetCode) // colored tile
			}
		}

		fmt.Print("    ") // gap between goal and current

		// Print current board row
		for j := 0; j < 3; j++ {
			tile := current[i][j]
			if tile.IsBlank {
				fmt.Print("   ")
			} else {
				fmt.Print(tile.Color + "   " + resetCode)
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

// findBlank returns the row and column of the blank tile in a 3x3 board.
// board: 2D slice of Tile pointers
// Returns row and column of blank tile, or (-1,-1) if not found
func findBlank(board [][]*Tile) (int, int) {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if board[i][j].IsBlank {
				return i, j
			}
		}
	}
	return -1, -1
}

// move attempts to move the blank tile in a specified direction.
// board: current 3x3 board
// dir: "w"=up, "s"=down, "a"=left, "d"=right
// Returns true if move was successful, false otherwise
func move(board [][]*Tile, dir string) bool {
	row, col := findBlank(board)
	newRow, newCol := row, col

	// Determine new position based on input
	switch dir {
	case "w":
		newRow--
	case "s":
		newRow++
	case "a":
		newCol--
	case "d":
		newCol++
	default:
		return false
	}

	// Check boundaries
	if newRow < 0 || newRow >= 3 || newCol < 0 || newCol >= 3 {
		return false
	}

	// Swap blank with target tile
	board[row][col], board[newRow][newCol] = board[newRow][newCol], board[row][col]
	return true
}

// checkWin checks if the current board matches the solved puzzle.
// Returns true if the board is solved, false otherwise.
func checkWin(board [][]*Tile) bool {
	expected := 0
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if board[i][j].CorrectPos != expected {
				return false
			}
			expected++
		}
	}
	return true
}

// playGame runs a single session of the color sliding puzzle.
// reader: buffered reader to capture user input
func playGame(reader *bufio.Reader) {
	tiles := createTiles()
	shuffleSolvable(tiles)          // generate random solvable goal
	goalBoard := createBoard(tiles) // goal board

	// Start current board as goal board and shuffle for realism
	currentBoard := createBoard(tiles)
	shuffleBoard(currentBoard, 50)

	for {
		// Display current and goal boards
		printBoard(currentBoard, goalBoard)

		// Check if player solved the puzzle
		if checkWin(currentBoard) {
			fmt.Println("Puzzle Solved!\n")
			break
		}

		// Prompt for user input
		fmt.Print("Move (or q to quit): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		// Handle quit: reveal solution
		if input == "q" {
			fmt.Println("\nYou quit the level. Here’s the solved puzzle:")
			printBoard(goalBoard, goalBoard)
			break
		}

		// Attempt to move blank tile; show error if invalid
		if !move(currentBoard, input) {
			fmt.Println("Invalid input! Use w/a/s/d and stay within bounds.")
		}
	}
}

// main initializes the game loop and handles replay logic.
func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\nColor Sliding Puzzle!!!\n")
	fmt.Println("Left: finished display goal")
	fmt.Println("Right: current board\n")
	fmt.Println("Use w/a/s/d to move the blank tile. Press q to quit.")

	for {
		playGame(reader)

		// Play-again input validation loop
		for {
			fmt.Print("Play again? (y/n): ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(strings.ToLower(input))

			if input == "y" {
				break
			} else if input == "n" {
				fmt.Println("\nThanks for playing! Goodbye!")
				return
			} else {
				// Catch any other invalid input
				fmt.Println("Invalid input! Please type 'y' to play again or 'n' to quit.")
			}
		}
	}
}
