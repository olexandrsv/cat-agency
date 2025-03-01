package common

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Validator struct {
	c        *gin.Context
	err      InvalidFormError
	parseErr error
}

func NewValidator(c *gin.Context) *Validator {
	return &Validator{
		c: c,
	}
}

func (v *Validator) parseForm() {
	if v.c.Request.MultipartForm != nil {
		return
	}
	if err := v.c.Request.ParseMultipartForm(10 << 20); err != nil {
		v.parseErr = errors.New("can't parse form")
	}
}

func (v *Validator) GetIntFromURL(name string) int {
	raw := v.c.Param(name)
	if raw == "" {
		v.err.URLFieldNotExists(name)
		return 0
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		v.err.WrongURLFieldType(name, "int")
	}

	return value
}

func (v *Validator) GetInt(name string) int {
	v.parseForm()
	if v.parseErr != nil {
		return 0
	}
	raw, ok := v.c.GetPostForm(name)
	if !ok {
		v.err.FieldNotExists(name)
		return 0
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		v.err.WrongFieldType(name, "int")
	}
	return value
}

func (v *Validator) GetString(name string) string {
	v.parseForm()
	if v.parseErr != nil {
		return ""
	}
	value, ok := v.c.GetPostForm(name)
	if !ok {
		v.err.FieldNotExists(name)
		return ""
	}
	return value
}

func (v *Validator) GetFloat64(name string) float64 {
	v.parseForm()
	if v.parseErr != nil {
		return 0
	}
	raw, ok := v.c.GetPostForm(name)
	if !ok {
		v.err.FieldNotExists(name)
		return 0
	}
	value, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		v.err.WrongFieldType(name, "float64")
	}
	return value
}

func (v *Validator) Err() error {
	if v.err.Err == nil {
		return nil
	}
	return v.err
}
