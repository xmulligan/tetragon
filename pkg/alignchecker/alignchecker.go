// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Tetragon
package alignchecker

import (
	"reflect"

	"github.com/isovalent/tetragon-oss/pkg/api/processapi"
	"github.com/isovalent/tetragon-oss/pkg/api/testapi"
	"github.com/isovalent/tetragon-oss/pkg/sensors/exec/execvemap"

	check "github.com/cilium/cilium/pkg/alignchecker"
)

// CheckStructAlignments checks whether size and offsets of the C and Go
// structs for the datapath match.
//
// C struct size info is extracted from the given ELF object file debug section
// encoded in DWARF.
//
// To find a matching C struct field, a Go field has to be tagged with
// `align:"field_name_in_c_struct". In the case of unnamed union field, such
// union fields can be referred with special tags - `align:"$union0"`,
// `align:"$union1"`, etc.
func CheckStructAlignments(path string) error {
	// Validate alignments of C and Go equivalent structs
	toCheck := map[string][]reflect.Type{
		// from perf_event_output
		"msg_exit":         {reflect.TypeOf(processapi.MsgExitEvent{})},
		"msg_test":         {reflect.TypeOf(testapi.MsgTestEvent{})},
		"execve_map_value": {reflect.TypeOf(execvemap.ExecveValue{})},
	}
	return check.CheckStructAlignments(path, toCheck)
}
