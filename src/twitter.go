package flex_eth

type TwitterAuth struct {
	bearer string
}

func (auth TwitterAuth) DebugText() string {
	return auth.bearer
}
