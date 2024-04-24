package queries

import (
	"context"
	"fmt"
	"log"

	"hw3/storage/internal/domain"
)

const insertMessageInfo = `INSERT INTO chat (id, nickname, message) 
	VALUES ($1, $2, $3) RETURNING nickname, message, message_time`

func (q *Queries) InsertMessage(ctx context.Context, messageInfo domain.MessageInfo) (domain.MessageInfo, error) {
	row := q.pool.QueryRow(ctx, insertMessageInfo, messageInfo.ID, messageInfo.Nickname, messageInfo.Message)
	var info domain.MessageInfo
	err := row.Scan(&info.Nickname, &info.Message, &info.Time)

	if err != nil {
		log.Println(err)
		return domain.MessageInfo{}, fmt.Errorf("error while inserting message")
	}

	return info, nil
}
