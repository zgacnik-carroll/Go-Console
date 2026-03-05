package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

// Tile struct represents a color tile
type Tile struct {
	Color      string
	CorrectPos int
	IsBlank    bool
}

// ANSI color codes for backgrounds
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

const resetCode = "\033[0m"

// Create a 3x3 board of colored tiles
func createTiles() []*Tile {
	tiles := []*Tile{}
	for i := 0; i < 8; i++ {
		tiles = append(tiles, &Tile{
			Color:      colorCodes[i%len(colorCodes)],
			CorrectPos: i,
			IsBlank:    false,
		})
	}
	tiles = append(tiles, &Tile{
		Color:      "",
		CorrectPos: 8,
		IsBlank:    true,
	})
	return tiles
}

// Convert 1D slice of tiles to 2D board
func createBoard(tiles []*Tile) [][]*Tile {
	board := make([][]*Tile, 3)
	index := 0
	for i := 0; i < 3; i++ {
		board[i] = make([]*Tile, 3)
		for j := 0; j < 3; j++ {
			board[i][j] = tiles[index]
			index++
		}
	}
	return board
}

// Count inversions in a 1D tile slice (ignoring blank)
func countInversions(tiles []*Tile) int {
	count := 0
	for i := 0; i < len(tiles); i++ {
		if tiles[i].IsBlank {
			continue
		}
		for j := i + 1; j < len(tiles); j++ {
			if tiles[j].IsBlank {
				continue
			}
			if tiles[i].CorrectPos > tiles[j].CorrectPos {
				count++
			}
		}
	}
	return count
}

// Shuffle tiles randomly and ensure solvable
func shuffleSolvable(tiles []*Tile) {
	rand.Seed(time.Now().UnixNano())
	for {
		// Fisher–Yates shuffle
		for i := len(tiles) - 1; i > 0; i-- {
			j := rand.Intn(i + 1)
			tiles[i], tiles[j] = tiles[j], tiles[i]
		}
		if countInversions(tiles)%2 == 0 {
			break // solvable
		}
	}
	// Update CorrectPos after shuffle
	for i, t := range tiles {
		t.CorrectPos = i
	}
}

// Shuffle board with random legal moves (optional, more natural movement)
func shuffleBoard(board [][]*Tile, moves int) {
	dirs := []string{"w", "a", "s", "d"}
	for i := 0; i < moves; i++ {
		move(board, dirs[rand.Intn(4)])
	}
}

// Print board with mini goal display on the left
func printBoard(current, goal [][]*Tile) {
	fmt.Println()
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			tile := goal[i][j]
			if tile.IsBlank {
				fmt.Print("   ")
			} else {
				fmt.Print(tile.Color + "   " + resetCode)
			}
		}
		fmt.Print("    ")
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

// Find blank tile coordinates
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

// Move blank tile
func move(board [][]*Tile, dir string) bool {
	row, col := findBlank(board)
	newRow, newCol := row, col

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

	if newRow < 0 || newRow >= 3 || newCol < 0 || newCol >= 3 {
		return false
	}

	board[row][col], board[newRow][newCol] = board[newRow][newCol], board[row][col]
	return true
}

// Check win condition
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

func playGame(reader *bufio.Reader) {
	tiles := createTiles()
	shuffleSolvable(tiles) // random and solvable
	goalBoard := createBoard(tiles)

	// Current board is a copy of goal, then shuffled slightly for realism
	currentBoard := createBoard(tiles)
	shuffleBoard(currentBoard, 50)

	for {
		printBoard(currentBoard, goalBoard)

		if checkWin(currentBoard) {
			fmt.Println("Puzzle Solved!\n")
			break
		}

		fmt.Print("Move (or q to quit): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "q" {
			fmt.Println("\nYou quit the level. Here’s the solved puzzle:")
			printBoard(goalBoard, goalBoard)
			break
		}

		if !move(currentBoard, input) {
			fmt.Println("Invalid input! Use w/a/s/d and stay within bounds.")
		}
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\nColor Sliding Puzzle!!!\n")
	fmt.Println("Left: finished display goal")
	fmt.Println("Right: current board\n")
	fmt.Println("Use w/a/s/d to move the blank tile. Press q to quit.")

	for {
		playGame(reader)

		// Ask user if they want to play again, with input validation
		for {
			fmt.Print("Play again? (y/n): ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(strings.ToLower(input))

			if input == "y" {
				fmt.Println()
				break // start a new game
			} else if input == "n" {
				fmt.Println("\nThanks for playing! Goodbye!")
				return
			} else {
				fmt.Println("Invalid input! Please type 'y' to play again or 'n' to quit.")
			}
		}
	}
}
