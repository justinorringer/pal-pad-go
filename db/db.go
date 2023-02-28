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

func getLine(line string) (l models.Line, err error) {
	// unmarshal the line into a Line struct
	l = models.Line{}
	err = json.Unmarshal([]byte(line), &l)

	if err != nil {
		// log error
		log.Printf("Error unmarshalling line: %s", err)
		return models.Line{}, err
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
func ProcessDraw(rc *RedisClient, data []byte) (err error) {
	// marshal the message into a Line struct
	line := models.Line{}
	err = json.Unmarshal(data, &line)

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

type SyncData struct {
	SketchID uuid.UUID
}

func ProcessSync(rc *RedisClient, data []byte) (err error) {
	syncData := SyncData{}
	err = json.Unmarshal(data, &syncData)

	if err != nil {
		// log error
		log.Printf("Error unmarshalling sync data: %s", err)
		return
	}

	// get the sketch from the database
	sketchTable, err := getSketch(rc, syncData.SketchID.String())

	if err != nil {
		// log error
		log.Printf("Error finding matching sketch: %s", err)
		return
	}

	sketch := models.Sketch{
		ID:    sketchTable.ID,
		Lines: []models.Line{},
	}

	// get the lines from the database
	for _, lineID := range sketchTable.LineID {
		line, err := rc.Get(lineID.String())

		if err != nil {
			// log error
			log.Printf("Error finding matching line: %s", err)
			continue
		}

		l, err := getLine(line)

		if err != nil {
			// log error
			log.Printf("Error finding matching line: %s", err)
			continue
		}

		sketch.Lines = append(sketch.Lines, l)

	}

	return err

}

type ClearData struct {
	SketchID uuid.UUID
}

func ProcessClear(rc *RedisClient, data []byte) (err error) {
	clearData := ClearData{}
	err = json.Unmarshal(data, &clearData)

	if err != nil {
		// log error
		log.Printf("Error unmarshalling clear data: %s", err)
		return
	}

	// get the sketch from the database
	sketchTable, err := getSketch(rc, clearData.SketchID.String())

	if err != nil {
		// log error
		log.Printf("Error finding matching sketch: %s", err)
		return
	}

	// delete the lines from the database
	for _, lineID := range sketchTable.LineID {
		err = rc.Del(lineID.String())

		if err != nil {
			// log error
			log.Printf("Error deleting line: %s", err)
			continue
		}
	}

	sketchTable = SketchTable{
		ID:     sketchTable.ID,
		LineID: []uuid.UUID{},
	}

	// update sketch
	err = updateSketch(rc, sketchTable)

	if err != nil {
		// log error
		log.Printf("Error updating sketch: %s", err)
		return
	}

	return
}
