package slack

import (
	"log"
)

// ReplaceEphemeral sends a message back visible only by the caller replacing the previous mesage.
func ReplaceEphemeral(text Text, sections ...*Block) *Message {
	return newReplaceMessage(responseTypeEphemeral, text, sections...)
}

// Ephemeral sends a message back visible only by the caller.
func Ephemeral(text Text, sections ...*Block) *Message {
	return newMessage(responseTypeEphemeral, text, sections...)
}

// InChannel sends a message back visible by everyone in the channel.
func InChannel(text Text, sections ...*Block) *Message {
	return newMessage(responseTypeInChannel, text, sections...)
}

// BadRequest sends an ephemeral message to the user indicating that the request was invalid.
func BadRequest(err error) *Message {
	return Ephemeral(Text(err.Error()), Section(PlainText(err.Error())))
}

// InternalServerError sends an ephemeral message to the user indicating that the request failed.
func InternalServerError(err error) *Message {
	log.Printf("[ERROR] %s", err)
	msg := "Sorry, that didn't work. Try again or contact the app administrator."
	return Ephemeral(Text(msg), Section(PlainText(msg)))
}
