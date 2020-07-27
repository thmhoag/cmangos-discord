package permissions

type Config struct {
	Cmds	map[string]*CmdPermission	`yaml: cmds`
}

type CmdPermission struct {
	AllowRoles 	[]string	`yaml:"allowRoles"`
	DenyRoles 	[]string	`yaml:"denyRoles"`
	AllowUsers	[]string	`yaml:"allowUsers"`
	DenyUsers	[]string	`yaml:"denyUsers"`
}