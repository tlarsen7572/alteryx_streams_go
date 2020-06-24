package twitter

import (
	"fmt"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/tlarsen7572/goalteryx/api"
	"github.com/tlarsen7572/goalteryx/output_connection"
	"github.com/tlarsen7572/goalteryx/recordblob"
	"github.com/tlarsen7572/goalteryx/recordinfo"
	"strings"
)

type Ii struct {
	ToolId  int
	Config  *Config
	Output  output_connection.OutputConnection
	inInfo  recordinfo.RecordInfo
	stream  *twitter.Stream
	outInfo recordinfo.RecordInfo
}

func (ii *Ii) Init(recordInfoIn string) bool {
	var err error
	ii.inInfo, err = recordinfo.FromXml(recordInfoIn)
	if err != nil {
		api.OutputMessage(ii.ToolId, api.Error, err.Error())
		return false
	}
	_, err = ii.inInfo.GetFieldByName(`Event`)
	if err != nil {
		api.OutputMessage(ii.ToolId, api.Error, `missing required 'Event' field`)
	}

	generator := recordinfo.NewGenerator()
	generator.AddInt64Field(`ID`, `Twitter`)
	generator.AddV_WStringField(`Text`, `Twitter`, 300)
	generator.AddV_WStringField(`Created At`, `Twitter`, 100)
	generator.AddV_WStringField(`User`, `Twitter`, 100)
	ii.outInfo = generator.GenerateRecordInfo()
	_ = ii.Output.Init(ii.outInfo)
	return true
}

func (ii *Ii) PushRecord(record recordblob.RecordBlob) bool {
	event, isNull, err := ii.inInfo.GetStringValueFrom(`Event`, record)
	if err != nil || isNull {
		api.OutputMessage(ii.ToolId, api.Error, `error or null value received`)
		return false
	}

	if event == `Start` {
		config := oauth1.NewConfig(ii.Config.ConsumerKey, ii.Config.ConsumerSecret)
		token := oauth1.NewToken(ii.Config.AccessToken, ii.Config.AccessTokenSecret)
		httpClient := config.Client(oauth1.NoContext, token)
		client := twitter.NewClient(httpClient)
		demux := twitter.NewSwitchDemux()
		demux.Tweet = ii.receiveTweet
		demux.Warning = ii.tweetWarning
		demux.StreamDisconnect = ii.tweetDisconnect
		var follow []string
		var track []string
		if ii.Config.Follow != `` {
			follow = strings.Split(ii.Config.Follow, `,`)
		}
		if ii.Config.Track != `` {
			track = strings.Split(ii.Config.Track, `,`)
		}
		params := &twitter.StreamFilterParams{
			StallWarnings: twitter.Bool(true),
			Follow:        follow,
			Track:         track,
		}
		ii.stream, err = client.Streams.Filter(params)
		if err != nil {
			api.OutputMessage(ii.ToolId, api.Error, err.Error())
			return false
		}
		go demux.HandleChan(ii.stream.Messages)
	}
	if event == `End` {
		ii.stream.Stop()
	}
	return true
}

func (ii *Ii) UpdateProgress(percent float64) {
	ii.Output.UpdateProgress(percent)
}

func (ii *Ii) Close() {
	api.OutputMessage(ii.ToolId, api.Complete, ``)
	ii.Output.Close()
}

func (ii *Ii) CacheSize() int {
	return 0
}

func (ii *Ii) receiveTweet(tweet *twitter.Tweet) {
	if tweet.RetweetedStatus != nil || tweet.QuotedStatus != nil {
		return
	}

	_ = ii.outInfo.SetIntField(`ID`, int(tweet.ID))
	text := tweet.Text
	if tweet.ExtendedTweet != nil {
		text = tweet.ExtendedTweet.FullText
	}
	_ = ii.outInfo.SetStringField(`Text`, text)
	_ = ii.outInfo.SetStringField(`Created At`, tweet.CreatedAt)
	_ = ii.outInfo.SetStringField(`User`, tweet.User.ScreenName)

	record, err := ii.outInfo.GenerateRecord()
	if err != nil {
		api.OutputMessage(ii.ToolId, api.Error, err.Error())
		return
	}
	ii.Output.PushRecord(record)
}

func (ii *Ii) tweetWarning(warning *twitter.StallWarning) {
	api.OutputMessage(ii.ToolId, api.Warning, warning.Message)
}

func (ii *Ii) tweetDisconnect(disconnect *twitter.StreamDisconnect) {
	api.OutputMessage(ii.ToolId, api.Warning, fmt.Sprintf(`disconnected because: %v`, disconnect.Reason))
}
