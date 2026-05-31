package common

import (
	"time"

	"github.com/nskforward/gate4/internal/domain/model"
	"github.com/nskforward/gate4/pkg/pb"
)

func ConvertInUsers(users []model.User) []*pb.User {
	result := make([]*pb.User, 0, len(users))
	for _, user := range users {
		result = append(result, ConvertInUser(user))
	}
	return result
}

func ConvertInUser(user model.User) *pb.User {
	return &pb.User{
		Id:      user.ID,
		Created: user.Created.Unix(),
	}
}

func ConvertOutUsers(users []*pb.User) []model.User {
	result := make([]model.User, 0, len(users))
	for _, user := range users {
		result = append(result, ConvertOutUser(user))
	}
	return result
}

func ConvertOutUser(user *pb.User) model.User {
	return model.User{
		ID:      user.Id,
		Created: time.Unix(user.Created, 0),
	}
}
