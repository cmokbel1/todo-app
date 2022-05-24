package crypto_test

import (
	"testing"

	"github.com/cmokbel1/todo-app/backend/crypto"
)

func TestCreateAndCompare(t *testing.T) {
	in := "password"
	out, err := crypto.CreateHash(in)
	if err != nil {
		t.Fatal(err)
	}

	if match, err := crypto.ComparePasswordAndHash(in, out); err != nil {
		t.Fatal(err)
	} else if !match {
		t.Fatalf("password %v and hash %v should match", in, out)
	}

	in = "pa$$word"
	if match, err := crypto.ComparePasswordAndHash(in, out); err != nil {
		t.Fatal(err)
	} else if match {
		t.Fatalf("password %v and hash %v should not match", in, out)
	}
}

func TestRandomString(t *testing.T) {
	set := make(map[string]struct{})
	for i := 0; i < 5000; i++ {
		str := crypto.RandomString()
		if _, ok := set[str]; ok {
			t.Fatalf("random string collision: %v", str)
		} else if len(str) != 16 {
			t.Fatalf("want random string len %v got %v", 16, len(str))
		}
		set[str] = struct{}{}
	}
}
