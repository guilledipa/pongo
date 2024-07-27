package main

import (
	"bytes"
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth   = 640
	screenHeight  = 480
	bolaSpeed     = 3
	paletaSpeed   = 6 // 6px per tick (by default Ebiten updates game state @ 60Hz)
	fontSize      = 24
	smallFontSize = fontSize / 2
)

var (
	arcadeFaceSource *text.GoTextFaceSource
)

type Objecto struct {
	X, Y, W, H int
}

type Paleta struct {
	Objecto
}

func (p *Paleta) MoveOnKeyPress() {
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		p.Y += paletaSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		p.Y -= paletaSpeed
	}
}

type Bola struct {
	Objecto
	dxdt, dydt int // Velocidad en x e y per tick
}

func (b *Bola) Move() {
	b.X += b.dxdt
	b.Y += b.dydt
}

type Game struct {
	paleta           Paleta
	bola             Bola
	score, highScore int
}

// Layout controls the size of the window
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

// Draw is called every frame and is used to draw stuff on screen
func (g *Game) Draw(screen *ebiten.Image) {
	// Paleta
	vector.DrawFilledRect(screen,
		float32(g.paleta.X), float32(g.paleta.Y),
		float32(g.paleta.W), float32(g.paleta.H),
		color.White, false,
	)
	// Bola
	vector.DrawFilledRect(screen,
		float32(g.bola.X), float32(g.bola.Y),
		float32(g.bola.W), float32(g.bola.H),
		color.White, false,
	)
	// Score
	score := fmt.Sprintf("Score: %d", g.score)
	op := &text.DrawOptions{}
	op.ColorScale.ScaleWithColor(color.White)
	op.GeoM.Translate(10, 10)
	text.Draw(screen, score, &text.GoTextFace{
		Source: arcadeFaceSource,
		Size:   smallFontSize,
	}, op)
	// High score
	highScore := fmt.Sprintf("High score: %d", g.highScore)
	op = &text.DrawOptions{}
	op.ColorScale.ScaleWithColor(color.White)
	op.GeoM.Translate(10, 30)
	text.Draw(screen, highScore, &text.GoTextFace{
		Source: arcadeFaceSource,
		Size:   smallFontSize,
	}, op)
}

// Update the game state - Called 60Hz
func (g *Game) Update() error {
	g.paleta.MoveOnKeyPress()
	g.bola.Move()
	g.CollideWithWall()
	g.CollideWithPaleta()
	return nil
}

func (g *Game) Reset() {
	g.bola.X = 0
	g.bola.Y = 0
	g.score = 0
}

func (g *Game) CollideWithWall() {
	if g.bola.X >= screenWidth {
		g.Reset()
	}
	if g.bola.X <= 0 {
		g.bola.dxdt = bolaSpeed
	}
	if g.bola.Y <= 0 {
		g.bola.dydt = bolaSpeed
	}
	if g.bola.Y >= screenHeight {
		g.bola.dydt = -bolaSpeed
	}
}

func (g *Game) CollideWithPaleta() {
	if g.bola.X >= g.paleta.X && g.bola.Y >= g.paleta.Y && g.bola.Y <= g.paleta.Y+g.paleta.H {
		g.bola.dxdt = -g.bola.dxdt
		g.score++
		if g.score > g.highScore {
			g.highScore = g.score
		}
	}
}

func init() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.ArcadeN_ttf))
	if err != nil {
		log.Fatal(err)
	}
	arcadeFaceSource = s
}

func main() {
	ebiten.SetWindowTitle("Pongo ~ Pong en Go")
	ebiten.SetWindowSize(screenWidth, screenHeight)

	paleta := Paleta{
		Objecto: Objecto{
			X: 600,
			Y: 200,
			W: 15,
			H: 100,
		},
	}
	bola := Bola{
		Objecto: Objecto{
			X: 600,
			Y: 200,
			W: 15,
			H: 100,
		},
	}
	g := &Game{
		paleta: paleta,
		bola:   bola,
	}
	err := ebiten.RunGame(g)
	if err != nil {
		log.Fatalf("Unable to run game: %v", err)
	}
}
