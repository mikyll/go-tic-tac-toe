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

// Utility to skip the terminal bell sound, when selecting prompt options
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

type menu struct {
	Entry string
}

type playerSymbols struct {
	P1, P2 string
}

type board struct {
	X, Y                                                 int
	X1Y1, X1Y2, X1Y3, X2Y1, X2Y2, X2Y3, X3Y1, X3Y2, X3Y3 string
}

func updateBoardElement(gameBoard *[9]board, el board) {
	for i := 0; i < len(gameBoard); i++ {
		gameBoard[i] = el
	}
}

func remove(slice []board, i int) []board {
	return append(slice[:i], slice[i+1:]...)
}

func swapPlayer(player string, ps playerSymbols) string {
	if player == ps.P1 {
		return ps.P2
	} else if player == ps.P2 {
		return ps.P1
	}
	return ""
}

func selectMove(gameBoard [9]board, updatedBoard board, player string, playerTurn *int, playerMove int, moveHistory []int) ([9]board, []board, []int) {
	newBoard := gameBoard
	choices := gameBoard[:]

	moveHistory[*playerTurn] = playerMove
	*playerTurn++

	player = swapPlayer(player, ps)

	updateBoardElement(&newBoard, updatedBoard)
	for i := 0; i < len(gameBoard); i++ {
		choices[i] = updatedBoard
	}

	choices[0].X1Y1 = player
	choices[0].X = 1
	choices[0].Y = 1
	choices[1].X1Y2 = player
	choices[1].X = 1
	choices[1].Y = 2
	choices[2].X1Y3 = player
	choices[2].X = 1
	choices[2].Y = 3
	choices[3].X2Y1 = player
	choices[3].X = 2
	choices[3].Y = 1
	choices[4].X2Y2 = player
	choices[4].X = 2
	choices[4].Y = 2
	choices[5].X2Y3 = player
	choices[5].X = 2
	choices[5].Y = 3
	choices[6].X3Y1 = player
	choices[6].X = 3
	choices[6].Y = 1
	choices[7].X3Y2 = player
	choices[7].X = 3
	choices[7].Y = 2
	choices[8].X3Y3 = player
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
	if gameBoard.X1Y1 != "" && gameBoard.X1Y1 != " " && gameBoard.X1Y1 == gameBoard.X1Y2 && gameBoard.X1Y1 == gameBoard.X1Y3 {
		return gameBoard.X1Y1
	}
	if gameBoard.X2Y1 != "" && gameBoard.X2Y1 != " " && gameBoard.X2Y1 == gameBoard.X2Y2 && gameBoard.X2Y1 == gameBoard.X2Y3 {
		return gameBoard.X2Y1
	}
	if gameBoard.X3Y1 != "" && gameBoard.X3Y1 != " " && gameBoard.X3Y1 == gameBoard.X3Y2 && gameBoard.X3Y1 == gameBoard.X3Y3 {
		return gameBoard.X3Y1
	}
	// columns
	if gameBoard.X1Y1 != "" && gameBoard.X1Y1 != " " && gameBoard.X1Y1 == gameBoard.X2Y1 && gameBoard.X1Y1 == gameBoard.X3Y1 {
		return gameBoard.X1Y1
	}
	if gameBoard.X1Y2 != "" && gameBoard.X1Y2 != " " && gameBoard.X1Y2 == gameBoard.X2Y2 && gameBoard.X1Y2 == gameBoard.X3Y2 {
		return gameBoard.X1Y2
	}
	if gameBoard.X1Y3 != "" && gameBoard.X1Y3 != " " && gameBoard.X1Y3 == gameBoard.X2Y3 && gameBoard.X1Y3 == gameBoard.X3Y3 {
		return gameBoard.X1Y3
	}
	// others
	if gameBoard.X1Y1 != "" && gameBoard.X1Y1 != " " && gameBoard.X1Y1 == gameBoard.X2Y2 && gameBoard.X1Y1 == gameBoard.X3Y3 {
		return gameBoard.X1Y1
	}
	if gameBoard.X3Y1 != "" && gameBoard.X3Y1 != " " && gameBoard.X3Y1 == gameBoard.X2Y2 && gameBoard.X3Y1 == gameBoard.X1Y3 {
		return gameBoard.X3Y1
	}
	return ""
}

func checkWinPrint(gameBoard board) string {
	// rows
	if gameBoard.X1Y1 != "" && gameBoard.X1Y1 == gameBoard.X1Y2 && gameBoard.X1Y1 == gameBoard.X1Y3 {
		p := gameBoard.X1Y1
		fmt.Printf("\n\n %v | %v | %v \n---+---+---\n   |   |   \n---+---+---\n   |   |   ", p, p, p)
		return gameBoard.X1Y1
	}
	if gameBoard.X2Y1 != "" && gameBoard.X2Y1 == gameBoard.X2Y2 && gameBoard.X2Y1 == gameBoard.X2Y3 {
		p := gameBoard.X2Y1
		fmt.Printf("\n\n   |   |   \n---+---+---\n %v | %v | %v \n---+---+---\n   |   |   ", p, p, p)
		return gameBoard.X2Y1
	}
	if gameBoard.X3Y1 != "" && gameBoard.X3Y1 == gameBoard.X3Y2 && gameBoard.X3Y1 == gameBoard.X3Y3 {
		p := gameBoard.X3Y1
		fmt.Printf("\n\n   |   |   \n---+---+---\n   |   |   \n---+---+---\n %v | %v | %v ", p, p, p)
		return gameBoard.X3Y1
	}
	// columns
	if gameBoard.X1Y1 != "" && gameBoard.X1Y1 == gameBoard.X2Y1 && gameBoard.X1Y1 == gameBoard.X3Y1 {
		p := gameBoard.X1Y1
		fmt.Printf("\n\n %v |   |   \n---+---+---\n %v |   |   \n---+---+---\n %v |   |   ", p, p, p)
		return gameBoard.X1Y1
	}
	if gameBoard.X1Y2 != "" && gameBoard.X1Y2 == gameBoard.X2Y2 && gameBoard.X1Y2 == gameBoard.X3Y2 {
		p := gameBoard.X1Y2
		fmt.Printf("\n\n   | %v |   \n---+---+---\n   | %v |   \n---+---+---\n   | %v |   ", p, p, p)
		return gameBoard.X1Y2
	}
	if gameBoard.X1Y3 != "" && gameBoard.X1Y3 == gameBoard.X2Y3 && gameBoard.X1Y3 == gameBoard.X3Y3 {
		p := gameBoard.X1Y3
		fmt.Printf("\n\n   |   | %v \n---+---+---\n   |   | %v \n---+---+---\n   |   | %v ", p, p, p)
		return gameBoard.X1Y3
	}
	// others
	if gameBoard.X1Y1 != "" && gameBoard.X1Y1 == gameBoard.X2Y2 && gameBoard.X1Y1 == gameBoard.X3Y3 {
		p := gameBoard.X1Y1
		fmt.Printf("\n\n %v |   |   \n---+---+---\n   | %v |   \n---+---+---\n   |   | %v ", p, p, p)
		return gameBoard.X1Y1
	}
	if gameBoard.X3Y1 != "" && gameBoard.X3Y1 == gameBoard.X2Y2 && gameBoard.X3Y1 == gameBoard.X1Y3 {
		p := gameBoard.X3Y1
		fmt.Printf("\n\n   |   | %v \n---+---+---\n   | %v |   \n---+---+---\n %v |   |   ", p, p, p)
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
func printWholeBoard(b [9]board) {
	fmt.Printf(`
[0]           [1]           [2]           [3]
 %s | %s | %s  #  %s | %s | %s  #  %s | %s | %s  #  %s | %s | %s
---+---+--- # ---+---+--- # ---+---+--- # ---+---+--- # 
 %s | %s | %s  #  %s | %s | %s  #  %s | %s | %s  #  %s | %s | %s
---+---+--- # ---+---+--- # ---+---+--- # ---+---+--- # 
 %s | %s | %s  #  %s | %s | %s  #  %s | %s | %s  #  %s | %s | %s
 `, b[0].X1Y1, b[0].X1Y2, b[0].X1Y3, b[1].X1Y1, b[1].X1Y2, b[1].X1Y3, b[2].X1Y1, b[2].X1Y2, b[2].X1Y3, b[3].X1Y1, b[3].X1Y2, b[3].X1Y3,
		b[0].X2Y1, b[0].X2Y2, b[0].X2Y3, b[1].X2Y1, b[1].X2Y2, b[1].X2Y3, b[2].X2Y1, b[2].X2Y2, b[2].X2Y3, b[3].X2Y1, b[3].X2Y2, b[3].X2Y3,
		b[0].X3Y1, b[0].X3Y2, b[0].X3Y3, b[1].X3Y1, b[1].X3Y2, b[1].X3Y3, b[2].X3Y1, b[2].X3Y2, b[2].X3Y3, b[3].X3Y1, b[3].X3Y2, b[3].X3Y3)
}

var ps = playerSymbols{P1: "X", P2: "O"}

func main() {
	/*// pointer test
	var pointer *int
	var pointed int

	pointer = &pointed
	fmt.Println("pointer: ", pointer)
	fmt.Println("*pointer:", *pointer)
	fmt.Println("&pointer:", &pointer)
	fmt.Println("pointed: ", pointed)
	fmt.Println("&pointed:", &pointed)

	return*/

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
	player := ps.P1
	opponent := ps.P2
	gameBoard := [9]board{
		{X: 1, Y: 1, X1Y1: " ", X1Y2: " ", X1Y3: " ", X2Y1: " ", X2Y2: " ", X2Y3: " ", X3Y1: " ", X3Y2: " ", X3Y3: " "},
		{X: 1, Y: 2, X1Y1: " ", X1Y2: " ", X1Y3: " ", X2Y1: " ", X2Y2: " ", X2Y3: " ", X3Y1: " ", X3Y2: " ", X3Y3: " "},
		{X: 1, Y: 3, X1Y1: " ", X1Y2: " ", X1Y3: " ", X2Y1: " ", X2Y2: " ", X2Y3: " ", X3Y1: " ", X3Y2: " ", X3Y3: " "},
		{X: 2, Y: 1, X1Y1: " ", X1Y2: " ", X1Y3: " ", X2Y1: " ", X2Y2: " ", X2Y3: " ", X3Y1: " ", X3Y2: " ", X3Y3: " "},
		{X: 2, Y: 2, X1Y1: " ", X1Y2: " ", X1Y3: " ", X2Y1: " ", X2Y2: " ", X2Y3: " ", X3Y1: " ", X3Y2: " ", X3Y3: " "},
		{X: 2, Y: 3, X1Y1: " ", X1Y2: " ", X1Y3: " ", X2Y1: " ", X2Y2: " ", X2Y3: " ", X3Y1: " ", X3Y2: " ", X3Y3: " "},
		{X: 3, Y: 1, X1Y1: " ", X1Y2: " ", X1Y3: " ", X2Y1: " ", X2Y2: " ", X2Y3: " ", X3Y1: " ", X3Y2: " ", X3Y3: " "},
		{X: 3, Y: 2, X1Y1: " ", X1Y2: " ", X1Y3: " ", X2Y1: " ", X2Y2: " ", X2Y3: " ", X3Y1: " ", X3Y2: " ", X3Y3: " "},
		{X: 3, Y: 3, X1Y1: " ", X1Y2: " ", X1Y3: " ", X2Y1: " ", X2Y2: " ", X2Y3: " ", X3Y1: " ", X3Y2: " ", X3Y3: " "},
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
		Label:    fmt.Sprintf("Turn %d, you play as %s. Choose your next move.", turnCounter, player),
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
		fmt.Println("Turn:", turnCounter)
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

		gameBoard, playerChoices, choicesHistory = selectMove(gameBoard, playerChoices[v], player, &turnCounter, v, choicesHistory)
		fmt.Println("after select: ")
		printWholeBoard(gameBoard)
		if turnCounter > 3 {
			win := checkWin(gameBoard[0])
			if win != "" {
				//printBoard(gameBoard[0])
				fmt.Printf("Player %s won.\n\n", win)
				return
			} else if turnCounter == 8 {
				fmt.Printf("Draw.\n\n")
				return
			}
		}

		// opponent random move
		opponentChoice = rand.Intn(len(playerChoices))
		gameBoard, playerChoices, choicesHistory = selectMove(gameBoard, playerChoices[opponentChoice], opponent, &turnCounter, opponentChoice, choicesHistory)
		if turnCounter > 3 {
			win := checkWin(gameBoard[0])
			if win != "" {
				//printBoard(gameBoard[0])
				fmt.Printf("Player %s won.\n\n", win)
				return
			} else if turnCounter == 8 {
				fmt.Printf("Draw.\n\n")
				return
			}
		}
	}
	//
	//
	//
	//
	/*gamePrompt := promptui.Select{
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
	fmt.Printf("Chosen option %d\n", v)*/
}
