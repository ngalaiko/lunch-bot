package users

type ID string

type User struct {
	ID   ID     `dynamodbav:"id" json:"id"`
	Name string `dynamodbav:"name" json:"name"`
}
