package contexts

import (
	"encoding/json"
	"github.com/MerinEREN/iiPackages/datastore/context"
)

// ContextWithValueOnly is used for page context request's response body.
type ContextWithValueOnly struct {
	Value string `json:"value"`
}

// GetLangValue returns context or contexts only with value of the corresponding language.
func GetLangValue(cs context.Contexts, lang string) (interface{}, error) {
	cwvos := make(map[string]ContextWithValueOnly)
	var err error
	for i, v := range cs {
		contextValues := make(map[string]string)
		err = json.Unmarshal(v.ValuesBS, &contextValues)
		if err != nil {
			return nil, err
		}
		cwvo := ContextWithValueOnly{contextValues[lang]}
		if len(cs) == 1 {
			return cwvo, err
		}
		cwvos[i] = cwvo
	}
	return cwvos, err
}
