package common

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"slices"

	"github.com/solutionchallenge/ondaum-server/pkg/utils"
)

type Inspection string

const (
	InspectionPHQ9 Inspection = "phq-9"
	InspectionGAD7 Inspection = "gad-7"
	InspectionPSS  Inspection = "pss"
)

var SupportedInspections = InspectionList{
	InspectionPHQ9,
	InspectionGAD7,
	InspectionPSS,
}

type InspectionList []Inspection

func (e *InspectionList) Scan(src interface{}) error {
	bytes, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("failed to type assert InspectionList")
	}
	var result []Inspection
	if err := json.Unmarshal(bytes, &result); err != nil {
		return err
	}
	*e = result
	return nil
}

func (e *InspectionList) Value() (driver.Value, error) {
	buf, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	return string(buf), nil
}

func (e *InspectionList) Validate() bool {
	for _, test := range *e {
		if !slices.Contains(SupportedInspections, test) {
			return false
		}
	}
	return true
}

func (e *InspectionList) ToString() string {
	result, _ := e.Value()
	return result.(string)
}

type InspectionPaper struct {
	Name      string               `json:"name"`
	Guides    string               `json:"guides"`
	Questions []InspectionQuestion `json:"questions"`
	Results   []InspectionResult   `json:"results"`
	Scoring   InspectionScoring    `json:"scoring"`
}

type InspectionQuestion struct {
	Index    int                `json:"index"`
	Question string             `json:"question"`
	Answers  []InspectionAnswer `json:"answers"`
}

type InspectionAnswer struct {
	Score int    `json:"score"`
	Text  string `json:"text"`
}

type InspectionResult struct {
	Min         int    `json:"min"`
	Max         int    `json:"max"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Critical    bool   `json:"critical"`
}

type InspectionScoring struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

func ReadInspectionPaperFrom(path string, root ...string) (InspectionPaper, error) {
	jsonData, err := utils.ReadFileFrom(path, root...)
	if err != nil {
		return InspectionPaper{}, err
	}
	var testPaper InspectionPaper
	if err := json.Unmarshal([]byte(jsonData), &testPaper); err != nil {
		return InspectionPaper{}, err
	}
	return testPaper, nil
}
