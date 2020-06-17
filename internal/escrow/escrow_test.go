// Copyright 2018 Square Inc.
//
// Use of this source code is governed by a GNU
// General Public License license version 3 that
// can be found in the LICENSE file.

package escrow

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPostCryptServer(t *testing.T) {
	expected := CryptServerData{
		Pass:      "2345.1234.6566.foo",
		Serialnum: "1234foobar",
		Hostname:  "testing.example.com",
		Username:  "tester",
	}

	mockCryptServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		if r.Method != "POST" {
			t.Errorf("expected 'POST' request, got '%s'", r.Method)
		}

		if r.URL.EscapedPath() != "/checkin/" {
			t.Errorf("expected request to '/checkin', got '%s'", r.URL.EscapedPath())
		}

		authusername, authpassword, _ := r.BasicAuth()
		if authusername != "" || authpassword != "" {
			t.Errorf("expected no basic auth got %q:%q", authusername, authpassword)
		}

		// cryptserver expects the following form data recovery_password, serial,
		// macname, username
		// see: https://github.com/grahamgilbert/Crypt-Server/blob/master/server/views.py#L442
		r.ParseForm()
		actual := r.Form.Get("recovery_password")
		if actual != expected.Pass {
			t.Errorf("expected 'recovery_password=%v' got %v", expected.Pass, actual)
		}

		actual = r.Form.Get("serial")
		if actual != expected.Serialnum {
			t.Errorf("expected 'serial=%v' got %v", expected.Serialnum, actual)
		}

		actual = r.Form.Get("macname")
		if actual != expected.Hostname {
			t.Errorf("expected 'macname=%v' got %v", expected.Hostname, actual)
		}

		actual = r.Form.Get("username")
		if actual != expected.Username {
			t.Errorf("expected 'username=%v' got %v", expected.Username, actual)
		}
	}))
	defer mockCryptServer.Close()

	endpoint := CryptServerInfo{
		Server: mockCryptServer.URL,
		URI:    "/checkin/",
	}

	resp, err := expected.PostCryptServer(endpoint)
	if err != nil {
		t.Errorf("errored posting escrow data to mock cryptserver with %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("errored posting escrow data to mock cryptserver expected 200 got %v", resp.StatusCode)
	}
}


func TestPostBasicAuthCryptServer(t *testing.T) {
	expected := CryptServerData{
		Pass:      "2345.1234.6566.foo",
		Serialnum: "1234foobar",
		Hostname:  "testing.example.com",
		Username:  "tester",
	}
	basicauth := []string{"DarthHelmet", "12345"}

	mockCryptServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		if r.Method != "POST" {
			t.Errorf("expected 'POST' request, got '%s'", r.Method)
		}

		if r.URL.EscapedPath() != "/checkin/" {
			t.Errorf("expected request to '/checkin', got '%s'", r.URL.EscapedPath())
		}

		authusername, authpassword, _ := r.BasicAuth()
		if authusername != basicauth[0] || authpassword != basicauth[1] {
			t.Errorf("expected username:password=%q, got '%s:%s'", basicauth, authusername, authpassword)
		}

		// cryptserver expects the following form data recovery_password, serial,
		// macname, username
		// see: https://github.com/grahamgilbert/Crypt-Server/blob/master/server/views.py#L442
		r.ParseForm()
		actual := r.Form.Get("recovery_password")
		if actual != expected.Pass {
			t.Errorf("expected 'recovery_password=%v' got %v", expected.Pass, actual)
		}

		actual = r.Form.Get("serial")
		if actual != expected.Serialnum {
			t.Errorf("expected 'serial=%v' got %v", expected.Serialnum, actual)
		}

		actual = r.Form.Get("macname")
		if actual != expected.Hostname {
			t.Errorf("expected 'macname=%v' got %v", expected.Hostname, actual)
		}

		actual = r.Form.Get("username")
		if actual != expected.Username {
			t.Errorf("expected 'username=%v' got %v", expected.Username, actual)
		}
	}))
	defer mockCryptServer.Close()

	endpoint := CryptServerInfo{
		Server: mockCryptServer.URL,
		URI:    "/checkin/",
		Username: basicauth[0],
		Password: basicauth[1],
	}

	resp, err := expected.PostCryptServer(endpoint)
	if err != nil {
		t.Errorf("errored posting escrow data to mock cryptserver with %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("errored posting escrow data to mock cryptserver expected 200 got %v", resp.StatusCode)
	}
}
