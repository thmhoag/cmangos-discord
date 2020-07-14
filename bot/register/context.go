package register

import "github.com/thmhoag/cmangos-discord/pkg/cmangos"

type Ctx interface {
	MangosClient() *cmangos.SoapClient
}