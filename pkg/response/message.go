package response

import "fmt"

type textBlockType string

const (
	textBlockTypePlainText textBlockType = "plain_text"
	textBlockTypeMarkdown  textBlockType = "mrkdwn"
)

// https://api.slack.com/reference/messaging/composition-objects#text
type textBlock struct {
	Type textBlockType `json:"type"`
	Text string        `json:"text"`
}

func Markdown(format string, a ...interface{}) *textBlock {
	return &textBlock{
		Type: textBlockTypeMarkdown,
		Text: fmt.Sprintf(format, a...),
	}
}

func PlainText(format string, a ...interface{}) *textBlock {
	return &textBlock{
		Type: textBlockTypePlainText,
		Text: fmt.Sprintf(format, a...),
	}
}

type Option struct {
	Text  *textBlock `json:"text"`
	Value string     `json:"value"`
}

type accessoryType string

const (
	accessoryTypeStaticSelect accessoryType = "static_select"
)

type accessory struct {
	ActionID    string        `json:"action_id"`
	Type        accessoryType `json:"type"`
	Placeholder *textBlock    `json:"placeholder"`
	Options     []*Option     `json:"options"`
}

type blockObjectType string

const (
	blockObjectTypeSection blockObjectType = "section"
	blockObjectTypeDivider blockObjectType = "divider"
)

type Block struct {
	Type      blockObjectType `json:"type"`
	Text      *textBlock      `json:"text,omitempty"`
	Fields    []*textBlock    `json:"fields,omitempty"`
	Accessory *accessory      `json:"accessory,omitempty"`
}

func Divider() *Block {
	return &Block{
		Type: blockObjectTypeDivider,
	}
}

type sectionOption func(*Block)

// https://api.slack.com/reference/block-kit/blocks#section
func Section(text *textBlock, fields ...*textBlock) *Block {
	b := &Block{
		Type:   blockObjectTypeSection,
		Text:   text,
		Fields: fields,
	}
	return b
}

// https://api.slack.com/reference/block-kit/blocks#section
func Select(text *textBlock, options ...sectionOption) *Block {
	b := Section(text)
	for _, apply := range options {
		apply(b)
	}
	return b
}

// https://api.slack.com/reference/block-kit/block-elements#static_select
func SelectOption(text *textBlock, value string) *Option {
	return &Option{
		Text:  text,
		Value: value,
	}
}

// https://api.slack.com/reference/block-kit/block-elements#static_select
func Static(placeholeer *textBlock, actionID string, options ...*Option) sectionOption {
	return func(b *Block) {
		b.Accessory = &accessory{
			ActionID:    actionID,
			Type:        accessoryTypeStaticSelect,
			Placeholder: placeholeer,
			Options:     options,
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
	ResponseType responseType `json:"response_type,omitempty"`
	Blocks       []*Block     `json:"blocks,omitempty"`
}

func newMessage(rt responseType, sections ...*Block) *message {
	return &message{
		ResponseType: rt,
		Blocks:       sections,
	}
}
