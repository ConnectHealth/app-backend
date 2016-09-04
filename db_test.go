package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPeople(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(1, 1, "")

	patients, err := getPatients()
	if err != nil {
		t.Fail()
	}
	assert.Equal(cap(patients), 1)
	assert.Equal(patients[0].FirstName, "kal")
}
