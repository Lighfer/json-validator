# Usage:
validator has only one function: Validate, it returns nil if JSON string is valid, otherwise reutrns the ErrJSON.
```Go
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
```