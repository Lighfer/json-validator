package json

import (
	"io/ioutil"
	"os"
	"testing"
)

// 1.json 2.json 3.json epub.json issue72.json issue74.json issue859.json comes from fastjson.

func TestValidate_1(t *testing.T) {
	f, err := os.Open("1.json")
	if err != nil {
		t.Fatal(err)
	}

	b, _ := ioutil.ReadAll(f)

	if err = Validate(string(b)); err != nil {
		t.Fatal(err)
	}
}

func TestValidate_2(t *testing.T) {
	f, err := os.Open("2.json")
	if err != nil {
		t.Fatal(err)
	}

	b, _ := ioutil.ReadAll(f)

	if err = Validate(string(b)); err != nil {
		t.Fatal(err)
	}
}

func TestValidate_3(t *testing.T) {
	f, err := os.Open("3.json")
	if err != nil {
		t.Fatal(err)
	}

	b, _ := ioutil.ReadAll(f)

	if err = Validate(string(b)); err != nil {
		t.Fatal(err)
	}
}

func TestValidate_epub(t *testing.T) {
	f, err := os.Open("epub.json")
	if err != nil {
		t.Fatal(err)
	}

	b, _ := ioutil.ReadAll(f)

	if err = Validate(string(b)); err != nil {
		t.Fatal(err)
	}
}

func TestValidate_issue72(t *testing.T) {
	f, err := os.Open("issue72.json")
	if err != nil {
		t.Fatal(err)
	}

	b, _ := ioutil.ReadAll(f)

	if err = Validate(string(b)); err != nil {
		t.Fatal(err)
	}
}

func TestValidate_issue74(t *testing.T) {
	f, err := os.Open("issue74.json")
	if err != nil {
		t.Fatal(err)
	}

	b, _ := ioutil.ReadAll(f)

	if err = Validate(string(b)); err != nil {
		t.Fatal(err)
	}
}

func TestValidate_issue859(t *testing.T) {
	f, err := os.Open("issue859.json")
	if err != nil {
		t.Fatal(err)
	}

	b, _ := ioutil.ReadAll(f)

	if err = Validate(string(b)); err != nil {
		t.Fatal(err)
	}
}
