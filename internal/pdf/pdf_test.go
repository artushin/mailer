package pdf

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestEnvelope(t *testing.T) {
	b := bytes.NewBuffer(nil)

	err := GenerateEnvelope(&Envelope{
		From: []string{
			"Mr. Foo Boar",
			"12345 Street Dr.",
			"Apt 4304",
			"New York, NY, 10000",
		},
		To: []string{
			"Mr. Foo Boar",
			"12345 Street Dr.",
			"Apt 4304",
			"New York, NY, 10000",
		},
	}, b)
	if err != nil {
		t.Fatal(err)
	}

	if err := ioutil.WriteFile("./envelope.pdf", b.Bytes(), 0775); err != nil {
		t.Fatal(err)
	}
}

func TestLetter(t *testing.T) {
	b := bytes.NewBuffer(nil)

	err := GenerateLetter(&Letter{
		From: []string{
			"Mr. Foo Boar",
			"12345 Street Dr.",
			"Apt 4304",
			"New York, NY, 10000",
		},
		Person:    "Dr Boar Foo",
		Signature: "Mr. Foo Bar",
	}, b)
	if err != nil {
		t.Fatal(err)
	}

	if err := ioutil.WriteFile("./letter.pdf", b.Bytes(), 0775); err != nil {
		t.Fatal(err)
	}
}
