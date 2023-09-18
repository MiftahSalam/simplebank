package gapi

import (
	db "simplebank/db/sqlc"
	"simplebank/pb"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func toPbUser(user db.User) *pb.User {
	return &pb.User{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
		CreatedAt:         timestamppb.New(user.CreatedAt),
	}
}