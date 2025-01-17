package authenticate

type AuthenticateRequest struct {
	FlowId    string `json:"flowId" binding:"required"`
	Signature string `json:"signature" binding:"required"`
}

type AuthenticatePayload struct {
	Token string `json:"token"`
}

type GetFlowIdPayload struct {
	Eula   string `json:"eula,omitempty"`
	FlowId string `json:"flowId"`
}

type web2AuthRequest struct {
	Token     string `json:"token"`
	Provider  string `json:"provider"`
	User_Type string `json:"type"`
}
