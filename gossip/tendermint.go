package gossip

func (svc *Service) initApp() {
	req := svc.config.Net.ChainInfo()
	_ = svc.abciApp.InitChain(req)
	// TODO: check response
}
