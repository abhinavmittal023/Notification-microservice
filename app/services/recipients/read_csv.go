package recipients

import (
	"encoding/csv"
	"io"
	"mime/multipart"
	"strconv"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers"
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
	_, err = reader.Read() //first row is for headers
	if err != nil {
		return nil, errors.Wrap(err, "Error reading from file")
	}

	for i := 0; ; i = i + 1 {
		record, err := reader.Read()
		if err == io.EOF {
			break // reached end of the file
		} else if err != nil {
			return nil, errors.Wrap(err, "Error reading from file")
		}
		var channelID uint64
		if record[4] != "" {
			channelID, err = strconv.ParseUint(record[4], 10, 64)
			if err != nil {
				return nil, errors.Wrap(err, "Error converting channelID to int")
			}
			recipients = append(recipients, serializers.RecipientInfo{RecipientUUID: record[0], Email: record[1], PushToken: record[2], WebToken: record[3], PreferredChannelID: uint64(channelID)})
		} else {
			recipients = append(recipients, serializers.RecipientInfo{RecipientUUID: record[0], Email: record[1], PushToken: record[2], WebToken: record[3]})
		}
	}
	return &recipients, nil
}
