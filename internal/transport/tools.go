package transport

import (
	"os"
	"path/filepath"
	"time"

	"github.com/nskforward/gate4/internal/brokers"
	"github.com/nskforward/gate4/internal/users"
	"github.com/nskforward/gate4/pkg/pb"
)

func convertInUsers(in []*pb.User) []*users.User {
	result := make([]*users.User, 0, len(in))
	for _, user := range in {
		result = append(result, convertInUser(user))
	}
	return result
}

func convertOutUsers(in []*users.User) []*pb.User {
	result := make([]*pb.User, 0, len(in))
	for _, user := range in {
		result = append(result, convertOutUser(user))
	}
	return result
}

func convertOutUser(user *users.User) *pb.User {
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

func convertInUser(user *pb.User) *users.User {
	secret := ""
	if user.Secret != nil {
		secret = *user.Secret
	}
	return &users.User{
		ID:        *user.Id,
		BrokerID:  brokers.BrokerID(user.BrokerId),
		AccountID: user.AccountId,
		Expires:   time.Unix(user.Expires, 0),
		Created:   time.Unix(user.Created, 0),
		Blocked:   user.Blocked,
		Secret:    secret,
	}
}

func UnixSocketPath() string {
	return filepath.Join(os.TempDir(), "gate4.sock")
}
