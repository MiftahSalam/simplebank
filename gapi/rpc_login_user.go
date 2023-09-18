package gapi

import (
	"context"
	"database/sql"
	db "simplebank/db/sqlc"
	"simplebank/pb"
	"simplebank/util"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	user, err := server.store.GetUser(ctx, req.GetUsername())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "username not found: %s", err)
		}

		return nil, status.Errorf(codes.Internal, "failed to get user: %s", err)
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "incorrect password: %s", err)
	}

	accessToken, payloadAccessToken, err := server.tokenManager.CreateToken(user.Username, server.config.TokenExpiredDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get access token: %s", err)
	}

	refreshToken, payloadRefreshToken, err := server.tokenManager.CreateToken(user.Username, server.config.RefreshTokenExpiredDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get refresh token: %s", err)
	}

	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           payloadRefreshToken.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		ExpiredAt:    payloadRefreshToken.ExpiredAt,
		UserAgent:    server.extractMetadata(ctx).UserAgent,
		ClientIp:     server.extractMetadata(ctx).ClientIP,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create session: %s", err)
	}

	return &pb.LoginUserResponse{
		User:                  toPbUser(user),
		SessionId:             session.ID.String(),
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiredAt:  timestamppb.New(payloadAccessToken.ExpiredAt),
		RefreshTokenExpiredAt: timestamppb.New(payloadRefreshToken.ExpiredAt),
	}, nil
}
