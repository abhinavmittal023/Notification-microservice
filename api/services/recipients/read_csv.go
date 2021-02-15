package recipients

import (
	"encoding/csv"
	"io"
	"mime/multipart"
	"strconv"

	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
	"github.com/pkg/errors"
)

// ReadCSV function reads the content of the CSV file into the recipient struct
func ReadCSV(csvFile *multipart.FileHeader) (*[]models.Recipient, error) {
	// Open the file
	var recipients []models.Recipient
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
		var id int
		if record[3] != "" {
			id, err = strconv.Atoi(record[3])
			if err != nil {
				return nil, errors.Wrap(err, "Error converting id to int")
			}
			recipients = append(recipients, models.Recipient{Email: record[0], PushToken: record[1], WebToken: record[2], PreferredChannelID: uint64(id)})
		} else {
			recipients = append(recipients, models.Recipient{Email: record[0], PushToken: record[1], WebToken: record[2]})
		}
	}
	return &recipients, nil
}
