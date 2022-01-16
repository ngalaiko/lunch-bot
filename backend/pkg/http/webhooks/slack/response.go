package slack

import (
	"log"
)

// ReplaceEphemeral sends a message back visible only by the caller replacing the previous mesage.
func ReplaceEphemeral(text string, sections ...*Block) *Message {
	return NewReplaceMessage(ResponseTypeEphemeral, text, sections...)
}

// Ephemeral sends a message back visible only by the caller.
func Ephemeral(text string, sections ...*Block) *Message {
	return NewMessage(ResponseTypeEphemeral, text, sections...)
}

// InChannel sends a message back visible by everyone in the channel.
func InChannel(text string, sections ...*Block) *Message {
	return NewMessage(ResponseTypeInChannel, text, sections...)
}

// BadRequest sends an ephemeral message to the user indicating that the request was invalid.
func BadRequest(err error) *Message {
	return Ephemeral(err.Error(), Section(PlainText(err.Error())))
}

// InternalServerError sends an ephemeral message to the user indicating that the request failed.
func InternalServerError(err error) *Message {
	log.Printf("[ERROR] %s", err)
	msg := "Sorry, that didn't work. Try again or contact the app administrator."
	return Ephemeral(msg, Section(PlainText(msg)))
}
