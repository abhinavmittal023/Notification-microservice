package configuration

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

var resp Configuration //variable to store final decoded data
var dbString string    //variable to store the dbString

//GetResp function returns the configuration struct
func GetResp() Configuration {
	return resp
}

//GetDBString function returns the database string
func GetDBString() string {
	return dbString
}

//Init Function unmarshalls the config.json into the struct
func Init() error {

	file, err := os.Open("./configuration/config.json") //opening json file
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

	dbString = "user=" + resp.Database.User + " password=" + resp.Database.Password + " dbname=" +
		resp.Database.DbName + " sslmode=" + resp.Database.SSLMode
	return nil
}
