package queries

import (
	"context"
	"fmt"
	"log"

	"2sem/hw1/homework/internal/domain"
)

const getMessagesFromLastQuery = `SELECT id, nickname, message, message_time 
	FROM chat ORDER BY message_time DESC LIMIT $1`

func (q *Queries) GetLastMessages(ctx context.Context, count int) ([]*domain.MessageInfo, error) {
	rows, err := q.pool.Query(ctx, getMessagesFromLastQuery, count)
	if err != nil {
		return nil, fmt.Errorf("can't select message by count: %w", err)
	}
	defer rows.Close()

	messageInfos := make([]*domain.MessageInfo, 0)
	for rows.Next() {
		info := &domain.MessageInfo{}
		if err := rows.Scan(&info.ID, &info.Nickname, &info.Message, &info.Time); err != nil {
			return nil, fmt.Errorf("can't scan info: %w", err)
		}

		messageInfos = append(messageInfos, info)
	}

	return messageInfos, nil
}

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
