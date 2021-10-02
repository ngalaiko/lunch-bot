package request

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/hetiansu5/urlquery"
)

func Parse(req events.APIGatewayProxyRequest, dist interface{}) error {
	body, err := getBody(req)
	if err != nil {
		return fmt.Errorf("failed to parse body: %w", err)
	}

	contentType := req.Headers["content-type"]
	switch contentType {
	case "application/json":
		return parseJSON(body, dist)
	case "application/x-www-form-urlencoded":
		return parseQuery(body, dist)
	default:
		return fmt.Errorf("don't know how to handle '%s", contentType)
	}
}

func getBody(req events.APIGatewayProxyRequest) ([]byte, error) {
	if !req.IsBase64Encoded {
		return []byte(req.Body), nil
	}
	body, err := decodeBase64([]byte(req.Body))
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}
	return body, nil
}

func parseJSON(in []byte, dist interface{}) error {
	if err := json.Unmarshal(in, dist); err != nil {
		return fmt.Errorf("failed to unmarshal request: %w", err)
	}
	return nil
}

func parseQuery(in []byte, dist interface{}) error {
	if err := urlquery.Unmarshal(in, dist); err != nil {
		return fmt.Errorf("failed to unmarshal request: %w", err)
	}
	return nil
}

func decodeBase64(in []byte) ([]byte, error) {
	enc := base64.StdEncoding
	dbuf := make([]byte, enc.DecodedLen(len(in)))
	n, err := enc.Decode(dbuf, in)
	return dbuf[:n], err
}
