package main

import (
	"fmt"
	"os"
	"syscall"

	"github.com/manifoldco/promptui"
	"golang.org/x/term"
)

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

func remove(slice []choices, i int) []choices {
	return append(slice[:i], slice[i+1:]...)
}

type choices struct {
	X, Y                                                 int
	X1Y1, X1Y2, X1Y3, X2Y1, X2Y2, X2Y3, X3Y1, X3Y2, X3Y3 string
}

func selectMove(oldBoard []choices, choices []choices, move int, player string) ([]choices, []choices) {
	newBoard := oldBoard
	newChoices := choices

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
	newChoices = remove(newChoices, move)

	return newBoard, newChoices
}

func printBoard(board []choices) {
	fmt.Printf(`
 %s | %s | %s
---+---+---
 %s | %s | %s
---+---+---
 %s | %s | %s
 `, board[0].X1Y1, board[1].X1Y2, board[2].X1Y3, board[3].X2Y1, board[4].X2Y2, board[5].X2Y3, board[6].X3Y1, board[7].X3Y2, board[8].X3Y3)
}

func main() {
	if !term.IsTerminal(int(syscall.Stdin)) {
		fmt.Println("Terminal is not interactive! Consider using flags or environment variables!")
		return
	}

	/*mainMenu := []menu{
		{Entry: "Single Player"},
		{Entry: "Multiplayer"},
		{Entry: "About"},
		{Entry: "Quit"},
	}

	singlePlayerMenu := []menu{
		{Entry: "Easy"},
		{Entry: "Hard"},
		{Entry: "About"},
		{Entry: "Back"},
	}

	mainMenuTemplate := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "\U000027A4  {{ .Entry | white }}",
		Inactive: "  {{ .Entry | white }} ",
		Selected: "\U00002023 {{ .Entry | white }}",
	}

	mainMenuPrompt := promptui.Select{
		Label:     "--- Main Menu ---",
		Items:     mainMenu,
		Templates: mainMenuTemplate,
		Size:      4,
	}

	i, _, err := mainMenuPrompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}
	fmt.Printf("Chosen option %d\n", i)*/

	// game
	player := " "
	oldBoard := []choices{
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
	player = "X"
	choices := []choices{
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
		Details: fmt.Sprintf(
			`
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
		Items:     choices,
		Templates: gameTemplate,
		Size:      4,
		Stdout:    &bellSkipper{},
	}

	i, _, err := gamePrompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}
	fmt.Printf("Chosen option %d\n", i)

	// update only the chosen element in every
	/*newBoard := oldBoard
	for j := 0; j < len(newBoard); j++ {
		newBoard[j] = choices[i]
		choices[j] = newBoard[j]
		// need to add the X to the choices
	}
	// remove the choice at #i
	choices = remove(choices, i)

	for k := 0; k < len(choices); k++ {
		fmt.Printf("(%s %s %s, %s %s %s, %s %s %s)\n", choices[k].X1Y1, choices[k].X1Y2, choices[k].X1Y3, choices[k].X2Y1, choices[k].X2Y2, choices[k].X2Y3, choices[k].X3Y1, choices[k].X3Y2, choices[k].X3Y3)
	}
	*/

	oldBoard, choices = selectMove(oldBoard, choices, i, player)

	/*fmt.Println("newBoard")
	printBoard(oldBoard)
	fmt.Println("choices")
	printBoard(choices)*/

	gamePrompt = promptui.Select{
		Label:     "",
		Items:     choices,
		Templates: gameTemplate,
		Size:      4,
		Stdout:    &bellSkipper{},
	}
	i, _, err = gamePrompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}
	fmt.Printf("Chosen option %d\n", i)
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
