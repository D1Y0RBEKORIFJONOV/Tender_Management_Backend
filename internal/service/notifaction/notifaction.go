package notifactionservice

import (
	"awesomeProject/internal/entity"
	redisrepository "awesomeProject/internal/infastructure/repository/redis"
	notificationusecase "awesomeProject/internal/usecase/notification"
	"context"
	"log/slog"
	"time"
)

type Notification struct {
	logger       *slog.Logger
	notification *notificationusecase.NotificationRepository
	redis        *redisrepository.RedisUserRepository
}

func NewNotification(logger *slog.Logger, notification *notificationusecase.NotificationRepository,
	redis *redisrepository.RedisUserRepository) *Notification {
	return &Notification{logger: logger, notification: notification, redis: redis}
}

func (n *Notification) CreateNotification(ctx context.Context, req *entity.CreateNotification) error {
	const op = "Service.CreateNotification"
	log := n.logger.With(slog.String("Method", op))
	log.Info("Create Notification")
	defer log.Info("Create Notification Completed")
	message := entity.MessageBid{
		CreateAt:   time.Now(),
		SenderName: "Tender Management ",
		Status:     "you are successfully connected to notification service",
	}

	err := n.notification.SaveNotification(ctx, &entity.Notification{})
	err = n.notification.SaveNotification(ctx, &entity.Notification{
		UserId:   req.UserId,
		Offset:   1,
		Messages: []entity.MessageBid{message},
	})
	if err != nil {
		log.Error("err", err.Error())
		return err
	}
	return nil
}

func (n *Notification) GetNotification(ctx context.Context, req *entity.GetNotificationReq) (*entity.GetNotificationResp, error) {
	const op = "Service.GetNotification"
	log := n.logger.With(slog.String("Method", op))
	log.Info("Get Notification")
	defer log.Info("Get Notification Completed")
	if req.Offset == 0 {
		offset, err := n.notification.GetOffsetNotification(ctx, req.UserId)
		if err != nil {
			log.Error("err", err.Error())
			return nil, err
		}
		messages, err := n.notification.GetNotification(ctx, &entity.GetNotificationReq{
			Offset: offset,
			UserId: req.UserId,
		})
		if err != nil {
			log.Error("err", err.Error())
			return nil, err
		}
		offSetUpdate := offset
		var result entity.GetNotificationResp
		for i := offset - 1; int(i) < len(messages); i++ {
			result.Messages = append(result.Messages, messages[i])
			offSetUpdate++
		}
		if len(result.Messages) != 0 {
			err = n.notification.UpdateOffsetNotification(ctx, req.UserId, offSetUpdate)
			if err != nil {
				log.Error("err", err.Error())
				return nil, err
			}
		}

		return &result, nil
	}
	messages, err := n.notification.GetNotification(ctx, &entity.GetNotificationReq{
		Offset: req.Offset,
		UserId: req.UserId,
	})
	if err != nil {
		log.Error("err", err.Error())
		return nil, err
	}
	result := entity.GetNotificationResp{
		Messages: messages,
	}
	return &result, nil
}

func (n *Notification) AddNotification(ctx context.Context, req *entity.AddNotificationReq) error {
	const op = "Service.AddNotification"
	log := n.logger.With(slog.String("Method", op))
	log.Info("Add Notification")
	defer log.Info("Add Notification Completed")
	messages := entity.MessageBid{
		CreateAt:   time.Now(),
		SenderName: req.CreateMessage.SenderName,
		Status:     req.CreateMessage.Status,
	}
	err := n.notification.AddNotification(ctx, req.UserId, messages)
	if err != nil {
		log.Error("err", err.Error())
		return err
	}
	err = n.redis.Check(ctx)
	if err != nil {
		log.Error("err", err.Error())
		return err
	}
	return nil
}
