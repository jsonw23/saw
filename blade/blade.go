package blade

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/TylerBrock/colorjson"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/fatih/color"
	sawConfig "github.com/jsonw23/saw/config"
)

// A Blade is a Saw execution instance
type Blade struct {
	config *sawConfig.Configuration
	aws    *config.Config
	output *sawConfig.OutputConfiguration
	cwl    *cloudwatchlogs.Client
}

// NewBlade creates a new Blade with CloudWatchLogs instance from provided config
func NewBlade(
	sawConfig *sawConfig.Configuration,
	awsConfig *config.Config,
	outputConfig *sawConfig.OutputConfiguration,
) *Blade {
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	blade := Blade{}

	blade.cwl = cloudwatchlogs.NewFromConfig(cfg)
	blade.config = sawConfig
	blade.output = outputConfig

	return &blade
}

// GetLogGroups gets the log groups from AWS given the blade configuration
func (b *Blade) GetLogGroups() ([]types.LogGroup, error) {
	logGroups := make([]types.LogGroup, 0)
	input := b.config.DescribeLogGroupsInput()
	for {
		output, err := b.cwl.DescribeLogGroups(context.TODO(), input)
		if err != nil {
			return nil, err
		}
		logGroups = append(logGroups, output.LogGroups...)
		input.NextToken = output.NextToken
		if input.NextToken == nil {
			break
		}
	}
	return logGroups, nil
}

// GetLogStreams gets the log streams from AWS given the blade configuration
func (b *Blade) GetLogStreams() ([]types.LogStream, error) {
	streams := make([]types.LogStream, 0)
	input := b.config.DescribeLogStreamsInput()
	for {
		output, err := b.cwl.DescribeLogStreams(context.TODO(), input)
		if err != nil {
			return nil, err
		}
		streams = append(streams, output.LogStreams...)
		input.NextToken = output.NextToken
		if input.NextToken == nil {
			break
		}
	}

	return streams, nil
}

// GetEvents gets events from AWS given the blade configuration
func (b *Blade) GetEvents() {
	formatter := b.output.Formatter()
	input := b.config.FilterLogEventsInput()
	for {
		output, err := b.cwl.FilterLogEvents(context.TODO(), input)
		if err != nil {
			fmt.Println("Error", err)
			os.Exit(2)
		}
		for _, event := range output.Events {
			if b.output.Pretty {
				fmt.Println(formatEvent(formatter, &event))
			} else {
				fmt.Println(*event.Message)
			}
		}

		input.NextToken = output.NextToken
		if input.NextToken == nil {
			os.Exit(0)
		}
	}
}

// StreamEvents continuously prints log events to the console
func (b *Blade) StreamEvents() {
	var lastSeenTime *int64
	var seenEventIDs map[string]bool
	formatter := b.output.Formatter()
	input := b.config.FilterLogEventsInput()

	clearSeenEventIds := func() {
		seenEventIDs = make(map[string]bool, 0)
	}

	addSeenEventIDs := func(id *string) {
		seenEventIDs[*id] = true
	}

	updateLastSeenTime := func(ts *int64) {
		if lastSeenTime == nil || *ts > *lastSeenTime {
			lastSeenTime = ts
			clearSeenEventIds()
		}
	}

	for {
		output, err := b.cwl.FilterLogEvents(context.TODO(), input)
		if err != nil {
			fmt.Println("Error", err)
			os.Exit(2)
		}

		for _, event := range output.Events {
			updateLastSeenTime(event.Timestamp)
			if _, seen := seenEventIDs[*event.EventId]; !seen {
				var message string
				if b.output.Raw {
					message = *event.Message
				} else {
					message = formatEvent(formatter, &event)
				}
				message = strings.TrimRight(message, "\n")
				fmt.Println(message)
				addSeenEventIDs(event.EventId)
			}
		}

		input.NextToken = output.NextToken
		if input.NextToken == nil {
			if lastSeenTime != nil {
				input.StartTime = lastSeenTime
			}
			time.Sleep(1 * time.Second)
		}
	}
}

// formatEvent returns a CloudWatch log event as a formatted string using the provided formatter
func formatEvent(formatter *colorjson.Formatter, event *types.FilteredLogEvent) string {
	red := color.New(color.FgRed).SprintFunc()
	white := color.New(color.FgWhite).SprintFunc()

	str := event.Message
	bytes := []byte(*event.Message)
	date := time.Unix(*event.Timestamp, 0)
	dateStr := date.Format(time.RFC3339)
	streamStr := aws.String(*event.LogStreamName)
	jl := map[string]interface{}{}

	if err := json.Unmarshal(bytes, &jl); err != nil {
		return fmt.Sprintf("[%s] (%s) %s", red(dateStr), white(streamStr), *str)
	}

	output, _ := formatter.Marshal(jl)
	return fmt.Sprintf("[%s] (%s) %s", red(dateStr), white(streamStr), output)
}
