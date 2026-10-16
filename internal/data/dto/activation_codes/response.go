package activationcodes

type ActivateResponse struct {
	RedirectTo string `json:"redirect_to,omitempty"`
	Error      string `json:"error,omitempty"`
	Expired    bool   `json:"expired,omitempty"`
}
