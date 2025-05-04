package common

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Emotion string

const (
	EmotionJoy      Emotion = "joy"
	EmotionSadness  Emotion = "sadness"
	EmotionAnger    Emotion = "anger"
	EmotionSurprise Emotion = "surprise"
	EmotionFear     Emotion = "fear"
	EmotionDisgust  Emotion = "disgust"
	EmotionNeutral  Emotion = "neutral"
)

type EmotionList []Emotion

func (e *EmotionList) Scan(src interface{}) error {
	bytes, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("failed to type assert EmotionList")
	}
	var result []Emotion
	if err := json.Unmarshal(bytes, &result); err != nil {
		return err
	}
	*e = result
	return nil
}

func (e EmotionList) Value() (driver.Value, error) {
	return json.Marshal(e)
}
