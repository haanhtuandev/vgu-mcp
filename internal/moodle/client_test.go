package moodle

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestGetEnrolledCourses(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/webservice/rest/server.php" {
			t.Fatalf("path = %q", request.URL.Path)
		}
		if err := request.ParseForm(); err != nil {
			t.Fatal(err)
		}
		want := url.Values{"userid": {"7"}, "wstoken": {"secret"}, "moodlewsrestformat": {"json"}, "wsfunction": {"core_enrol_get_users_courses"}}
		if request.Form.Encode() != want.Encode() {
			t.Fatalf("form = %q, want %q", request.Form.Encode(), want.Encode())
		}
		_, _ = writer.Write([]byte(`[{"id":1,"fullname":"Algorithms","shortname":"ALG","categoryid":2}]`))
	}))
	defer server.Close()

	courses, err := NewClient(server.URL, "secret").GetEnrolledCourses(context.Background(), 7)
	if err != nil {
		t.Fatal(err)
	}
	if len(courses) != 1 || courses[0].FullName != "Algorithms" {
		t.Fatalf("courses = %#v", courses)
	}
}

func TestGetSiteInfoReturnsMoodleError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte(`{"exception":"moodle_exception","errorcode":"invalidtoken","message":"Invalid token"}`))
	}))
	defer server.Close()

	_, err := NewClient(server.URL, "secret").GetSiteInfo(context.Background())
	if err == nil || !strings.Contains(err.Error(), "invalidtoken") {
		t.Fatalf("error = %v", err)
	}
}
