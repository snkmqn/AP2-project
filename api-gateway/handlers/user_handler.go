package handlers

import (
	"context"
	"fmt"
	"google.golang.org/grpc/metadata"
	"log"
	"net/http"

	userpb "proto/generated/ecommerce/user"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	Client userpb.UserServiceClient
}

func NewUserHandler(client userpb.UserServiceClient) *UserHandler {
	return &UserHandler{Client: client}
}

func (h *UserHandler) RegisterUser(c *gin.Context) {
	var req userpb.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	resp, err := h.Client.RegisterUser(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *UserHandler) LoginUserHandler(c *gin.Context) {
	var req userpb.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	grpcReq := &userpb.LoginRequest{
		Email: req.GetEmail(),
		Password: req.GetPassword(),
	}

	resp, err := h.Client.LoginUser(context.Background(), grpcReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to login user: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": resp.GetToken(), "username": resp.GetUsername(), "email": resp.GetEmail()})
}

func (h *UserHandler) RetrieveProfileHandler(c *gin.Context) {

	userId := c.Param("user_id")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}
	log.Printf("Attempting to retrieve user profile with user_id: %s", userId)


	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
		return
	}

	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	md := metadata.Pairs("authorization", "Bearer "+token)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	resp, err := h.Client.RetrieveProfile(ctx, &userpb.RetrieveProfileRequest{UserId: userId})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to retrieve user profile: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"username":  resp.GetUsername(),
		"email": resp.GetEmail(),
	})
}


