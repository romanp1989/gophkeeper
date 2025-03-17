package crypto

import (
	"bytes"
	"testing"
)

func TestDeriveKey(t *testing.T) {
	type testCase struct {
		name       string
		password   string
		salt       string
		wantKeyLen int
		wantErr    bool
	}

	testCases := []testCase{
		{
			name:       "valid_password_and_salt",
			password:   "secretPassword",
			salt:       "mySalt",
			wantKeyLen: 32,
			wantErr:    false,
		},
		{
			name:       "empty_password",
			password:   "",
			salt:       "salt",
			wantKeyLen: 32,
			wantErr:    true,
		},
		{
			name:       "empty_salt",
			password:   "password",
			salt:       "",
			wantKeyLen: 32,
			wantErr:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			key, err := DeriveKey(tc.password, tc.salt)

			if (err != nil) != tc.wantErr {
				t.Fatalf("DeriveKey() error = %v, wantErr = %v", err, tc.wantErr)
			}
			if err == nil && len(key) != tc.wantKeyLen {
				t.Errorf("DeriveKey() length = %d, want %d", len(key), tc.wantKeyLen)
			}
		})
	}
}

func TestEncrypt(t *testing.T) {
	type testCase struct {
		name      string
		plaintext string
		key       []byte
		wantErr   bool
	}

	testKey := bytes.Repeat([]byte{0xAA}, 32)

	testCases := []testCase{
		{
			name:      "valid_key_32_bytes",
			plaintext: "Hello, world!",
			key:       testKey,
			wantErr:   false,
		},
		{
			name:      "invalid_key_17_bytes",
			plaintext: "SecretMessage",
			key:       bytes.Repeat([]byte{0xBB}, 17),
			wantErr:   true,
		},
		{
			name:      "invalid_key_8_bytes",
			plaintext: "TooShortKey",
			key:       bytes.Repeat([]byte{0xCC}, 8),
			wantErr:   true,
		},
		{
			name:      "invalid_key_33_bytes",
			plaintext: "TooLongKey",
			key:       bytes.Repeat([]byte{0xDD}, 33),
			wantErr:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ciphertext, err := Encrypt(tc.plaintext, tc.key)

			if (err != nil) != tc.wantErr {
				t.Fatalf("Encrypt() error = %v, wantErr = %v", err, tc.wantErr)
			}

			if !tc.wantErr && ciphertext == "" {
				t.Error("Encrypt() returned empty ciphertext, expected non-empty")
			}
		})
	}
}

func TestDecrypt(t *testing.T) {
	type testCase struct {
		name        string
		ciphertext  string
		key         []byte
		shouldMatch string
		wantErr     bool
	}

	validKey := bytes.Repeat([]byte{0xAA}, 32)

	encryptedHello, err := Encrypt("Hello, world!", validKey)
	if err != nil {
		t.Fatalf("unable to prepare ciphertext for tests: %v", err)
	}

	testCases := []testCase{
		{
			name:        "valid_key_and_ciphertext",
			ciphertext:  encryptedHello,
			key:         validKey,
			shouldMatch: "Hello, world!",
			wantErr:     false,
		},
		{
			name:        "wrong_key",
			ciphertext:  encryptedHello,
			key:         bytes.Repeat([]byte{0xBB}, 32),
			shouldMatch: "",
			wantErr:     true,
		},
		{
			name:        "invalid_hex_string",
			ciphertext:  "not_hex_data",
			key:         validKey,
			shouldMatch: "",
			wantErr:     true,
		},
		{
			name:        "ciphertext_too_short",
			ciphertext:  "1234ABCD",
			key:         validKey,
			shouldMatch: "",
			wantErr:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			plaintext, err := Decrypt(tc.ciphertext, tc.key)
			if (err != nil) != tc.wantErr {
				t.Fatalf("Decrypt() error = %v, wantErr = %v", err, tc.wantErr)
			}
			if !tc.wantErr && plaintext != tc.shouldMatch {
				t.Errorf("Decrypt() got = %s, want = %s", plaintext, tc.shouldMatch)
			}
		})
	}
}
