package configuration

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

var resp Configuration //variable to store final decoded data

//GetResp function returns the configuration struct
func GetResp() Configuration {
	return resp
}

//Init Function marshalls the config.json into the struct
func Init() error {

	file, err := os.Open("./configuration/config.json") //opening jason file
	if err != nil {
		return err
	}
	defer file.Close()

	bt, err := ioutil.ReadAll(file) //reading it using ioutil
	if err != nil {
		return err
	}

	err = json.Unmarshal(bt, &resp) //decoding from bytes of encodings
	if err != nil {
		return err
	}
	return nil
}
