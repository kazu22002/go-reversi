package main

import (
	"fmt"
	"github.com/kazu22002/go-reversi/reversi/cell"
	"github.com/kazu22002/go-reversi/reversi/game"
	"github.com/kazu22002/go-reversi/reversi/player"
	"github.com/nsf/termbox-go"
	"strings"
	"sync"
)

func main() {
	fmt.Println("\n############# REVERSI #############")

	playerBlack, playerWhite := selectPlay()

	party := playGame(playerBlack, playerWhite )

	fmt.Println("\n############# Result #############")
	fmt.Println("")
	fmt.Println(game.Render(party))
	fmt.Println(result(party))
	fmt.Println("\n############# Thank you #############")
}

func playGame(playerBlack, playerWhite player.Player) game.Game{

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	gameCh := make(chan game.Game)
	keyCh := make(chan termbox.Key)

	go drawLoop(gameCh)
	go keyEventLoop(keyCh)

	party := game.New([]player.Player{playerBlack, playerWhite})

	party = controller(party, gameCh, keyCh)
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	return party
}

func selectPlay() (player.Player, player.Player) {
	var playType string

	fmt.Println("1: Black. human , White. CPU")
	fmt.Println("2: Black. CPU , White. human")
	fmt.Println("3: Black. human , White. human")
	fmt.Print("select : ")
	fmt.Scanf("%s", &playType)

	if playType == "2" {
		return player.New(false, cell.TypeBlack), player.New(true, cell.TypeWhite)
	} else if playType == "3" {
		return player.New(true, cell.TypeBlack), player.New(true, cell.TypeWhite)
	}

	return player.New(true, cell.TypeBlack), player.New(false, cell.TypeWhite)
}

func result(resultGame game.Game) string {
	black, white := game.Result(resultGame)
	return fmt.Sprintf("Black: %d, White: %d", black, white)
}

var mu sync.Mutex

//画面描画
func drawLoop(sch chan game.Game) {
	for {
		st := <-sch
		mu.Lock()
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

		// todo text decoration underline
		//screen := game.Render(st)
		//fmt.Println(screen)
		drawString(game.Render(st))

		termbox.Flush()
		mu.Unlock()
	}
}

func drawString(str string) {
	lines := strings.Split(str, "\n")
	for i := 0; i < len(lines); i++ {
		drawLine(lines[i], i)
	}
}

func drawLine(str string, y int){
	runes := []rune(str)
	for i := 0; i < len(runes); i++ {
		termbox.SetCell(i, y, runes[i], termbox.ColorDefault, termbox.ColorDefault)
	}
}

//ゲームメイン処理
func controller(gameMain game.Game, gameCh chan game.Game, keyCh chan termbox.Key) game.Game {
	gameMain, _ = game.PlayTurn(gameMain)
	gameCh <- gameMain

	for {
		if game.IsFinished(gameMain){
			return gameMain
		}

		select {
		case key := <-keyCh: //キーイベント
			mu.Lock()
			switch key {
			case termbox.KeyEsc, termbox.KeyCtrlC: //ゲーム終了
				mu.Unlock()
				return gameMain
			case termbox.KeyArrowLeft: //ひだり
				gameMain, _ = game.EventLeft(gameMain)
				break
			case termbox.KeyArrowRight: //みぎ
				gameMain, _ = game.EventRight(gameMain)
				break
			case termbox.KeyEnter: //決定
				gameMain, _ = game.EventEnter(gameMain)
				gameMain, _ = game.PlayTurn(gameMain)
				break
			}
			mu.Unlock()

			gameCh <- gameMain
			break
		//case <-stateCh: //ゲームイベント
		//	mu.Lock()
		//	mu.Unlock()
		//	break
		default:
			break
		}
	}
}

//キーイベント
func keyEventLoop(kch chan termbox.Key) {
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			kch <- ev.Key
		default:
		}
	}
}
