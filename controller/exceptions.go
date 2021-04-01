package controller

import "fmt"

type DetectionError struct {
	msg string
}

func (e *DetectionError) Error() string {  
    return fmt.Sprintf("unable to detect controller: %s", e.msg)
}

type UnexpectedError struct {
	msg string
}

func (e *UnexpectedError) Error() string {  
    return fmt.Sprintf("unexpected error: %s", e.msg)
}

type UnexpectedResponse struct {
	msg string
}

func (e *UnexpectedResponse) Error() string {  
    return fmt.Sprintf("unexpected response: %s", e.msg)
}