package webauthn

import (
	"github.com/NHAS/webauthn/protocol"
)

// Credential contains all needed information about a WebAuthn credential for storage.
type Credential struct {
	// A probabilistically-unique byte sequence identifying a public key credential source and its authentication assertions.
	ID []byte

	// The public key portion of a Relying Party-specific credential key pair, generated by an authenticator and returned to
	// a Relying Party at registration time (see also public key credential). The private key portion of the credential key
	// pair is known as the credential private key. Note that in the case of self attestation, the credential key pair is also
	// used as the attestation key pair, see self attestation for details.
	PublicKey []byte

	// The attestation format used (if any) by the authenticator when creating the credential.
	AttestationType string

	// The transport types the authenticator supports.
	Transport []protocol.AuthenticatorTransport

	// The commonly stored flags.
	Flags CredentialFlags

	// The Authenticator information for a given certificate.
	Authenticator Authenticator
}

type CredentialFlags struct {
	// Flag UP indicates the users presence.
	UserPresent bool

	// Flag UV indicates the user performed verification.
	UserVerified bool

	// Flag BE indicates the credential is able to be backed up and/or sync'd between devices. This should NEVER change.
	BackupEligible bool

	// Flag BS indicates the credential has been backed up and/or sync'd. This value can change but it's recommended
	// that RP's keep track of this value.
	BackupState bool
}

// Descriptor converts a Credential into a protocol.CredentialDescriptor.
func (c Credential) Descriptor() (descriptor protocol.CredentialDescriptor) {
	return protocol.CredentialDescriptor{
		Type:            protocol.PublicKeyCredentialType,
		CredentialID:    c.ID,
		Transport:       c.Transport,
		AttestationType: c.AttestationType,
	}
}

// MakeNewCredential will return a credential pointer on successful validation of a registration response.
func MakeNewCredential(c *protocol.ParsedCredentialCreationData) (*Credential, error) {
	newCredential := &Credential{
		ID:              c.Response.AttestationObject.AuthData.AttData.CredentialID,
		PublicKey:       c.Response.AttestationObject.AuthData.AttData.CredentialPublicKey,
		AttestationType: c.Response.AttestationObject.Format,
		Transport:       c.Response.Transports,
		Flags: CredentialFlags{
			UserPresent:    c.Response.AttestationObject.AuthData.Flags.HasUserPresent(),
			UserVerified:   c.Response.AttestationObject.AuthData.Flags.HasUserVerified(),
			BackupEligible: c.Response.AttestationObject.AuthData.Flags.HasBackupEligible(),
			BackupState:    c.Response.AttestationObject.AuthData.Flags.HasBackupState(),
		},
		Authenticator: Authenticator{
			AAGUID:     c.Response.AttestationObject.AuthData.AttData.AAGUID,
			SignCount:  c.Response.AttestationObject.AuthData.Counter,
			Attachment: c.AuthenticatorAttachment,
		},
	}

	return newCredential, nil
}
