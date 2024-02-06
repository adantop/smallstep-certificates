package wire

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseUserID(t *testing.T) {
	ok := `{"name": "Alice Smith", "domain": "wire.com", "handle": "wireapp://%40alice_wire@wire.com"}`
	failJSON := `{"name": }`
	emptyHandle := `{"name": "Alice Smith", "domain": "wire.com", "handle": ""}`
	emptyName := `{"name": "", "domain": "wire.com", "handle": "wireapp://%40alice_wire@wire.com"}`
	emptyDomain := `{"name": "Alice Smith", "domain": "", "handle": "wireapp://%40alice_wire@wire.com"}`
	tests := []struct {
		name       string
		data       []byte
		wantWireID UserID
		wantErr    bool
	}{
		{name: "ok", data: []byte(ok), wantWireID: UserID{Name: "Alice Smith", Domain: "wire.com", Handle: "wireapp://%40alice_wire@wire.com"}},
		{name: "fail/json", data: []byte(failJSON), wantErr: true},
		{name: "fail/empty-handle", data: []byte(emptyHandle), wantErr: true},
		{name: "fail/empty-name", data: []byte(emptyName), wantErr: true},
		{name: "fail/empty-domain", data: []byte(emptyDomain), wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotWireID, err := ParseUserID(tt.data)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.wantWireID, gotWireID)
		})
	}
}

func TestParseDeviceID(t *testing.T) {
	ok := `{"name": "device", "domain": "wire.com", "client-id": "wireapp://CzbfFjDOQrenCbDxVmgnFw!594930e9d50bb175@wire.com", "handle": "wireapp://%40alice_wire@wire.com"}`
	failJSON := `{"name": }`
	emptyHandle := `{"name": "device", "domain": "wire.com", "client-id": "wireapp://CzbfFjDOQrenCbDxVmgnFw!594930e9d50bb175@wire.com", "handle": ""}`
	emptyName := `{"name": "", "domain": "wire.com", "client-id": "wireapp://CzbfFjDOQrenCbDxVmgnFw!594930e9d50bb175@wire.com", "handle": "wireapp://%40alice_wire@wire.com"}`
	emptyDomain := `{"name": "device", "domain": "", "client-id": "wireapp://CzbfFjDOQrenCbDxVmgnFw!594930e9d50bb175@wire.com", "handle": "wireapp://%40alice_wire@wire.com"}`
	emptyClientID := `{"name": "device", "domain": "wire.com", "client-id": "", "handle": "wireapp://%40alice_wire@wire.com"}`
	tests := []struct {
		name       string
		data       []byte
		wantWireID DeviceID
		wantErr    bool
	}{
		{name: "ok", data: []byte(ok), wantWireID: DeviceID{Name: "device", Domain: "wire.com", ClientID: "wireapp://CzbfFjDOQrenCbDxVmgnFw!594930e9d50bb175@wire.com", Handle: "wireapp://%40alice_wire@wire.com"}},
		{name: "fail/json", data: []byte(failJSON), wantErr: true},
		{name: "fail/empty-handle", data: []byte(emptyHandle), wantErr: true},
		{name: "fail/empty-name", data: []byte(emptyName), wantErr: true},
		{name: "fail/empty-domain", data: []byte(emptyDomain), wantErr: true},
		{name: "fail/empty-client-id", data: []byte(emptyClientID), wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotWireID, err := ParseDeviceID(tt.data)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.wantWireID, gotWireID)
		})
	}
}

func TestParseClientID(t *testing.T) {
	tests := []struct {
		name        string
		clientID    string
		want        ClientID
		expectedErr error
	}{
		{name: "ok", clientID: "wireapp://CzbfFjDOQrenCbDxVmgnFw!594930e9d50bb175@wire.com", want: ClientID{Scheme: "wireapp", Username: "CzbfFjDOQrenCbDxVmgnFw", DeviceID: "594930e9d50bb175", Domain: "wire.com"}},
		{name: "fail/uri", clientID: "bla", expectedErr: errors.New(`invalid Wire client ID URI "bla": error parsing bla: scheme is missing`)},
		{name: "fail/scheme", clientID: "not-wireapp://bla.com", expectedErr: errors.New(`invalid Wire client ID scheme "not-wireapp"; expected "wireapp"`)},
		{name: "fail/username", clientID: "wireapp://user@wire.com", expectedErr: errors.New(`invalid Wire client ID username "user"`)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseClientID(tt.clientID)
			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
