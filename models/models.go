package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// sketch contains an array of lines (as json)
type Sketch struct {
	ID    uuid.UUID `json:"id"`
	Lines []Line    `json:"lines"`
}

type Line struct {
	ID         uuid.UUID `json:"id"`
	SketchID   uuid.UUID `json:"sketch_id"`
	Timestamp  time.Time `json:"timestamp"` // when line is finished
	UserID     int       `json:"user_id"`
	BrushIndex int       `json:"brush_index"`
	Points     []Point   `json:"points"`
	Color      string    `json:"color"`
}

// limiting the coordinates to 1920x1080 canvas,
// so we can use 16 bit integers
type Point struct {
	X       int16 `json:"x"`
	Y       int16 `json:"y"`
	Size    int16 `json:"size"`
	Opacity int16 `json:"opacity"`
}
