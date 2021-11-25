package slack

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

type Text string

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

type Option struct {
	Text  *TextBlock `json:"text"`
	Value string     `json:"value"`
}

type accessoryType string

const (
	accessoryTypeButton accessoryType = "button"
)

type accessory struct {
	Type     accessoryType `json:"type"`
	Text     *TextBlock    `json:"text"`
	ActionID string        `json:"action_id"`
	Value    string        `json:"value"`
}

type blockObjectType string

const (
	blockObjectTypeSection blockObjectType = "section"
	blockObjectTypeDivider blockObjectType = "divider"
)

type Block struct {
	Type      blockObjectType `json:"type"`
	Text      *TextBlock      `json:"text,omitempty"`
	Fields    []*TextBlock    `json:"fields,omitempty"`
	Accessory *accessory      `json:"accessory,omitempty"`
}

func Divider() *Block {
	return &Block{
		Type: blockObjectTypeDivider,
	}
}

type sectionOption func(*Block)

// https://api.slack.com/reference/block-kit/blocks#section
func SectionFields(fields []*TextBlock, options ...sectionOption) *Block {
	b := Section(nil, fields...)
	for _, apply := range options {
		apply(b)
	}
	return b
}

// https://api.slack.com/reference/block-kit/blocks#section
func Section(text *TextBlock, fields ...*TextBlock) *Block {
	return &Block{
		Type:   blockObjectTypeSection,
		Text:   text,
		Fields: fields,
	}
}

// https://api.slack.com/reference/block-kit/block-elements#button
func WithButton(text *TextBlock, actionID, value string) sectionOption {
	return func(b *Block) {
		b.Accessory = &accessory{
			Text:     text,
			Type:     accessoryTypeButton,
			ActionID: actionID,
			Value:    value,
		}
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
	ReplaceOriginal bool         `json:"replace_original"`
	ResponseType    responseType `json:"response_type,omitempty"`
	Text            Text         `json:"text,omitempty"` // Text used in notifications as a fallback for Blocks
	Blocks          []*Block     `json:"blocks,omitempty"`
}

func newReplaceMessage(rt responseType, text Text, sections ...*Block) *message {
	return &message{
		ReplaceOriginal: true,
		ResponseType:    rt,
		Text:            text,
		Blocks:          sections,
	}
}

func newMessage(rt responseType, text Text, sections ...*Block) *message {
	return &message{
		ResponseType: rt,
		Text:         text,
		Blocks:       sections,
	}
}
