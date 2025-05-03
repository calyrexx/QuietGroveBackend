package errorspkg

import "fmt"

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

type ErrInvalidInterval uint8

func (e ErrInvalidInterval) Error() string {
	return fmt.Sprintf("invalid interval [%d]", e)
}

func (err ErrConstructorDependencies) Error() string {
	return fmt.Sprintf("constructor [%s] got not correct dependency [%s] is [%s]", err.constructor, err.dependency, err.state)
}

func NewErrConstructorDependencies(constructor, dependency, state string) error {
	return ErrConstructorDependencies{constructor: constructor, dependency: dependency, state: state}
}

type ErrSomethingIsEmpty string

func (e ErrSomethingIsEmpty) Error() string {
	return fmt.Sprintf("[%s] is empty", string(e))
}

type ErrUnitIsNil struct {
	unit string
}

func (err ErrUnitIsNil) Error() string {
	return fmt.Sprintf("is nil: %s", err.unit)
}
func NewErrUnitIsNil(unit string) error {
	return ErrUnitIsNil{unit: unit}
}

type ErrCronFunc struct {
	err error
}

func (e ErrCronFunc) Error() string {
	return fmt.Sprintf("CronFunc - error: %s", e.err)
}

func NewErrCronFunc(err error) error {
	return ErrCronFunc{err: err}
}
