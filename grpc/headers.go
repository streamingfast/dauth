package grpc

import (
	"context"

	"github.com/streamingfast/dauth"
	"google.golang.org/grpc/metadata"
)

func TransferIdentityToOutgoingContext(ctx context.Context) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx
	}

	var (
		userId   string
		apiKeyId string
		ip       string
	)

	if len(md.Get(dauth.SFHeaderUserID)) > 0 {
		userId = md.Get(dauth.SFHeaderUserID)[0]
	}

	if len(md.Get(dauth.SFHeaderApiKeyID)) > 0 {
		apiKeyId = md.Get(dauth.SFHeaderApiKeyID)[0]
	}

	if len(md.Get(dauth.SFHeaderIP)) > 0 {
		ip = md.Get(dauth.SFHeaderIP)[0]
	}

	return metadata.AppendToOutgoingContext(
		ctx,
		dauth.SFHeaderUserID, userId,
		dauth.SFHeaderApiKeyID, apiKeyId,
		dauth.SFHeaderIP, ip,
	)
}
