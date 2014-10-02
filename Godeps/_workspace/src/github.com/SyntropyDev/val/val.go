package val

import (
	"errors"
	"fmt"
	"strings"
)

type Validator struct {
	fVals []fieldValidator
}

type fieldValidator struct {
	k string
	v interface{}
	f ValFunc
}

func New() *Validator {
	return &Validator{fVals: []fieldValidator{}}
}

func (this *Validator) Add(k string, v interface{}, funcs ...ValFunc) {
	for _, f := range funcs {
		fVal := fieldValidator{k: k, v: v, f: f}
		this.fVals = append(this.fVals, fVal)
	}
}

func (this *Validator) Validate() (bool, map[string]error) {
	errs := map[string]error{}
	for _, fVal := range this.fVals {
		if err := fVal.f(fVal.k, fVal.v); nil != err {
			errs[fVal.k] = err
		}
	}
	valid := len(errs) == 0
	return valid, errs
}

func ErrorFromMap(errMap map[string]error) error {
	errs := []string{}
	for key, err := range errMap {
		nErr := fmt.Sprintf("%s - %s", key, err)
		errs = append(errs, nErr)
	}
	return errors.New(strings.Join(errs, ","))
}
