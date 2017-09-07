package main

import (
	"fmt"
	"math/rand"
	"path"
	"sync"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type pipe struct {
	texture *sdl.Texture
	speed   int32

	mutex sync.RWMutex
	pipes []*sdl.Rect
}

func newPipe(r *sdl.Renderer, assetDir string) (*pipe, error) {
	texture, err := img.LoadTexture(r, path.Join(assetDir, "pipe/pipe.png"))
	if err != nil {
		return nil, fmt.Errorf("could not load pipe texture: %v", err)
	}
	pipe := &pipe{texture: texture, speed: 2}
	// spawn pipes in a seperate goroutine
	go func() {
		// spawn a new pipe once ever 1.2 seconds
		const pipeSpawnDuration = 1200 * time.Millisecond
		tick := time.Tick(pipeSpawnDuration)

		// don't seed the global rand
		random := rand.New(rand.NewSource(time.Now().Unix()))

		for range tick {
			h := 100 + int32(random.Intn(300))
			rect := &sdl.Rect{X: 800, Y: (600 - h) - 80, W: 50, H: h}

			if random.Float32() > 0.5 {
				// flip the pipe upside down
				rect.Y = 0
			}
			pipe.mutex.Lock()
			pipe.pipes = append(pipe.pipes, rect)
			pipe.mutex.Unlock()
		}
	}()
	return pipe, nil
}

// draw renders the internal state of the pipes.
func (p *pipe) draw(r *sdl.Renderer) error {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	for _, pipe := range p.pipes {
		flip := sdl.FLIP_NONE
		if pipe.Y == 0 {
			// pipe should be upside down
			flip = sdl.FLIP_VERTICAL
		}
		if err := r.CopyEx(p.texture, nil, pipe, 0, nil, flip); err != nil {
			return fmt.Errorf("could not copy pipe: %v", err)
		}
	}
	return nil
}

// update updates the internal state of the pipes.
func (p *pipe) update() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	rem := p.pipes[:0]
	for _, pipe := range p.pipes {
		pipe.X -= p.speed
		if pipe.X+pipe.W > 0 {
			// only keep pipes that are still on the screen
			rem = append(rem, pipe)
		}
	}
	p.pipes = rem
}

// intersect returns true iff the bird intersects with any pipes.
func (p *pipe) intersect(b *bird) bool {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	for _, pipe := range p.pipes {
		if b.intersect(pipe) {
			return true
		}
	}
	return false
}

// hasPassed returns true iff the bird has passed
// the closest pipe.
func (p *pipe) hasPassed(b *bird) bool {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	if len(p.pipes) == 0 {
		return false
	}

	// index 0 is always the pipe closest to the bird.
	return b.hasPassed(p.pipes[0])
}

// restart rewinds the pipe object to its origin state.
func (p *pipe) restart() {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.pipes = p.pipes[:0]
}

// destroy releases resources associated with the pipe object.
// Do not use the object after calling this function.
func (p *pipe) destroy() {
	p.texture.Destroy()
}
