package openid4vp

import (
	"encoding/base64"
	"fmt"

	"github.com/fxamacker/cbor/v2"
	"github.com/kokukuma/identity-credential-api-demo/protocol"
)

// TODO: session transcript: 9.1.5.1 Session transcript

const BROWSER_HANDOVER_V1 = "BrowserHandoverv1"

type OriginInfo struct {
	Cat     int     `json:"cat"`
	Type    int     `json:"type"`
	Details Details `json:"details"`
}

type Details struct {
	BaseURL string `json:"baseUrl"`
}

func generateBrowserSessionTranscript(nonce []byte, origin string, requesterIdHash []byte) ([]byte, error) {
	originInfo := OriginInfo{
		Cat:  1,
		Type: 1,
		Details: Details{
			BaseURL: origin,
		},
	}
	originInfoBytes, err := cbor.Marshal(originInfo)
	if err != nil {
		return nil, fmt.Errorf("error encoding origin info: %v", err)
	}

	// Create the final CBOR array
	browserHandover := []interface{}{
		nil, // DeviceEngagementBytes
		nil, // EReaderKeyBytes
		[]interface{}{ // BrowserHandover
			BROWSER_HANDOVER_V1,
			nonce,
			originInfoBytes,
			requesterIdHash,
		},
	}

	transcript, err := cbor.Marshal(browserHandover)
	if err != nil {
		return nil, fmt.Errorf("error encoding transcript: %v", err)
	}

	return transcript, nil
}

func generateOID4VPSessionTranscript(nonce []byte, clientID, responseURI, apu string) ([]byte, error) {
	mdocGeneratedNonce, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(apu)
	if err != nil {
		return nil, fmt.Errorf("failed to decode mdocGeneratedNonce")
	}

	clientIdToHash, err := cbor.Marshal([]interface{}{clientID, mdocGeneratedNonce})
	if err != nil {
		return nil, err
	}

	// Create responseUriToHash
	responseUriToHash, err := cbor.Marshal([]interface{}{responseURI, mdocGeneratedNonce})
	if err != nil {
		return nil, err
	}
	clientIdHash := protocol.Digest(clientIdToHash, "SHA-256")
	responseURIHash := protocol.Digest(responseUriToHash, "SHA-256")

	// Create the final CBOR array
	oid4vpHandover := []interface{}{
		nil, // DeviceEngagementBytes
		nil, // EReaderKeyBytes
		[]interface{}{ // OID4VPHandover
			clientIdHash,
			responseURIHash,
			nonce,
		},
	}

	transcript, err := cbor.Marshal(oid4vpHandover)
	if err != nil {
		return nil, fmt.Errorf("error encoding transcript: %v", err)
	}

	return transcript, nil
}
