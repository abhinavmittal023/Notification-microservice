package constants

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

var resp Constants //variable to store final decoded data

//GetConstants function returns the configuration struct
func GetConstants() Constants {
	return resp
}

//Init Function unmarshalls the constants.json into the struct
func Init() error {

	file, err := os.Open("./constants/constants.json") //opening json file
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
