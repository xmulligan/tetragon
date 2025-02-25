// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Tetragon

package tracing

import (
	"context"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"testing"

	"github.com/isovalent/tetragon-oss/api/v1/fgs"
	ec "github.com/isovalent/tetragon-oss/pkg/eventchecker"
	"github.com/isovalent/tetragon-oss/pkg/kernels"
	"github.com/isovalent/tetragon-oss/pkg/observer"
	"github.com/isovalent/tetragon-oss/pkg/testutils"
	"github.com/stretchr/testify/assert"

	_ "github.com/isovalent/tetragon-oss/pkg/sensors/exec"
)

func TestKprobeNSChanges(t *testing.T) {
	if !kernels.MinKernelVersion("5.3.0") {
		t.Skip("matchNamespaceChanges requires at least 5.3.0 version")
	}

	var doneWG, readyWG sync.WaitGroup
	defer doneWG.Wait()

	ctx, cancel := context.WithTimeout(context.Background(), cmdWaitTime)
	defer cancel()

	testBin := testContribPath("namespace-tester/test_ns")
	testCmd := exec.CommandContext(ctx, testBin)
	testPipes, err := testutils.NewCmdBufferedPipes(testCmd)
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		logOut(t, "stdout> ", testPipes.StdoutRd)
	}()

	go func() {
		logOut(t, "stderr> ", testPipes.StderrRd)
	}()

	// makeSpecFile creates a new spec file bsed on the template, and the provided arguments
	makeSpecFile := func(pid string) string {
		data := map[string]string{
			"MatchedPID":   pid,
			"NamespacePID": "false",
		}
		specName, err := testutils.GetSpecFromTemplate("nschanges.yaml.tmpl", data)
		if err != nil {
			t.Fatal(err)
		}
		return specName
	}

	pidStr := strconv.Itoa(os.Getpid())
	specFname := makeSpecFile(pidStr)
	t.Logf("pid is %s and spec file is %s", pidStr, specFname)

	obs, err := observer.GetDefaultObserverWithFile(t, specFname, fgsLib)
	if err != nil {
		t.Fatalf("GetDefaultObserverWithFile error: %s", err)
	}
	observer.LoopEvents(ctx, t, &doneWG, &readyWG, obs)
	readyWG.Wait()

	if err := testCmd.Start(); err != nil {
		t.Fatal(err)
	}

	if err := testCmd.Wait(); err != nil {
		t.Fatalf("command failed with %s. Context error: %s", err, ctx.Err())
	}

	writeArg0 = ec.GenericArgFileChecker(ec.StringMatchAlways(), ec.SuffixStringMatch("strange.txt"), ec.FullStringMatch(""))
	writeArg1 = ec.GenericArgBytesCheck([]byte("testdata\x00"))
	writeArg2 = ec.GenericArgSizeCheck(9)
	kpChecker := ec.NewKprobeChecker().
		WithFunctionName("__x64_sys_write").
		WithArgs([]ec.GenericArgChecker{writeArg0, writeArg1, writeArg2}).
		WithAction(fgs.KprobeAction_KPROBE_ACTION_POST)
	checker := ec.NewOrderedMultiResponseChecker(
		ec.NewKprobeEventChecker().
			HasKprobe(kpChecker).
			HasKprobe(kpChecker).
			End(),
	)
	err = observer.JsonTestCheck(t, &checker)
	assert.NoError(t, err)
}

func TestKprobeCapChanges(t *testing.T) {
	if !kernels.MinKernelVersion("5.3.0") {
		t.Skip("matchCapabilityChanges requires at least 5.3.0 version")
	}

	var doneWG, readyWG sync.WaitGroup
	defer doneWG.Wait()

	ctx, cancel := context.WithTimeout(context.Background(), cmdWaitTime)
	defer cancel()

	testBin := testContribPath("capabilities-tester/test_caps")
	testCmd := exec.CommandContext(ctx, testBin)
	testPipes, err := testutils.NewCmdBufferedPipes(testCmd)
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		logOut(t, "stdout> ", testPipes.StdoutRd)
	}()

	go func() {
		logOut(t, "stderr> ", testPipes.StderrRd)
	}()

	// makeSpecFile creates a new spec file bsed on the template, and the provided arguments
	makeSpecFile := func(pid string) string {
		data := map[string]string{
			"MatchedPID":   pid,
			"NamespacePID": "false",
		}
		specName, err := testutils.GetSpecFromTemplate("capchanges.yaml.tmpl", data)
		if err != nil {
			t.Fatal(err)
		}
		return specName
	}

	pidStr := strconv.Itoa(os.Getpid())
	specFname := makeSpecFile(pidStr)
	t.Logf("pid is %s and spec file is %s", pidStr, specFname)

	obs, err := observer.GetDefaultObserverWithFile(t, specFname, fgsLib)
	if err != nil {
		t.Fatalf("GetDefaultObserverWithFile error: %s", err)
	}
	observer.LoopEvents(ctx, t, &doneWG, &readyWG, obs)
	readyWG.Wait()

	if err := testCmd.Start(); err != nil {
		t.Fatal(err)
	}

	if err := testCmd.Wait(); err != nil {
		t.Fatalf("command failed with %s. Context error: %s", err, ctx.Err())
	}

	writeArg0 = ec.GenericArgFileChecker(ec.StringMatchAlways(), ec.SuffixStringMatch("strange.txt"), ec.FullStringMatch(""))
	writeArg1 = ec.GenericArgBytesCheck([]byte("testdata\x00"))
	writeArg2 = ec.GenericArgSizeCheck(9)
	kpChecker := ec.NewKprobeChecker().
		WithFunctionName("__x64_sys_write").
		WithArgs([]ec.GenericArgChecker{writeArg0, writeArg1, writeArg2}).
		WithAction(fgs.KprobeAction_KPROBE_ACTION_POST)
	checker := ec.NewOrderedMultiResponseChecker(
		ec.NewKprobeEventChecker().
			HasKprobe(kpChecker).
			HasKprobe(kpChecker).
			End(),
	)
	err = observer.JsonTestCheck(t, &checker)
	assert.NoError(t, err)
}
