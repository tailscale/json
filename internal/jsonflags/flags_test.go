// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package jsonflags

import "testing"

func TestFlags(t *testing.T) {
	type Check struct{ want Flags }
	type Join struct{ in Flags }
	type Set struct{ in Bools }
	type Clear struct{ in Bools }
	type Get struct {
		in     Bools
		want   bool
		wantOk bool
	}

	calls := []any{
		Get{in: AllowDuplicateNames, want: false, wantOk: false},
		Set{in: AllowDuplicateNames | 0},
		Get{in: AllowDuplicateNames, want: false, wantOk: true},
		Set{in: AllowDuplicateNames | 1},
		Get{in: AllowDuplicateNames, want: true, wantOk: true},
		Check{want: Flags{Presence: uint64(AllowDuplicateNames), Value: uint64(AllowDuplicateNames)}},
		Get{in: AllowInvalidUTF8, want: false, wantOk: false},
		Set{in: AllowInvalidUTF8 | 1},
		Get{in: AllowInvalidUTF8, want: true, wantOk: true},
		Set{in: AllowInvalidUTF8 | 0},
		Get{in: AllowInvalidUTF8, want: false, wantOk: true},
		Get{in: AllowDuplicateNames, want: true, wantOk: true},
		Check{want: Flags{Presence: uint64(AllowDuplicateNames | AllowInvalidUTF8), Value: uint64(AllowDuplicateNames)}},
		Set{in: AllowDuplicateNames | AllowInvalidUTF8 | 0},
		Check{want: Flags{Presence: uint64(AllowDuplicateNames | AllowInvalidUTF8), Value: uint64(0)}},
		Set{in: AllowDuplicateNames | AllowInvalidUTF8 | 0},
		Check{want: Flags{Presence: uint64(AllowDuplicateNames | AllowInvalidUTF8), Value: uint64(0)}},
		Set{in: AllowDuplicateNames | AllowInvalidUTF8 | 1},
		Check{want: Flags{Presence: uint64(AllowDuplicateNames | AllowInvalidUTF8), Value: uint64(AllowDuplicateNames | AllowInvalidUTF8)}},
		Join{in: Flags{Presence: 0, Value: 0}},
		Check{want: Flags{Presence: uint64(AllowDuplicateNames | AllowInvalidUTF8), Value: uint64(AllowDuplicateNames | AllowInvalidUTF8)}},
		Join{in: Flags{Presence: uint64(Expand | AllowInvalidUTF8), Value: uint64(AllowDuplicateNames)}},
		Check{want: Flags{Presence: uint64(Expand | AllowDuplicateNames | AllowInvalidUTF8), Value: uint64(AllowDuplicateNames)}},
		Clear{in: AllowDuplicateNames | AllowInvalidUTF8},
		Check{want: Flags{Presence: uint64(Expand), Value: uint64(0)}},
		Set{in: AllowInvalidUTF8 | Deterministic | IgnoreStructErrors | 1},
		Set{in: Expand | StringifyNumbers | RejectFloatOverflow | 0},
		Check{want: Flags{Presence: uint64(AllowInvalidUTF8 | Deterministic | IgnoreStructErrors | Expand | StringifyNumbers | RejectFloatOverflow), Value: uint64(AllowInvalidUTF8 | Deterministic | IgnoreStructErrors)}},
		Clear{in: ^AllCoderFlags},
		Check{want: Flags{Presence: uint64(AllowInvalidUTF8 | Expand), Value: uint64(AllowInvalidUTF8)}},
	}
	var fs Flags
	for i, call := range calls {
		switch call := call.(type) {
		case Join:
			fs.Join(call.in)
		case Set:
			fs.Set(call.in)
		case Clear:
			fs.Clear(call.in)
		case Get:
			if got, gotOk := fs.GetOk(call.in); got != call.want || gotOk != call.wantOk {
				t.Fatalf("%d: GetOk = (%v, %v), want (%v, %v)", i, got, gotOk, call.want, call.wantOk)
			}
		case Check:
			if fs != call.want {
				t.Fatalf("%d: got %x, want %x", i, fs, call.want)
			}
		}
	}
}
