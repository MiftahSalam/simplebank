package api

import (
	"errors"
	"fmt"
	"net/http"
	db "simplebank/db/sqlc"
	"time"

	"github.com/gin-gonic/gin"
)

type renewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type renewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiredAt time.Time `json:"access_token_expired_at"`
}

func (server *Server) renewAccessToken(ctx *gin.Context) {
	var req renewAccessTokenRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	refreshTokenPayload, err := server.tokenManager.VerifyToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	session, err := server.store.GetSession(ctx, refreshTokenPayload.ID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if session.IsBlocked {
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("session blocked")))
			return
		}
	}

	if session.Username != refreshTokenPayload.Username {
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("incorrect session user")))
			return
		}
	}

	if session.RefreshToken != req.RefreshToken {
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("mismatch session token")))
			return
		}
	}

	if time.Now().After(session.ExpiredAt) {
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("expired session")))
			return
		}
	}

	accessToken, payloadAccessToken, err := server.tokenManager.CreateToken(refreshTokenPayload.Username, server.config.TokenExpiredDuration)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, renewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiredAt: payloadAccessToken.ExpiredAt,
	})

}
