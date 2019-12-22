package gouldian

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Output struct {
	status  int
	headers map[string]string
	body    string
}

func (err Output) Error() string {
	return err.body
}

func (out *Output) Json(val interface{}) *Output {
	body, _ := json.Marshal(val)
	out.body = string(body)
	return out
	// return &Output{200, map[string]string{}, string(body)}
}

func (out *Output) Text(text string) *Output {
	out.body = text
	return out
	// return &Output{200, map[string]string{}, text}
}

func Ok() *Output {
	return &Output{200, map[string]string{}, ""}
}

func Success(status int) *Output {
	return &Output{status, map[string]string{}, ""}
}

func (out *Output) With(head string, value string) *Output {
	out.headers[head] = value
	return out
}

//
type Issue struct {
	Type    string      `json:"type"`
	Status  int         `json:"status"`
	Title   string      `json:"title"`
	Details interface{} `json:"details"`
}

func (err Issue) Error() string {
	return fmt.Sprintf(strconv.Itoa(err.Status) + ": " + err.Title)
}

func Unauthorized(spec interface{}) Issue {
	return Issue{"https://httpstatuses.com/401", 401, "Unauthorized", spec}
}
