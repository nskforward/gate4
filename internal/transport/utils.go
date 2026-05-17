package transport

import (
	"github.com/nskforward/gate4/pkg/pb"
	"github.com/nskforward/gate4/pkg/types"
)

func convertPositions(info *types.AccountInfo) *pb.GetPositionsResponse {
	positions := make([]*pb.Position, 0, len(info.Positions))
	for _, pos := range info.Positions {
		positions = append(positions, convertPosition(pos))
	}
	return &pb.GetPositionsResponse{
		Positions: positions,
	}
}

func convertPosition(in types.Position) *pb.Position {
	return &pb.Position{
		Symbol: in.Symbol,
		Price:  in.Price,
		Size:   in.Size,
		Profit: in.Profit,
	}
}

func convertQuote(in types.Quote) *pb.QuoteStreamResponse {
	return &pb.QuoteStreamResponse{
		Symbol:    in.Symbol,
		Timestamp: in.Timestamp,
		Ask:       in.Ask,
		Bid:       in.Bid,
	}
}

func convertSessions(in []types.Session) *pb.GetScheduleResponse {
	sessions := make([]*pb.ScheduleSession, 0, len(in))
	for _, sess := range in {
		sessions = append(sessions, &pb.ScheduleSession{
			Type:  string(sess.Type),
			Start: sess.Start,
			End:   sess.End,
		})
	}
	return &pb.GetScheduleResponse{
		Sessions: sessions,
	}
}

func convertAccountTrade(in types.AccountTrade) *pb.AccountTradeResponse {
	return &pb.AccountTradeResponse{
		Symbol:    in.Symbol,
		Timestamp: in.Timestamp,
		AccountId: in.AccountID,
		Price:     in.Price,
		Size:      in.Size,
	}
}

func convertAsset(in types.AssetInfo) *pb.GetAssetResponse {
	return &pb.GetAssetResponse{
		Symbol:      in.Symbol,
		Description: in.Description,
		Decimals:    in.Decimals,
		MinStep:     in.MinStep,
		LotSize:     in.LotSize,
		Currency:    in.Currency,
	}
}
