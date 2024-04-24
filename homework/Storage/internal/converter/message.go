package converter

import (
	"hw3/storage/internal/domain"
	"hw3/storage/internal/transport/dto"
)

type MessageConverter struct {
}

func (c *MessageConverter) MapDtoToDomain(dto dto.MessageInfoDTO) domain.MessageInfo {
	return domain.MessageInfo{
		Nickname: dto.Nickname,
		Message:  dto.Message,
		Time:     dto.Time,
	}
}
