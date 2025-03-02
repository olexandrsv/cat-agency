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

func (e *CustomError) append(err error) {
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

type InternalError struct {
	CustomError
}

func NewInternalError(err error) InternalError {
	return InternalError{
		CustomError: newCustomError(err),
	}
}

type NoRowsError struct {
	CustomError
}

func NewNoRowsError(err error) NoRowsError {
	return NoRowsError{
		CustomError: newCustomError(err),
	}
}

type HTTPRequestError struct {
	CustomError
}

func NewHTTPRequestError(err error) HTTPRequestError {
	return HTTPRequestError{
		CustomError: newCustomError(err),
	}
}

type JSONError struct {
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

func NewInvalidBreedError(breed string) InvalidBreedError {
	return InvalidBreedError{
		CustomError: newCustomError(fmt.Errorf("invalid breed: %s", breed)),
	}
}

type InvalidFormError struct {
	CustomError
}

func NewInvalidFieldError(err error) InvalidFormError {
	return InvalidFormError{
		CustomError: newCustomError(err),
	}
}

func (e *InvalidFormError) URLFieldNotExists(field string) {
	e.append(fmt.Errorf("url field: %s doesn't exists", field))
}

func (e *InvalidFormError) WrongURLFieldType(field, expectedType string) {
	e.append(fmt.Errorf("can't convert url field: %s to type: %s", field, expectedType))
}

func (e *InvalidFormError) FieldNotExists(field string) {
	e.append(fmt.Errorf("field: %s doesn't exists in form", field))
}

func (e *InvalidFormError) WrongFieldType(field, expectedType string) {
	e.append(fmt.Errorf("can't convert form field: %s to type: %s", field, expectedType))
}

func WriteError(c *gin.Context, err error) {
	log.Println(err)
	code := http.StatusInternalServerError
	message := "Internal error"

	switch err.(type) {
	case InvalidFormError:
		code = http.StatusBadRequest
		message = fmt.Sprintf("Invalid data in form:\n%s", err.Error())
	case InvalidBreedError, MissionAssignedError, WrongCatIDError, TargetCompletedError, MissionCompletedError,
		ManyMissionTargetsError, JSONParseError, FewMissionTargetsError:
		code = http.StatusBadRequest
		message = err.Error()
	case NoRowsError:
		code = http.StatusBadRequest
		message = "Invalid data"
	}

	c.JSON(code, message)
}

type MissionAssignedError struct {
	CustomError
}

func NewMissionAssignedError() MissionAssignedError {
	return MissionAssignedError{
		CustomError: newCustomError(errors.New("mission is already assigned to cat")),
	}
}

type WrongCatIDError struct {
	CustomError
}

func NewWrongCatIDError(catID int) WrongCatIDError {
	return WrongCatIDError{
		CustomError: newCustomError(fmt.Errorf("no cat with id: %d", catID)),
	}
}

type TargetCompletedError struct {
	CustomError
}

func NewTargetCompletedError(targetID int) TargetCompletedError {
	return TargetCompletedError{
		CustomError: newCustomError(fmt.Errorf("target with id: %d has been already completed", targetID)),
	}
}

type MissionCompletedError struct {
	CustomError
}

func NewMissionCompletedError(missionID int) MissionCompletedError {
	return MissionCompletedError{
		CustomError: newCustomError(fmt.Errorf("mission with id: %d has been already completed", missionID)),
	}
}

type FewMissionTargetsError struct {
	CustomError
}

func NewFewMissionTargetsError() FewMissionTargetsError {
	return FewMissionTargetsError{
		CustomError: newCustomError(errors.New("mission have only 1 target")),
	}
}

type ManyMissionTargetsError struct {
	CustomError
}

func NewManyMissionTargetsError() ManyMissionTargetsError {
	return ManyMissionTargetsError{
		CustomError: newCustomError(errors.New("mission already have 3 targets")),
	}
}

type JSONParseError struct {
	CustomError
}

func NewJSONParseError(err error) JSONParseError {
	return JSONParseError{
		CustomError: newCustomError(err),
	}
}

type TargetsDublicateError struct {
	CustomError
}

func NewTargetsDublicateError() TargetsDublicateError {
	return TargetsDublicateError{
		CustomError: newCustomError(errors.New("mission targets should be unique")),
	}
}
