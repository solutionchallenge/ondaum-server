package common

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"slices"
)

type Feature string

/*
* `"escalate_crisis"`: 사용자가 심각한 위기 상황(자살/자해 위험)임을 감지했을 때, 즉시 모든 대화를 중단하고 이 액션을 반환합니다. (애플리케이션은 이 액션을 받아 전문가 연결 안내 등 비상 대응 절차를 수행해야 함)
* `"suggest_test_phq9"`: 사용자에게 간이 우울증 검사(PHQ-9) 기능 사용을 제안할 때 사용합니다.
* `"suggest_test_gad7"`: 사용자에게 간이 불안 증상 검사(PHQ-9) 기능 사용을 제안할 때 사용합니다.
* `"suggest_test_pss"`: 사용자에게 간이 스트레스 척도 검사(PHQ-9) 기능 사용을 제안할 때 사용합니다.
* `"end_conversation"`: 사용자와 협의 하에 대화를 종료한 경우에 사용합니다.
 */
const (
	FeatureEscalateCrisis  Feature = "escalate_crisis"
	FeatureSuggestTestPHQ9 Feature = "suggest_test_phq9"
	FeatureSuggestTestGAD7 Feature = "suggest_test_gad7"
	FeatureSuggestTestPSS  Feature = "suggest_test_pss"
	FeatureEndConversation Feature = "end_conversation"
)

var SupportedFeatures = FeatureList{
	FeatureEscalateCrisis,
	FeatureSuggestTestPHQ9,
	FeatureSuggestTestGAD7,
	FeatureSuggestTestPSS,
	FeatureEndConversation,
}

type FeatureList []Feature

func (e *FeatureList) Scan(src interface{}) error {
	bytes, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("failed to type assert FeatureList")
	}
	var result []Feature
	if err := json.Unmarshal(bytes, &result); err != nil {
		return err
	}
	*e = result
	return nil
}

func (e *FeatureList) Value() (driver.Value, error) {
	buf, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	return string(buf), nil
}

func (e *FeatureList) Validate() bool {
	for _, feature := range *e {
		if !slices.Contains(SupportedFeatures, feature) {
			return false
		}
	}
	return true
}

func (e *FeatureList) ToString() string {
	result, _ := e.Value()
	return result.(string)
}
