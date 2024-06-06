package webauthn

type defaultUser struct {
	id          []byte
	credentials []Credential
}

var _ User = (*defaultUser)(nil)

func (user *defaultUser) WebAuthnID() []byte {
	return user.id
}

func (user *defaultUser) WebAuthnName() string {
	return "newUser"
}

func (user *defaultUser) WebAuthnDisplayName() string {
	return "New User"
}

func (user *defaultUser) WebAuthnIcon() string {
	return "https://pics.com/avatar.png"
}

func (user *defaultUser) WebAuthnCredential(ID []byte) *Credential {
	return &Credential{}
}

func (user *defaultUser) WebAuthnCredentials() []*Credential {
	return []*Credential{}
}
