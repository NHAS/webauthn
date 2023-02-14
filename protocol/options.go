package protocol

import (
	"github.com/go-webauthn/webauthn/protocol/webauthncose"
)

type CredentialCreation struct {
	Response PublicKeyCredentialCreationOptions `json:"publicKey"`
}

type CredentialAssertion struct {
	Response PublicKeyCredentialRequestOptions `json:"publicKey"`
}

// In order to create a Credential via create(), the caller specifies a few parameters in a PublicKeyCredentialCreationOptions object.
// See §5.4. Options for Credential Creation https://www.w3.org/TR/webauthn/#dictionary-makecredentialoptions
type PublicKeyCredentialCreationOptions struct {
	Challenge              URLEncodedBase64         `json:"challenge"`
	RelyingParty           RelyingPartyEntity       `json:"rp"`
	User                   UserEntity               `json:"user"`
	Parameters             []CredentialParameter    `json:"pubKeyCredParams,omitempty"`
	AuthenticatorSelection AuthenticatorSelection   `json:"authenticatorSelection,omitempty"`
	Timeout                int                      `json:"timeout,omitempty"`
	CredentialExcludeList  []CredentialDescriptor   `json:"excludeCredentials,omitempty"`
	Extensions             AuthenticationExtensions `json:"extensions,omitempty"`
	Attestation            ConveyancePreference     `json:"attestation,omitempty"`
}

// The PublicKeyCredentialRequestOptions dictionary supplies get() with the data it needs to generate an assertion.
// Its challenge member MUST be present, while its other members are OPTIONAL.
// See §5.5. Options for Assertion Generation https://www.w3.org/TR/webauthn/#assertion-options
type PublicKeyCredentialRequestOptions struct {
	Challenge          URLEncodedBase64            `json:"challenge"`
	Timeout            int                         `json:"timeout,omitempty"`
	RelyingPartyID     string                      `json:"rpId,omitempty"`
	AllowedCredentials []CredentialDescriptor      `json:"allowCredentials,omitempty"`
	UserVerification   UserVerificationRequirement `json:"userVerification,omitempty"` // Default is "preferred"
	Extensions         AuthenticationExtensions    `json:"extensions,omitempty"`
}

// This dictionary contains the attributes that are specified by a caller when referring to a public
// key credential as an input parameter to the create() or get() methods. It mirrors the fields of
// the PublicKeyCredential object returned by the latter methods.
// See §5.10.3. Credential Descriptor https://www.w3.org/TR/webauthn/#credential-dictionary
type CredentialDescriptor struct {
	// The valid credential types.
	Type CredentialType `json:"type"`
	// CredentialID The ID of a credential to allow/disallow
	CredentialID URLEncodedBase64 `json:"id"`
	// The authenticator transports that can be used
	Transport []AuthenticatorTransport `json:"transports,omitempty"`

	// The AttestationType from the Credential. Used internally only.
	AttestationType string `json:"-"`
}

// CredentialParameter is the credential type and algorithm
// that the relying party wants the authenticator to create
type CredentialParameter struct {
	Type      CredentialType                       `json:"type"`
	Algorithm webauthncose.COSEAlgorithmIdentifier `json:"alg"`
}

// This enumeration defines the valid credential types.
// It is an extension point; values can be added to it in the future, as
// more credential types are defined. The values of this enumeration are used
// for versioning the Authentication Assertion and attestation structures according
// to the type of the authenticator.
// See §5.10.3. Credential Descriptor https://www.w3.org/TR/webauthn/#credentialType
type CredentialType string

const (
	// PublicKeyCredentialType - Currently one credential type is defined, namely "public-key".
	PublicKeyCredentialType CredentialType = "public-key"
)

// AuthenticationExtensions - referred to as AuthenticationExtensionsClientInputs in the
// spec document, this member contains additional parameters requesting additional processing
// by the client and authenticator.
// This is currently under development
type AuthenticationExtensions map[string]interface{}

// WebAuthn Relying Parties may use the AuthenticatorSelectionCriteria dictionary to specify their requirements
// regarding authenticator attributes. See §5.4.4. Authenticator Selection Criteria
// https://www.w3.org/TR/webauthn/#authenticatorSelection
type AuthenticatorSelection struct {
	// AuthenticatorAttachment If this member is present, eligible authenticators are filtered to only
	// authenticators attached with the specified AuthenticatorAttachment enum
	AuthenticatorAttachment AuthenticatorAttachment `json:"authenticatorAttachment,omitempty"`
	// RequireResidentKey this member describes the Relying Party's requirements regarding resident
	// credentials. If the parameter is set to true, the authenticator MUST create a client-side-resident
	// public key credential source when creating a public key credential.
	RequireResidentKey *bool `json:"requireResidentKey,omitempty"`

	// ResidentKey this member describes the Relying Party's requirements regarding resident
	// credentials per Webauthn Level 2.
	ResidentKey ResidentKeyRequirement `json:"residentKey,omitempty"`

	// UserVerification This member describes the Relying Party's requirements regarding user verification for
	// the create() operation. Eligible authenticators are filtered to only those capable of satisfying this
	// requirement.
	UserVerification UserVerificationRequirement `json:"userVerification,omitempty"`
}

// ConveyancePreference is the type representing the Webauthn Attestation Conveyance Preference Enumeration.
//
// WebAuthn Relying Parties may use AttestationConveyancePreference to specify their preference regarding attestation
// conveyance during credential generation.
//
// See §5.4.7 https://www.w3.org/TR/webauthn-2/#enum-attestation-convey
type ConveyancePreference string

const (
	// PreferNoAttestation is a ConveyancePreference value.
	//
	// This value indicates that the Relying Party is not interested in authenticator attestation. For example, in order
	// to potentially avoid having to obtain user consent to relay identifying information to the Relying Party, or to
	// save a round trip to an Attestation CA or Anonymization CA.
	//
	// This is the default value.
	//
	// See https://www.w3.org/TR/webauthn-2/#dom-attestationconveyancepreference-none
	PreferNoAttestation ConveyancePreference = "none"

	// PreferIndirectAttestation is a ConveyancePreference value.
	//
	// This value indicates that the Relying Party prefers an attestation conveyance yielding verifiable attestation
	// statements, but allows the client to decide how to obtain such attestation statements. The client MAY replace the
	// authenticator-generated attestation statements with attestation statements generated by an Anonymization CA, in
	// order to protect the user’s privacy, or to assist Relying Parties with attestation verification in a
	// heterogeneous ecosystem.
	//
	// Note: There is no guarantee that the Relying Party will obtain a verifiable attestation statement in this case.
	// For example, in the case that the authenticator employs self attestation.
	//
	// See https://www.w3.org/TR/webauthn-2/#dom-attestationconveyancepreference-indirect
	PreferIndirectAttestation ConveyancePreference = "indirect"

	// PreferDirectAttestation is a ConveyancePreference value.
	//
	// This value indicates that the Relying Party wants to receive the attestation statement as generated by the
	// authenticator.
	//
	// See https://www.w3.org/TR/webauthn-2/#dom-attestationconveyancepreference-direct
	PreferDirectAttestation ConveyancePreference = "direct"

	// PreferEnterpriseAttestation is a ConveyancePreference value.
	//
	// This value indicates that the Relying Party wants to receive an attestation statement that may include uniquely
	// identifying information. This is intended for controlled deployments within an enterprise where the organization
	// wishes to tie registrations to specific authenticators. User agents MUST NOT provide such an attestation unless
	// the user agent or authenticator configuration permits it for the requested RP ID.
	//
	// If permitted, the user agent SHOULD signal to the authenticator (at invocation time) that enterprise
	// attestation is requested, and convey the resulting AAGUID and attestation statement, unaltered, to the Relying
	// Party.
	//
	// See https://www.w3.org/TR/webauthn-2/#dom-attestationconveyancepreference-enterprise
	PreferEnterpriseAttestation ConveyancePreference = "enterprise"
)

func (a *PublicKeyCredentialRequestOptions) GetAllowedCredentialIDs() [][]byte {
	var allowedCredentialIDs = make([][]byte, len(a.AllowedCredentials))
	for i, credential := range a.AllowedCredentials {
		allowedCredentialIDs[i] = credential.CredentialID
	}
	return allowedCredentialIDs
}

type Extensions interface{}

type ServerResponse struct {
	Status  ServerResponseStatus `json:"status"`
	Message string               `json:"errorMessage"`
}

type ServerResponseStatus string

const (
	StatusOk     ServerResponseStatus = "ok"
	StatusFailed ServerResponseStatus = "failed"
)
