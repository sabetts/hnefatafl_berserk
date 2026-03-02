package main

import (
	"bytes"
	_ "embed"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
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

	//go:embed ai_assets/blood.png
	Blood_png []byte

	//go:embed ai_assets/oldblood.png
	OldBlood_png []byte

	//go:embed ai_assets/button.png
	Button_png []byte

)

var (
	throneImg    *ebiten.Image
	defenderImg  *ebiten.Image
	knightImg    *ebiten.Image
	kingImg      *ebiten.Image
	attackerImg  *ebiten.Image
	commanderImg *ebiten.Image
	boardImg     *ebiten.Image
	bloodImg     *ebiten.Image
	oldBloodImg  *ebiten.Image

	// buttons
	buttonImg  *ebiten.Image

	// Font
	mplusFaceSource *text.GoTextFaceSource
	mplusNormalFace *text.GoTextFace
	mplusBigFace    *text.GoTextFace
)

func GetImgForPiece(p board.Piece) *ebiten.Image {
	switch p {
	case board.Empty:
		// Nothing to draw
		return nil
	case board.Defender:
		return defenderImg
	case board.Knight:
		return knightImg
	case board.King:
		return kingImg
	case board.Attacker:
		return attackerImg
	case board.Commander:
		return commanderImg
	default:
		log.Fatal("Unknown piece ", p)
	}
	return nil
}

const screenWidth = 1024
const screenHeight = 768

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

type UIElement struct {
	x,y int
	w,h int
	hovering bool
	mousedown bool
	clicked bool
	mx, my int
}

func (e *UIElement) IsInside(x, y int) bool {
	return x >= e.x && y >= e.y && x < e.x+e.w && y < e.y+e.h
}

func (e *UIElement) UpdateMouse(pressed, click bool, x, y int) {
	if e.IsInside(x,y) {
		e.clicked = click
		e.hovering = !pressed
		e.mousedown = pressed
		e.mx = x - e.x
		e.my = y - e.y
	} else {
		e.hovering = false
		e.mousedown = false
		e.clicked = false
	}
}

type Button struct {
	UIElement
	label string
}

type Updateable interface {
	Update()
}

type Drawable interface {
	Draw(screen *ebiten.Image)
}

func (b *Button) Draw(screen *ebiten.Image) {
	var bg color.RGBA
	var fg color.RGBA
	if b.mousedown {
		bg = color.RGBA{0x99,0x99,0x99,0xff}
		fg = color.RGBA{0xFF,0xFF,0xFF,0xff}
	} else if b.hovering {
		bg = color.RGBA{0x99,0x99,0x99,0xff}
		fg = color.RGBA{0xCC,0xCC,0xCC,0xff}
	} else {
		bg = color.RGBA{0x99,0x99,0x99,0xff}
		fg = color.RGBA{0x00,0x00,0x00,0xff}
	}

	vector.FillRect(screen, float32(b.x), float32(b.y), float32(b.w), float32(b.h), bg, false)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(b.x), float64(b.y))
	screen.DrawImage(buttonImg, op)

	top := &text.DrawOptions{}
	top.PrimaryAlign = text.AlignCenter
	top.SecondaryAlign = text.AlignCenter
	top.GeoM.Translate(float64(b.x+b.w/2), float64(b.y + b.h/2))
	top.ColorScale.ScaleWithColor(fg)
	text.Draw(screen, b.label, mplusBigFace, top)
}

type Game struct {
	UIElement
	board           board.Board
	history         []board.Board
	future          []board.Board
	selectionActive bool
	selectedTile    board.Coord
	selectionMoves  []board.Move
	allCaptures     []board.Capture
	// piece movement
	movingActive bool
	movingPercent float32
}

func NewGame() *Game {
	g := &Game{}
	g.board = board.NewBoard()
	return g
}

func (g *Game) HandleDone() {
	if !g.board.LastMove.Berserk {
		return
	}

	g.history = append(g.history, board.CopyBoard(g.board))
	g.board.MakeMove(board.Move{
		From: board.Coord{0,0},
		To: board.Coord{0,0},
	})
	g.selectionActive = false
	g.movingActive = false
}

func (g *Game) HandleUndo() {
	if len(g.history) <= 0 {
		return
	}

	g.future = append(g.future, g.board)
	g.board = g.history[len(g.history)-1]
	g.history = g.history[0:len(g.history)-1]

	if g.board.LastMove.Berserk {
		g.selectionActive = true
		g.selectedTile = g.board.LastMove.To
		g.selectionMoves = g.board.GetValidMoves(g.selectedTile, g.board.LastMove.Berserk)
	} else {
		g.selectionActive = false
	}
	g.movingActive = false
}

func (g *Game) HandleClick() {
	if !g.clicked {
		return
	}
	fmt.Println("click", g.clicked, g.mx, g.my)

	coord, ok := PixelToCoord(&g.board, g.mx, g.my)
	if !ok {
		return
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
			g.selectionMoves = g.board.GetValidMoves(coord, g.board.LastMove.Berserk)
		}
	} else if targetMoveIdx >= 0 {
		g.allCaptures = append(g.allCaptures, g.board.LastMove.Captures...)
		g.history = append(g.history, board.CopyBoard(g.board))
		clear(g.future)
		g.board.MakeMove(g.selectionMoves[targetMoveIdx])
		if g.board.LastMove.Berserk {
			g.selectedTile = g.board.LastMove.To
			g.selectionMoves = g.board.GetValidMoves(g.selectedTile, g.board.LastMove.Berserk)
		} else {
			g.selectionActive = false
		}
		g.movingActive = true
		g.movingPercent = 0
	}

	// blood stain for captures. but they fade into faded blood stains
	// that last hte rest of the game. by the end the board looks
	// gruisomely bloody.


	fmt.Println("selection", g.selectionActive, g.selectedTile, len(g.selectionMoves))
}

func (g *Game) Update() {
	g.HandleClick()

	if g.movingActive {
		g.movingPercent += 0.1
		if g.movingPercent >= 1 {
			g.movingActive = false
		}
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	//bs := boardImg.Bounds().Size()
	{
		// op.GeoM.Translate(screenWidth/2, screenHeight/2)
		op := &ebiten.DrawImageOptions{}
		screen.DrawImage(boardImg, op)
	}

	// Draw the old captures
	for _, c := range g.allCaptures {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(boardOffsetX, boardOffsetY)
		op.GeoM.Translate(float64(c.Coord.X*tileWidth), float64(c.Coord.Y*tileHeight))
		op.GeoM.Translate(float64(tileWidth)/2, float64(tileHeight)/2)
		s := oldBloodImg.Bounds().Size()
		op.GeoM.Translate(-float64(s.X)/2, -float64(s.Y)/2)
		screen.DrawImage(oldBloodImg, op)
	}

	// Draw the most recent captures
	for _, c := range g.board.LastMove.Captures {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(boardOffsetX, boardOffsetY)
		op.GeoM.Translate(float64(c.Coord.X*tileWidth), float64(c.Coord.Y*tileHeight))
		op.GeoM.Translate(float64(tileWidth)/2, float64(tileHeight)/2)
		var img *ebiten.Image
		if g.movingActive {
			img = GetImgForPiece(c.Piece)
		} else {
			img = bloodImg
		}
		s := img.Bounds().Size()
		op.GeoM.Translate(-float64(s.X)/2, -float64(s.Y)/2)
		screen.DrawImage(img, op)
	}

	// Now draw all the stationary pieces
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
			// Skip the moving piece.
			if g.movingActive && g.board.LastMove.To.X == x && g.board.LastMove.To.Y == y {
				continue
			}
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(boardOffsetX, boardOffsetY)
			op.GeoM.Translate(float64(x*tileWidth), float64(y*tileHeight))
			op.GeoM.Translate(float64(tileWidth)/2, float64(tileHeight)/2)
			p := g.board.Squares[g.board.Idx(board.Coord{x, y})]
			img := GetImgForPiece(p)
			if img != nil {
				s := img.Bounds().Size()
				op.GeoM.Translate(-float64(s.X)/2, -float64(s.Y)/2)
				screen.DrawImage(img, op)
			}
		}
	}

	// Draw the moving piece
	if  g.movingActive {
		op := &ebiten.DrawImageOptions{}
		img := GetImgForPiece(g.board.LastMove.Piece)

		dx := g.movingPercent * float32((g.board.LastMove.To.X - g.board.LastMove.From.X) * tileWidth)
		dy := g.movingPercent * float32((g.board.LastMove.To.Y - g.board.LastMove.From.Y) * tileHeight)

		op.GeoM.Translate(boardOffsetX, boardOffsetY)
		op.GeoM.Translate(float64(g.board.LastMove.From.X*tileWidth), float64(g.board.LastMove.From.Y*tileHeight))
		op.GeoM.Translate(float64(tileWidth)/2, float64(tileHeight)/2)
		op.GeoM.Translate(float64(dx), float64(dy))
		s := img.Bounds().Size()
		op.GeoM.Translate(-float64(s.X)/2, -float64(s.Y)/2)
		screen.DrawImage(img, op)
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

type UI struct {
	Game *Game
	BackupGame *Game
	Done Button
	Undo Button
	Restart Button
	Quit Button
}


func (ui *UI) Update() error {
	pressed := ebiten.IsMouseButtonPressed(ebiten.MouseButton0)
	click := inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0)
	x, y := ebiten.CursorPosition()

	ui.Game.UpdateMouse(pressed, click, x,y)
	ui.Done.UpdateMouse(pressed, click, x,y)
	ui.Undo.UpdateMouse(pressed, click, x,y)
	ui.Restart.UpdateMouse(pressed, click, x,y)
	ui.Quit.UpdateMouse(pressed, click, x,y)

	ui.Game.Update()

	if ui.Done.clicked {
		ui.Game.HandleDone()
	}

	if ui.Undo.clicked {
		// Special-case undo the game restart.
		if len(ui.Game.history) == 0 && len(ui.Game.future) == 0 {
			ui.Game = ui.BackupGame
			ui.BackupGame = nil
		} else {
			ui.Game.HandleUndo()
		}
	}

	if ui.Restart.clicked {
		ui.BackupGame = ui.Game
		ui.Game = NewGame()
		// FIXME: bleh
		boardSize := boardImg.Bounds().Size()
		ui.Game.x = 0
		ui.Game.y = 0
		ui.Game.w = boardSize.X
		ui.Game.h = boardSize.Y
	}

	if ui.Quit.clicked {
		return ebiten.Termination
	}

	return nil
}

func (ui *UI) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x80, 0x80, 0xc0, 0xff})
	vector.FillRect(screen, 2000, 100, 500, 200, color.White, false)

	ui.Game.Draw(screen)

	ui.Done.Draw(screen)
	ui.Undo.Draw(screen)
	ui.Restart.Draw(screen)
	ui.Quit.Draw(screen)
}

func (ui *UI) Layout(outsideWidth, outsideHeight int) (int, int) {
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

	img, _, err = image.Decode(bytes.NewReader(Blood_png))
	if err != nil {
		return err
	}
	bloodImg = ebiten.NewImageFromImage(img)

	img, _, err = image.Decode(bytes.NewReader(OldBlood_png))
	if err != nil {
		return err
	}
	oldBloodImg = ebiten.NewImageFromImage(img)

	img, _, err = image.Decode(bytes.NewReader(Button_png))
	if err != nil {
		return err
	}
	buttonImg = ebiten.NewImageFromImage(img)

	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		return err
	}
	mplusFaceSource = s

	mplusNormalFace = &text.GoTextFace{
		Source: mplusFaceSource,
		Size:   24,
	}
	mplusBigFace = &text.GoTextFace{
		Source: mplusFaceSource,
		Size:   32,
	}

	return nil
}

func NewUI() *UI {
	ui := &UI{}
	ui.Game = NewGame()

	boardSize := boardImg.Bounds().Size()
	buttonSize := buttonImg.Bounds().Size()
	pad := 5

	ui.Game.x = 0
	ui.Game.y = 0
	ui.Game.w = boardSize.X
	ui.Game.h = boardSize.Y

	ui.Done = Button{
		UIElement: UIElement{
			x: boardSize.X,
			y: 0,
			w: buttonSize.X,
			h: buttonSize.Y,
		},
		label: "Done",
	}
	ui.Undo = Button{
		UIElement: UIElement{
			x: boardSize.X,
			y: (buttonSize.Y + pad)*1,
			w: buttonSize.X,
			h: buttonSize.Y,
		},
		label: "Undo",
	}
	ui.Restart = Button{
		UIElement: UIElement{
			x: boardSize.X,
			y: (buttonSize.Y+pad)*2,
			w: buttonSize.X,
			h: buttonSize.Y,
		},
		label: "Restart",
	}
	ui.Quit = Button{
		UIElement: UIElement{
			x: boardSize.X,
			y: (buttonSize.Y+pad)*3,
			w: buttonSize.X,
			h: buttonSize.Y,
		},
		label: "Quit",
	}

	return ui
}

func main() {
	err := loadImages()
	if err != nil {
		log.Fatal(err)
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)

	ui := NewUI()

	if err := ebiten.RunGame(ui); err != nil {
		log.Fatal(err)
	}
}
