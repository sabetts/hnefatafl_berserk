package main

import (
	"bytes"
	_ "embed"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image"
	_ "image/png"
	"image/color"
	"log"
	//"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"fmt"

	"berzerk/board"
)

// Image assets
var (
	//go:embed ai_assets/throne.png
	Throne_png []byte

	//go:embed ai_assets/defender.png
	Defender_png []byte

	//go:embed ai_assets/knight.png
	Knight_png []byte

	//go:embed ai_assets/king.png
	King_png []byte

	//go:embed ai_assets/attacker.png
	Attacker_png []byte

	//go:embed ai_assets/commander.png
	Commander_png []byte

	//go:embed ai_assets/board_final.png
	Board_png []byte
)

var (
	throneImg    *ebiten.Image
	defenderImg  *ebiten.Image
	knightImg    *ebiten.Image
	kingImg      *ebiten.Image
	attackerImg  *ebiten.Image
	commanderImg *ebiten.Image
	boardImg     *ebiten.Image
)

const screenWidth = 800
const screenHeight = 800

const boardOffsetX = 138
const boardOffsetY = 138

const tileWidth = 106
const tileHeight = 107

func PixelToCoord(b *board.Board, x, y int) (board.Coord, bool) {
	x -= boardOffsetX
	x /= tileWidth
	y -= boardOffsetY
	y /= tileHeight
	if x < 0 || x >= b.Size ||
		y < 0 || y >= b.Size {
		return board.Coord{0, 0}, false
	}
	return board.Coord{x, y}, true
}


type Game struct {
	board           board.Board
	selectionActive bool
	selectedTile    board.Coord
	selectionMoves  []board.Move
}


func (g *Game) Update() error {
	click := inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0)
	if !click {
		return nil
	}

	x, y := ebiten.CursorPosition()

	fmt.Println("click", click, x, y)

	coord, ok := PixelToCoord(&g.board, x, y)
	if !ok {
		return nil
	}

	fmt.Println("coord", coord)

	piece := g.board.PieceAt(coord)
	fmt.Println("piece", piece)

	targetMoveIdx := -1
	if g.selectionActive {
		for i, move := range(g.selectionMoves) {
			if coord == move.To {
				targetMoveIdx = i
				break
			}
		}
	}

	fmt.Println("target", targetMoveIdx)

	// They clicked a piece
	if (g.board.Turn == board.TurnAttacker && board.IsAttackerSide(piece)) ||
		(g.board.Turn == board.TurnDefender && board.IsDefenderSide(piece)) {
		if g.selectionActive && g.selectedTile == coord {
			g.selectionActive = false
		} else {
			g.selectionActive = true
			g.selectedTile = coord
			g.selectionMoves = g.board.GetValidMoves(coord)
		}
	} else if targetMoveIdx >= 0 {
		g.board.MakeMove(g.selectionMoves[targetMoveIdx])
		g.selectionActive = false
	}

	fmt.Println("selection", g.selectionActive, g.selectedTile, len(g.selectionMoves))

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	// op.GeoM.Translate(screenWidth/2, screenHeight/2)
	screen.DrawImage(boardImg, op)
	for y := range g.board.Size {
		// Throne needs to be drawn after rows above it have been
		// drawn because it's a bit tall.
		if y == g.board.Size/2 {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(boardOffsetX, boardOffsetY)
			op.GeoM.Translate(float64(g.board.Size/2*tileWidth), float64(g.board.Size/2*tileHeight))
			op.GeoM.Translate(float64(tileWidth)/2, float64(tileHeight))
			s := throneImg.Bounds().Size()
			op.GeoM.Translate(-float64(s.X)/2, -float64(s.Y))
			screen.DrawImage(throneImg, op)
		}
		for x := range g.board.Size {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(boardOffsetX, boardOffsetY)
			op.GeoM.Translate(float64(x*tileWidth), float64(y*tileHeight))
			op.GeoM.Translate(float64(tileWidth)/2, float64(tileHeight)/2)
			p := g.board.Squares[g.board.Idx(board.Coord{x, y})]
			var img *ebiten.Image
			switch p {
			case board.Empty:
				// Nothing to draw
			case board.Defender:
				img = defenderImg
			case board.Knight:
				img = knightImg
			case board.King:
				img = kingImg
			case board.Attacker:
				img = attackerImg
			case board.Commander:
				img = commanderImg
			default:
				log.Fatal("Unknown piece ", p)
			}
			if img != nil {
				s := img.Bounds().Size()
				op.GeoM.Translate(-float64(s.X)/2, -float64(s.Y)/2)
				screen.DrawImage(img, op)
			}
		}
	}
	// Draw the selection and valid moves
	if g.selectionActive {
		vector.StrokeCircle(
			screen,
			float32(boardOffsetX + tileWidth/2 + g.selectedTile.X * tileWidth),
			float32(boardOffsetY + tileHeight/2 + g.selectedTile.Y * tileHeight),
			tileWidth/2,
			10,
			color.RGBA{0xff, 0xff, 0x00, 0xff},
			true,
		)

		for _, m := range g.selectionMoves {
			var fillColor color.RGBA
			if len(m.Captures) > 0 {
				fillColor = color.RGBA{0xff, 0x00, 0x00, 0xff}
			} else {
				fillColor = color.RGBA{0xff, 0xff, 0x00, 0xff}
			}
			vector.FillCircle(
				screen,
				float32(boardOffsetX + tileWidth/2 + tileWidth * m.To.X),
				float32(boardOffsetY + tileHeight/2 + tileHeight * m.To.Y),
				tileWidth/8,
				fillColor,
				true,
			)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth * 2, screenHeight * 2
}

func loadImages() error {
	img, _, err := image.Decode(bytes.NewReader(Throne_png))
	if err != nil {
		log.Fatal("Throne: ", err)
	}
	throneImg = ebiten.NewImageFromImage(img)

	img, _, err = image.Decode(bytes.NewReader(Defender_png))
	if err != nil {
		return err
	}
	defenderImg = ebiten.NewImageFromImage(img)

	img, _, err = image.Decode(bytes.NewReader(Knight_png))
	if err != nil {
		return err
	}
	knightImg = ebiten.NewImageFromImage(img)

	img, _, err = image.Decode(bytes.NewReader(King_png))
	if err != nil {
		return err
	}
	kingImg = ebiten.NewImageFromImage(img)

	img, _, err = image.Decode(bytes.NewReader(Attacker_png))
	if err != nil {
		return err
	}
	attackerImg = ebiten.NewImageFromImage(img)

	img, _, err = image.Decode(bytes.NewReader(Commander_png))
	if err != nil {
		return err
	}
	commanderImg = ebiten.NewImageFromImage(img)

	img, _, err = image.Decode(bytes.NewReader(Board_png))
	if err != nil {
		return err
	}
	boardImg = ebiten.NewImageFromImage(img)

	return nil
}

func main() {
	err := loadImages()
	if err != nil {
		log.Fatal(err)
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)

	game := &Game{}
	game.board = board.NewBoard()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
