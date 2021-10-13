package response

import "fmt"

type textBlockType string

const (
	textBlockTypePlainText textBlockType = "plain_text"
	textBlockTypeMarkdown  textBlockType = "mrkdwn"
)

// https://api.slack.com/reference/messaging/composition-objects#text
type TextBlock struct {
	Type textBlockType `json:"type"`
	Text string        `json:"text"`
}

func Markdown(format string, a ...interface{}) *TextBlock {
	return &TextBlock{
		Type: textBlockTypeMarkdown,
		Text: fmt.Sprintf(format, a...),
	}
}

func PlainText(format string, a ...interface{}) *TextBlock {
	return &TextBlock{
		Type: textBlockTypePlainText,
		Text: fmt.Sprintf(format, a...),
	}
}

type blockObjectType string

const (
	blockObjectTypeSection blockObjectType = "section"
	blockObjectTypeDivider blockObjectType = "divider"
)

// https://api.slack.com/reference/block-kit/blocks#section
type SectionBlock struct {
	Type   blockObjectType `json:"type"`
	Text   *TextBlock      `json:"text,omitempty"`
	Fields []*TextBlock    `json:"fields,omitempty"`
}

func Divider() *SectionBlock {
	return &SectionBlock{
		Type: blockObjectTypeDivider,
	}
}

func Section(text *TextBlock, fields ...*TextBlock) *SectionBlock {
	return &SectionBlock{
		Type:   blockObjectTypeSection,
		Text:   text,
		Fields: fields,
	}
}

// https://api.slack.com/interactivity/slash-commands#responding_immediate_response
type responseType string

const (
	responseTypeEphemeral responseType = "ephemeral"
	responseTypeInChannel responseType = "in_channel"
)

// slash command response message
// https://api.slack.com/interactivity/slash-commands#responding_to_commands
type message struct {
	ResponseType responseType    `json:"response_type,omitempty"`
	Blocks       []*SectionBlock `json:"blocks,omitempty"`
}

func newMessage(rt responseType, sections ...*SectionBlock) *message {
	return &message{
		ResponseType: rt,
		Blocks:       sections,
	}
}
