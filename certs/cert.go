package certs

import "embed"

// Cert переменная с файловой системой для TLS сертификатов
//
//go:embed "*.pem"
var Cert embed.FS
