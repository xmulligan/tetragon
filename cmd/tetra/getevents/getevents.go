// Copyright (C) Isovalent, Inc. - All Rights Reserved.
//
// NOTICE: All information contained herein is, and remains the property of
// Isovalent Inc and its suppliers, if any. The intellectual and technical
// concepts contained herein are proprietary to Isovalent Inc and its suppliers
// and may be covered by U.S. and Foreign Patents, patents in process, and are
// protected by trade secret or copyright law.  Dissemination of this information
// or reproduction of this material is strictly forbidden unless prior written
// permission is obtained from Isovalent Inc.

package getevents

import (
	"context"
	"encoding/json"
	"errors"
	"os"

	"github.com/isovalent/tetragon-oss/api/v1/fgs"
	"github.com/isovalent/tetragon-oss/cmd/tetra/common"
	"github.com/isovalent/tetragon-oss/pkg/encoder"
	"github.com/isovalent/tetragon-oss/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func getEvents(ctx context.Context, client fgs.FineGuidanceSensorsClient) {
	stream, err := client.GetEvents(ctx, &fgs.GetEventsRequest{})
	if err != nil {
		logger.GetLogger().WithError(err).Fatal("Failed to call GetEvents")
	}
	var eventEncoder encoder.EventEncoder
	if viper.GetString(common.KeyOutput) == "compact" {
		colorMode := encoder.ColorMode(viper.GetString(common.KeyColor))
		eventEncoder = encoder.NewCompactEncoder(os.Stdout, colorMode)
	} else {
		eventEncoder = json.NewEncoder(os.Stdout)
	}
	for {
		res, err := stream.Recv()
		if err != nil {
			if !errors.Is(err, context.Canceled) && status.Code(err) != codes.Canceled {
				logger.GetLogger().WithError(err).Fatal("Failed to receive events")
			}
			return
		}
		if err = eventEncoder.Encode(res); err != nil {
			logger.GetLogger().WithError(err).WithField("event", res).Debug("Failed to encode event")
		}
	}
}

func New() *cobra.Command {
	cmd := cobra.Command{
		Use:   "getevents",
		Short: "Print events",
		Run: func(cmd *cobra.Command, args []string) {
			common.CliRun(getEvents)
		},
	}

	flags := cmd.Flags()
	flags.StringP("output", "o", "json", "Output format. json or compact")
	flags.String("color", "auto", "Colorize compact output. auto, always, or never")
	viper.BindPFlags(flags)
	return &cmd
}
