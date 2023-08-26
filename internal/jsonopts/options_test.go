// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package jsonopts_test

import (
	"reflect"
	"testing"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/internal/jsonflags"
	. "github.com/go-json-experiment/json/internal/jsonopts"
	"github.com/go-json-experiment/json/jsontext"
)

func makeFlags(f ...jsonflags.Bools) (fs jsonflags.Flags) {
	for _, f := range f {
		fs.Set(f)
	}
	return fs
}

func TestCopyCoderOptions(t *testing.T) {
	got := &Struct{
		Flags:        makeFlags(jsonflags.Indent|jsonflags.AllowInvalidUTF8|0, jsonflags.Expand|jsonflags.AllowDuplicateNames|jsonflags.Unmarshalers|1),
		CoderValues:  CoderValues{Indent: "    "},
		ArshalValues: ArshalValues{Unmarshalers: "something"},
	}
	src := &Struct{
		Flags:        makeFlags(jsonflags.Indent|jsonflags.Deterministic|jsonflags.Marshalers|1, jsonflags.Expand|0),
		CoderValues:  CoderValues{Indent: "\t"},
		ArshalValues: ArshalValues{Marshalers: "something"},
	}
	want := &Struct{
		Flags:        makeFlags(jsonflags.AllowInvalidUTF8|jsonflags.Expand|0, jsonflags.Indent|jsonflags.AllowDuplicateNames|jsonflags.Unmarshalers|1),
		CoderValues:  CoderValues{Indent: "\t"},
		ArshalValues: ArshalValues{Unmarshalers: "something"},
	}
	got.CopyCoderOptions(src)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("CopyCoderOptions:\n\tgot:  %+v\n\twant: %+v", got, want)
	}
}

func TestJoin(t *testing.T) {
	tests := []struct {
		in   Options
		want *Struct
	}{{
		in:   jsonflags.AllowInvalidUTF8 | 1,
		want: &Struct{Flags: makeFlags(jsonflags.AllowInvalidUTF8 | 1)},
	}, {
		in: jsonflags.Expand | 0,
		want: &Struct{
			Flags: makeFlags(jsonflags.AllowInvalidUTF8|1, jsonflags.Expand|0)},
	}, {
		in: Indent("\t"), // implicitly sets Expand=true
		want: &Struct{
			Flags:       makeFlags(jsonflags.AllowInvalidUTF8 | jsonflags.Expand | jsonflags.Indent | 1),
			CoderValues: CoderValues{Indent: "\t"},
		},
	}, {
		in: &Struct{
			Flags: makeFlags(jsonflags.Expand|jsonflags.EscapeForJS|0, jsonflags.AllowInvalidUTF8|1),
		},
		want: &Struct{
			Flags:       makeFlags(jsonflags.AllowInvalidUTF8|jsonflags.Indent|1, jsonflags.Expand|jsonflags.EscapeForJS|0),
			CoderValues: CoderValues{Indent: "\t"},
		},
	}, {
		in: &DefaultOptionsV1, want: &DefaultOptionsV1, // v1 fully replaces before
	}, {
		in: &DefaultOptionsV2, want: &DefaultOptionsV2}, // v2 fully replaces before
	}
	got := new(Struct)
	for i, tt := range tests {
		got.Join(tt.in)
		if !reflect.DeepEqual(got, tt.want) {
			t.Fatalf("%d: Join:\n\tgot:  %+v\n\twant: %+v", i, got, tt.want)
		}
	}
}

func TestGet(t *testing.T) {
	opts := &Struct{
		Flags:        makeFlags(jsonflags.Indent|jsonflags.Deterministic|jsonflags.Marshalers|1, jsonflags.Expand|0),
		CoderValues:  CoderValues{Indent: "\t"},
		ArshalValues: ArshalValues{Marshalers: new(json.Marshalers)},
	}
	if v, ok := json.GetOption(nil, jsontext.AllowDuplicateNames); v || ok {
		t.Errorf("GetOption(..., AllowDuplicateNames) = (%v, %v), want (false, false)", v, ok)
	}
	if v, ok := json.GetOption(jsonflags.AllowInvalidUTF8|0, jsontext.AllowDuplicateNames); v || ok {
		t.Errorf("GetOption(..., AllowDuplicateNames) = (%v, %v), want (false, false)", v, ok)
	}
	if v, ok := json.GetOption(jsonflags.AllowDuplicateNames|0, jsontext.AllowDuplicateNames); v || !ok {
		t.Errorf("GetOption(..., AllowDuplicateNames) = (%v, %v), want (false, true)", v, ok)
	}
	if v, ok := json.GetOption(jsonflags.AllowDuplicateNames|1, jsontext.AllowDuplicateNames); !v || !ok {
		t.Errorf("GetOption(..., AllowDuplicateNames) = (%v, %v), want (true, true)", v, ok)
	}
	if v, ok := json.GetOption(Indent(""), jsontext.AllowDuplicateNames); v || ok {
		t.Errorf("GetOption(..., AllowDuplicateNames) = (%v, %v), want (false, false)", v, ok)
	}
	if v, ok := json.GetOption(Indent(" "), jsontext.WithIndent); v != " " || !ok {
		t.Errorf(`GetOption(..., WithIndent) = (%q, %v), want (" ", true)`, v, ok)
	}
	if v, ok := json.GetOption(jsonflags.AllowDuplicateNames|1, jsontext.WithIndent); v != "" || ok {
		t.Errorf(`GetOption(..., WithIndent) = (%q, %v), want ("", false)`, v, ok)
	}
	if v, ok := json.GetOption(opts, jsontext.AllowDuplicateNames); v || ok {
		t.Errorf("GetOption(..., AllowDuplicateNames) = (%v, %v), want (false, false)", v, ok)
	}
	if v, ok := json.GetOption(opts, json.Deterministic); !v || !ok {
		t.Errorf("GetOption(..., Deterministic) = (%v, %v), want (true, true)", v, ok)
	}
	if v, ok := json.GetOption(opts, jsontext.Expand); v || !ok {
		t.Errorf("GetOption(..., Expand) = (%v, %v), want (false, true)", v, ok)
	}
	if v, ok := json.GetOption(opts, jsontext.AllowInvalidUTF8); v || ok {
		t.Errorf("GetOption(..., AllowInvalidUTF8) = (%v, %v), want (false, false)", v, ok)
	}
	if v, ok := json.GetOption(opts, jsontext.WithIndent); v != "\t" || !ok {
		t.Errorf(`GetOption(..., WithIndent) = (%q, %v), want ("\t", true)`, v, ok)
	}
	if v, ok := json.GetOption(opts, jsontext.WithIndentPrefix); v != "" || ok {
		t.Errorf(`GetOption(..., WithIndentPrefix) = (%q, %v), want ("", false)`, v, ok)
	}
	if v, ok := json.GetOption(opts, json.WithMarshalers); v == nil || !ok {
		t.Errorf(`GetOption(..., WithMarshalers) = (%v, %v), want (non-nil, true)`, v, ok)
	}
	if v, ok := json.GetOption(opts, json.WithUnmarshalers); v != nil || ok {
		t.Errorf(`GetOption(..., WithUnmarshalers) = (%v, %v), want (nil, false)`, v, ok)
	}
}

var sink struct {
	Bool       bool
	String     string
	Marshalers *json.Marshalers
}

func BenchmarkGetBool(b *testing.B) {
	b.ReportAllocs()
	opts := json.DefaultOptionsV2()
	for i := 0; i < b.N; i++ {
		sink.Bool, sink.Bool = json.GetOption(opts, jsontext.AllowDuplicateNames)
	}
}

func BenchmarkGetIndent(b *testing.B) {
	b.ReportAllocs()
	opts := json.DefaultOptionsV2()
	for i := 0; i < b.N; i++ {
		sink.String, sink.Bool = json.GetOption(opts, jsontext.WithIndent)
	}
}

func BenchmarkGetMarshalers(b *testing.B) {
	b.ReportAllocs()
	opts := json.JoinOptions(json.DefaultOptionsV2(), json.WithMarshalers(nil))
	for i := 0; i < b.N; i++ {
		sink.Marshalers, sink.Bool = json.GetOption(opts, json.WithMarshalers)
	}
}
