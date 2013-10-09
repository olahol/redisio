package redisio

import (
	"io/ioutil"
	"os"
	"testing"
	"testing/quick"
)

func TestRequest(t *testing.T) {
	f := func(request []string) bool {
		if len(request) < 1 {
			return true
		}

		temp, err := ioutil.TempFile("", "test")

		if err != nil {
			return false
		}

		name := temp.Name()

		rd1 := NewWriter(temp)
		err = rd1.WriteRequest(request)
        rd1.Flush()

		if err != nil {
			return false
		}

		temp.Close()

		temp, err = os.Open(name)

		if err != nil {
			return false
		}

		rd2 := NewReader(temp)
		reply, err := rd2.ReadRequest()

		if err != nil {
			return false
		}

		temp.Close()

		for i := range reply {
			if reply[i] != request[i] {
				return false
			}
		}

		os.Remove(name)

		return true
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func TestReplyBulk(t *testing.T) {
	f := func(bulk string) bool {
		temp, err := ioutil.TempFile("", "test")

		if err != nil {
			return false
		}

		name := temp.Name()

		rd1 := NewWriter(temp)
		err = rd1.WriteBulk(bulk)
        rd1.Flush()

		if err != nil {
			return false
		}

		temp.Close()

		temp, err = os.Open(name)

		if err != nil {
			return false
		}

		rd2 := NewReader(temp)
		_, reply, err := rd2.ReadReply()

		if err != nil {
			return false
		}

		temp.Close()

		if reply[0] != bulk {
			return false
		}

		os.Remove(name)

		return true
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}
