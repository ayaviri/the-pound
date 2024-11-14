package db

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"
)

type NotificationType string

const (
	TypeTreat  NotificationType = "treat"
	TypeRebark NotificationType = "rebark"
	TypePaw    NotificationType = "paw"
	TypeFollow NotificationType = "follow"
)

type BaseNotificationPayload struct {
	FromDogId       string `json:"from_dog_id"`
	FromDogUsername string `json:"from_dog_username"`
}

type BarkRelatedNotificationPayload struct {
	BaseNotificationPayload
	BarkId string `json:"bark_id"`
	Bark   string `json:"bark"`
}

type TreatNotificationPayload struct {
	BarkRelatedNotificationPayload
}

type RebarkNotificationPayload struct {
	BarkRelatedNotificationPayload
}

type PawNotificationPayload struct {
	// The bark information contained in this notification is that of
	// the PAW, not the parent bark
	BarkRelatedNotificationPayload
}

type FollowNotificationPayload struct {
	BaseNotificationPayload
	IsApproved bool `json:"is_approved"`
}

type Notification struct {
	Id           string           `json:"id"`
	IsRead       bool             `json:"is_read"`
	ToDogId      string           `json:"to_dog_id"`
	Type         NotificationType `json:"type"`
	Payload      interface{}      `json:"payload"`
	CreationDate time.Time        `json:"created_at"`
}

func GetUnreadNotifications(
	e DBExecutor,
	dogId string,
	count uint,
	offset uint,
) ([]Notification, error) {
	query := `select id, is_read, to_dog_id, type, payload, created_at from notification
where to_dog_id = $1 and is_read = false order by created_at desc limit $2 offset $3`
	var rows *sql.Rows
	rows, err = e.Query(query, dogId, count, offset)

	if err != nil {
		return []Notification{}, nil
	}

	return ConstructNotificationsFromRows(rows)
}

func SetNotificationToRead(e DBExecutor, notificationId string) error {
	statement := "update notification set is_read = true where id = $1"
	_, err = e.Exec(statement, notificationId)

	return err
}

func ConstructNotificationsFromRows(r *sql.Rows) ([]Notification, error) {
	notifications := make([]Notification, 0)

	for r.Next() {
		n, err := constructNotificationFromRow(r)

		if err != nil {
			return []Notification{}, err
		}

		notifications = append(notifications, n)
	}

	return notifications, nil
}

// Scans the next row from the given set and constructs a notification from it
func constructNotificationFromRow(r *sql.Rows) (Notification, error) {
	var payloadStr string
	n := Notification{}
	err = r.Scan(
		&n.Id,
		&n.IsRead,
		&n.ToDogId,
		&n.Type,
		&payloadStr,
		&n.CreationDate,
	)

	if err != nil {
		return n, err
	}

	payloadStruct, err := notificationFromType(n.Type)

	if err != nil {
		return n, err
	}

	err = json.Unmarshal([]byte(payloadStr), &payloadStruct)

	if err != nil {
		return n, err
	}

	n.Payload = payloadStruct

	return n, nil
}

func notificationFromType(nType NotificationType) (interface{}, error) {
	switch nType {
	case TypeTreat:
		return TreatNotificationPayload{}, nil
	case TypeRebark:
		return RebarkNotificationPayload{}, nil
	case TypePaw:
		return PawNotificationPayload{}, nil
	case TypeFollow:
		return FollowNotificationPayload{}, nil
	default:
		return nil, errors.New("Unrecognised notification type")
	}

}
