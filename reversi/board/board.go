package board

import (
	"errors"
	"github.com/kazu22002/go-reversi/reversi/cell"
	"github.com/kazu22002/go-reversi/reversi/matrix"
	"github.com/kazu22002/go-reversi/reversi/vector"
)

type Board [][]uint8

func New(xSize uint8, ySize uint8) Board {
	board := Board{}
	var y uint8
	for y = 0; y < ySize; y++ {
		board = append(board, make([]uint8, xSize, xSize))
	}
	return board
}

func IsValidBoardSize(xSize int, ySize int) bool {
	return xSize%2 == 0 && ySize%2 == 0
}

func InitCells(board Board) (Board, error) {
	xSize, ySize := matrix.GetSize(board)
	if !IsValidBoardSize(xSize, ySize) {
		return board, errors.New("Invalid board Size, x/y dim must be even to place departure cells")
	}
	return DrawCells(GetDepartureCells(board), board), nil
}

func GetDepartureCells(board Board) []cell.Cell {

	xSize, ySize := matrix.GetSize(board)

	xMiddle := uint8(xSize / 2)
	yMiddle := uint8(ySize / 2)

	return []cell.Cell{
		cell.New(xMiddle, yMiddle, cell.TypeBlack),
		cell.New(xMiddle-1, yMiddle-1, cell.TypeBlack),
		cell.New(xMiddle-1, yMiddle, cell.TypeWhite),
		cell.New(xMiddle, yMiddle-1, cell.TypeWhite),
	}

}

func Render(board Board, cellProposals cell.Cell) string {

	renderMatrix := [][]string{}

	for yPos, row := range board {
		renderMatrix = append(renderMatrix, make([]string, len(row)))
		for xPos, cellType := range row {
			_, proposalCellIdx := FindCell(uint8(xPos), uint8(yPos), cellProposals)
			if proposalCellIdx != -1 {
				renderMatrix[yPos][xPos] = cell.GetSymbol(cellProposals.CellType)
			} else {
				renderMatrix[yPos][xPos] = cell.GetSymbol(cellType)
			}
		}
	}

	return matrix.Render(renderMatrix)
}

func Result(board Board) (int, int){
	black := 0
	white := 0
	for _, row := range board {
		for _, cellType := range row {
			if cellType == cell.TypeBlack {
				black++
			} else if cellType == cell.TypeWhite {
				white++
			}
		}
	}
	return black, white
}

func IsFull(board Board) bool {
	for _, ySlice := range board {
		for _, cellType := range ySlice {
			if cellType == cell.TypeEmpty {
				return false
			}
		}
	}
	return true
}

func DrawCells(cells []cell.Cell, board Board) Board {
	newBoard := board
	for _, cell := range cells {
		newBoard[cell.Y][cell.X] = cell.CellType
	}
	return newBoard
}

func GetCellType(xPos uint8, yPos uint8, board Board) uint8 {
	if !(uint8(len(board)-1) >= yPos && uint8(len(board[yPos])-1) >= xPos) {
		return cell.TypeEmpty
	}
	return board[yPos][xPos]
}

func GetFlippedCellsFromCellChange(cellChange cell.Cell, board Board) []cell.Cell {

	if GetCellType(cellChange.X, cellChange.Y, board) != cell.TypeEmpty {
		return []cell.Cell{}
	}

	var flippedCells []cell.Cell

	for _, directionnalVector := range vector.GetDirectionnalVectors() {
		flippedInDirection := GetFlippedCellsForCellChangeAndDirectionVector(cellChange, directionnalVector, board)
		flippedCells = append(flippedCells, flippedInDirection...)
	}

	return flippedCells

}

func GetFlippedCellsForCellChangeAndDirectionVector(cellChange cell.Cell, directionVector vector.Vector, board Board) []cell.Cell {

	flippedCells := []cell.Cell{}

	var localCellType uint8
	localCellPosition := vector.Vector{int(cellChange.X), int(cellChange.Y)}
	reverseCellType := cell.GetReverseCellType(cellChange.CellType)

	for {
		localCellPosition = vector.VectorAdd(localCellPosition, directionVector)
		localCellType = GetCellType(uint8(localCellPosition.X), uint8(localCellPosition.Y), board)
		if localCellType != reverseCellType {
			break
		}
		flippedCell := cell.New(uint8(localCellPosition.X), uint8(localCellPosition.Y), cellChange.CellType)
		flippedCells = append(flippedCells, flippedCell)
	}

	if localCellType == cellChange.CellType && len(flippedCells) > 0 {
		return flippedCells
	}

	return []cell.Cell{}

}

func IsLegalCellChange(cellChange cell.Cell, board Board) bool {
	return len(GetFlippedCellsFromCellChange(cellChange, board)) > 0
}

func GetLegalCellChangesForCellType(cellType uint8, board Board) []cell.Cell {

	legalCellChanges := []cell.Cell{}

	for y, row := range board {
		for x, _ := range row {
			cellChange := cell.Cell{uint8(x), uint8(y), cellType}
			if IsLegalCellChange(cellChange, board) {
				legalCellChanges = append(legalCellChanges, cellChange)
			}
		}
	}

	return legalCellChanges

}

func GetCellDistribution(board Board) map[uint8]uint8 {
	dist := map[uint8]uint8{cell.TypeEmpty: uint8(0), cell.TypeBlack: uint8(0), cell.TypeWhite: uint8(0)}
	for _, row := range board {
		for _, cellType := range row {
			dist[cellType]++
		}
	}
	return dist
}

func FindCell(x uint8, y uint8, cellSelect cell.Cell) (cell.Cell, int) {
	if cellSelect.X == x && cellSelect.Y == y {
		return cellSelect, 0
	}
	return cell.Cell{}, -1
}
