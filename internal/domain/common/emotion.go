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
	buf, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	return string(buf), nil
}

func (e *EmotionList) Validate() bool {
	for _, emotion := range *e {
		if !slices.Contains(SupportedEmotions, emotion) {
			return false
		}
	}
	return true
}

func (e *EmotionList) ToString() string {
	result, _ := e.Value()
	return result.(string)
}

type EmotionRate struct {
	Emotion Emotion `json:"emotion"`
	Rate    float64 `json:"rate"`
}

func (e *EmotionRate) Scan(src interface{}) error {
	bytes, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("failed to type assert EmotionList")
	}
	var result EmotionRate
	if err := json.Unmarshal(bytes, &result); err != nil {
		return err
	}
	*e = result
	return nil
}

func (e *EmotionRate) Value() (driver.Value, error) {
	buf, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	return string(buf), nil
}

func (e *EmotionRate) Validate() bool {
	return slices.Contains(SupportedEmotions, e.Emotion)
}

func (e *EmotionRate) ToString() string {
	result, _ := e.Value()
	return result.(string)
}

type EmotionRateList []EmotionRate

func (e *EmotionRateList) Scan(src interface{}) error {
	bytes, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("failed to type assert EmotionRateList")
	}
	var result []EmotionRate
	if err := json.Unmarshal(bytes, &result); err != nil {
		return err
	}
	*e = result
	return nil
}

func (e *EmotionRateList) Value() (driver.Value, error) {
	buf, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	return string(buf), nil
}

func (e *EmotionRateList) Validate() bool {
	for _, emotionRate := range *e {
		if !emotionRate.Validate() {
			return false
		}
	}
	return true
}

func (e *EmotionRateList) ToString() string {
	result, _ := e.Value()
	return result.(string)
}

func (e *EmotionRateList) GetDominant() Emotion {
	if len(*e) == 0 {
		return EmotionNeutral
	}
	var maxRate float64
	var dominant Emotion
	for _, emotionRate := range *e {
		if emotionRate.Rate > maxRate {
			maxRate = emotionRate.Rate
			dominant = emotionRate.Emotion
		}
	}
	return dominant
}

type EmotionCount struct {
	Emotion Emotion `json:"emotion"`
	Count   int64   `json:"count"`
}

func (e *EmotionCount) Scan(src interface{}) error {
	bytes, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("failed to type assert EmotionCount")
	}
	var result EmotionCount
	if err := json.Unmarshal(bytes, &result); err != nil {
		return err
	}
	*e = result
	return nil
}

func (e *EmotionCount) Value() (driver.Value, error) {
	buf, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	return string(buf), nil
}

func (e *EmotionCount) Validate() bool {
	return slices.Contains(SupportedEmotions, e.Emotion)
}

func (e *EmotionCount) ToString() string {
	result, _ := e.Value()
	return result.(string)
}

type EmotionCountList []EmotionCount

func (e *EmotionCountList) Scan(src interface{}) error {
	bytes, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("failed to type assert EmotionCountList")
	}
	var result []EmotionCount
	if err := json.Unmarshal(bytes, &result); err != nil {
		return err
	}
	*e = result
	return nil
}

func (e *EmotionCountList) Value() (driver.Value, error) {
	buf, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	return string(buf), nil
}

func (e *EmotionCountList) Validate() bool {
	for _, emotionCount := range *e {
		if !emotionCount.Validate() {
			return false
		}
	}
	return true
}

func (e *EmotionCountList) ToString() string {
	result, _ := e.Value()
	return result.(string)
}

func (e *EmotionCountList) GetDominant() Emotion {
	if len(*e) == 0 {
		return EmotionNeutral
	}
	var maxCount int64
	var dominant Emotion
	for _, emotionCount := range *e {
		if emotionCount.Count > maxCount {
			maxCount = emotionCount.Count
			dominant = emotionCount.Emotion
		}
	}
	return dominant
}
