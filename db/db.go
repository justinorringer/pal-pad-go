package db

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/justinorringer/pal-pad-go/models"
)

type SketchTable struct {
	ID     uuid.UUID
	LineID []uuid.UUID
}

func getSketch(rc *RedisClient, id string) (sketchTable SketchTable, err error) {
	// get the sketch from the database
	sketch, err := rc.Get(id)

	if err != nil {
		// log error
		log.Printf("Error getting sketch: %s", err)
		return
	}

	// unmarshal the sketch into a Sketch struct
	sketchTable = SketchTable{}
	err = json.Unmarshal([]byte(sketch), &sketchTable)

	if err != nil {
		// log error
		log.Printf("Error unmarshalling sketch: %s", err)
		return
	}

	return
}

func updateSketch(rc *RedisClient, sketchTable SketchTable) (err error) {
	// marshal the sketch into a json string
	json, err := json.Marshal(sketchTable)

	if err != nil {
		// log error
		log.Printf("Error marshalling sketch: %s", err)
		return
	}

	// save the sketch to the database
	err = rc.Set(sketchTable.ID.String(), string(json))

	if err != nil {
		// log error
		log.Printf("Error setting sketch: %s", err)
		return
	}

	return
}

func updateLine(rc *RedisClient, line models.Line) (err error) {
	// marshal the line into a json string
	json, err := json.Marshal(line)

	if err != nil {
		// log error
		log.Printf("Error marshalling line: %s", err)
		return
	}

	// save the line to the database
	err = rc.Set(line.ID.String(), string(json))

	if err != nil {
		// log error
		log.Printf("Error setting line: %s", err)
		return
	}

	return
}

// Save the line to the database
func ProcessMessage(rc *RedisClient, message []byte) (err error) {
	// marshal the message into a Line struct
	line := models.Line{}
	err = json.Unmarshal(message, &line)

	if err != nil {
		// log error
		log.Printf("Error unmarshalling line: %s", err)
		return
	}

	// get the sketch from the database
	sketchTable, err := getSketch(rc, line.SketchID.String())

	if err != nil {
		// log error
		log.Printf("Error finding matching sketch: %s", err)
		return
	}

	// add the line to the sketch
	sketchTable.LineID = append(sketchTable.LineID, line.ID)

	// save the sketch to the database
	err = updateSketch(rc, sketchTable)

	if err != nil {
		// log error
		log.Printf("Error updating sketch: %s", err)
		return
	}

	// save the line to the database
	err = updateLine(rc, line)

	if err != nil {
		// log error
		log.Printf("Error updating line: %s", err)
		return
	}

	return
}
