package common

import (
	"database/sql/driver"
	"encoding/json"
	"slices"

	"github.com/solutionchallenge/ondaum-server/pkg/utils"
)

type Diagnosis string

const (
	DiagnosisPHQ9 Diagnosis = "phq-9"
	DiagnosisGAD7 Diagnosis = "gad-7"
	DiagnosisPSS  Diagnosis = "pss"
)

var DiagnosisDescriptions = map[Diagnosis]string{
	DiagnosisPHQ9: "Patient Health Questionnaire-9",
	DiagnosisGAD7: "Generalized Anxiety Disorder-7",
	DiagnosisPSS:  "Perceived Stress Scale",
}

var DiagnosisFilepaths = map[Diagnosis]string{
	DiagnosisPHQ9: "resource/diagnosis/phq-9-en.json",
	DiagnosisGAD7: "resource/diagnosis/gad-7-en.json",
	DiagnosisPSS:  "resource/diagnosis/pss-en.json",
}

var SupportedDiagnoses = DiagnosisList{
	DiagnosisPHQ9,
	DiagnosisGAD7,
	DiagnosisPSS,
}

type DiagnosisList []Diagnosis

func (e *DiagnosisList) Scan(src interface{}) error {
	bytes, ok := src.([]byte)
	if !ok {
		return utils.NewError("failed to type assert DiagnosisList")
	}
	var result []Diagnosis
	if err := json.Unmarshal(bytes, &result); err != nil {
		return err
	}
	*e = result
	return nil
}

func (e *DiagnosisList) Value() (driver.Value, error) {
	buf, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	return string(buf), nil
}

func (e *DiagnosisList) Validate() bool {
	for _, test := range *e {
		if !slices.Contains(SupportedDiagnoses, test) {
			return false
		}
	}
	return true
}

func (e *DiagnosisList) ToString() string {
	result, _ := e.Value()
	return result.(string)
}

type DiagnosisPaper struct {
	Name      string              `json:"name"`
	Guides    string              `json:"guides"`
	Questions []DiagnosisQuestion `json:"questions"`
	Results   []DiagnosisResult   `json:"results"`
	Scoring   DiagnosisScoring    `json:"scoring"`
}

type DiagnosisQuestion struct {
	Index    int               `json:"index"`
	Question string            `json:"question"`
	Answers  []DiagnosisAnswer `json:"answers"`
}

type DiagnosisAnswer struct {
	Score int    `json:"score"`
	Text  string `json:"text"`
}

type DiagnosisResult struct {
	Min         int    `json:"min"`
	Max         int    `json:"max"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Critical    bool   `json:"critical"`
}

type DiagnosisScoring struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

func ReadDiagnosisPaperFrom(path string, root ...string) (DiagnosisPaper, error) {
	jsonData, err := utils.ReadFileFrom(path, root...)
	if err != nil {
		return DiagnosisPaper{}, err
	}
	var testPaper DiagnosisPaper
	if err := json.Unmarshal([]byte(jsonData), &testPaper); err != nil {
		return DiagnosisPaper{}, err
	}
	return testPaper, nil
}
