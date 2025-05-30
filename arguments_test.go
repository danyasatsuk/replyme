package replyme

import (
	"github.com/go-faker/faker/v4"
	"testing"
)

func TestArgument(t *testing.T) {
	name := faker.Word()
	value := faker.Sentence()

	arg := Argument{
		Name:  name,
		value: value,
	}

	arg.setValue("test")

	if arg.GetValue() != "test" {
		t.Fatal("value mismatch")
	}
}
