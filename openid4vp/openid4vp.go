package openid4vp

import (
	"encoding/base64"
	"fmt"

	doc "github.com/kokukuma/mdoc-verifier/document"
)

var (
	b64 = base64.URLEncoding.WithPadding(base64.StdPadding)
)

// https://openid.net/specs/openid-4-verifiable-presentations-1_0.html

type AuthorizationRequest struct {
	ClientID               string                 `json:"client_id"`
	ClientIDScheme         string                 `json:"client_id_scheme"`
	ResponseType           string                 `json:"response_type"`
	Nonce                  string                 `json:"nonce"`
	PresentationDefinition PresentationDefinition `json:"presentation_definition"`
	ResponseURI            string                 `json:"response_uri"`
	ResponseMode           string                 `json:"response_mode"`
	Scope                  string                 `json:"scope"`
	State                  string                 `json:"state"`
	ClientMetadata         ClientMetadata         `json:"client_metadata"`
}

type ClientMetadata struct {
	AuthorizationEncryptedResopnseAlg string   `json:"authorization_encrypted_response_alg"`
	AuthorizationEncryptedResopnseEnc string   `json:"authorization_encrypted_response_enc"`
	IDTokenEncryptedResponseAlg       string   `json:"id_token_encrypted_response_alg"`
	IDTokenEncryptedResponseEnc       string   `json:"id_token_encrypted_response_enc"`
	JwksURI                           string   `json:"jwks_uri"`
	SubjectSyntaxTypesSupported       []string `json:"subject_syntax_types_supported"`
	IDTokenSignedResponseAlg          string   `json:"id_token_signed_response_alg"`
}

type PresentationDefinition struct {
	ID               string            `json:"id"`
	InputDescriptors []InputDescriptor `json:"input_descriptors"`
}

type InputDescriptor struct {
	ID          string      `json:"id"`
	Format      Format      `json:"format"`
	Constraints Constraints `json:"constraints"`
}

type Constraints struct {
	LimitDisclosure string      `json:"limit_disclosure"`
	Fields          []PathField `json:"fields"`
}

type Format struct {
	MsoMdoc MsoMdoc `json:"mso_mdoc,omitempty"`
}

type MsoMdoc struct {
	Alg []string `json:"alg"`
}

type OpenID4VPData struct {
	VPToken                string                 `json:"vp_token"`
	State                  string                 `json:"state"`
	PresentationSubmission PresentationSubmission `json:"presentation_submission"`

	// ?
	APV string
	APU string
}
type PresentationSubmission struct {
	ID            string      `json:"id"`
	DefinitionID  string      `json:"definition_id"`
	DescriptorMap interface{} `json:"descriptor_map"`
}

type PathField struct {
	Path           []string `json:"path"`
	IntentToRetain bool     `json:"intent_to_retain"`
}

func FormatFields(ns doc.NameSpace, retain bool, ids ...doc.ElementIdentifier) []PathField {
	result := []PathField{}

	for _, id := range ids {
		result = append(result, PathField{
			Path:           []string{fmt.Sprintf("$['%s']['%s']", ns, id)},
			IntentToRetain: retain,
		})
	}
	return result
}

func CreatePresentationDefinition() PresentationDefinition {
	return PresentationDefinition{
		ID: "mDL-request-demo",
		InputDescriptors: []InputDescriptor{
			{
				ID: "eu.europa.ec.eudi.pid.1",
				Format: Format{
					MsoMdoc: MsoMdoc{
						Alg: []string{"ES256"},
					},
				},
				Constraints: Constraints{
					LimitDisclosure: "required",
					Fields: FormatFields(
						doc.EUDIPID1, true,
						doc.EudiFamilyName,
					),
				},
			},
			// {
			// 	ID: "org.iso.18013.5.1.mDL",
			// 	Format: Format{
			// 		MsoMdoc: MsoMdoc{
			// 			Alg: []string{"ES256"},
			// 		},
			// 	},
			// 	Constraints: Constraints{
			// 		LimitDisclosure: "required",
			// 		Fields: FormatFields(
			// 			doc.ISO1801351, true,
			// 			doc.FamilyName,
			// 			doc.GivenName,
			// 		),
			// 	},
			// },
		},
	}
}

func CreateClientMetadata() ClientMetadata {
	return ClientMetadata{
		AuthorizationEncryptedResopnseAlg: "ECDH-ES",
		AuthorizationEncryptedResopnseEnc: "A128CBC-HS256",
		IDTokenEncryptedResponseAlg:       "RSA-OAEP-256",
		IDTokenEncryptedResponseEnc:       "A128CBC-HS256",
		JwksURI:                           "https://fido-kokukuma.jp.ngrok.io/wallet/jwks.json",
		SubjectSyntaxTypesSupported:       []string{"urn:ietf:params:oauth:jwk-thumbprint"},
		IDTokenSignedResponseAlg:          "RS256",
	}
}
