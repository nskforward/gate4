package common

import (
	"time"

	"github.com/nskforward/gate4/internal/domain/model"
	"github.com/nskforward/gate4/pkg/pb"
)

func ConvertOutTokens(tokens []*pb.Token) []model.Token {
	result := make([]model.Token, 0, len(tokens))
	for _, token := range tokens {
		result = append(result, ConvertOutToken(token))
	}
	return result
}

func ConvertOutToken(token *pb.Token) model.Token {
	return model.Token{
		ID:      token.Id,
		UserID:  token.UserId,
		Created: time.Unix(token.Created, 0),
		Expires: time.Unix(token.Expires, 0),
	}
}

func ConvertInTokens(tokens []model.Token) []*pb.Token {
	result := make([]*pb.Token, 0, len(tokens))
	for _, token := range tokens {
		result = append(result, ConvertInToken(token))
	}
	return result
}

func ConvertInToken(token model.Token) *pb.Token {
	return &pb.Token{
		Id:      token.ID,
		UserId:  token.UserID,
		Created: token.Created.Unix(),
		Expires: token.Expires.Unix(),
	}
}
