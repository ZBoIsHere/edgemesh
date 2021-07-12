package gateway

func (gw *EdgeGateway) Run() {
	gw.APIConn.Init()
}
