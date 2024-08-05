package utils

import (
	"encoding/json"
	"errors"
	"strconv"
)

type SubjectID int

func (s *SubjectID) UnmarshalJSON(data []byte) error {
	var intID int
	var stringID string

	if err := json.Unmarshal(data, &intID); err == nil {
		*s = SubjectID(intID)
		return nil
	}

	if err := json.Unmarshal(data, &stringID); err == nil {
		id, err := strconv.Atoi(stringID)
		if err != nil {
			return err
		}
		*s = SubjectID(id)
		return nil
	}

	return errors.New("invalid subject ID")
}
