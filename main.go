package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Tile struct {
	Lines        []string
	CorrectIndex int
}

func createTiles() []*Tile {

	image := []string{
		"#########",
		"#       #",
		"#       #",
		"#       #",
		"#       #",
		"#       #",
		"#       #",
		"#       #",
		"#########",
	}

	tileSize := 3
	tiles := []*Tile{}
	index := 0

	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {

			var lines []string

			for i := 0; i < tileSize; i++ {

				line := image[row*tileSize+i]
				start := col * tileSize
				end := start + tileSize

				lines = append(lines, line[start:end])
			}

			tile := &Tile{
				Lines:        lines,
				CorrectIndex: index,
			}

			tiles = append(tiles, tile)
			index++
		}
	}

	// Make last tile blank
	tiles[8] = &Tile{
		Lines:        nil,
		CorrectIndex: 8,
	}

	return tiles
}

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

	// Swap last two tiles to create puzzle state
	board[2][1], board[2][2] =
		board[2][2], board[2][1]

	return board
}

func printBoard(board [][]*Tile) {
	fmt.Println()

	for row := 0; row < 3; row++ {

		for line := 0; line < 3; line++ {

			for col := 0; col < 3; col++ {

				tile := board[row][col]

				if tile.Lines == nil {
					fmt.Print("    ") // 4 spaces for blank tile width
				} else {
					fmt.Print(tile.Lines[line])
				}
			}

			fmt.Println()
		}
	}

	fmt.Println()
}

func checkWin(board [][]*Tile) bool {
	expected := 0

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {

			if board[i][j].CorrectIndex != expected {
				return false
			}
			expected++
		}
	}

	return true
}

func findBlank(board [][]*Tile) (int, int) {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if board[i][j].Lines == nil {
				return i, j
			}
		}
	}
	return -1, -1
}

func move(board [][]*Tile, dir string) {

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
		fmt.Println("Invalid input. Please type w, a, s, or d.")
		return
	}

	// Boundary check
	if newRow < 0 || newRow > 2 || newCol < 0 || newCol > 2 {
		return
	}

	// Swap POINTERS
	board[row][col], board[newRow][newCol] =
		board[newRow][newCol], board[row][col]
}

func main() {

	reader := bufio.NewReader(os.Stdin)

	tiles := createTiles()
	board := createBoard(tiles)

	fmt.Println("ASCII Sliding Puzzle")
	fmt.Println("Use w/a/s/d to move bottom-right tile. q to quit.")

	for {

		printBoard(board)

		if checkWin(board) {
			fmt.Println("Puzzle Solved!")
			return
		}

		fmt.Print("Move: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "q" {
			return
		}

		move(board, input)
	}
}
