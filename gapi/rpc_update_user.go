package gapi

import (
	"context"
	db "simplebank/db/sqlc"
	"simplebank/pb"
	"simplebank/util"
	"simplebank/validation"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	authPayload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, unauthorizedError(err)
	}

	violations := validateUpdateUserRequest(req)
	if violations != nil {
		return nil, invalidArgument(violations)
	}

	if authPayload.Username != req.Username {
		return nil, status.Errorf(codes.PermissionDenied, "cannot update user info")
	}

	arg := db.UpdateUserParams{
		Username: req.GetUsername(),
		FullName: pgtype.Text{
			String: req.GetFullName(),
			Valid:  req.FullName != nil,
		},
		Email: pgtype.Text{
			String: req.GetEmail(),
			Valid:  req.Email != nil,
		},
	}

	if req.Password != nil {
		hashedPassword, err := util.HashPassword(req.GetPassword())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
		}

		arg.HashedPassword = pgtype.Text{
			String: hashedPassword,
			Valid:  req.Password != nil,
		}

		arg.PasswordChangedAt = pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		}
	}

	user, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		if db.ErrorCode(err) == db.UniqueKeyViolation {
			return nil, status.Errorf(codes.AlreadyExists, "%s", err)
		} else {
			return nil, status.Errorf(codes.Internal, "failed to Update user: %s", err)
		}
	}

	return &pb.UpdateUserResponse{
		User: toPbUser(user),
	}, nil

}

func validateUpdateUserRequest(req *pb.UpdateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validation.ValidateUsername(req.Username); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if req.Password != nil {
		if err := validation.ValidatePassword(req.GetPassword()); err != nil {
			violations = append(violations, fieldViolation("password", err))
		}
	}
	if req.Email != nil {
		if err := validation.ValidateEmail(req.GetEmail()); err != nil {
			violations = append(violations, fieldViolation("email", err))
		}
	}
	if req.FullName != nil {
		if err := validation.ValidateFullname(req.GetUsername()); err != nil {
			violations = append(violations, fieldViolation("full_name", err))
		}
	}

	return
}
