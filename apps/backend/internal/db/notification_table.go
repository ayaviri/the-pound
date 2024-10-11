package db

//
// import (
// 	"database/sql"
// 	"encoding/json"
// 	"errors"
// 	"time"
// )
//
// /*
// What kinds of notifications do I want there to be ?
// For now, let's just do treats and follows/follow requests
//
// Treats have:
// - Dog that gave the treat
// - Bark that the treat was given to
//
// Follows/follow requests have:
// - Dog that followed/requested to follow
//
// ReBarks have:
// - Dog that rebarked
// - Bark that was rebarked
// */
// type Notification struct {
// }
//
// func GetNotifications(
// 	e DBExecutor,
// 	dogId string,
// 	count uint,
// 	offset uint,
// ) ([]Notification, error) {
// 	query := `select type, payload, created_at from notification
// where for_dog_id = $1 order by created_at desc limit $2 offset $3`
// 	var rows *sql.Rows
// 	rows, err = e.Query(query, dogId, count, offset)
//
// 	if err != nil {
// 		return []Notification{}, nil
// 	}
//
// 	var notifications []Notification
//
// 	for rows.Next() {
// 		var notificationType string
// 		var payload string
// 		var createdAt time.Time
// 		err = rows.Scan(&notificationType, &payload, &createdAt)
//
// 		if err != nil {
// 			return []Notification{}, err
// 		}
//
// 		switch notificationType {
// 		case "treat":
// 			err = json.Unmarshal([]byte(payload), &TreatNotification{})
//
// 			if err != nil {
// 				return []Notification{}, err
// 			}
// 		default:
// 			return []Notification{}, errors.New("Notification with invalid type found")
// 		}
// 	}
// 	return notifications, nil
// }
