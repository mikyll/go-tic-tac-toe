package main

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"golang.org/x/term"
)

const MAIN_MENU = -1
const SINGLEPLAYER_MENU = -2
const MULTIPLAYER_MENU = -3
const ROOM_MENU = -4
const ABOUT_MENU = -4
const SINGLEPLAYER = 0
const MULTIPLAYER = 1
const ABOUT = 2
const QUIT = 3
const EASY = 0
const HARD = 1
const BACK_SINGLEPLAYER = 2
const LOCAL = 0
const LAN = 1
const BACK_MULTIPLAYER = 2
const CREATE_ROOM = 0
const JOIN_ROOM = 1
const BACK_ROOM = 2
const BACK_ABOUT = 0
const MODE_SP_EASY = 0
const MODE_SP_HARD = 1
const MODE_MP_LOCAL = 2
const MODE_MP_LAN = 3

// Utility to skip the terminal bell sound, when selecting prompt options
type bellSkipper struct{}

// Write implements an io.WriterCloser over os.Stderr, but it skips the terminal bell character.
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

func selectMove(gameBoard [9]board, updatedBoard board, player string, playerTurn int, playerMove int, moveHistory []int) ([9]board, []board, []int) {
	newBoard := gameBoard
	choices := gameBoard[:]

	moveHistory[playerTurn] = playerMove

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
	if playerTurn != 8 {
		for i := 0; moveHistory[i] != -1; i++ {
			choices = remove(choices, moveHistory[i])
		}
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
	if gameBoard.X1Y1 != "" && gameBoard.X1Y1 != " " && gameBoard.X1Y1 == gameBoard.X1Y2 && gameBoard.X1Y1 == gameBoard.X1Y3 {
		p := gameBoard.X1Y1
		fmt.Printf("\n\n %v | %v | %v \n---+---+---\n   |   |   \n---+---+---\n   |   |   ", p, p, p)
		return gameBoard.X1Y1
	}
	if gameBoard.X2Y1 != "" && gameBoard.X2Y1 != " " && gameBoard.X2Y1 == gameBoard.X2Y2 && gameBoard.X2Y1 == gameBoard.X2Y3 {
		p := gameBoard.X2Y1
		fmt.Printf("\n\n   |   |   \n---+---+---\n %v | %v | %v \n---+---+---\n   |   |   ", p, p, p)
		return gameBoard.X2Y1
	}
	if gameBoard.X3Y1 != "" && gameBoard.X3Y1 != " " && gameBoard.X3Y1 == gameBoard.X3Y2 && gameBoard.X3Y1 == gameBoard.X3Y3 {
		p := gameBoard.X3Y1
		fmt.Printf("\n\n   |   |   \n---+---+---\n   |   |   \n---+---+---\n %v | %v | %v ", p, p, p)
		return gameBoard.X3Y1
	}
	// columns
	if gameBoard.X1Y1 != "" && gameBoard.X1Y1 != " " && gameBoard.X1Y1 == gameBoard.X2Y1 && gameBoard.X1Y1 == gameBoard.X3Y1 {
		p := gameBoard.X1Y1
		fmt.Printf("\n\n %v |   |   \n---+---+---\n %v |   |   \n---+---+---\n %v |   |   ", p, p, p)
		return gameBoard.X1Y1
	}
	if gameBoard.X1Y2 != "" && gameBoard.X1Y2 != " " && gameBoard.X1Y2 == gameBoard.X2Y2 && gameBoard.X1Y2 == gameBoard.X3Y2 {
		p := gameBoard.X1Y2
		fmt.Printf("\n\n   | %v |   \n---+---+---\n   | %v |   \n---+---+---\n   | %v |   ", p, p, p)
		return gameBoard.X1Y2
	}
	if gameBoard.X1Y3 != "" && gameBoard.X1Y3 != " " && gameBoard.X1Y3 == gameBoard.X2Y3 && gameBoard.X1Y3 == gameBoard.X3Y3 {
		p := gameBoard.X1Y3
		fmt.Printf("\n\n   |   | %v \n---+---+---\n   |   | %v \n---+---+---\n   |   | %v ", p, p, p)
		return gameBoard.X1Y3
	}
	// others
	if gameBoard.X1Y1 != "" && gameBoard.X1Y1 != " " && gameBoard.X1Y1 == gameBoard.X2Y2 && gameBoard.X1Y1 == gameBoard.X3Y3 {
		p := gameBoard.X1Y1
		fmt.Printf("\n\n %v |   |   \n---+---+---\n   | %v |   \n---+---+---\n   |   | %v ", p, p, p)
		return gameBoard.X1Y1
	}
	if gameBoard.X3Y1 != "" && gameBoard.X3Y1 != " " && gameBoard.X3Y1 == gameBoard.X2Y2 && gameBoard.X3Y1 == gameBoard.X1Y3 {
		p := gameBoard.X3Y1
		fmt.Printf("\n\n   |   | %v \n---+---+---\n   | %v |   \n---+---+---\n %v |   |   ", p, p, p)
		return gameBoard.X3Y1
	}
	return ""
}

// utilities
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
[0]         // [1]         // [2]         // [3]         // [4]         // [5]         // [6]         // [7]         // [8]
 %s | %s | %s  //  %s | %s | %s  //  %s | %s | %s  //  %s | %s | %s  //  %s | %s | %s  //  %s | %s | %s  //  %s | %s | %s  //  %s | %s | %s  //  %s | %s | %s 
---+---+--- // ---+---+--- // ---+---+--- // ---+---+--- // ---+---+--- // ---+---+--- // ---+---+--- // ---+---+--- // ---+---+--- 
 %s | %s | %s  //  %s | %s | %s  //  %s | %s | %s  //  %s | %s | %s  //  %s | %s | %s  //  %s | %s | %s  //  %s | %s | %s  //  %s | %s | %s  //  %s | %s | %s 
---+---+--- // ---+---+--- // ---+---+--- // ---+---+--- // ---+---+--- // ---+---+--- // ---+---+--- // ---+---+--- // ---+---+--- 
 %s | %s | %s  //  %s | %s | %s  //  %s | %s | %s  //  %s | %s | %s  //  %s | %s | %s  //  %s | %s | %s  //  %s | %s | %s  //  %s | %s | %s  //  %s | %s | %s 

 `, b[0].X1Y1, b[0].X1Y2, b[0].X1Y3, b[1].X1Y1, b[1].X1Y2, b[1].X1Y3, b[2].X1Y1, b[2].X1Y2, b[2].X1Y3, b[3].X1Y1, b[3].X1Y2, b[3].X1Y3, b[4].X1Y1, b[4].X1Y2, b[4].X1Y3, b[5].X1Y1, b[5].X1Y2, b[5].X1Y3, b[6].X1Y1, b[6].X1Y2, b[6].X1Y3, b[7].X1Y1, b[7].X1Y2, b[7].X1Y3, b[8].X1Y1, b[8].X1Y2, b[8].X1Y3,
		b[0].X2Y1, b[0].X2Y2, b[0].X2Y3, b[1].X2Y1, b[1].X2Y2, b[1].X2Y3, b[2].X2Y1, b[2].X2Y2, b[2].X2Y3, b[3].X2Y1, b[3].X2Y2, b[3].X2Y3, b[4].X2Y1, b[4].X2Y2, b[4].X2Y3, b[5].X2Y1, b[5].X2Y2, b[5].X2Y3, b[6].X2Y1, b[6].X2Y2, b[6].X2Y3, b[7].X2Y1, b[7].X2Y2, b[7].X2Y3, b[8].X2Y1, b[8].X2Y2, b[8].X2Y3,
		b[0].X3Y1, b[0].X3Y2, b[0].X3Y3, b[1].X3Y1, b[1].X3Y2, b[1].X3Y3, b[2].X3Y1, b[2].X3Y2, b[2].X3Y3, b[3].X3Y1, b[3].X3Y2, b[3].X3Y3, b[4].X3Y1, b[4].X3Y2, b[4].X3Y3, b[5].X3Y1, b[5].X3Y2, b[5].X3Y3, b[6].X3Y1, b[6].X3Y2, b[6].X3Y3, b[7].X3Y1, b[7].X3Y2, b[7].X3Y3, b[8].X3Y1, b[8].X3Y2, b[8].X3Y3)
}

func pressAnyKey(message string) {
	// switch stdin into 'raw' mode
	fmt.Println(message)
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println(err)
		return
	}

	b := make([]byte, 1)
	_, err = os.Stdin.Read(b)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("the char %q was hit", string(b[0]))
	term.Restore(int(os.Stdin.Fd()), oldState)
}

var ps = playerSymbols{P1: "X", P2: "O"}

func game(mode int, first bool, c net.Conn) {
	var v int
	var err error

	turnCounter := 0
	choicesHistory := []int{-1, -1, -1, -1, -1, -1, -1, -1, -1}
	// NB: player1 is always X, player2 is always O
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
		Active:   "> ({{ .X | cyan }}, {{ .Y | green }})", // nice alternative: \U000027A4
		Inactive: " ({{ .X | cyan }}, {{ .Y | green }})",
		Selected: "> ({{ .X | cyan }}, {{ .Y | green }})",
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
	for {
		gameTemplate.Label = fmt.Sprintf("Turn %d, you play as %s. Choose your next move.", turnCounter, player)
		gamePrompt := promptui.Select{
			Label:     "",
			Items:     playerChoices,
			Templates: gameTemplate,
			Size:      4,
			Stdout:    &bellSkipper{},
		}

		switch mode {
		case MODE_SP_EASY, MODE_SP_HARD, MODE_MP_LOCAL:
			v, _, err = gamePrompt.Run()
			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Chosen option %d\n", v)
		case MODE_MP_LAN:
			if first {
				v, _, err = gamePrompt.Run()
				if err != nil {
					fmt.Printf("Prompt failed %v\n", err)
					os.Exit(1)
				}
				fmt.Printf("Chosen option %d\n", v)

				// send to client
				fmt.Fprintf(c, strconv.Itoa(v)+"\n")
			} else {
				// no prompt & receive move from client
				printBoard(gameBoard[0])
				fmt.Println("Waiting for opponent move...")
				message, _ := bufio.NewReader(c).ReadString('\n')
				v, err = strconv.Atoi(strings.TrimSuffix(message, "\n"))
				if err != nil {
					fmt.Println(err)
					return
				}

				fmt.Println("Turn", strconv.Itoa(turnCounter), "chose"+strconv.Itoa(v))
			}
		}
		gameBoard, playerChoices, choicesHistory = selectMove(gameBoard, playerChoices[v], player, turnCounter, v, choicesHistory)
		if turnCounter > 3 {
			win := checkWin(gameBoard[0])
			if win != "" {
				printBoard(gameBoard[0])
				fmt.Printf("Player %s won.\n\n", win)
				return
			} else if turnCounter == 8 {
				printBoard(gameBoard[0])
				fmt.Printf("Draw.\n\n")
				return
			}
		}
		turnCounter++

		// opponent random move
		switch mode {
		case MODE_SP_EASY:
			opponentChoice = rand.Intn(len(playerChoices))
		case MODE_SP_HARD:
			// TO-DO
			return
		case MODE_MP_LOCAL:
			gameTemplate.Label = fmt.Sprintf("Turn %d, you play as %s. Choose your next move.", turnCounter, opponent)
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
				os.Exit(1)
			}
			fmt.Printf("Chosen option %d\n", v)
			opponentChoice = v
		case MODE_MP_LAN:
			gameTemplate.Label = fmt.Sprintf("Turn %d, you play as %s. Choose your next move.", turnCounter, opponent)
			gamePrompt := promptui.Select{
				Label:     "",
				Items:     playerChoices,
				Templates: gameTemplate,
				Size:      4,
				Stdout:    &bellSkipper{},
			}

			if !first {
				v, _, err = gamePrompt.Run()
				if err != nil {
					fmt.Printf("Prompt failed %v\n", err)
					os.Exit(1)
				}
				fmt.Printf("Chosen option %d\n", v)

				// send to client
				fmt.Fprintf(c, strconv.Itoa(v)+"\n")
			} else {
				// no prompt & receive move from client
				printBoard(gameBoard[0])
				fmt.Println("Waiting for opponent move...")
				message, _ := bufio.NewReader(c).ReadString('\n')
				v, err = strconv.Atoi(strings.TrimSuffix(message, "\n"))
				if err != nil {
					fmt.Println(err)
					return
				}

				fmt.Println("Turn", strconv.Itoa(turnCounter), "chose"+strconv.Itoa(v))
			}
			opponentChoice = v
		}
		gameBoard, playerChoices, choicesHistory = selectMove(gameBoard, playerChoices[opponentChoice], opponent, turnCounter, opponentChoice, choicesHistory)
		if turnCounter > 3 {
			win := checkWin(gameBoard[0])
			if win != "" {
				printBoard(gameBoard[0])
				fmt.Printf("Player %s won.\n\n", win)
				return
			} else if turnCounter == 8 {
				printBoard(gameBoard[0])
				fmt.Printf("Draw.\n\n")
				return
			}
		}
		turnCounter++
	}
}

func validateIP(address string) error {
	if net.ParseIP(address) == nil {
		return errors.New("invalid IP address")
	}
	return nil
}

func main() {
	if !term.IsTerminal(int(syscall.Stdin)) {
		fmt.Println("Terminal is not interactive! Consider using flags or environment variables!")
		return
	}

	rand.Seed(time.Now().UnixNano())

	var state int
	var err error
	var ip string

	// main menu entries init
	mainMenu := []menu{
		{Entry: "Single Player"},
		{Entry: "Multiplayer"},
		{Entry: "About"},
		{Entry: "Quit"},
	}

	singlePlayerMenu := []menu{
		{Entry: "Easy"},
		{Entry: "Hard (coming soon)"},
		{Entry: "Back"},
	}

	multiPlayerMenu := []menu{
		{Entry: "Local"},
		{Entry: "LAN"},
		{Entry: "Back"},
	}

	roomMenu := []menu{
		{Entry: "Create Game"},
		{Entry: "Join Game"},
		{Entry: "Back"},
	}

	mainMenuTemplate := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "> {{ .Entry | cyan }}",
		Inactive: " {{ .Entry | white }} ",
		Selected: "> {{ .Entry | white }}",
	}

	mainMenuPrompt := promptui.Select{
		Label:     "-------- Main Menu ---------",
		Items:     mainMenu,
		Templates: mainMenuTemplate,
		Size:      4,
		Stdout:    &bellSkipper{},
	}

	singlePlayerTemplate := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "> {{ .Entry | cyan }}",
		Inactive: " {{ .Entry | white }} ",
		Selected: "> {{ .Entry | white }}",
	}

	singlePlayerPrompt := promptui.Select{
		Label:     "------ Single Player -------",
		Items:     singlePlayerMenu,
		Templates: singlePlayerTemplate,
		Size:      3,
		Stdout:    &bellSkipper{},
	}

	multiPlayerTemplate := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "> {{ .Entry | cyan }}",
		Inactive: " {{ .Entry | white }} ",
		Selected: "> {{ .Entry | white }}",
	}

	multiPlayerPrompt := promptui.Select{
		Label:     "------ Multi Player --------",
		Items:     multiPlayerMenu,
		Templates: multiPlayerTemplate,
		Size:      3,
		Stdout:    &bellSkipper{},
	}

	roomTemplate := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "> {{ .Entry | cyan }}",
		Inactive: " {{ .Entry | white }} ",
		Selected: "> {{ .Entry | white }}",
	}

	roomPrompt := promptui.Select{
		Label:     "------ Multi Player --------",
		Items:     roomMenu,
		Templates: roomTemplate,
		Size:      3,
		Stdout:    &bellSkipper{},
	}

	ipTemplate := &promptui.PromptTemplates{
		Prompt:  "{{ . }}",
		Valid:   "{{ . | green }}",
		Invalid: "{{ . | red }}",
		Success: "Connecting to ",
	}

	ipPrompt := promptui.Prompt{
		Label:     "Connect to IP address: ",
		Templates: ipTemplate,
		Validate:  validateIP,
		Stdout:    &bellSkipper{},
	}

	state = MAIN_MENU
	// main menu loop
	for state == MAIN_MENU {
		state, _, err = mainMenuPrompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			os.Exit(1)
		}
		switch state {
		case SINGLEPLAYER:
			state = SINGLEPLAYER_MENU
			// singleplayer loop
			for state == SINGLEPLAYER_MENU {
				state, _, err = singlePlayerPrompt.Run()
				if err != nil {
					fmt.Printf("Prompt failed %v\n", err)
					os.Exit(1)
				}
				switch state {
				case EASY:
					game(MODE_SP_EASY, true, nil)
					pressAnyKey("Press any key to continue ... ")
					state = MAIN_MENU
				case HARD:
					// TO-DO
					state = MAIN_MENU
				case BACK_SINGLEPLAYER:
					state = BACK_SINGLEPLAYER
				}
				state = MAIN_MENU
			}
			state = MAIN_MENU
		case MULTIPLAYER:
			state = MULTIPLAYER_MENU
			// multiplayer loop
			for state == MULTIPLAYER_MENU {
				state, _, err = multiPlayerPrompt.Run()
				if err != nil {
					fmt.Printf("Prompt failed %v\n", err)
					os.Exit(1)
				}
				switch state {
				case LOCAL:
					game(MODE_MP_LOCAL, true, nil)
					pressAnyKey("Press any key to continue ... ")
					state = MAIN_MENU
				case LAN:
					// room loop
					state = ROOM_MENU
					for state == ROOM_MENU {
						state, _, err = roomPrompt.Run()
						if err != nil {
							fmt.Printf("Prompt failed %v\n", err)
							os.Exit(1)
						}
						switch state {
						case CREATE_ROOM:
							// TO-DO: listening on IP:port ...
							address := "localhost:4000"
							l, err := net.Listen("tcp", address)
							if err != nil {
								fmt.Println(err)
								return
							}
							fmt.Println("[Server] Waiting for connections on", address)

							// Listen for Client connection
							c, err := l.Accept()
							if err != nil {
								fmt.Println(err)
								return
							}
							fmt.Println("[Server] Connection accepted, picking player turns...")

							/*// auxiliary function
							firstTurn := func() bool {
								return rand.Intn(2) == 1
							}*/

							// turn picking
							turn := rand.Intn(2) == 1
							if turn {
								fmt.Println("[Server] You play as", ps.P1)
							} else {
								fmt.Println("[Server] You play as", ps.P2)
							}

							// Send turn to Client
							fmt.Fprintf(c, strconv.FormatBool(!turn)+"\n")

							// Start game
							game(MODE_MP_LAN, turn, c)

							l.Close()
							// NB if using a game() function it could be nice to use defer to close the listen(?)
							state = ROOM_MENU
						case JOIN_ROOM:
							// TO-DO: enter&validate IP ...
							ip, err = ipPrompt.Run()
							if err != nil {
								fmt.Printf("Prompt failed %v\n", err)
								return
							}

							fmt.Println("[Client] Connection attempt to", ip+":4000")

							// connect to Server
							c, err := net.Dial("tcp", ip+":4000")
							if err != nil {
								fmt.Println(err)
								return
							}

							// Receive turn
							message, _ := bufio.NewReader(c).ReadString('\n')
							turn, err := strconv.ParseBool(strings.TrimSuffix(message, "\n"))
							if err != nil {
								fmt.Println(err)
								return
							}

							//fmt.Println("[Client] Message Received:", message) // test

							if turn {
								fmt.Println("[Client] You play as", ps.P1)
							} else {
								fmt.Println("[Client] You play as", ps.P2)
							}

							// Start game
							game(MODE_MP_LAN, turn, c)

							state = ROOM_MENU
						case BACK_ROOM:
							state = MULTIPLAYER_MENU
						}
					}
					state = MULTIPLAYER_MENU
				case BACK_MULTIPLAYER:
					state = BACK_MULTIPLAYER
				}
			}
			state = MAIN_MENU
		case ABOUT:
			// TO-DO
			fmt.Printf("---------- About -----------\nTic-Tac-Toe\nApp developed by ")
			color.Cyan("Michele Righi")

			fmt.Printf("\nIf you like it, consider leaving a star on GitHub!\nLink: https://github.com/mikyll/go-tic-tac-toe\n\n\n")
			pressAnyKey("Press any key to go back ... ")
			state = MAIN_MENU
		case QUIT:
			return
		}
	}
}
