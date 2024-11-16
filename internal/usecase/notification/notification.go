package notificationusecase

import (
	"awesomeProject/internal/entity"
	"context"
)

type notification interface {
	CreateNotification(ctx context.Context, req *entity.CreateNotification) error
	GetNotification(ctx context.Context, req *entity.GetNotificationReq) (*entity.GetNotificationResp, error)
	AddNotification(ctx context.Context, req *entity.AddNotificationReq) error
}

type NotificationUseCase struct {
	notification notification
}

func NewNotificationUseCase(ntf notification) *NotificationUseCase {
	return &NotificationUseCase{ntf}
}

func (n *NotificationUseCase) CreateNotification(ctx context.Context, req *entity.CreateNotification) error {
	return n.notification.CreateNotification(ctx, req)
}
func (n *NotificationUseCase) GetNotification(ctx context.Context, req *entity.GetNotificationReq) (*entity.GetNotificationResp, error) {
	return n.notification.GetNotification(ctx, req)
}

func (n *NotificationUseCase) AddNotification(ctx context.Context, req *entity.AddNotificationReq) error {
	return n.notification.AddNotification(ctx, req)
}
