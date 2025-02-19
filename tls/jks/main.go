package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/binary"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"time"
)

const (
	MagicNumber    = 0xFEEDFEED
	VERSION        = 2
	PrivateKeyTag  = 1
	TrustedCertTag = 2
)

type JKSEntry struct {
	Tag        int32
	Alias      string
	Timestamp  int64
	PrivateKey []byte
	CertChain  [][]byte
}

func main() {
	if len(os.Args) > 1 {
		// 인증서 파일 경로가 주어진 경우
		certPath := os.Args[1]
		convertCertToJKS(certPath)
	} else {
		// 인수가 없는 경우 새로운 인증서 생성
		generateNewJKS()
	}
}

func convertCertToJKS(certPath string) {
	// 인증서 파일 읽기
	certPEMData, err := os.ReadFile(certPath)
	if err != nil {
		panic(fmt.Sprintf("인증서 파일을 읽을 수 없습니다: %v", err))
	}

	// PEM 디코딩
	block, _ := pem.Decode(certPEMData)
	if block == nil {
		panic("PEM 형식의 인증서가 아닙니다")
	}

	// 인증서 파싱
	_, err = x509.ParseCertificate(block.Bytes)
	if err != nil {
		panic(fmt.Sprintf("인증서 파싱 실패: %v", err))
	}

	// JKS 엔트리 생성 (인증서만 포함)
	entry := JKSEntry{
		Tag:       TrustedCertTag, // 인증서만 있는 경우 TRUSTED_CERT_TAG 사용
		Alias:     "imported-cert",
		Timestamp: time.Now().UnixNano() / int64(time.Millisecond),
		CertChain: [][]byte{block.Bytes},
	}

	// JKS 파일 생성
	createJKSFile(entry)
	fmt.Println("인증서가 JKS 파일로 변환되었습니다: keystore.jks")
}

func generateNewJKS() {
	// 개인키 생성
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	// 인증서 템플릿 생성
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName:         "localhost",
			Organization:       []string{"Organization"},
			OrganizationalUnit: []string{"Dev"},
			Country:            []string{"KR"},
			Province:           []string{"Seoul"},
			Locality:           []string{"Seoul"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	// 자체 서명된 인증서 생성
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		panic(err)
	}

	// 개인키를 PEM 형식으로 저장
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	err = os.WriteFile("private.key", privateKeyPEM, 0600)
	if err != nil {
		panic(fmt.Sprintf("개인키 저장 실패: %v", err))
	}

	// 인증서를 PEM 형식으로 저장
	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: derBytes,
	})
	err = os.WriteFile("certificate.pem", certPEM, 0644)
	if err != nil {
		panic(fmt.Sprintf("인증서 저장 실패: %v", err))
	}

	// JKS 엔트리 생성
	entry := JKSEntry{
		Tag:        PrivateKeyTag,
		Alias:      "server-alias",
		Timestamp:  time.Now().UnixNano() / int64(time.Millisecond),
		PrivateKey: x509.MarshalPKCS1PrivateKey(privateKey),
		CertChain:  [][]byte{derBytes},
	}

	// JKS 파일 생성
	createJKSFile(entry)
	fmt.Println("새로운 JKS 파일이 생성되었습니다: keystore.jks")
	fmt.Printf("- JKS 파일: %s\n", filepath.Join(".", "keystore.jks"))
	fmt.Printf("- 인증서: %s\n", filepath.Join(".", "certificate.pem"))
	fmt.Printf("- 개인키: %s\n", filepath.Join(".", "private.key"))
}

func createJKSFile(entry JKSEntry) {
	var buf bytes.Buffer

	// Magic number 와 버전 쓰기
	if err := binary.Write(&buf, binary.BigEndian, uint32(MagicNumber)); err != nil {
		panic(err)
	}
	if err := binary.Write(&buf, binary.BigEndian, uint32(VERSION)); err != nil {
		panic(err)
	}

	// 엔트리 개수 쓰기
	if err := binary.Write(&buf, binary.BigEndian, uint32(1)); err != nil {
		panic(err)
	}

	// 엔트리 데이터 쓰기
	writeEntry(&buf, entry)

	// 파일에 저장
	err := os.WriteFile("keystore.jks", buf.Bytes(), 0600)
	if err != nil {
		panic(err)
	}
}

func writeEntry(buf *bytes.Buffer, entry JKSEntry) {
	// 태그 쓰기
	if err := binary.Write(buf, binary.BigEndian, entry.Tag); err != nil {
		panic(err)
	}

	// Alias 쓰기
	if err := binary.Write(buf, binary.BigEndian, uint16(len(entry.Alias))); err != nil {
		panic(err)
	}
	buf.WriteString(entry.Alias)

	// Timestamp 쓰기
	if err := binary.Write(buf, binary.BigEndian, entry.Timestamp); err != nil {
		panic(err)
	}

	if entry.Tag == PrivateKeyTag {
		// 개인키 쓰기
		if err := binary.Write(buf, binary.BigEndian, uint32(len(entry.PrivateKey))); err != nil {
			panic(err)
		}
		buf.Write(entry.PrivateKey)
	}

	// 인증서 체인 길이 쓰기
	if err := binary.Write(buf, binary.BigEndian, uint32(len(entry.CertChain))); err != nil {
		panic(err)
	}

	// 인증서 체인 데이터 쓰기
	for _, cert := range entry.CertChain {
		if err := binary.Write(buf, binary.BigEndian, uint32(len(cert))); err != nil {
			panic(err)
		}
		buf.Write(cert)
	}
}
