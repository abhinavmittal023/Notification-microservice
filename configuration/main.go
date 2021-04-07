package configuration

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

var resp Configuration // variable to store final decoded data

// GetResp function returns the configuration struct
func GetResp() Configuration {
	return resp
}

// Init Function unmarshalls the config.json into the struct
func Init() error {

	cwd, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "Unable to get working directory")
	}

	for ; strings.Split(cwd, "/")[len(strings.Split(cwd, "/"))-1] != "notifications-microservice"; cwd = filepath.Dir(cwd) {
	}

	file, err := os.Open(fmt.Sprintf("%s/configuration/config.json", cwd)) // opening json file
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
