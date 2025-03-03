// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Hubble

package observe

import (
	"fmt"
	"time"

	"github.com/cilium/hubble/pkg/defaults"
	hubprinter "github.com/cilium/hubble/pkg/printer"
	hubtime "github.com/cilium/hubble/pkg/time"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	selectorOpts struct {
		all          bool
		last         uint64
		since, until string
		follow       bool
		first        uint64
	}

	formattingOpts struct {
		output string

		timeFormat string

		enableIPTranslation bool
		nodeName            bool
		numeric             bool
		color               string
	}

	otherOpts struct {
		ignoreStderr    bool
		printRawFilters bool
	}

	printer *hubprinter.Printer

	// selector flags
	selectorFlags = pflag.NewFlagSet("selectors", pflag.ContinueOnError)
	// generic formatting flags, available to `hubble observe`, including sub-commands.
	formattingFlags = pflag.NewFlagSet("Formatting", pflag.ContinueOnError)
	// other flags
	otherFlags = pflag.NewFlagSet("other", pflag.ContinueOnError)
)

func init() {
	selectorFlags.BoolVar(&selectorOpts.all, "all", false, "Get all flows stored in Hubble's buffer")
	selectorFlags.Uint64Var(&selectorOpts.last, "last", 0, fmt.Sprintf("Get last N flows stored in Hubble's buffer (default %d)", defaults.FlowPrintCount))
	selectorFlags.Uint64Var(&selectorOpts.first, "first", 0, "Get first N flows stored in Hubble's buffer")
	selectorFlags.BoolVarP(&selectorOpts.follow, "follow", "f", false, "Follow flows output")
	selectorFlags.StringVar(&selectorOpts.since,
		"since", "",
		fmt.Sprintf(`Filter flows since a specific date. The format is relative (e.g. 3s, 4m, 1h43,, ...) or one of:
  StampMilli:             %s
  YearMonthDay:           %s
  YearMonthDayHour:       %s
  YearMonthDayHourMinute: %s
  RFC3339:                %s
  RFC3339Milli:           %s
  RFC3339Micro:           %s
  RFC3339Nano:            %s
  RFC1123Z:               %s
 `,
			time.StampMilli,
			hubtime.YearMonthDay,
			hubtime.YearMonthDayHour,
			hubtime.YearMonthDayHourMinute,
			time.RFC3339,
			hubtime.RFC3339Milli,
			hubtime.RFC3339Micro,
			time.RFC3339Nano,
			time.RFC1123Z,
		),
	)
	selectorFlags.StringVar(&selectorOpts.until,
		"until", "",
		fmt.Sprintf(`Filter flows until a specific date. The format is relative (e.g. 3s, 4m, 1h43,, ...) or one of:
  StampMilli:             %s
  YearMonthDay:           %s
  YearMonthDayHour:       %s
  YearMonthDayHourMinute: %s
  RFC3339:                %s
  RFC3339Milli:           %s
  RFC3339Micro:           %s
  RFC3339Nano:            %s
  RFC1123Z:               %s
 `,
			time.StampMilli,
			hubtime.YearMonthDay,
			hubtime.YearMonthDayHour,
			hubtime.YearMonthDayHourMinute,
			time.RFC3339,
			hubtime.RFC3339Milli,
			hubtime.RFC3339Micro,
			time.RFC3339Nano,
			time.RFC1123Z,
		),
	)

	formattingFlags.StringVarP(
		&formattingOpts.output, "output", "o", "compact",
		`Specify the output format, one of:
  compact:  Compact output
  dict:     Each flow is shown as KEY:VALUE pair
  jsonpb:   JSON encoded GetFlowResponse according to proto3's JSON mapping
  json:     Alias for jsonpb
  table:    Tab-aligned columns
`)
	formattingFlags.BoolVarP(&formattingOpts.nodeName, "print-node-name", "", false, "Print node name in output")
	formattingFlags.StringVar(
		&formattingOpts.timeFormat, "time-format", "StampMilli",
		fmt.Sprintf(`Specify the time format for printing. This option does not apply to the json and jsonpb output type. One of:
  StampMilli:             %s
  YearMonthDay:           %s
  YearMonthDayHour:       %s
  YearMonthDayHourMinute: %s
  RFC3339:                %s
  RFC3339Milli:           %s
  RFC3339Micro:           %s
  RFC3339Nano:            %s
  RFC1123Z:               %s
 `,
			time.StampMilli,
			hubtime.YearMonthDay,
			hubtime.YearMonthDayHour,
			hubtime.YearMonthDayHourMinute,
			time.RFC3339,
			hubtime.RFC3339Milli,
			hubtime.RFC3339Micro,
			time.RFC3339Nano,
			time.RFC1123Z,
		),
	)

	otherFlags.BoolVarP(&otherOpts.ignoreStderr,
		"silent-errors", "s", false,
		"Silently ignores errors and warnings")
	otherFlags.BoolVar(&otherOpts.printRawFilters,
		"print-raw-filters", false,
		"Print allowlist/denylist filters and exit without sending the request to Hubble server")
}

// New observer command.
func New(vp *viper.Viper) *cobra.Command {
	observeCmd := newObserveCmd(vp)
	flowsCmd := newFlowsCmd(vp)

	observeCmd.AddCommand(
		newAgentEventsCommand(vp),
		newDebugEventsCommand(vp),
		flowsCmd,
	)

	return observeCmd
}
