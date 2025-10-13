package constants

// JobType defines the type for our job constants.
type JobType string

// Defines all the background job types in the system.
const (
	JobTypeDigitalSign            JobType = "DIGITAL_SIGN_NOTE"
	JobTypeUploadNote             JobType = "UPLOAD_NOTE"
	JobTypeUploadNoteWithData     JobType = "UPLOAD_NOTE_WITH_DATA"
	JobTypeUploadNoteWithDataDNCN JobType = "UPLOAD_NOTE_WITH_DATA_DNCN"
	JobTypeAccfin                 JobType = "ACCFIN"
	// Add other job types here as you need them, for example:
	// JobTypeSendWelcomeEmail JobType = "SEND_WELCOME_EMAIL"
)
