package common

import (
	"fmt"
	"log"
	"net/http"

	"errors"

	"github.com/gin-gonic/gin"
)

type CustomError struct {
	Err error
}

func (e *CustomError) append(err error){
	if e.Err == nil {
		e.Err = err
		return
	}
	e.Err = errors.Join(e.Err, err)
}

func (e CustomError) Error() string {
	return fmt.Sprintf("%+v", e.Err)
}

func newCustomError(err error) CustomError {
	return CustomError{
		Err: err,
	}
}

type DatabaseError struct {
	CustomError
}

func NewDatabseError(err error) DatabaseError {
	return DatabaseError{
		CustomError: newCustomError(err),
	}
}

type InternalError struct{
	CustomError
}

func NewInternalError(err error) InternalError {
	return InternalError{
		CustomError: newCustomError(err),
	}
}

type NoRowsError struct{
	CustomError
}

func NewNoRowsError(err error) NoRowsError {
	return NoRowsError{
		CustomError: newCustomError(err),
	}
}

type HTTPRequestError struct{
	CustomError
}

func NewHTTPRequestError(err error) HTTPRequestError {
	return HTTPRequestError{
		CustomError: newCustomError(err),
	}
}

type JSONError struct{
	CustomError
}

func NewJSONError(err error) JSONError {
	return JSONError{
		CustomError: newCustomError(err),
	}
}

type InvalidBreedError struct {
	CustomError
}

func NewInvalidBreedError(breed string) InvalidBreedError{
	return InvalidBreedError{
		CustomError: newCustomError(fmt.Errorf("invalid breed: %s", breed)),
	}
}

type InvalidFormError struct{
	CustomError
}

func NewInvalidFieldError(err error) InvalidFormError {
	return InvalidFormError{
		CustomError: newCustomError(err),
	}
}

func (e *InvalidFormError) URLFieldNotExists(field string){
	e.append(fmt.Errorf("url field: %s doesn't exists", field))
}

func (e *InvalidFormError) WrongURLFieldType(field, expectedType string) {
	e.append(fmt.Errorf("can't convert url field: %s to type: %s", field, expectedType))
}

func (e *InvalidFormError) FieldNotExists(field string){
	e.append(fmt.Errorf("field: %s doesn't exists in form", field))
}

func (e *InvalidFormError) WrongFieldType(field, expectedType string) {
	e.append(fmt.Errorf("can't convert form field: %s to type: %s", field, expectedType))
}

func WriteError(c *gin.Context, err error){
	log.Println(err)
	code := http.StatusInternalServerError
	message := "Internal error"

	switch err.(type){
	case InvalidFormError:
		code = http.StatusBadRequest
		message = fmt.Sprintf("Invalid data in form:\n%s", err.Error())
	case InvalidBreedError:
		code = http.StatusBadRequest
		message = err.Error()
	}

	c.JSON(code, message)
}