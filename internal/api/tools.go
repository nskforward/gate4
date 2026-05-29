package api

import (
	"time"

	"github.com/nskforward/gate4/internal/domain/users"
	"github.com/nskforward/gate4/pkg/pb"
)

func ConvertInUsers(in []*pb.User) []users.User {
	result := make([]users.User, 0, len(in))
	for _, user := range in {
		result = append(result, ConvertInUser(user))
	}
	return result
}

func ConvertOutUsers(in []users.User) []*pb.User {
	result := make([]*pb.User, 0, len(in))
	for _, user := range in {
		result = append(result, ConvertOutUser(user))
	}
	return result
}

func ConvertOutUser(user users.User) *pb.User {
	return &pb.User{
		BrokerId:  string(user.BrokerID),
		AccountId: user.AccountID,
		Expires:   user.Expires.Unix(),
		Created:   user.Created.Unix(),
		Id:        &user.ID,
		Blocked:   user.Blocked,
		Secret:    new("********"),
	}
}

func ConvertInUser(user *pb.User) users.User {
	secret := ""
	if user.Secret != nil {
		secret = *user.Secret
	}
	return users.User{
		ID:        *user.Id,
		BrokerID:  users.BrokerID(user.BrokerId),
		AccountID: user.AccountId,
		Expires:   time.Unix(user.Expires, 0),
		Created:   time.Unix(user.Created, 0),
		Blocked:   user.Blocked,
		Secret:    secret,
	}
}
