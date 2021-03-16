package logwrapper

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Event stores messages to log later, from our standard interface
type Event struct {
	id      int
	message string
}

// StandardLogger enforces specific log message formats
type StandardLogger struct {
	*logrus.Logger
}

// NewLogger initializes the standard logger
func NewLogger() *StandardLogger {

	var baseLogger = logrus.New()

	var standardLogger = &StandardLogger{baseLogger}

	standardLogger.Formatter = &logrus.JSONFormatter{}
	return standardLogger
}

// OpenFile function opens or creates the log file
func OpenFile() (*os.File, error) {
	var filename string = "logfile.log"
	// Create the log file if doesn't exist. And append to it if it already exists.
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	return f, err
}

// Declare variables to store log messages as new Events
var (
	internalServerError = Event{1, "Internal Server Error due to: %s"}
	entityAdded         = Event{2, "%s added successfully"}
	entityDeleted       = Event{3, "%s deleted successfully"}
	entityUpdated       = Event{4, "%s updated successfully"}
)

// EntityAdded is a standard info message for addition of entity
func (l *StandardLogger) EntityAdded(argumentName string) {
	l.Infof(entityAdded.message, argumentName)
}

// InternalServerError is a standard error message for internal server error
func (l *StandardLogger) InternalServerError(argumentName string) {
	l.Errorf(internalServerError.message, argumentName)
}

// EntityDeleted is a standard info message for deletion of entity
func (l *StandardLogger) EntityDeleted(argumentName string) {
	l.Infof(entityDeleted.message, argumentName)
}

// EntityUpdated is a standard info message for updation of entity
func (l *StandardLogger) EntityUpdated(argumentName string) {
	l.Infof(entityUpdated.message, argumentName)
}
