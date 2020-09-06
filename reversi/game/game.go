package game

import (
	"github.com/kazu22002/go-reversi/reversi/board"
	"github.com/kazu22002/go-reversi/reversi/cell"
	"github.com/kazu22002/go-reversi/reversi/player"
)

type Game struct {
	Board           board.Board
	Players         []player.Player
	CurrPlayerIndex uint8
	SelectPosition  int
	EnablePosition  []cell.Cell
}

func New(players []player.Player) Game {
	gameBoard, _ := board.InitCells(board.New(8, 8))
	enableCell := []cell.Cell{}
	return Game{
		gameBoard,
		players,
		0,
		0,
		enableCell,
	}
}

func Render(game Game) string {
	if len(game.EnablePosition) > 0  {
		return board.Render(game.Board, game.EnablePosition[game.SelectPosition])
	}

	return board.Render(game.Board, cell.Cell{})
}

func Result(game Game) (int, int) {
	return board.Result(game.Board)
}

func IsFinished(game Game) bool {
	return board.IsFull(game.Board)
}

func GetCurrentPlayer(game Game) player.Player {
	return game.Players[game.CurrPlayerIndex]
}

func GetScore(game Game) map[player.Player]uint8 {
	dist := board.GetCellDistribution(game.Board)
	score := make(map[player.Player]uint8, 2)
	for _, player := range game.Players {
		score[player] = dist[player.CellType]
	}
	return score
}

func SwitchPlayer(game Game) Game {

	newGame := game
	if newGame.CurrPlayerIndex == 0 {
		newGame.CurrPlayerIndex = 1
	} else {
		newGame.CurrPlayerIndex = 0
	}
	return newGame
}

func CanPlayerChangeCells(player player.Player, currentGame Game) bool {
	return len(board.GetLegalCellChangesForCellType(player.CellType, currentGame.Board)) > 0
}

//func RenderAskBoard(game Game) string {
//	currentPlayer := GetCurrentPlayer(game)
//	legalCellChanges := board.GetLegalCellChangesForCellType(currentPlayer.CellType, game.Board)
//	return board.Render(game.Board, legalCellChanges)
//}

func PlayTurn(currentGame Game) (Game, error) {
	if IsFinished(currentGame){
		return currentGame, nil
	}

	if !CanPlayerChangeCells(GetCurrentPlayer(currentGame), currentGame) {
		newGame := SwitchPlayer(currentGame)
		return PlayTurn(newGame)
	}

	currentGame.SelectPosition = 0
	currentPlayer := GetCurrentPlayer(currentGame)
	newGame := currentGame
	if !currentPlayer.HumanPlayer {
		newGame = askForCellEnable(newGame)
		// todo play ai
		newGame, _ = CellChange(newGame, newGame.EnablePosition[0])
		return PlayTurn(newGame)
	}
	newGame = askForCellEnable(newGame)

	return newGame, nil
}

func askForCellEnable(game Game) Game{
	currentPlayer := GetCurrentPlayer(game)
	legalCellChanges := board.GetLegalCellChangesForCellType(currentPlayer.CellType, game.Board)

	game.EnablePosition = legalCellChanges
	return game
}

func CellChange(currentGame Game, cellChange cell.Cell) (Game, error) {
	cellChangesFromChoice := append(board.GetFlippedCellsFromCellChange(cellChange, currentGame.Board), cellChange)
	currentGame.Board = board.DrawCells(cellChangesFromChoice, currentGame.Board)

	return SwitchPlayer(currentGame), nil
}

func EventLeft(currentGame Game) (Game, error){
	if currentGame.SelectPosition <= 0 {
		currentGame.SelectPosition = len(currentGame.EnablePosition) - 1
	} else {
		currentGame.SelectPosition--
	}
	return currentGame, nil
}

func EventRight(currentGame Game) (Game, error){
	if len(currentGame.EnablePosition) <= ( currentGame.SelectPosition + 1 ) {
		currentGame.SelectPosition = 0
	} else {
		currentGame.SelectPosition++
	}
	return currentGame, nil
}

func EventEnter(currentGame Game) (Game, error){
	return CellChange(currentGame, currentGame.EnablePosition[currentGame.SelectPosition])
}
