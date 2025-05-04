package common

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"slices"
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

var SupportedEmotions = EmotionList{
	EmotionJoy,
	EmotionSadness,
	EmotionAnger,
	EmotionSurprise,
	EmotionFear,
	EmotionDisgust,
	EmotionNeutral,
}

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

func (e *EmotionList) Value() (driver.Value, error) {
	return json.Marshal(e)
}

func (e *EmotionList) Validate() bool {
	for _, emotion := range *e {
		if !slices.Contains(SupportedEmotions, emotion) {
			return false
		}
	}
	return true
}
