package dauth

import (
	"context"
	"google.golang.org/grpc/metadata"
)

const (
	SFHeaderUserID   string = "x-sf-user-id"
	SFHeaderApiKeyID string = "x-sf-api-key-id"
	SFHeaderIP       string = "x-real-ip"
)

func FromContext(ctx context.Context) (userID string, apiKeyId string, ip string) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return
	}

	userID = getHeader(SFHeaderUserID, md)
	apiKeyId = getHeader(SFHeaderApiKeyID, md)
	ip = getHeader(SFHeaderIP, md)
	return
}

func getHeader(key string, md metadata.MD) string {
	if len(md.Get(key)) > 0 {
		return md.Get(key)[0]
	}
	return ""
}
