package configuration

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

var resp Configuration // variable to store final decoded data

// GetResp function returns the configuration struct
func GetResp() Configuration {
	return resp
}

// Init Function unmarshalls the config.json into the struct
func Init() error {

	file, err := os.Open("./configuration/config.json") // opening json file
	if err != nil {
		return errors.Wrap(err, "Unable to open the config file")
	}
	defer file.Close()

	bt, err := ioutil.ReadAll(file) // reading it using ioutil
	if err != nil {
		return errors.Wrap(err, "Unable to read the config file")
	}

	err = json.Unmarshal(bt, &resp) // decoding from bytes of encodings
	if err != nil {
		return errors.Wrap(err, "Unable to unmarshall data from config file")
	}

	return nil
}
