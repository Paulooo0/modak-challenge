package notification

import "github.com/google/uuid"

type SendNotificationRequest struct {
	UserID  uuid.UUID `json:"user_id" binding:"required"`
	Type    string    `json:"type" binding:"required"`
	Message string    `json:"message" binding:"required"`
}

type StatusResponse struct {
	Status string `json:"status"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
