package forms

const (
	FormButtonIdPrefix   = "form-"
	FormModalIdPrefix    = "modal-"
	ModAcceptBtnIdPrefix = "form-mod-accept-"
	ModRejectBtnIdPrefix = "form-mod-reject-"
	ModBanBtnIdPrefix    = "form-mod-ban-"
)

type Forms struct {
	// TODO: mod actions cache instance
}

type UserInput struct {
	InputId string
	Value   string
}
