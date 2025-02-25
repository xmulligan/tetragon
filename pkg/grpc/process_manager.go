// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Tetragon

package grpc

import (
	"fmt"
	"sync"

	"github.com/cilium/hubble/pkg/cilium"
	"github.com/isovalent/tetragon-oss/api/v1/fgs"
	"github.com/isovalent/tetragon-oss/pkg/api/processapi"
	"github.com/isovalent/tetragon-oss/pkg/api/readyapi"
	"github.com/isovalent/tetragon-oss/pkg/api/testapi"
	"github.com/isovalent/tetragon-oss/pkg/api/tracingapi"
	"github.com/isovalent/tetragon-oss/pkg/dns"
	"github.com/isovalent/tetragon-oss/pkg/eventcache"
	"github.com/isovalent/tetragon-oss/pkg/execcache"
	"github.com/isovalent/tetragon-oss/pkg/grpc/exec"
	"github.com/isovalent/tetragon-oss/pkg/grpc/test"
	"github.com/isovalent/tetragon-oss/pkg/grpc/tracing"
	"github.com/isovalent/tetragon-oss/pkg/logger"
	"github.com/isovalent/tetragon-oss/pkg/metrics"
	"github.com/isovalent/tetragon-oss/pkg/reader/node"
	"github.com/isovalent/tetragon-oss/pkg/sensors"
	"github.com/isovalent/tetragon-oss/pkg/server"
	"github.com/sirupsen/logrus"
)

type execProcess interface {
	HandleExecveMessage(*processapi.MsgExecveEventUnix) *fgs.GetEventsResponse
	HandleExitMessage(*processapi.MsgExitEventUnix) *fgs.GetEventsResponse
}

var (
	tracingGrpc *tracing.Grpc
	execGrpc    execProcess
)

// ProcessManager maintains a cache of processes from fgs exec events.
type ProcessManager struct {
	eventCache *eventcache.Cache
	execCache  *execcache.Cache
	nodeName   string
	Server     *server.Server
	// synchronize access to the listeners map.
	mux               sync.Mutex
	listeners         map[server.Listener]struct{}
	ciliumState       *cilium.State
	enableProcessCred bool
	enableProcessNs   bool
	enableEventCache  bool
	enableCilium      bool
	dns               *dns.Cache
}

// NewProcessManager returns a pointer to an initialized ProcessManager struct.
func NewProcessManager(
	ciliumState *cilium.State,
	manager *sensors.Manager,
	enableProcessCred bool,
	enableProcessNs bool,
	enableEventCache bool,
	enableCilium bool,
) (*ProcessManager, error) {
	var err error

	pm := &ProcessManager{
		nodeName:          node.GetNodeNameForExport(),
		ciliumState:       ciliumState,
		listeners:         make(map[server.Listener]struct{}),
		enableProcessCred: enableProcessCred,
		enableProcessNs:   enableProcessNs,
		enableEventCache:  enableEventCache,
		enableCilium:      enableCilium,
	}

	pm.dns, err = dns.NewCache()
	if err != nil {
		return nil, fmt.Errorf("failed to create DNS cache %w", err)
	}
	pm.Server = server.NewServer(pm, manager)
	pm.eventCache = eventcache.New(pm.Server, pm.dns)
	pm.execCache = execcache.New(pm.Server, pm.dns)

	tracingGrpc = tracing.New(ciliumState, pm.dns, pm.eventCache, enableCilium, enableProcessCred, enableProcessNs)
	execGrpc = exec.New(pm.execCache, pm.eventCache, enableProcessCred, enableProcessNs)

	logger.GetLogger().WithField("enableCilium", enableCilium).WithFields(logrus.Fields{
		"enableEventCache":  enableEventCache,
		"enableProcessCred": enableProcessCred,
		"enableProcessNs":   enableProcessNs,
	}).Info("Starting process manager")
	return pm, nil
}

// Notify implements Listener.Notify.
func (pm *ProcessManager) Notify(event interface{}) error {
	var processedEvent *fgs.GetEventsResponse
	switch msg := event.(type) {
	case *readyapi.MsgFGSReady:
		// pass
	case *processapi.MsgExecveEventUnix:
		processedEvent = execGrpc.HandleExecveMessage(msg)
	case *processapi.MsgExitEventUnix:
		processedEvent = execGrpc.HandleExitMessage(msg)
	case *tracingapi.MsgGenericKprobeUnix:
		processedEvent = tracingGrpc.HandleGenericKprobeMessage(msg)
	case *tracingapi.MsgGenericTracepointUnix:
		processedEvent = tracingGrpc.HandleGenericTracepointMessage(msg)
	case *testapi.MsgTestEventUnix:
		processedEvent = test.HandleTestMessage(msg)

	default:
		logger.GetLogger().WithField("event", event).Warnf("unhandled event of type %T", msg)
		metrics.ErrorCount.WithLabelValues(string(metrics.UnhandledEvent)).Inc()
		return nil
	}
	if processedEvent != nil {
		pm.NotifyListener(event, processedEvent)
	}
	return nil
}

// Close implements Listener.Close.
func (pm *ProcessManager) Close() error {
	return nil
}

func (pm *ProcessManager) AddListener(listener server.Listener) {
	logger.GetLogger().WithField("getEventsListener", listener).Debug("Adding a getEventsListener")
	pm.mux.Lock()
	defer pm.mux.Unlock()
	pm.listeners[listener] = struct{}{}
}

func (pm *ProcessManager) RemoveListener(listener server.Listener) {
	logger.GetLogger().WithField("getEventsListener", listener).Debug("Removing a getEventsListener")
	pm.mux.Lock()
	defer pm.mux.Unlock()
	delete(pm.listeners, listener)
}

func (pm *ProcessManager) NotifyListener(original interface{}, processed *fgs.GetEventsResponse) {
	pm.mux.Lock()
	defer pm.mux.Unlock()
	for l := range pm.listeners {
		l.Notify(processed)
	}
	metrics.ProcessEvent(original, processed)
}
