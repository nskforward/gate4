package broker

import (
	"github.com/nskforward/gate4/pkg/pb"
	"github.com/nskforward/gate4/pkg/types"
)

func toPbSessions(in []types.Session) []*pb.ScheduleSession {
	result := make([]*pb.ScheduleSession, 0, len(in))
	for _, sess := range in {
		result = append(result, toPbSession(&sess))
	}
	return result
}

func toPbSession(in *types.Session) *pb.ScheduleSession {
	if in == nil {
		return nil
	}
	return &pb.ScheduleSession{
		Type:  string(in.Type),
		Start: in.Start,
		End:   in.End,
	}
}

func toPbQuote(in types.Quote) *pb.QuoteStreamResponse {
	return &pb.QuoteStreamResponse{
		Symbol:    in.Symbol,
		Timestamp: in.Timestamp,
		Ask:       in.Ask,
		Bid:       in.Bid,
	}
}
