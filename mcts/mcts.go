package mcts

import (
	"fmt"
	"log"
	"time"
	"math/rand"
	"berzerk/board"
)

func GetPossibleMoves(b board.Board) []board.Move{
	moves := make([]board.Move, 0, 100)
	if b.LastMove.Berserk {
		done := board.Move{
			From: board.Coord{0,0},
			To: board.Coord{0,0},
		}
		moves = append(moves, done)
	}
	for y := range b.Size {
		for x := range b.Size {
			p := b.PieceAtXY(x,y)
			if p == board.Empty {
				continue
			}
			if board.IsAttackerSide(p) && b.Turn == board.TurnDefender {
				continue
			}
			if board.IsDefenderSide(p) && b.Turn == board.TurnAttacker {
				continue
			}
			m := b.GetValidMoves(board.Coord{X:x, Y:y}, b.LastMove.Berserk)
			moves = append(moves, m...)
		}
	}
	return moves
}

func PickMoveCaptureAggressively(r *rand.Rand, moves []board.Move) int {
	caps := make([]int, 0, 10)
	for i, m := range moves {
		if len(m.Captures) > 0 {
			caps = append(caps, i)
		}
	}
	if len(caps) > 0 {
		i := rand.Intn(len(caps))
		choice := caps[i]
		return choice
	} else {
		if len(moves) == 0 {
			return -1
		}
		i := rand.Intn(len(moves))
		return i
	}
}

func PickMoveRandom(r *rand.Rand, moves []board.Move) int {
	i := r.Intn(len(moves))
	return i
}

func Search(b board.Board) int {
	//rand.New(rand.NewSource(123)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i:= 0; i<5000; i++ {
		moves := GetPossibleMoves(b)
		//choice := PickMoveRandom(r, moves)
		choice := PickMoveCaptureAggressively(r, moves)
		if choice < 0 {
			fmt.Println("")
			b.Print()
			log.Fatal("no moves!")
		}

		b.MakeMove(moves[choice])
		if b.IsGameOver() {
			return int(b.Turn) * 2 - 1
		}
	}
	return 0
}

func RateMoves (count int, b board.Board) ([]board.Move, []float32) {
	moves := GetPossibleMoves(b)
	ratings := make([]float32, len(moves))
	for mIdx, move := range moves {
		newb := board.CopyBoard(b)
		newb.MakeMove(move)
		// Play a bunch of games from this position and tally up the
		// wins/losses.
		results := make([]float32, count)
		for i := 0; i<count; i++ {
			results[i] = float32(Search(board.CopyBoard(newb)))
		}
		var sum float32 = 0
		for _, v := range results {
			sum += float32(v)
		}
		fmt.Println(results)
		ratings[mIdx] = sum / float32(count)
		// Make sure 1 always means the best and -1 means the worst.
		if b.Turn == board.TurnDefender {
			ratings[mIdx] *= -1
		}
	}
	return moves, ratings
}

func RateMovesParallel (count int, b board.Board) ([]board.Move, []float32) {
	moves := GetPossibleMoves(b)
	ratings := make([]float32, len(moves))
	for mIdx, move := range moves {
		newb := board.CopyBoard(b)
		newb.MakeMove(move)
		// Play a bunch of games from this position and tally up the
		// wins/losses.
		results := make([]float32, count)
		msgs := make(chan int)
		for i := 0; i<count; i++ {
			go func() { msgs <- Search(board.CopyBoard(newb)) }()
		}
		for i := 0; i<count; i++ {
			results[i] = float32(<- msgs)
		}
		var sum float32 = 0
		for _, v := range results {
			sum += float32(v)
		}
		// fmt.Println(results)
		ratings[mIdx] = sum / float32(count)
		// Make sure 1 always means the best and -1 means the worst.
		if b.Turn == board.TurnDefender {
			ratings[mIdx] *= -1
		}
	}
	return moves, ratings
}
