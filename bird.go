package main

import (
	"fmt"
	"math"
	"path"
	"sync"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	gravity    = 0.12
	jumpSpeed  = 5
	angleBoost = 10
)

type bird struct {
	frames []*sdl.Texture

	mutex           sync.RWMutex
	frame           uint64
	x, y, w, h      int32
	velocity, angle float64
	dead            bool
	timeOnDeath     time.Time
}

func newBird(r *sdl.Renderer, assetDir string) (*bird, error) {
	frames := make([]*sdl.Texture, 10)
	for i := 1; i <= len(frames); i++ {
		texture, err := img.LoadTexture(r, path.Join(assetDir, fmt.Sprintf("bird/bird_frame_%d.png", i)))
		if err != nil {
			return nil, fmt.Errorf("could not load bird texture: %v", err)
		}
		frames[i-1] = texture
	}
	return &bird{
		frames: frames,
		x:      10,
		y:      windowHeight / 2,
		w:      50,
		h:      43,
	}, nil
}

// draw renders the internal state of the bird.
func (b *bird) draw(r *sdl.Renderer) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	rect := &sdl.Rect{X: b.x, Y: windowHeight - b.y - b.h/2, W: b.w, H: b.h}
	var i uint64
	if b.dead {
		// only the last two frames are eligible,
		// since they are the frames indicating death.
		i = (uint64(len(b.frames)) - (b.frame / 10 % 2)) - 1
	} else {
		// all frames eligible except the last two.
		i = b.frame / 10 % (uint64(len(b.frames)) - 2)
	}

	if err := r.CopyEx(b.frames[i], nil, rect, b.angle, nil, sdl.FLIP_NONE); err != nil {
		return fmt.Errorf("could not copy bird frame %d: %v", i, err)
	}
	b.frame++
	return nil
}

// update updates the internal state of the bird.
func (b *bird) update() {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if b.dead {
		return
	}
	// fall velocity
	b.y -= int32(b.velocity)
	// bird angle to show falling or flying up
	b.angle = math.Min(25, b.velocity*angleBoost)
	atTop := b.y+b.h/2 > windowHeight
	atBottom := b.y-b.h/2 < 80

	if atTop || atBottom {
		// prevent the bird to get out from the top and bottom
		b.y += int32(b.velocity)
		if atBottom {
			// make the bird angle straight
			b.angle = 0
		}
	}
	b.velocity += gravity
}

func (b *bird) handleKeyEvent(event sdl.Event) {
	switch v := event.(type) {
	case *sdl.KeyDownEvent:
		if v.Keysym.Sym == sdl.K_SPACE {
			// jump up
			b.velocity = -jumpSpeed
		}
	}
}

// intersect returns true iff the bird intersects with the given rect.
func (b *bird) intersect(r *sdl.Rect) bool {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	rect := sdl.Rect{X: b.x, Y: windowHeight - b.y - b.h/2, W: b.w, H: b.h}
	return rect.HasIntersection(r)
}

// hasPassed returns true iff the bird has passed the given rect.
func (b *bird) hasPassed(r *sdl.Rect) bool {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	return b.x == r.X+b.w
}

// restart rewinds the bird object to its origin state.
// Only useable if the bird is dead.
func (b *bird) restart() {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if !b.dead {
		return
	}
	b.y = 300
	b.velocity = 0
	b.dead = false
}

// destroy releases resources associated with the bird object.
// Do not use the object after calling this function.
func (b *bird) destroy() {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	for i := range b.frames {
		b.frames[i].Destroy()
	}
	b.frames = nil
}

// onCollision sets the bird as dead.
func (b *bird) onCollision() {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if b.dead {
		return
	}
	b.dead = true
	b.timeOnDeath = time.Now()
}

// deadSince returns the duration since the bird died.
func (b *bird) deadSince() time.Duration {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	if !b.dead {
		return 0
	}
	return time.Since(b.timeOnDeath)
}
