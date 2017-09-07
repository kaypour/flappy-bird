package main

import (
	"fmt"
	"path"
	"runtime"
	"strconv"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	lingerOnDeath = 2 * time.Second
)

type scene struct {
	renderer *sdl.Renderer
	bg       *sdl.Texture
	font     *ttf.Font
	music    *mix.Music
	bird     *bird
	pipe     *pipe
	frame    uint64
	score    float64
}

func newScene(r *sdl.Renderer, assetDir string) (*scene, error) {
	bg, err := img.LoadTexture(r, path.Join(assetDir, "background/background.png"))
	if err != nil {
		return nil, fmt.Errorf("could not load background texture: %v", err)
	}
	font, err := ttf.OpenFont(path.Join(assetDir, "font/font.ttf"), 40)
	if err != nil {
		return nil, fmt.Errorf("could not load font texture: %v", err)
	}
	music, err := mix.LoadMUS(path.Join(assetDir, "sound/intro.ogg"))
	if err != nil {
		return nil, fmt.Errorf("could not load music file: %v", err)
	}
	bird, err := newBird(r, assetDir)
	if err != nil {
		return nil, err
	}
	pipe, err := newPipe(r, assetDir)
	if err != nil {
		return nil, err
	}

	return &scene{
		renderer: r,
		bg:       bg,
		font:     font,
		music:    music,
		bird:     bird,
		pipe:     pipe,
	}, nil
}

func (s *scene) run() error {
	if err := s.music.Play(-1); err != nil {
		return fmt.Errorf("could not play music: %v", err)
	}
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	events := make(chan sdl.Event)
	out := make(chan error)
	go func() {
		tick := time.Tick(10 * time.Millisecond)
		for {
			select {
			case event := <-events:
				switch event.(type) {
				case *sdl.QuitEvent:
					out <- nil
					return
				default:
					s.bird.handleKeyEvent(event)
				}
			case <-tick:
				if err := s.draw(); err != nil {
					out <- err
					return
				}
				if s.bird.dead {
					if s.bird.deadSince() > lingerOnDeath {
						s.score = 0
						s.bird.restart()
						s.pipe.restart()
					}
					continue
				}
				s.update()
			}
		}
	}()
	for {
		select {
		case err := <-out:
			return err
		case events <- sdl.WaitEventTimeout(10):
		}
	}
}

func (s *scene) draw() error {
	if err := s.renderer.Clear(); err != nil {
		return fmt.Errorf("could not clear renderer: %v", err)
	}
	if err := s.renderer.Copy(s.bg, nil, nil); err != nil {
		return fmt.Errorf("could not copy background: %v", err)
	}
	if err := s.bird.draw(s.renderer); err != nil {
		return err
	}
	if err := s.pipe.draw(s.renderer); err != nil {
		return err
	}
	if err := s.drawScore(); err != nil {
		return err
	}
	s.renderer.Present()
	return nil
}

func (s *scene) drawScore() error {
	score := strconv.FormatFloat(s.score, 'f', 0, 64)
	font, err := s.font.RenderUTF8_Solid(score, sdl.Color{})
	if err != nil {
		return fmt.Errorf("could not render score: %v", err)
	}
	defer font.Free()

	var clip sdl.Rect
	font.GetClipRect(&clip)
	t, err := s.renderer.CreateTextureFromSurface(font)
	if err != nil {
		return fmt.Errorf("could not create score texture: %v", err)
	}
	defer t.Destroy()

	// align the score in the middle of the screen.
	r := &sdl.Rect{
		X: (windowWidth / 2) - (clip.W / 2),
		Y: 10,
		W: clip.W,
		H: clip.H,
	}
	if err := s.renderer.Copy(t, nil, r); err != nil {
		return fmt.Errorf("could not copy score texture: %v", err)
	}
	return nil
}

func (s *scene) update() {
	s.bird.update()
	s.pipe.update()
	if s.pipe.intersect(s.bird) {
		s.bird.onCollision()
		return
	}
	if s.pipe.hasPassed(s.bird) {
		s.score++
	}
}

// destroy releases resources associated with the scene object.
// This function should be called upon exit.
func (s *scene) destroy() {
	s.bg.Destroy()
	s.font.Close()
	s.music.Free()
	s.bird.destroy()
	s.pipe.destroy()
}
