package wallet

import (
	"encoding/json"
	"os"
)

func ReadKeysFile(path string) (*File, error) {

	path = defaultKeysFile(path)

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(file)
	decoder.DisallowUnknownFields()
	decodedFile := &keysFileJSON{}
	err = decoder.Decode(&decodedFile)
	if err != nil {
		return nil, err
	}

	keysFile := &File{
		path: path,
	}
	err = keysFile.fromJSON(decodedFile)
	if err != nil {
		return nil, err
	}

	return keysFile, nil
}
