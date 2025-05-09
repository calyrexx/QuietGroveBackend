package errorspkg

import (
	"fmt"
	"time"
)

type ErrViperReadInConfig struct {
	errorMsg error
}

func (e *ErrViperReadInConfig) Error() string {
	return fmt.Sprintf("viper.ReadInConfig (conf) - error: %v", e.errorMsg)
}

func NewErrViperReadInConfig(errorMsg error) *ErrViperReadInConfig {
	return &ErrViperReadInConfig{errorMsg: errorMsg}
}

type ErrReadConfigViper struct {
	section  string
	errorMsg error
}

func (e *ErrReadConfigViper) Error() string {
	return fmt.Sprintf("read config in section %s - error: %v", e.section, e.errorMsg)
}

func NewErrReadConfigViper(section string, errorMsg error) *ErrReadConfigViper {
	return &ErrReadConfigViper{section: section, errorMsg: errorMsg}
}

type ErrConstructorDependencies struct {
	constructor string
	dependency  string
	state       string
}

func (err ErrConstructorDependencies) Error() string {
	return fmt.Sprintf("constructor [%s] got not correct dependency [%s] is [%s]", err.constructor, err.dependency, err.state)
}

func NewErrConstructorDependencies(constructor, dependency, state string) error {
	return ErrConstructorDependencies{constructor: constructor, dependency: dependency, state: state}
}

type ErrRepoFailed struct {
	operation string
	method    string
	errorMsg  error
}

func (err ErrRepoFailed) Error() string {
	return fmt.Sprintf("operation [%s] in method [%s] failed: %s", err.operation, err.method, err.errorMsg)
}

func NewErrRepoFailed(operation, method string, err error) error {
	return &ErrRepoFailed{
		operation: operation,
		method:    method,
		errorMsg:  err,
	}
}

type ErrRepoNotFound struct {
	unit   string
	id     string
	method string
}

func (err ErrRepoNotFound) Error() string {
	return fmt.Sprintf("[%s] with id [%s] in method [%s] not found", err.unit, err.id, err.method)
}

func NewErrRepoNotFound(unit, id, method string) error {
	return &ErrRepoNotFound{
		unit:   unit,
		id:     id,
		method: method,
	}
}

type ErrHouseUnavailable struct {
	HouseID  int
	CheckIn  time.Time
	CheckOut time.Time
}

func (err ErrHouseUnavailable) Error() string {
	return fmt.Sprintf("house [%d] unavailable, checkIn: %v, checkOut: %v", err.HouseID, err.CheckIn, err.CheckOut)
}

func NewErrHouseUnavailable(houseID int, checkIn, checkOut time.Time) error {
	return &ErrHouseUnavailable{
		HouseID:  houseID,
		CheckIn:  checkIn,
		CheckOut: checkOut,
	}
}
