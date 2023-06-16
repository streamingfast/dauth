package dauth

import (
	"context"

	"google.golang.org/grpc/metadata"
)

const (
	SFHeaderUserID   string = "X-SF-USER-ID"
	SFHeaderApiKeyID string = "X-SF-API-KEY-ID"
	SFHeaderIP       string = "X-Real-IP"
)

type Identity struct {
	UserID   string
	ApiKeyID string
	IP       string
}

func TransferAuthHeadersToOutgoingContext(ctx context.Context) context.Context {
	identity := GetAuthInfoFromIncomingContext(ctx)
	return metadata.AppendToOutgoingContext(
		ctx,
		SFHeaderUserID, identity.UserID,
		SFHeaderApiKeyID, identity.ApiKeyID,
		SFHeaderIP, identity.IP,
	)
}

func GetAuthInfoFromIncomingContext(ctx context.Context) Identity {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return Identity{}
	}

	var userId, apiKeyId, ip string

	if len(md.Get(SFHeaderUserID)) > 0 {
		userId = md.Get(SFHeaderUserID)[0]
	}

	if len(md.Get(SFHeaderApiKeyID)) > 0 {
		apiKeyId = md.Get(SFHeaderApiKeyID)[0]
	}

	if len(md.Get(SFHeaderIP)) > 0 {
		ip = md.Get(SFHeaderIP)[0]
	}

	return Identity{
		UserID:   userId,
		ApiKeyID: apiKeyId,
		IP:       ip,
	}
}
