package bot

import "github.com/thmhoag/cmangos-discord/pkg/cmangos"

type globalCtx struct {
	mangosClient cmangos.SoapClient
}

func (c *globalCtx) MangosClient() cmangos.SoapClient {
	return c.mangosClient
}