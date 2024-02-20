package gapi

import (
	"context"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	db "pet-bank/db/sqlc"
	"pet-bank/pb"
	"pet-bank/val"
	"time"
)

func (server *Server) VerifyEmail(ctx context.Context, req *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error) {
	violations := validateVerifyEmailRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	verifyEmail, err := server.store.GetVerifyEmail(ctx, req.GetEmailId())

	if verifyEmail.SecretCode != req.GetSecretCode() {
		return nil, status.Errorf(codes.InvalidArgument, "invalid secret code: %s", err)
	}

	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to get verification: %s", err)
	}

	if verifyEmail.IsUsed {
		return nil, status.Errorf(codes.DeadlineExceeded, "code already verified: %s", err)
	}

	if time.Now().After(verifyEmail.ExpiredAt) {
		return nil, status.Errorf(codes.DeadlineExceeded, "code already expiried: %s", err)
	}

	result, err := server.store.VerifyEmailTx(ctx, db.VerifyEmailTxParams{
		EmailID: verifyEmail.ID,
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to verify email: %s", err)
	}

	return &pb.VerifyEmailResponse{IsVerified: result.User.IsEmailVerified}, nil
}

func validateVerifyEmailRequest(req *pb.VerifyEmailRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateEmailId(req.GetEmailId()); err != nil {
		violations = append(violations, fieldViolation("email_id", err))
	}

	if err := val.ValidateSecretCode(req.GetSecretCode()); err != nil {
		violations = append(violations, fieldViolation("secret_code", err))
	}

	return
}
