package board

type Piece byte

const (
	Empty     = 0
	Defender  = 1
	Knight    = 2
	King      = 3
	Attacker  = 4
	Commander = 5
)

func IsAttackerSide(p Piece) bool {
	return p == Attacker || p == Commander
}

func IsDefenderSide(p Piece) bool {
	return p == Defender || p == Knight || p == King
}

func IsSameSide(a, b Piece) bool {
	return (IsAttackerSide(a) && IsAttackerSide(b) ||
		IsDefenderSide(a) && IsDefenderSide(b))
}

func IsOpposingSide(a, b Piece) bool {
	return (IsAttackerSide(a) && IsDefenderSide(b) ||
		IsDefenderSide(a) && IsAttackerSide(b))
}

type Coord struct {
	X int
	Y int
}

type Capture struct {
	Piece Piece
	Coord Coord
}

type Move struct {
	Piece    Piece
	From     Coord
	To       Coord
	Captures []Capture
	// Is this the beginning or continuation of a chain of berserk
	// moves?
	Berserk  bool
}

type Turn byte

const (
	TurnAttacker = 0
	TurnDefender = 1
)

// The Berzerk rule works by allowing the player to play move than one
// move in a row. So if the last move ended with a capture and there
// are still more berzerk moves that can be chained together, then the
// Turn remains unchanged.
type Board struct {
	Size    int
	Squares []Piece
	LastMove Move
	Turn     Turn
}

func CopyBoard(b Board) Board {
	dest := Board{}
	dest.Size = b.Size
	dest.Squares = make([]Piece, b.Size*b.Size)
	_ = copy(dest.Squares, b.Squares)

	// FIXME: Separate function?
	dest.LastMove.Piece = b.LastMove.Piece
	dest.LastMove.From = b.LastMove.From
	dest.LastMove.To = b.LastMove.To
	dest.LastMove.Captures = make([]Capture, len(b.LastMove.Captures))
	_ = copy(dest.LastMove.Captures, b.LastMove.Captures)
	dest.LastMove.Berserk = b.LastMove.Berserk

	dest.Turn = b.Turn

	return dest
}

func (b *Board) Idx(c Coord) int {
	return c.Y*b.Size + c.X
}

func NewBoard() Board {
	// Berserk Hnefetafl is played on an 11x11 board. In principle
	// other sizes are possible.
	size := 11
	b := Board{}
	b.Size = size
	b.Squares = make([]Piece, size*size)
	//b.LastMove = 
	b.Turn = TurnAttacker

	// Place attacker pieces
	for i := 3; i < size-3; i++ {
		b.Squares[b.Idx(Coord{i, 0})] = Attacker
		b.Squares[b.Idx(Coord{i, size - 1})] = Attacker
		b.Squares[b.Idx(Coord{0, i})] = Attacker
		b.Squares[b.Idx(Coord{size - 1, i})] = Attacker
	}
	b.Squares[b.Idx(Coord{1, size / 2})] = Commander
	b.Squares[b.Idx(Coord{size - 2, size / 2})] = Commander
	b.Squares[b.Idx(Coord{size / 2, 1})] = Commander
	b.Squares[b.Idx(Coord{size / 2, size - 2})] = Commander

	// Place defender pieces
	b.Squares[b.Idx(Coord{size / 2, size / 2})] = King

	b.Squares[b.Idx(Coord{size/2 - 1, size/2 - 1})] = Defender
	b.Squares[b.Idx(Coord{size/2 - 1, size / 2})] = Defender
	b.Squares[b.Idx(Coord{size/2 - 1, size/2 + 1})] = Defender

	b.Squares[b.Idx(Coord{size/2 + 1, size/2 - 1})] = Knight
	b.Squares[b.Idx(Coord{size/2 + 1, size / 2})] = Defender
	b.Squares[b.Idx(Coord{size/2 + 1, size/2 + 1})] = Defender

	b.Squares[b.Idx(Coord{size / 2, size/2 - 1})] = Defender
	b.Squares[b.Idx(Coord{size / 2, size/2 + 1})] = Defender

	b.Squares[b.Idx(Coord{size/2 - 2, size / 2})] = Defender
	b.Squares[b.Idx(Coord{size/2 + 2, size / 2})] = Defender
	b.Squares[b.Idx(Coord{size / 2, size/2 - 2})] = Defender
	b.Squares[b.Idx(Coord{size / 2, size/2 + 2})] = Defender

	return b
}

func (b *Board) IsCornerCoord(c Coord) bool {
	if c.X == 0 && c.Y == 0 {
		return true
	}
	if c.X == 0 && c.Y == b.Size-1 {
		return true
	}
	if c.X == b.Size-1 && c.Y == 0 {
		return true
	}
	if c.X == b.Size-1 && c.Y == b.Size-1 {
		return true
	}
	return false
}

func (b *Board) IsRestrictedCoord(c Coord) bool {
	if b.IsCornerCoord(c) {
		return true
	}
	if c.X == b.Size/2 && c.Y == b.Size/2 {
		return true
	}
	return false
}

func (b *Board) ValidCoord(c Coord) bool {
	return c.X >= 0 && c.X < b.Size && c.Y >= 0 && c.Y < b.Size
}

func (b *Board) PieceAt(c Coord) Piece {
	return b.Squares[c.Y*b.Size+c.X]
}

func (b *Board) PieceAtXY(x, y int) Piece {
	return b.Squares[y*b.Size+x]
}

// Identifies the coordinates of any pieces that would be captured by
// placing an `aggressor` piece at the coordinate.
func (b *Board) GetSandwichCaptures(aggressor Piece, c Coord) []Capture {
	eastCoord := Coord{c.X + 1, c.Y}
	eastCoord2 := Coord{c.X + 2, c.Y}

	westCoord := Coord{c.X - 1, c.Y}
	westCoord2 := Coord{c.X - 2, c.Y}

	southCoord := Coord{c.X, c.Y + 1}
	southCoord2 := Coord{c.X, c.Y + 2}

	northCoord := Coord{c.X, c.Y - 1}
	northCoord2 := Coord{c.X, c.Y - 2}

	coords := [4][2]Coord{
		{eastCoord, eastCoord2},
		{westCoord, westCoord2},
		{southCoord, southCoord2},
		{northCoord, northCoord2},
	}

	captures := make([]Capture, 0, 4)
	for _, pair := range coords {
		if !b.ValidCoord(pair[0]) {
			continue
		}
		if !b.ValidCoord(pair[1]) {
			continue
		}
		middle := b.PieceAt(pair[0])
		other := b.PieceAt(pair[1])

		var isCapture = false
		if IsAttackerSide(aggressor) && IsAttackerSide(other) &&
			(middle == Defender || middle == Knight) {
			isCapture = true
		} else if IsDefenderSide(aggressor) && IsDefenderSide(other) &&
			(middle == Attacker || middle == Commander) {
			isCapture = true
		} else if IsAttackerSide(aggressor) && other == Empty && b.IsRestrictedCoord(pair[1]) &&
			(middle == Defender || middle == Knight) {
			isCapture = true
		} else if IsDefenderSide(aggressor) && other == Empty && b.IsRestrictedCoord(pair[1]) &&
			(middle == Attacker || middle == Commander) {
			isCapture = true
		}

		if isCapture {
			captures = append(captures, Capture{Piece: middle, Coord: pair[0]})
		}
	}

	return captures
}

// Identifies if placing the piece at the `c` Coord captures the king
// in a 4 piece capture.
func (b *Board) GetKingFourWayCapture(aggressor Piece, c Coord) []Capture {
	type capture struct {
		middle Coord
		others [3]Coord
	}

	coords := [4]capture{
		// east
		capture{
			middle: Coord{c.X + 1, c.Y},
			others: [3]Coord{
				Coord{c.X + 2, c.Y},
				Coord{c.X + 1, c.Y + 1},
				Coord{c.X + 1, c.Y - 1},
			},
		},
		// west
		capture{
			middle: Coord{c.X - 1, c.Y},
			others: [3]Coord{
				Coord{c.X - 2, c.Y},
				Coord{c.X - 1, c.Y + 1},
				Coord{c.X - 1, c.Y - 1},
			},
		},
		// south
		capture{
			middle: Coord{c.X, c.Y + 1},
			others: [3]Coord{
				Coord{c.X, c.Y + 2},
				Coord{c.X + 1, c.Y + 1},
				Coord{c.X - 1, c.Y + 1},
			},
		},
		// north
		capture{
			middle: Coord{c.X, c.Y - 1},
			others: [3]Coord{
				Coord{c.X, c.Y - 2},
				Coord{c.X + 1, c.Y - 1},
				Coord{c.X - 1, c.Y - 1},
			},
		},
	}

	captures := make([]Capture, 0, 4)
    if !IsAttackerSide(aggressor) {
		return captures
	}
outerLoop:
	for _, cap := range coords {
		if !b.ValidCoord(cap.middle) {
			continue
		}
		middle := b.PieceAt(cap.middle)
		if middle != King {
			continue
		}
		for i := range(3) {
			if !b.ValidCoord(cap.others[i]) {
				continue outerLoop
			}
			o := b.PieceAt(cap.others[i])
			if o == Empty && !b.IsRestrictedCoord(cap.others[i]) {
				continue outerLoop
			}
			if !IsAttackerSide(o) {
				continue outerLoop
			}
		}

		// The king is captured!!
		captures = append(captures, Capture{Piece: middle, Coord: cap.middle})
	}

	return captures
}

// Identifies if placing the piece at the `c` Coord captures the king
// in a 2 piece capture. The King can be captured between 2 Commanders
// or a Commander and an empty restricted square.
func (b *Board) GetKingTwoWayCapture(aggressor Piece, c Coord) []Capture {
	type capture struct {
		middle Coord
		other Coord
	}

	coords := [4]capture{
		// east
		capture{
			middle: Coord{c.X + 1, c.Y},
			other: Coord{c.X + 2, c.Y},
		},
		// west
		capture{
			middle: Coord{c.X - 1, c.Y},
			other: Coord{c.X - 2, c.Y},
		},
		// south
		capture{
			middle: Coord{c.X, c.Y + 1},
			other: Coord{c.X, c.Y + 2},
		},
		// north
		capture{
			middle: Coord{c.X, c.Y - 1},
			other: Coord{c.X, c.Y - 2},
		},
	}

	captures := make([]Capture, 0, 4)
    if aggressor != Commander {
		return captures
	}
	for _, cap := range coords {
		if !b.ValidCoord(cap.middle) {
			continue
		}
		middle := b.PieceAt(cap.middle)
		if middle != King {
			continue
		}
		if !b.ValidCoord(cap.other) {
			continue
		}
		o := b.PieceAt(cap.other)
		if o == Commander || (o == Empty && b.IsRestrictedCoord(cap.other)) {
			// The king is captured!!
			captures = append(captures, Capture{Piece: middle, Coord: cap.middle})
		}
	}

	return captures
}

func (b *Board) GetAllCaptures(aggressor Piece, c Coord) []Capture {
	c1 := b.GetSandwichCaptures(aggressor, c)
	c2 := b.GetKingTwoWayCapture(aggressor, c)
	c3 := b.GetKingFourWayCapture(aggressor, c)

	return append(append(c1, c2...), c3...)
}

func (b *Board) GetValidMovesDefender(aggressor Piece, c Coord) []Move {
	moves := make([]Move, 0, b.Size*b.Size)
	// west of
	for i := c.X + 1; i < b.Size; i++ {
		dest := Coord{i, c.Y}
		if b.IsRestrictedCoord(dest) {
			continue
		}
		destPiece := b.PieceAt(dest)
		if destPiece != Empty {
			break
		}
		m := Move{
			Piece:    aggressor,
			From:     c,
			To:       dest,
			Captures: b.GetSandwichCaptures(aggressor, dest),
		}
		moves = append(moves, m)
	}

	// east of
	for i := c.X - 1; i >= 0; i-- {
		dest := Coord{i, c.Y}
		if b.IsRestrictedCoord(dest) {
			continue
		}
		destPiece := b.PieceAt(dest)
		if destPiece != Empty {
			break
		}
		m := Move{
			Piece:    aggressor,
			From:     c,
			To:       dest,
			Captures: b.GetSandwichCaptures(aggressor, dest),
		}
		moves = append(moves, m)
	}

	// south of
	for i := c.Y + 1; i < b.Size; i++ {
		dest := Coord{c.X, i}
		if b.IsRestrictedCoord(dest) {
			continue
		}
		destPiece := b.PieceAt(dest)
		if destPiece != Empty {
			break
		}
		m := Move{
			Piece:    aggressor,
			From:     c,
			To:       dest,
			Captures: b.GetSandwichCaptures(aggressor, dest),
		}
		moves = append(moves, m)
	}

	// north of
	for i := c.Y - 1; i >= 0; i-- {
		dest := Coord{c.X, i}
		if b.IsRestrictedCoord(dest) {
			continue
		}
		destPiece := b.PieceAt(dest)
		if destPiece != Empty {
			break
		}
		m := Move{
			Piece:    aggressor,
			From:     c,
			To:       dest,
			Captures: b.GetSandwichCaptures(aggressor, dest),
		}
		moves = append(moves, m)
	}

	return moves
}

func (b *Board) GetValidMovesKnight(c Coord) []Move {
	// A Knight moves like a basic defender.
	moves := b.GetValidMovesDefender(Knight, c)

	// But it can also jump over basic attacker pieces.
	eastCoord := Coord{c.X + 1, c.Y}
	eastCoord2 := Coord{c.X + 2, c.Y}
	westCoord := Coord{c.X - 1, c.Y}
	westCoord2 := Coord{c.X - 2, c.Y}
	southCoord := Coord{c.X, c.Y + 1}
	southCoord2 := Coord{c.X, c.Y + 2}
	northCoord := Coord{c.X, c.Y - 1}
	northCoord2 := Coord{c.X, c.Y - 2}

	jumps := [4][2]Coord{
		{eastCoord, eastCoord2},
		{westCoord, westCoord2},
		{southCoord, southCoord2},
		{northCoord, northCoord2},
	}

	for _, j := range jumps {
		over := j[0]
		land := j[1]
		if b.ValidCoord(over) && b.ValidCoord(land) && b.PieceAt(over) == Attacker && b.PieceAt(land) == Empty && !b.IsRestrictedCoord(land) {
			// It's possible to land after jumping and capture pieces in
			// the regular sandwhich manner.
			caps := b.GetSandwichCaptures(Knight, land)
			caps = append(caps, Capture{Piece: b.PieceAt(over), Coord:over})
			m := Move{
				Piece:    Knight,
				From:     c,
				To:       land,
				Captures: caps,
			}
			moves = append(moves, m)
		}
	}

	return moves
}

func (b *Board) GetValidMovesKing(c Coord) []Move {
	moves := make([]Move, 0, b.Size*b.Size)

	// Unlike all other pieces, A King can occupy *any* square on the
	// board.
	// west of
	for i := c.X + 1; i < b.Size; i++ {
		dest := Coord{i, c.Y}
		destPiece := b.PieceAt(dest)
		if destPiece != Empty {
			break
		}
		m := Move{
			Piece:    King,
			From:     c,
			To:       dest,
			Captures: b.GetSandwichCaptures(King, dest),
		}
		moves = append(moves, m)
	}

	// east of
	for i := c.X - 1; i >= 0; i-- {
		dest := Coord{i, c.Y}
		destPiece := b.PieceAt(dest)
		if destPiece != Empty {
			break
		}
		m := Move{
			Piece:    King,
			From:     c,
			To:       dest,
			Captures: b.GetSandwichCaptures(King, dest),
		}
		moves = append(moves, m)
	}

	// south of
	for i := c.Y + 1; i < b.Size; i++ {
		dest := Coord{c.X, i}
		destPiece := b.PieceAt(dest)
		if destPiece != Empty {
			break
		}
		m := Move{
			Piece:    King,
			From:     c,
			To:       dest,
			Captures: b.GetSandwichCaptures(King, dest),
		}
		moves = append(moves, m)
	}

	// north of
	for i := c.Y - 1; i >= 0; i-- {
		dest := Coord{c.X, i}
		destPiece := b.PieceAt(dest)
		if destPiece != Empty {
			break
		}
		m := Move{
			Piece:    King,
			From:     c,
			To:       dest,
			Captures: b.GetSandwichCaptures(King, dest),
		}
		moves = append(moves, m)
	}

	// Kings can jump but don't capture the jumped piece.
	eastCoord := Coord{c.X + 1, c.Y}
	eastCoord2 := Coord{c.X + 2, c.Y}
	westCoord := Coord{c.X - 1, c.Y}
	westCoord2 := Coord{c.X - 2, c.Y}
	southCoord := Coord{c.X, c.Y + 1}
	southCoord2 := Coord{c.X, c.Y + 2}
	northCoord := Coord{c.X, c.Y - 1}
	northCoord2 := Coord{c.X, c.Y - 2}
	jumps := [4][2]Coord{
		{eastCoord, eastCoord2},
		{westCoord, westCoord2},
		{southCoord, southCoord2},
		{northCoord, northCoord2},
	}

	for _, j := range jumps {
		over := j[0]
		land := j[1]
		if b.ValidCoord(over) && b.ValidCoord(land) && b.PieceAt(over) == Attacker && b.PieceAt(land) == Empty && (b.IsRestrictedCoord(c) || b.IsRestrictedCoord(land)){
			// It's possible to land after jumping and capture pieces in
			// the regular sandwhich manner.
			caps := b.GetSandwichCaptures(King, land)
			m := Move{
				Piece:    King,
				From:     c,
				To:       land,
				Captures: caps,
			}
			moves = append(moves, m)
		}
	}

	return moves
}

func (b *Board) GetValidMovesAttacker(aggressor Piece, c Coord) []Move {
	moves := make([]Move, 0, b.Size*b.Size)
	// west of
	for i := c.X + 1; i < b.Size; i++ {
		dest := Coord{i, c.Y}
		if b.IsRestrictedCoord(dest) {
			continue
		}
		destPiece := b.PieceAt(dest)
		if destPiece != Empty {
			break
		}
		m := Move{
			Piece:    aggressor,
			From:     c,
			To:       dest,
			Captures: b.GetAllCaptures(aggressor, dest),
		}
		moves = append(moves, m)
	}

	// east of
	for i := c.X - 1; i >= 0; i-- {
		dest := Coord{i, c.Y}
		if b.IsRestrictedCoord(dest) {
			continue
		}
		destPiece := b.PieceAt(dest)
		if destPiece != Empty {
			break
		}
		m := Move{
			Piece:    aggressor,
			From:     c,
			To:       dest,
			Captures: b.GetAllCaptures(aggressor, dest),
		}
		moves = append(moves, m)
	}

	// south of
	for i := c.Y + 1; i < b.Size; i++ {
		dest := Coord{c.X, i}
		if b.IsRestrictedCoord(dest) {
			continue
		}
		destPiece := b.PieceAt(dest)
		if destPiece != Empty {
			break
		}
		m := Move{
			Piece:    aggressor,
			From:     c,
			To:       dest,
			Captures: b.GetAllCaptures(aggressor, dest),
		}
		moves = append(moves, m)
	}

	// north of
	for i := c.Y - 1; i >= 0; i-- {
		dest := Coord{c.X, i}
		if b.IsRestrictedCoord(dest) {
			continue
		}
		destPiece := b.PieceAt(dest)
		if destPiece != Empty {
			break
		}
		m := Move{
			Piece:    aggressor,
			From:     c,
			To:       dest,
			Captures: b.GetAllCaptures(aggressor, dest),
		}
		moves = append(moves, m)
	}

	return moves
}

func (b *Board) GetValidMovesCommander(c Coord) []Move {
	moves := make([]Move, 0, b.Size*b.Size)
	// A Commander moves like a basic attacker
	moves = append(moves, b.GetValidMovesAttacker(Commander, c)...)

	// But it can also jump over basic defender pieces.
	type jump struct {
		over Coord
		land Coord
	}

	// Commanders can jump but don't capture the jumped piece.
	jumps := [4]jump{
		// east
		jump{
			over: Coord{c.X + 1, c.Y},
			land: Coord{c.X + 2, c.Y},
		},
		// west
		jump{
			Coord{c.X - 1, c.Y},
			Coord{c.X - 2, c.Y},
		},
		// south
		jump{
			Coord{c.X, c.Y + 1},
			Coord{c.X, c.Y + 2},
		},
		// north
		jump{
			Coord{c.X, c.Y - 1},
			Coord{c.X, c.Y - 2},
		},
	}

	for _, j := range jumps {
		if b.ValidCoord(j.over) && b.ValidCoord(j.land) && b.PieceAt(j.over) == Defender && b.PieceAt(j.land) == Empty && !b.IsRestrictedCoord(j.land) {
			// It's possible to land after jumping and capture pieces
			// in the regular sandwhich manner.
			caps := b.GetAllCaptures(Commander, j.land)
			m := Move{
				Piece:    Commander,
				From:     c,
				To:       j.land,
				Captures: caps,
			}
			moves = append(moves, m)
		}
	}

	return moves
}

// Return valid moves for the piece at `c`
func (b *Board) GetValidMoves(c Coord, berserk bool) []Move {
	piece := b.PieceAt(c)
	if piece == Empty {
		return nil
	}

	// temporarily remove the piece at its current location to prevent
	// the capture logic thinking it can sandwhich capture with
	// itself. I *think* this would only happen for jumping moves.
	i := b.Idx(c)
	tmp := b.Squares[i]
	b.Squares[i] = Empty

	moves := make([]Move, 0, b.Size*b.Size)
	switch piece {
	case Defender:
		moves = append(moves, b.GetValidMovesDefender(Defender, c)...)
	case Knight:
		moves = append(moves, b.GetValidMovesKnight(c)...)
	case King:
		moves = append(moves, b.GetValidMovesKing(c)...)
	case Attacker:
		moves = append(moves, b.GetValidMovesAttacker(Attacker, c)...)
	case Commander:
		moves = append(moves, b.GetValidMovesCommander(c)...)
	}

	// Filter out valid berserk moves.
	// Note: The king may finish a berserk run with a winning move to a corner square.
	if berserk {
		berzerk_moves := make([]Move, 0)
			for _, m := range moves {
				if len(m.Captures) > 0 ||
					(piece == King && b.IsCornerCoord(m.To)) {
					berzerk_moves = append(berzerk_moves, m)
				}
			}
		moves = berzerk_moves
	}

	// Put the piece back.
	b.Squares[i] = tmp

	return moves
}

// Note, this function does not verify that move is valid. It is
// expected that move was picked from the slice returned by
// GetValidMoves. TODO: handle berzerk moves.
func (b *Board) MakeMove(move Move) {
	// Not moving the piece is considered a pass. This is how a player
	// stops a chain of berserk moves.
	if move.From == move.To {
		if !b.LastMove.Berserk {
			// a pass is not allowed!
			return
		}
		// Amend the last move.
		b.LastMove.Berserk = false
		if b.Turn == TurnAttacker {
			b.Turn = TurnDefender
		} else {
			b.Turn = TurnAttacker
		}
		return
	}

	from := b.Idx(move.From)
	to := b.Idx(move.To)
	b.Squares[to] = b.Squares[from]
	b.Squares[from] = Empty

	for _, c := range move.Captures {
		idx := b.Idx(c.Coord)
		b.Squares[idx] = Empty
	}

	// Check whether a berserk follow-up is possible
	if len(move.Captures) > 0 {
		berserk_moves := b.GetValidMoves(move.To, true)
		if len(berserk_moves) > 0 {
			move.Berserk = true
		}
	}

	// If berzerk moves are possible then it's still the current
	// player's turn.
	if !move.Berserk {
		if b.Turn == TurnAttacker {
			b.Turn = TurnDefender
		} else {
			b.Turn = TurnAttacker
		}
	}

	b.LastMove = move
}
