package contents

import (
	"encoding/json"
	"github.com/MerinEREN/iiPackages/datastore/content"
)

// GetLangValue sends the requested language value of the contents only.
func GetLangValue(cs content.Contents, lang string) (map[string]string, error) {
	contentsClient := make(map[string]string)
	var err error
	for i, v := range cs {
		contentValues := make(map[string]string)
		err = json.Unmarshal(v.ValuesBS, &contentValues)
		if err != nil {
			return nil, err
		}
		contentsClient[i] = contentValues[lang]
	}
	return contentsClient, err
}
