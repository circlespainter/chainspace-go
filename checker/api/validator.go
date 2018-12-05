package api

import (
	"errors"
	"fmt"
)

// Validator ...
type Validator struct{}

// TransactionValidator ...
type TransactionValidator interface {
	Validate(tx *Transaction) error
}

// Validate ...
func (v Validator) Validate(tx *Transaction) error {
	for _, tc := range tx.Traces {
		err := v.ValidateTrace(&tc, tx.Mappings)
		if err != nil {
			return err
		}
	}

	return nil
}

// ValidateTrace ...
func (v Validator) ValidateTrace(tc *Trace, mappings map[string]interface{}) error {
	for _, d := range tc.Dependencies {
		t := Trace(d)
		err := v.ValidateTrace(&t, mappings)
		if err != nil {
			return err
		}
	}

	if len(tc.InputObjectVersionIDs) <= 0 {
		return errors.New("Missing input version ID")
	}

	for _, v := range tc.InputObjectVersionIDs {
		_, ok := mappings[v]
		if !ok {
			return fmt.Errorf("Missing object mapping for key [%v]", v)
		}
	}

	for _, v := range tc.InputReferenceVersionIDs {
		_, ok := mappings[v]
		if !ok {
			return fmt.Errorf("Missing reference mapping for key [%v]", v)
		}
	}

	return nil
}