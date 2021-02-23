package recipients

import (
	"encoding/csv"
	"io"
	"mime/multipart"
	"strings"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
	"github.com/pkg/errors"
)

// ReadCSV function reads the content of the CSV file into the recipient struct
func ReadCSV(csvFile *multipart.FileHeader) (*[]serializers.RecipientInfo, error) {
	// Open the file
	var recipients []serializers.RecipientInfo
	recordFile, err := csvFile.Open()
	if err != nil {
		return nil, errors.Wrap(err, "Error opening file")
	}

	// Setup the reader
	reader := csv.NewReader(recordFile)

	// Read the records
	record, err := reader.Read() //first row is for headers
	if err != nil {
		return nil, errors.Wrap(err, "Error reading from file")
	}
	if len(record)<5 || record[0]!="ID" || record[1]!= "Email" || record[2]!= "PushToken" || record[3]!= "WebToken" || record[4]!= "PreferredChannelType"{
		return nil, errors.New("Incorrect Headers Format")
	}
	for i := 0; ; i = i + 1 {
		record, err = reader.Read()
		if err == io.EOF {
			break // reached end of the file
		} else if err != nil {
			return nil, errors.Wrap(err, "Error reading from file")
		}
		if record[4] != "" {
			channelType := record[4]
			channelTypeInt := constants.ChannelTypeToInt(channelType)
			if channelTypeInt == 0{
				return nil, errors.New("Error converting channel type to int")
			}
			recipients = append(recipients, serializers.RecipientInfo{RecipientID: strings.ToLower(record[0]), Email: strings.ToLower(record[1]), PushToken: record[2], WebToken: record[3], ChannelType: channelTypeInt})
		} else {
			recipients = append(recipients, serializers.RecipientInfo{RecipientID: strings.ToLower(record[0]), Email: strings.ToLower(record[1]), PushToken: record[2], WebToken: record[3]})
		}
	}
	return &recipients, nil
}
