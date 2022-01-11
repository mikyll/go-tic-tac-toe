package main

import (
	"fmt"
	"math/rand"
	"os"
	"syscall"
	"time"

	"github.com/manifoldco/promptui"
	"golang.org/x/term"
)

type menu struct {
	Entry string
}

type board struct {
	X, Y                                                 int
	X1Y1, X1Y2, X1Y3, X2Y1, X2Y2, X2Y3, X3Y1, X3Y2, X3Y3 string
}

type bellSkipper struct{}

// Write implements an io.WriterCloser over os.Stderr, but it skips the terminal
// bell character.
func (bs *bellSkipper) Write(b []byte) (int, error) {
	const charBell = 7 // c.f. readline.CharBell
	if len(b) == 1 && b[0] == charBell {
		return 0, nil
	}
	return os.Stderr.Write(b)
}

// Close implements an io.WriterCloser over os.Stderr.
func (bs *bellSkipper) Close() error {
	return os.Stderr.Close()
}

func remove(slice []board, i int) []board {
	return append(slice[:i], slice[i+1:]...)
}

func selectMove(gameBoard []board, updatedBoard board, player *string, playerTurn *int, playerMove int, moveHistory []int) ([]board, []board, []int) {
	newBoard := gameBoard
	choices := gameBoard

	moveHistory[*playerTurn] = playerMove
	*playerTurn++

	if *player == "X" {
		*player = "O"
	} else if *player == "O" {
		*player = "X"
	}

	for i := 0; i < len(gameBoard); i++ {
		newBoard[i] = updatedBoard
		choices[i] = updatedBoard
	}

	choices[0].X1Y1 = *player
	choices[0].X = 1
	choices[0].Y = 1
	choices[1].X1Y2 = *player
	choices[1].X = 1
	choices[1].Y = 2
	choices[2].X1Y3 = *player
	choices[2].X = 1
	choices[2].Y = 3
	choices[3].X2Y1 = *player
	choices[3].X = 2
	choices[3].Y = 1
	choices[4].X2Y2 = *player
	choices[4].X = 2
	choices[4].Y = 2
	choices[5].X2Y3 = *player
	choices[5].X = 2
	choices[5].Y = 3
	choices[6].X3Y1 = *player
	choices[6].X = 3
	choices[6].Y = 1
	choices[7].X3Y2 = *player
	choices[7].X = 3
	choices[7].Y = 2
	choices[8].X3Y3 = *player
	choices[8].X = 3
	choices[8].Y = 3

	// remove choices already picked
	for i := 0; moveHistory[i] != -1; i++ {
		choices = remove(choices, moveHistory[i])
	}

	return newBoard, choices, moveHistory
}

func checkWin(gameBoard board) string {
	// rows
	if gameBoard.X1Y1 == gameBoard.X1Y2 && gameBoard.X1Y1 == gameBoard.X1Y3 {
		return gameBoard.X1Y1
	}
	if gameBoard.X2Y1 == gameBoard.X2Y2 && gameBoard.X2Y1 == gameBoard.X2Y3 {
		return gameBoard.X2Y1
	}
	if gameBoard.X3Y1 == gameBoard.X3Y2 && gameBoard.X3Y1 == gameBoard.X3Y3 {
		return gameBoard.X3Y1
	}
	// columns
	if gameBoard.X1Y1 == gameBoard.X2Y1 && gameBoard.X1Y1 == gameBoard.X3Y1 {
		return gameBoard.X1Y1
	}
	if gameBoard.X1Y2 == gameBoard.X2Y2 && gameBoard.X1Y2 == gameBoard.X3Y2 {
		return gameBoard.X1Y2
	}
	if gameBoard.X1Y3 == gameBoard.X2Y3 && gameBoard.X1Y3 == gameBoard.X3Y3 {
		return gameBoard.X1Y3
	}
	// others
	if gameBoard.X1Y1 == gameBoard.X2Y2 && gameBoard.X1Y1 == gameBoard.X3Y3 {
		return gameBoard.X1Y1
	}
	if gameBoard.X3Y1 == gameBoard.X2Y2 && gameBoard.X3Y1 == gameBoard.X1Y3 {
		return gameBoard.X3Y1
	}
	return ""
}

func printBoard(b board) {
	fmt.Printf(`
 %s | %s | %s
---+---+---
 %s | %s | %s
---+---+---
 %s | %s | %s
 `, b.X1Y1, b.X1Y2, b.X1Y3, b.X2Y1, b.X2Y2, b.X2Y3, b.X3Y1, b.X3Y2, b.X3Y3)
}

func main() {
	if !term.IsTerminal(int(syscall.Stdin)) {
		fmt.Println("Terminal is not interactive! Consider using flags or environment variables!")
		return
	}

	rand.Seed(time.Now().UnixMilli())

	var v int
	var err error
	//var gameEnd bool = false

	// main menu
	mainMenu := []menu{
		{Entry: "Single Player"},
		{Entry: "Multiplayer"},
		{Entry: "About"},
		{Entry: "Quit"},
	}

	/*singlePlayerMenu := []menu{
		{Entry: "Easy"},
		{Entry: "Hard"},
		{Entry: "About"},
		{Entry: "Back"},
	}*/

	mainMenuTemplate := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "\U000027A4  {{ .Entry | cyan }}",
		Inactive: "  {{ .Entry | white }} ",
		Selected: "\U000027A4 {{ .Entry | white }}",
	}

	mainMenuPrompt := promptui.Select{
		Label:     "-------- Main Menu ---------",
		Items:     mainMenu,
		Templates: mainMenuTemplate,
		Size:      4,
		Stdout:    &bellSkipper{},
	}

	v = -1
	for i := 0; v == -1; i++ {
		v, _, err = mainMenuPrompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}
		switch v {
		case 1:
			v = -1
		case 2:
			v = -1
		case 3:
			return
		}
	}

	// game
	turnCounter := 0
	choicesHistory := []int{-1, -1, -1, -1, -1, -1, -1, -1, -1}
	player := "X"
	gameBoard := []board{
		{X: 1, Y: 1, X1Y1: " ", X1Y2: " ", X1Y3: " ", X2Y1: " ", X2Y2: " ", X2Y3: " ", X3Y1: " ", X3Y2: " ", X3Y3: " "},
		{X: 1, Y: 2, X1Y1: " ", X1Y2: " ", X1Y3: " ", X2Y1: " ", X2Y2: " ", X2Y3: " ", X3Y1: " ", X3Y2: " ", X3Y3: " "},
		{X: 1, Y: 3, X1Y1: " ", X1Y2: " ", X1Y3: " ", X2Y1: " ", X2Y2: " ", X2Y3: " ", X3Y1: " ", X3Y2: " ", X3Y3: " "},
		{X: 2, Y: 1, X1Y1: " ", X1Y2: " ", X1Y3: " ", X2Y1: " ", X2Y2: " ", X2Y3: " ", X3Y1: " ", X3Y2: " ", X3Y3: " "},
		{X: 2, Y: 2, X1Y1: " ", X1Y2: " ", X1Y3: " ", X2Y1: " ", X2Y2: " ", X2Y3: " ", X3Y1: " ", X3Y2: " ", X3Y3: " "},
		{X: 2, Y: 3, X1Y1: " ", X1Y2: " ", X1Y3: " ", X2Y1: " ", X2Y2: " ", X2Y3: " ", X3Y1: " ", X3Y2: " ", X3Y3: " "},
		{X: 3, Y: 1, X1Y1: " ", X1Y2: " ", X1Y3: " ", X2Y1: " ", X2Y2: " ", X2Y3: " ", X3Y1: " ", X3Y2: " ", X3Y3: " "},
		{X: 3, Y: 2, X1Y1: " ", X1Y2: " ", X1Y3: " ", X2Y1: " ", X2Y2: " ", X2Y3: " ", X3Y1: " ", X3Y2: " ", X3Y3: " "},
		{X: 3, Y: 3, X1Y1: " ", X1Y2: " ", X1Y3: " ", X2Y1: " ", X2Y2: " ", X2Y3: " ", X3Y1: " ", X3Y2: " ", X3Y3: player},
	}
	playerChoices := []board{
		{X: 1, Y: 1, X1Y1: player, X1Y2: " ", X1Y3: " ", X2Y1: " ", X2Y2: " ", X2Y3: " ", X3Y1: " ", X3Y2: " ", X3Y3: " "},
		{X: 1, Y: 2, X1Y1: " ", X1Y2: player, X1Y3: " ", X2Y1: " ", X2Y2: " ", X2Y3: " ", X3Y1: " ", X3Y2: " ", X3Y3: " "},
		{X: 1, Y: 3, X1Y1: " ", X1Y2: " ", X1Y3: player, X2Y1: " ", X2Y2: " ", X2Y3: " ", X3Y1: " ", X3Y2: " ", X3Y3: " "},
		{X: 2, Y: 1, X1Y1: " ", X1Y2: " ", X1Y3: " ", X2Y1: player, X2Y2: " ", X2Y3: " ", X3Y1: " ", X3Y2: " ", X3Y3: " "},
		{X: 2, Y: 2, X1Y1: " ", X1Y2: " ", X1Y3: " ", X2Y1: " ", X2Y2: player, X2Y3: " ", X3Y1: " ", X3Y2: " ", X3Y3: " "},
		{X: 2, Y: 3, X1Y1: " ", X1Y2: " ", X1Y3: " ", X2Y1: " ", X2Y2: " ", X2Y3: player, X3Y1: " ", X3Y2: " ", X3Y3: " "},
		{X: 3, Y: 1, X1Y1: " ", X1Y2: " ", X1Y3: " ", X2Y1: " ", X2Y2: " ", X2Y3: " ", X3Y1: player, X3Y2: " ", X3Y3: " "},
		{X: 3, Y: 2, X1Y1: " ", X1Y2: " ", X1Y3: " ", X2Y1: " ", X2Y2: " ", X2Y3: " ", X3Y1: " ", X3Y2: player, X3Y3: " "},
		{X: 3, Y: 3, X1Y1: " ", X1Y2: " ", X1Y3: " ", X2Y1: " ", X2Y2: " ", X2Y3: " ", X3Y1: " ", X3Y2: " ", X3Y3: player},
	}

	gameTemplate := &promptui.SelectTemplates{
		Label:    fmt.Sprintf("You play as %s. Choose your next move.", player),
		Active:   "\U000027A4  ({{ .X | cyan }}, {{ .Y | green }})",
		Inactive: "  ({{ .X | cyan }}, {{ .Y | green }})",
		Selected: "\U000027A4 ({{ .X | cyan }}, {{ .Y | green }})",
		Details: fmt.Sprintf(`
----------- Game -----------

         {{ .X1Y1 }} | {{ .X1Y2 }} | {{ .X1Y3 }}  {{"1" | green}}
        ---+---+---
         {{ .X2Y1 }} | {{ .X2Y2 }} | {{ .X2Y3 }}  {{"2" | green}}
        ---+---+---
         {{ .X3Y1 }} | {{ .X3Y2 }} | {{ .X3Y3 }}  {{"3" | green}}
         {{"1   2   3" | cyan}}

----------------------------
Selected Move: %s in ({{ .X | cyan }}, {{ .Y | green }})`, player),
	}

	opponentChoice := -1
	// roll for who begins
	for {
		gamePrompt := promptui.Select{
			Label:     "",
			Items:     playerChoices,
			Templates: gameTemplate,
			Size:      4,
			Stdout:    &bellSkipper{},
		}

		v, _, err = gamePrompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}
		fmt.Printf("Chosen option %d\n", v)

		printBoard(playerChoices[0])
		gameBoard, playerChoices, choicesHistory = selectMove(gameBoard, playerChoices[v], &player, &turnCounter, v, choicesHistory)
		printBoard(gameBoard[0])
		printBoard(gameBoard[1])
		printBoard(gameBoard[2])
		printBoard(playerChoices[0])
		if turnCounter > 3 {
			win := checkWin(gameBoard[0])
			if win != "" {
				printBoard(gameBoard[0])
				fmt.Printf("Player %s won.\n\n", win)
				return
			} else if turnCounter == 8 {
				fmt.Printf("Draw.\n\n")
				return
			}
		}

		// opponent random move
		opponentChoice = rand.Intn(len(playerChoices))
		gameBoard, playerChoices, choicesHistory = selectMove(gameBoard, playerChoices[opponentChoice], &player, &turnCounter, opponentChoice, choicesHistory)
		if turnCounter > 3 {
			win := checkWin(gameBoard[0])
			if win != "" {
				printBoard(gameBoard[0])
				fmt.Printf("Player %s won.\n\n", win)
				return
			} else if turnCounter == 8 {
				fmt.Printf("Draw.\n\n")
				return
			}
		}
	}
	gamePrompt := promptui.Select{
		Label:     "",
		Items:     playerChoices,
		Templates: gameTemplate,
		Size:      4,
		Stdout:    &bellSkipper{},
	}

	v, _, err = gamePrompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}
	fmt.Printf("Chosen option %d\n", v)

	gameBoard, playerChoices, choicesHistory = selectMove(gameBoard, playerChoices[v], &player, &turnCounter, v, choicesHistory)

	// opponent random move
	tmp := rand.Intn(len(playerChoices))
	gameBoard, playerChoices, choicesHistory = selectMove(gameBoard, playerChoices[tmp], &player, &turnCounter, tmp, choicesHistory)

	gamePrompt = promptui.Select{
		Label:     "",
		Items:     playerChoices,
		Templates: gameTemplate,
		Size:      4,
		Stdout:    &bellSkipper{},
	}
	v, _, err = gamePrompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}
	fmt.Printf("Chosen option %d\n", v)

	// NB: fix the selectMove(): we cannot use the index in choices, since they're not anymore in range [0, 8] but less
	gameBoard, playerChoices, choicesHistory = selectMove(gameBoard, playerChoices[v], &player, &turnCounter, v, choicesHistory)

	// opponent random move
	tmp = rand.Intn(len(playerChoices))
	gameBoard, playerChoices, choicesHistory = selectMove(gameBoard, playerChoices[tmp], &player, &turnCounter, tmp, choicesHistory)

	gamePrompt = promptui.Select{
		Label:     "",
		Items:     playerChoices,
		Templates: gameTemplate,
		Size:      4,
		Stdout:    &bellSkipper{},
	}
	v, _, err = gamePrompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}
	fmt.Printf("Chosen option %d\n", v)
}
