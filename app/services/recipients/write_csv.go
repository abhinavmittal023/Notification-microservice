package recipients

import (
	"bytes"
	"encoding/csv"
	"strconv"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers"
	"github.com/pkg/errors"
)

// CreateCSV function makes a csv file out of array of struct of recipient serializer
func CreateCSV(recipients *[]serializers.RecipientInfo) (*bytes.Buffer, error) {
	fileBytes := &bytes.Buffer{}
	fileWriter := csv.NewWriter(fileBytes)

	if err := fileWriter.Write([]string{"ID", "Email", "PushToken", "WebToken", "PreferredChannelID"}); err != nil {
		return nil, errors.Wrap(err, "error writing record to csv")
	}

	for _, recipient := range *recipients {
		var record []string
		record = append(record, strconv.FormatUint(recipient.ID, 10))
		record = append(record, recipient.Email)
		record = append(record, recipient.PushToken)
		record = append(record, recipient.WebToken)
		record = append(record, strconv.FormatUint(uint64(recipient.PreferredChannelType), 10))
		if err := fileWriter.Write(record); err != nil {
			return nil, errors.Wrap(err, "error writing record to csv")
		}
	}
	fileWriter.Flush()

	if err := fileWriter.Error(); err != nil {
		return nil, err
	}
	return fileBytes, nil
}
