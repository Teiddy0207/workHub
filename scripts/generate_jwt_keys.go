package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"
)

func main() {
	// Generate RSA key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Printf("Error generating private key: %v\n", err)
		return
	}

	// Encode private key to PEM
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	// Encode public key to PEM (PKIX format)
	publicKeyDER, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		fmt.Printf("Error marshaling public key: %v\n", err)
		return
	}

	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyDER,
	})

	// Encode to base64
	privateKeyBase64 := base64.StdEncoding.EncodeToString(privateKeyPEM)
	publicKeyBase64 := base64.StdEncoding.EncodeToString(publicKeyPEM)

	// Create config file content
	configContent := fmt.Sprintf(`app_name: "WorkHub"
node_env: "development"

webHttp:
  http_host: "localhost"
  http_address: 8080

webSocket:
  ws_host: "localhost"
  ws_address: 8081
  ws_allow: true

postgres:
  host: "localhost"
  port: 5432
  username: "postgres"
  password: "password"
  db_name: "workhub"
  max_idle_conns: 10
  max_open_conns: 100
  conn_max_lifetime: 3600

redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0
  pool_size: 10
  read_timeout: 3
  write_timeout: 3
  dial_timeout: 5
  timeout: 3

logger:
  encoding: "json"
  level: "info"
  zap_type: "development"
  disable_caller: false
  disable_stacktrace: false
  log_file: true
  payload: true

jwt:
  signing_method: "RS256"
  private_key: "%s"
  public_key: "%s"
  issuer: "WorkHub"
  refresh_token_expire: 604800  # 7 days
  long_token_expire: 86400     # 1 day
  short_token_expire: 3600     # 1 hour
  is_refresh_token: true
  validate_password: true
  len_token: 32
  token_expire: 3600

session_secret_key: "your-session-secret-key"
cookie_secret_key: "your-cookie-secret-key"
csrf_secret_key: "your-csrf-secret-key"

email:
  smtp_host: "smtp.gmail.com"
  smtp_port: 587
  username: "your-email@gmail.com"
  password: "your-app-password"
  from_name: "WorkHub"

aws:
  region: "us-east-1"
  bucket: "your-bucket"
  access_key_id: "your-access-key"
  secret_access_key: "your-secret-key"
  public_endpoint: "https://your-bucket.s3.amazonaws.com"
  api_endpoint: "https://s3.amazonaws.com"
  endpoint: "https://s3.amazonaws.com"
`, privateKeyBase64, publicKeyBase64)

	// Write config file
	err = os.WriteFile("config.yaml", []byte(configContent), 0644)
	if err != nil {
		fmt.Printf("Error writing config file: %v\n", err)
		return
	}

	fmt.Println("✅ JWT keys generated and config.yaml created successfully!")
	fmt.Println("📁 Private key (base64):", privateKeyBase64[:50]+"...")
	fmt.Println("📁 Public key (base64):", publicKeyBase64[:50]+"...")
	fmt.Println("📄 Config file: config.yaml")
}
