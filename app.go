package main

import (
	"fmt"
	"os"
	"syscall"

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

func selectMove(gameBoard []board, updatedBoard board, player string, playerTurn int, playerMove int, moveHistory []int) ([]board, []board, []int) {
	newBoard := gameBoard
	choices := gameBoard
	var nextPlayer string

	moveHistory[playerTurn] = playerMove

	if player == "X" {
		nextPlayer = "O"
	} else if player == "O" {
		nextPlayer = "X"
	}

	for i := 0; i < len(gameBoard); i++ {
		newBoard[i] = updatedBoard
		choices[i] = updatedBoard
	}

	choices[0].X1Y1 = nextPlayer
	choices[0].X = 1
	choices[0].Y = 1
	choices[1].X1Y2 = nextPlayer
	choices[1].X = 1
	choices[1].Y = 2
	choices[2].X1Y3 = nextPlayer
	choices[2].X = 1
	choices[2].Y = 3
	choices[3].X2Y1 = nextPlayer
	choices[3].X = 2
	choices[3].Y = 1
	choices[4].X2Y2 = nextPlayer
	choices[4].X = 2
	choices[4].Y = 2
	choices[5].X2Y3 = nextPlayer
	choices[5].X = 2
	choices[5].Y = 3
	choices[6].X3Y1 = nextPlayer
	choices[6].X = 3
	choices[6].Y = 1
	choices[7].X3Y2 = nextPlayer
	choices[7].X = 3
	choices[7].Y = 2
	choices[8].X3Y3 = nextPlayer
	choices[8].X = 3
	choices[8].Y = 3

	// remove choices already picked
	for i := 0; moveHistory[i] != -1; i++ {
		choices = remove(choices, moveHistory[i])
	}

	return newBoard, choices, moveHistory
}

/*func doOpponent(oldBoard []board, choices []board, move int, player string) ([]board, []board) {
	newBoard := oldBoard
	newChoices := choices
}*/

/*func selectMove(oldBoard []board, move int, player string) ([]board, []board) {
	newBoard := oldBoard
	// creo un nuovo array per le scelte, che corrisponde alla
	// newChoices := newBoard

	//playerChoice := string(choices[move].X + choices[move].Y)
	// add player in choices
	for i := 0; i < len(oldBoard); i++ {
		newBoard[i] = choices[move]
		newChoices[i] = choices[move]
		switch i {
		case 0:
			if choices[0].X1Y1 == " " {
				newChoices[0].X1Y1 = player
				newChoices[0].X = 1
				newChoices[0].Y = 1
			}
		case 1:
			if choices[1].X1Y2 == " " {
				newChoices[1].X1Y2 = player
				newChoices[1].X = 1
				newChoices[1].Y = 2
			}
		case 2:
			if choices[2].X1Y3 == " " {
				newChoices[2].X1Y3 = player
				newChoices[2].X = 1
				newChoices[2].Y = 3
			}
		case 3:
			if choices[3].X2Y1 == " " {
				newChoices[3].X2Y1 = player
				newChoices[3].X = 2
				newChoices[3].Y = 1
			}
		case 4:
			if choices[4].X2Y2 == " " {
				newChoices[4].X2Y2 = player
				newChoices[4].X = 2
				newChoices[4].Y = 2
			}
		case 5:
			if choices[5].X2Y3 == " " {
				newChoices[5].X2Y3 = player
				newChoices[5].X = 2
				newChoices[5].Y = 3
			}
		case 6:
			if choices[6].X3Y1 == " " {
				newChoices[6].X3Y1 = player
				newChoices[6].X = 3
				newChoices[6].Y = 1
			}
		case 7:
			if choices[7].X3Y2 == " " {
				newChoices[7].X3Y2 = player
				newChoices[7].X = 3
				newChoices[7].Y = 2
			}
		case 8:
			if choices[8].X3Y3 == " " {
				newChoices[8].X3Y3 = player
				newChoices[8].X = 3
				newChoices[8].Y = 3
			}
		}
	}

	for i := 0; i < len(newChoices); i++ {
		currChoice := strconv.Itoa(newChoices[i].X) + strconv.Itoa(newChoices[i].Y)
		switch currChoice {
		case "11":
			newChoices[i].X = 1
			newChoices[i].Y = 1
			if newChoices[i].X1Y1 == " " {
				newChoices[i].X1Y1 = player
			}
		case "12":
			newChoices[i].X = 1
			newChoices[i].Y = 2
			if newChoices[i].X1Y2 == " " {
				newChoices[i].X1Y2 = player
			}
		case "13":
			newChoices[i].X = 1
			newChoices[i].Y = 3
			if newChoices[i].X1Y3 == " " {
				newChoices[i].X1Y3 = player
			}
		case "21":
			newChoices[i].X = 2
			newChoices[i].Y = 1
			if newChoices[i].X2Y1 == " " {
				newChoices[i].X2Y1 = player
			}
		case "22":
			newChoices[i].X = 2
			newChoices[i].Y = 2
			if newChoices[i].X2Y2 == " " {
				newChoices[i].X2Y2 = player
			}
		case "23":
			newChoices[i].X = 2
			newChoices[i].Y = 3
			if newChoices[i].X2Y3 == " " {
				newChoices[i].X2Y3 = player
			}
		case "31":
			newChoices[i].X = 3
			newChoices[i].Y = 1
			if newChoices[i].X3Y1 == " " {
				newChoices[i].X3Y1 = player
			}
		case "32":
			newChoices[i].X = 3
			newChoices[i].Y = 2
			if newChoices[i].X3Y2 == " " {
				newChoices[i].X3Y2 = player
			}
		case "33":
			newChoices[i].X = 3
			newChoices[i].Y = 3
			if newChoices[i].X3Y3 == " " {
				newChoices[i].X3Y3 = player
			}
		}
	}
	newChoices = remove(newChoices, move)

	return newBoard, newChoices
}*/

func printBoard(b []board) {
	fmt.Printf(`
 %s | %s | %s
---+---+---
 %s | %s | %s
---+---+---
 %s | %s | %s
 `, b[0].X1Y1, b[1].X1Y2, b[2].X1Y3, b[3].X2Y1, b[4].X2Y2, b[5].X2Y3, b[6].X3Y1, b[7].X3Y2, b[8].X3Y3)
}

func main() {
	if !term.IsTerminal(int(syscall.Stdin)) {
		fmt.Println("Terminal is not interactive! Consider using flags or environment variables!")
		return
	}

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

         {{ .X1Y1 }} | {{ .X1Y2}} | {{ .X1Y3 }}  {{"1" | green}}
        ---+---+---
         {{ .X2Y1 }} | {{ .X2Y2}} | {{ .X2Y3 }}  {{"2" | green}}
        ---+---+---
         {{ .X3Y1 }} | {{ .X3Y2}} | {{ .X3Y3 }}  {{"3" | green}}
         {{"1   2   3" | cyan}}

----------------------------
Selected Move: %s in ({{ .X | cyan }}, {{ .Y | green }})`, player),
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

	gameBoard, playerChoices, choicesHistory = selectMove(gameBoard, playerChoices[v], player, turnCounter, v, choicesHistory)
	//_, playerChoices = selectMove(oldBoard, playerChoices, v, player)

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
	/*_, playerChoices = selectMove(oldBoard, playerChoices, v, player)

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
	//
	//
	//
	//
	//
	//
	//
	//
	//
	//
	//
	//
	//
	//
	//
	//
	//
	//
	//
	//
	/*peppers := []pepper{
		{Name: "Bell Pepper", HeatUnit: 0, Peppers: 0},
		{Name: "Banana Pepper", HeatUnit: 100, Peppers: 1},
		{Name: "Poblano", HeatUnit: 1000, Peppers: 2},
		{Name: "Jalapeño", HeatUnit: 3500, Peppers: 3},
		{Name: "Aleppo", HeatUnit: 10000, Peppers: 4},
		{Name: "Tabasco", HeatUnit: 30000, Peppers: 5},
		{Name: "Malagueta", HeatUnit: 50000, Peppers: 6},
		{Name: "Habanero", HeatUnit: 100000, Peppers: 7},
		{Name: "Red Savina Habanero", HeatUnit: 350000, Peppers: 8},
		{Name: "Dragon’s Breath", HeatUnit: 855000, Peppers: 9},
	}
	templates := &promptui.SelectTemplates{
			Label:    "{{ . }}?",
			Active:   "\U0001F336 {{ .Name | cyan }} ({{ .HeatUnit | red }})",
			Inactive: "  {{ .Name | cyan }} ({{ .HeatUnit | red }})",
			Selected: "\U0001F336 {{ .Name | red | cyan }}",
			Details: `
	--------- Pepper ----------
	{{ "Name:" | faint }}	{{ .Name }}
	{{ "Heat Unit:" | faint }}	{{ .HeatUnit }}
	{{ "Peppers:" | faint }}	{{ .Peppers }}`,
		}

		searcher := func(input string, index int) bool {
			pepper := peppers[index]
			name := strings.Replace(strings.ToLower(pepper.Name), " ", "", -1)
			input = strings.Replace(strings.ToLower(input), " ", "", -1)

			return strings.Contains(name, input)
		}

		prompt := promptui.Select{
			Label:     "Spicy Level",
			Items:     peppers,
			Templates: templates,
			Size:      4,
			Searcher:  searcher,
		}

		i, _, err := prompt.Run()

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		fmt.Printf("You choose number %d: %s\n", i+1, peppers[i].Name)*/
}

/*type menu struct {
	Entry string
}

type board struct {
	X1Y1, X1Y2, X1Y3, X2Y1, X2Y2, X2Y3, X3Y1, X3Y2, X3Y3 string
}

type choices struct {
	X, Y  int
	Board board
}

type pepper struct {
	Name     string
	HeatUnit int
	Peppers  int
}*/

/*func main() {
	prompt := promptui.Select{
		Label: "Select Day",
		Items: []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday",
			"Saturday", "Sunday"},
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	fmt.Printf("You choose %q\n", result)

	// Create a tic-tac-toe board.
	board := [][]string{
		[]string{"_", "_", "_"},
		[]string{"_", "_", "_"},
		[]string{"_", "_", "_"},
	}

	// The players take turns.
	board[0][0] = "X"
	board[2][2] = "O"
	board[1][2] = "X"
	board[1][0] = "O"
	board[0][2] = "X"

	for i := 0; i < len(board); i++ {
		fmt.Printf("%s\n", strings.Join(board[i], " "))
	}
}*/
