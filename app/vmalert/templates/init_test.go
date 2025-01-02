package templates

import (
	"fmt"
	"math"
	"net/url"
	"os"
	"strings"
	"testing"
	textTpl "text/template"
)

func TestMain(m *testing.M) {
	if err := Init([]string{}, nil, url.URL{}); err != nil {
		fmt.Println("failed to load template for test")
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func TestTemplateFuncs_Match(t *testing.T) {
	funcs := templateFuncs()
	// check "match" func
	matchFunc := funcs["match"].(func(pattern, s string) (bool, error))
	if _, err := matchFunc("invalid[regexp", "abc"); err == nil {
		t.Fatalf("expecting non-nil error on invalid regexp")
	}
	ok, err := matchFunc("abc", "def")
	if err != nil {
		t.Fatalf("unexpected error")
	}
	if ok {
		t.Fatalf("unexpected match")
	}
	ok, err = matchFunc("a.+b", "acsdb")
	if err != nil {
		t.Fatalf("unexpected error")
	}
	if !ok {
		t.Fatalf("unexpected mismatch")
	}
}

func TestTemplateFuncs_Formatting(t *testing.T) {
	f := func(funcName string, p any, resultExpected string) {
		t.Helper()

		funcs := templateFuncs()
		v := funcs[funcName]
		fLocal := v.(func(s any) (string, error))
		result, err := fLocal(p)
		if err != nil {
			t.Fatalf("unexpected error for %s(%f): %s", funcName, p, err)
		}
		if result != resultExpected {
			t.Fatalf("unexpected result for %s(%f); got\n%s\nwant\n%s", funcName, p, result, resultExpected)
		}
	}

	f("humanize1024", float64(0), "0")
	f("humanize1024", math.Inf(0), "+Inf")
	f("humanize1024", math.NaN(), "NaN")
	f("humanize1024", float64(127087), "124.1ki")
	f("humanize1024", float64(130137088), "124.1Mi")
	f("humanize1024", float64(133260378112), "124.1Gi")
	f("humanize1024", float64(136458627186688), "124.1Ti")
	f("humanize1024", float64(139733634239168512), "124.1Pi")
	f("humanize1024", float64(143087241460908556288), "124.1Ei")
	f("humanize1024", float64(146521335255970361638912), "124.1Zi")
	f("humanize1024", float64(150037847302113650318245888), "124.1Yi")
	f("humanize1024", float64(153638755637364377925883789312), "1.271e+05Yi")

	f("humanize", float64(127087), "127.1k")
	f("humanize", float64(136458627186688), "136.5T")

	f("humanizeDuration", 1, "1s")
	f("humanizeDuration", 0.2, "200ms")
	f("humanizeDuration", 42000, "11h 40m 0s")
	f("humanizeDuration", 16790555, "194d 8h 2m 35s")

	f("humanizePercentage", 1, "100%")
	f("humanizePercentage", 0.8, "80%")
	f("humanizePercentage", 0.015, "1.5%")

	f("humanizeTimestamp", 1679055557, "2023-03-17 12:19:17 +0000 UTC")
}

func TestTemplateFuncs_StringConversion(t *testing.T) {
	f := func(funcName, s, resultExpected string) {
		t.Helper()

		funcs := templateFuncs()
		v := funcs[funcName]
		fLocal := v.(func(s string) string)
		result := fLocal(s)
		if result != resultExpected {
			t.Fatalf("unexpected result for %s(%q); got\n%s\nwant\n%s", funcName, s, result, resultExpected)
		}
	}

	f("title", "foo bar", "Foo Bar")
	f("toUpper", "foo", "FOO")
	f("toLower", "FOO", "foo")
	f("pathEscape", "foo/bar\n+baz", "foo%2Fbar%0A+baz")
	f("queryEscape", "foo+bar\n+baz", "foo%2Bbar%0A%2Bbaz")
	f("jsonEscape", `foo{bar="baz"}`+"\n + 1", `"foo{bar=\"baz\"}\n + 1"`)
	f("quotesEscape", `foo{bar="baz"}`+"\n + 1", `foo{bar=\"baz\"}\n + 1`)
	f("htmlEscape", "foo < 10\nabc", "foo &lt; 10\nabc")
	f("crlfEscape", "foo\nbar\rx", `foo\nbar\rx`)
	f("stripPort", "foo", "foo")
	f("stripPort", "foo:1234", "foo")
	f("stripDomain", "foo.bar.baz", "foo")
	f("stripDomain", "foo.bar:123", "foo:123")
}

func TestTemplatesLoad_Success(t *testing.T) {
	f := func(pathPatterns []string, expectedTmpl *textTpl.Template) {
		t.Helper()

		masterTmplOrig := masterTmpl
		defer func() {
			masterTmpl = masterTmplOrig
		}()
		masterTmpl = nil

		if err := LoadTemplateFile(pathPatterns); err != nil {
			t.Fatalf("cannot load templates: %s", err)
		}
		if !isTemplatesTheSame(masterTmpl, expectedTmpl) {
			t.Fatalf("unexpected template\ngot\n%+v\nwant\n%+v", masterTmpl, expectedTmpl)
		}
	}

	// non existing path
	pathPatterns := []string{
		"templates/non-existing/good-*.tpl",
		"templates/absent/good-*.tpl",
	}
	expectedTmpl := textTpl.Must(newTemplate("").Parse(""))
	f(pathPatterns, expectedTmpl)

	// existing path
	pathPatterns = []string{
		"templates/test/good0-*.tpl",
	}
	expectedTmpl = textTpl.Must(newTemplate("").Parse(`
		{{- define "good0-test.tpl" -}}{{- end -}}
		{{- define "test.0" -}}
			{{ printf "Hello %s!" externalURL }}
		{{- end -}}
		{{- define "test.2" -}}
			{{ printf "Hello %s!" externalURL }}
		{{- end -}}
		{{- define "test.3" -}}
			{{ printf "Hello %s!" externalURL }}
		{{- end -}}
	`))
	f(pathPatterns, expectedTmpl)

	// template update
	pathPatterns = []string{
		"templates/other/nested/good0-*.tpl",
	}
	expectedTmpl = textTpl.Must(newTemplate("").Parse(`
		{{- define "good0-test.tpl" -}}{{- end -}}
		{{- define "test.0" -}}
			{{ printf "Hello %s!" externalURL }}
		{{- end -}}
		{{- define "test.1" -}}
			{{ printf "Hello %s!" externalURL }}
		{{- end -}}
		{{- define "test.3" -}}
			{{ printf "Hello %s!" externalURL }}
		{{- end -}}
	`))
	f(pathPatterns, expectedTmpl)
}

func TestTemplatesLoad_Failure(t *testing.T) {
	f := func(pathPatterns []string, expectedErrStr string) {
		t.Helper()

		err := LoadTemplateFile(pathPatterns)
		if err == nil {
			t.Fatalf("expecting non-nil error")
		}

		errStr := err.Error()
		if !strings.Contains(errStr, expectedErrStr) {
			t.Fatalf("the returned error %q doesn't contain %q", errStr, expectedErrStr)
		}
	}

	// load template with syntax error
	f([]string{
		"templates/other/nested/bad0-*.tpl",
		"templates/test/good0-*.tpl",
	}, "failed to parse template glob")
}

func TestTemplatesReload(t *testing.T) {
	masterTmplOrig := masterTmpl
	defer func() {
		masterTmpl = masterTmplOrig
	}()
	masterTmpl = nil

	// load with non existing path
	pathPatterns := []string{
		"templates/non-existing/good-*.tpl",
		"templates/absent/good-*.tpl",
	}
	if err := LoadTemplateFile(pathPatterns); err != nil {
		t.Fatalf("cannot load templates: %s", err)
	}
	tpl1 := GetCurrentTmpl()

	// reload with existing path
	pathPatterns = []string{
		"templates/test/good0-*.tpl",
	}
	if err := LoadTemplateFile(pathPatterns); err != nil {
		t.Fatalf("cannot load templates: %s", err)
	}
	tpl2 := GetCurrentTmpl()
	if isTemplatesTheSame(tpl1, tpl2) {
		t.Fatalf("tpl1 should be different from tpl2")
	}

	// reload the same path
	pathPatterns = []string{
		"templates/test/good0-*.tpl",
	}
	if err := LoadTemplateFile(pathPatterns); err != nil {
		t.Fatalf("cannot load templates: %s", err)
	}
	tpl3 := GetCurrentTmpl()
	if !isTemplatesTheSame(tpl2, tpl3) || tpl2.Name() != tpl3.Name() {
		t.Fatalf("tpl3 should be the same as tpl2")
	}
}

func TestIsTemplateTheSame(t *testing.T) {
	f := func(tmpl1, tmpl2 *textTpl.Template, isTheSame bool) {
		t.Helper()

		if isTemplatesTheSame(tmpl1, tmpl2) != isTheSame {
			t.Fatalf("unexpected result for isTemplatesTheSame")
		}
	}

	tmpl1 := textTpl.Must(newTemplate("t1").Parse("{{- define \"test\" -}}{{- end -}}"))
	tmpl2 := textTpl.Must(newTemplate("t2").Parse("{{- define \"test\" -}}{{- end -}}"))
	f(tmpl1, tmpl2, true)

	tmpl1, _ = tmpl1.Parse("{{- define \"test2\" -}}{{- end -}}")
	f(tmpl1, tmpl2, false)

	tmpl2, _ = tmpl2.Parse("{{- define \"test3\" -}}{{- end -}}")
	f(tmpl1, tmpl2, false)
}
