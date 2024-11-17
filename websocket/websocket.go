package websocket

import (
	"awesomeProject/internal/config"
	"awesomeProject/internal/entity"
	notificationusecase "awesomeProject/internal/usecase/notification"
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type MessageBid struct {
	Message string
	Bid     *entity.Bid
}
type MessageTender struct {
	Message string
	Tender  *entity.Tender
}
type Server struct {
	clients      map[string]*websocket.Conn
	clientsLock  sync.Mutex
	redisClient  *redis.Client
	notification *notificationusecase.NotificationUseCase
}

func NewServer(notification *notificationusecase.NotificationUseCase) *Server {
	cfg := config.New()
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisURL,
		Password: "",
		DB:       0,
	})
	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Conn errr: %v", err)
	}
	server := &Server{
		notification: notification,
		clients:      make(map[string]*websocket.Conn),
		redisClient:  redisClient,
	}
	go server.listenForNotifications()
	return server
}

func (s *Server) HandlerNotification(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	id := r.Header.Get("id")
	if id == "" {
		log.Println("id  not found ")
		return
	}

	s.clientsLock.Lock()
	s.clients[id] = conn
	s.clientsLock.Unlock()
	res, err := s.notification.GetNotification(context.Background(), &entity.GetNotificationReq{
		UserId: id,
	})

	if err != nil {
		log.Println(err)
		return
	}
	for i := 0; i < len(res.Messages); i++ {
		err = s.writeMessage(res.Messages[i], conn)
		if err != nil {
			log.Println(err)
		}
	}
	defer func() {
		s.clientsLock.Lock()
		delete(s.clients, id)
		s.clientsLock.Unlock()
	}()
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println("conn err:", err)
			break
		}
	}
}

func (s *Server) listenForNotifications() {
	pubsub := s.redisClient.Subscribe(context.Background(), "notifications")
	defer pubsub.Close()
	for {
		_, err := pubsub.ReceiveMessage(context.Background())
		if err != nil {
			log.Println("conn err:", err)
			continue
		}
		s.clientsLock.Lock()
		for id, conn := range s.clients {
			res, err := s.notification.GetNotification(context.Background(), &entity.GetNotificationReq{
				UserId: id,
			})
			if err != nil {
				log.Println("conn err", id, ":", err)
				continue
			}

			if len(res.Messages) == 0 {
				continue
			}

			for _, msg := range res.Messages {
				err = s.writeMessage(msg, conn)
				if err != nil {
					log.Println(err)
				}
			}
		}
		s.clientsLock.Unlock()
	}
}

func (s *Server) writeMessage(messages entity.MessageBid, conn *websocket.Conn) error {
	if messages.SenderName == "bid" {
		var message MessageBid
		message.Message = "You have been requested"
		err := json.Unmarshal([]byte(messages.Status), &message.Bid)
		if err != nil {
			log.Println(err)
		}
		err = conn.WriteJSON(message)
	} else if messages.SenderName == "tender" {
		var message MessageTender
		message.Message = "Congratulations, you have won this tender!"
		err := json.Unmarshal([]byte(messages.Status), &message.Tender)
		if err != nil {
			log.Println(err)
		}
		err = conn.WriteJSON(message)
		if err != nil {
			log.Println(err)
		}
	} else {
		err := conn.WriteMessage(websocket.TextMessage, []byte(messages.Status))
		if err != nil {
			log.Println(err)
		}
	}
	return nil
}
