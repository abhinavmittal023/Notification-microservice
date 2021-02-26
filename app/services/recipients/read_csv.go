package recipients

import (
	"encoding/csv"
	"fmt"
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

	if len(record) != len(constants.CSVHeaders()) {
		return nil, errors.New("Invalid number of headers provided")
	}

	csvHeaders := map[string]int{}
	for i, header := range constants.CSVHeaders() {
		csvHeaders[header] = i
	}

	for i, gotHeader := range record {
		_, present := csvHeaders[gotHeader]
		if !present {
			return nil, fmt.Errorf("Invalid Header %s", gotHeader)
		}
		csvHeaders[gotHeader] = i
	}

	for i := 0; ; i = i + 1 {
		record, err = reader.Read()
		if err == io.EOF {
			break // reached end of the file
		} else if err != nil {
			return nil, errors.Wrap(err, "Error reading from file")
		}
		if record[csvHeaders[constants.CSVHeaders()[4]]] != "" {
			channelType := record[csvHeaders[constants.CSVHeaders()[4]]]
			channelTypeInt := constants.ChannelTypeToInt(channelType)
			if channelTypeInt == 0 {
				return nil, fmt.Errorf("Invalid Channel type %s", channelType)
			}
			recipients = append(recipients, serializers.RecipientInfo{RecipientID: strings.ToLower(record[csvHeaders[constants.CSVHeaders()[0]]]), Email: strings.ToLower(record[csvHeaders[constants.CSVHeaders()[1]]]), PushToken: record[csvHeaders[constants.CSVHeaders()[2]]], WebToken: record[csvHeaders[constants.CSVHeaders()[3]]], ChannelType: channelTypeInt})
		} else {
			recipients = append(recipients, serializers.RecipientInfo{RecipientID: strings.ToLower(record[csvHeaders[constants.CSVHeaders()[0]]]), Email: strings.ToLower(record[csvHeaders[constants.CSVHeaders()[1]]]), PushToken: record[csvHeaders[constants.CSVHeaders()[2]]], WebToken: record[csvHeaders[constants.CSVHeaders()[3]]]})
		}
	}
	return &recipients, nil
}
