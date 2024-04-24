package converter

import (
	"time"

	"hw3/chat-service/internal/domain"
	"hw3/chat-service/internal/transport/dto"
)

type MessageConverter struct {
}

func (c *MessageConverter) MapRequestToDomain(dto dto.MessageInfoRequest) domain.MessageInfo {
	return domain.MessageInfo{
		Nickname: dto.Nickname,
		Message:  dto.Message,
	}
}

func (c *MessageConverter) MapDomainToResponse(domain domain.MessageInfo) dto.MessageInfoResponse {
	return dto.MessageInfoResponse{
		Nickname: domain.Nickname,
		Message:  domain.Message,
		Time:     domain.Time,
	}
}

func (c *MessageConverter) MapRequestToResponse(request dto.MessageInfoRequest) dto.MessageInfoResponse {
	return dto.MessageInfoResponse{
		Nickname: request.Nickname,
		Message:  request.Message,
		Time:     time.Now(),
	}
}

func (c *MessageConverter) MapSliceDomainToResponse(domains []domain.MessageInfo) []dto.MessageInfoResponse {
	dtos := make([]dto.MessageInfoResponse, len(domains))

	for i := 0; i < len(domains); i++ {
		dtos[i] = c.MapDomainToResponse(domains[i])
	}

	return dtos
}
