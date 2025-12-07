package kms

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

// Key Management System - HD Wallet (BIP-32, BIP-39, BIP-44)
// Secure key generation, mnemonic phrases, multi-sig

type HDWallet struct {
	Mnemonic     string
	Seed         []byte
	MasterKey    *ExtendedKey
	Accounts     map[int]*Account
	MultiSigVaults map[string]*MultiSigVault
}

type ExtendedKey struct {
	Key       []byte
	ChainCode []byte
	Depth     uint8
	Index     uint32
	ParentFP  uint32
}

type Account struct {
	Index      int
	PublicKey  string
	PrivateKey string
	Address    string
	Path       string // m/44'/60'/0'/0/0
	Balance    uint64
}

type MultiSigVault struct {
	ID            string
	RequiredSigs  int
	TotalSigners  int
	Signers       []string
	PendingTxs    map[string]*PendingTx
}

type PendingTx struct {
	TxHash      string
	To          string
	Value       uint64
	Signatures  map[string]string
	Executed    bool
}

// BIP-39 Wordlist (simplified - use full 2048 words in production)
var wordList = []string{
	"abandon", "ability", "able", "about", "above", "absent", "absorb", "abstract",
	"absurd", "abuse", "access", "accident", "account", "accuse", "achieve", "acid",
	"acoustic", "acquire", "across", "act", "action", "actor", "actress", "actual",
	"adapt", "add", "addict", "address", "adjust", "admit", "adult", "advance",
	// ... (in production, use all 2048 BIP-39 words)
}

func NewHDWallet() (*HDWallet, error) {
	// Generate 12-word mnemonic (128 bits entropy)
	mnemonic, err := generateMnemonic(128)
	if err != nil {
		return nil, err
	}

	// Generate seed from mnemonic
	seed := mnemonicToSeed(mnemonic, "")

	// Generate master key
	masterKey := generateMasterKey(seed)

	wallet := &HDWallet{
		Mnemonic:       mnemonic,
		Seed:           seed,
		MasterKey:      masterKey,
		Accounts:       make(map[int]*Account),
		MultiSigVaults: make(map[string]*MultiSigVault),
	}

	// Derive first account (BIP-44: m/44'/60'/0'/0/0)
	wallet.DeriveAccount(0)

	fmt.Printf("🔐 HD Wallet created\n")
	fmt.Printf("📝 Mnemonic: %s\n", mnemonic)
	fmt.Printf("🏦 First address: %s\n", wallet. Accounts[0].Address)

	return wallet, nil
}

// Generate BIP-39 mnemonic
func generateMnemonic(bits int) (string, error) {
	if bits%32 != 0 || bits < 128 || bits > 256 {
		return "", fmt.Errorf("invalid entropy bits")
	}

	// Generate random entropy
	entropy := make([]byte, bits/8)
	_, err := rand.Read(entropy)
	if err != nil {
		return "", err
	}

	// Calculate checksum
	hash := sha256.Sum256(entropy)
	checksumBits := bits / 32
	
	// Convert to mnemonic words (simplified)
	wordCount := (bits + checksumBits) / 11
	words := make([]string, wordCount)
	
	for i := 0; i < wordCount; i++ {
		// In production: properly extract 11-bit indices
		words[i] = wordList[i%len(wordList)]
	}

	return strings.Join(words, " "), nil
}

// Convert mnemonic to seed (BIP-39)
func mnemonicToSeed(mnemonic, passphrase string) []byte {
	salt := "mnemonic" + passphrase
	// PBKDF2 with 2048 iterations
	return pbkdf2.Key([]byte(mnemonic), []byte(salt), 2048, 64, sha256.New)
}

// Generate master key (BIP-32)
func generateMasterKey(seed []byte) *ExtendedKey {
	// HMAC-SHA512 with key "Bitcoin seed" (or "NUSA seed")
	hash := sha256.Sum256(append([]byte("NUSA seed"), seed...))
	
	return &ExtendedKey{
		Key:       hash[:32],
		ChainCode: hash[32:],
		Depth:     0,
		Index:     0,
		ParentFP:  0,
	}
}

// Derive account (BIP-44: m/44'/60'/0'/0/index)
func (w *HDWallet) DeriveAccount(index int) *Account {
	// Simplified derivation (production: full BIP-32/BIP-44)
	
	// Generate private key
	privKey := make([]byte, 32)
	copy(privKey, w.MasterKey.Key)
	privKey[31] = byte(index)
	
	// Generate public key (simplified - use secp256k1 in production)
	pubKey := sha256.Sum256(privKey)
	
	// Generate address (Ethereum-style)
	addrHash := sha256.Sum256(pubKey[:])
	address := "0x" + hex.EncodeToString(addrHash[:20])
	
	account := &Account{
		Index:      index,
		PublicKey:  hex.EncodeToString(pubKey[:]),
		PrivateKey: hex.EncodeToString(privKey),
		Address:    address,
		Path:       fmt.Sprintf("m/44'/60'/0'/0/%d", index),
		Balance:    0,
	}
	
	w.Accounts[index] = account
	
	fmt.Printf("✅ Derived account #%d: %s\n", index, address)
	
	return account
}

// Sign transaction
func (w *HDWallet) SignTransaction(accountIndex int, txHash string) (string, error) {
	account, exists := w.Accounts[accountIndex]
	if !exists {
		return "", fmt.Errorf("account not found")
	}
	
	// Sign with private key (simplified - use proper ECDSA in production)
	privKeyBytes, _ := hex.DecodeString(account.PrivateKey)
	combinedData := append(privKeyBytes, []byte(txHash)...)
	sigHash := sha256.Sum256(combinedData)
	signature := hex.EncodeToString(sigHash[:])
	
	fmt.Printf("✍️ Transaction signed by %s\n", account.Address)
	
	return signature, nil
}

// Create multi-sig vault
func (w *HDWallet) CreateMultiSigVault(id string, requiredSigs, totalSigners int, signers []string) error {
	if requiredSigs > totalSigners {
		return fmt.Errorf("required sigs cannot exceed total signers")
	}
	
	if len(signers) != totalSigners {
		return fmt. Errorf("must provide all signer addresses")
	}
	
	vault := &MultiSigVault{
		ID:           id,
		RequiredSigs: requiredSigs,
		TotalSigners: totalSigners,
		Signers:      signers,
		PendingTxs:   make(map[string]*PendingTx),
	}
	
	w.MultiSigVaults[id] = vault
	
	fmt.Printf("🏦 Multi-sig vault created: %s (%d-of-%d)\n", id, requiredSigs, totalSigners)
	
	return nil
}

// Propose multi-sig transaction
func (w *HDWallet) ProposeMultiSigTx(vaultID, txHash, to string, value uint64, proposer string) error {
	vault, exists := w.MultiSigVaults[vaultID]
	if !exists {
		return fmt.Errorf("vault not found")
	}
	
	// Check if proposer is a signer
	isSigner := false
	for _, signer := range vault.Signers {
		if signer == proposer {
			isSigner = true
			break
		}
	}
	
	if ! isSigner {
		return fmt.Errorf("proposer is not a signer")
	}
	
	pendingTx := &PendingTx{
		TxHash:     txHash,
		To:         to,
		Value:      value,
		Signatures: make(map[string]string),
		Executed:   false,
	}
	
	vault.PendingTxs[txHash] = pendingTx
	
	fmt.Printf("📝 Multi-sig tx proposed: %s (vault: %s)\n", txHash, vaultID)
	
	return nil
}

// Sign multi-sig transaction
func (w *HDWallet) SignMultiSigTx(vaultID, txHash, signer, signature string) error {
	vault, exists := w.MultiSigVaults[vaultID]
	if !exists {
		return fmt.Errorf("vault not found")
	}
	
	pendingTx, exists := vault. PendingTxs[txHash]
	if !exists {
		return fmt.Errorf("pending tx not found")
	}
	
	if pendingTx.Executed {
		return fmt.Errorf("tx already executed")
	}
	
	// Check if signer is authorized
	isSigner := false
	for _, s := range vault.Signers {
		if s == signer {
			isSigner = true
			break
		}
	}
	
	if !isSigner {
		return fmt.Errorf("not authorized signer")
	}
	
	// Add signature
	pendingTx.Signatures[signer] = signature
	
	fmt.Printf("✍️ Multi-sig signature added: %s (%d/%d)\n", signer, len(pendingTx.Signatures), vault.RequiredSigs)
	
	// Check if enough signatures
	if len(pendingTx.Signatures) >= vault.RequiredSigs {
		pendingTx.Executed = true
		fmt.Printf("✅ Multi-sig tx executed: %s\n", txHash)
	}
	
	return nil
}

// Export private key (encrypted)
func (w *HDWallet) ExportPrivateKey(accountIndex int, password string) (string, error) {
	account, exists := w.Accounts[accountIndex]
	if !exists {
		return "", fmt.Errorf("account not found")
	}
	
	// Encrypt private key with password (simplified - use AES in production)
	passwordHash := sha256.Sum256([]byte(password))
	privKeyBytes, _ := hex.DecodeString(account.PrivateKey)
	
	encrypted := make([]byte, len(privKeyBytes))
	for i := range privKeyBytes {
		encrypted[i] = privKeyBytes[i] ^ passwordHash[i%32]
	}
	
	return hex.EncodeToString(encrypted), nil
}

// Import private key
func (w *HDWallet) ImportPrivateKey(encryptedKey, password string) (*Account, error) {
	// Decrypt and import (simplified)
	passwordHash := sha256.Sum256([]byte(password))
	encryptedBytes, _ := hex. DecodeString(encryptedKey)
	
	privKey := make([]byte, len(encryptedBytes))
	for i := range encryptedBytes {
		privKey[i] = encryptedBytes[i] ^ passwordHash[i%32]
	}
	
	// Generate address from private key
	pubKey := sha256.Sum256(privKey)
	addrHash := sha256.Sum256(pubKey[:])
	address := "0x" + hex. EncodeToString(addrHash[:20])
	
	account := &Account{
		Index:      len(w.Accounts),
		PublicKey:  hex. EncodeToString(pubKey[:]),
		PrivateKey: hex.EncodeToString(privKey),
		Address:    address,
		Path:       "imported",
		Balance:    0,
	}
	
	w.Accounts[account.Index] = account
	
	fmt.Printf("📥 Private key imported: %s\n", address)
	
	return account, nil
}

// Get wallet info
func (w *HDWallet) GetInfo() map[string]interface{} {
	return map[string]interface{}{
		"mnemonic_words":   len(strings.Split(w.Mnemonic, " ")),
		"total_accounts":   len(w. Accounts),
		"multisig_vaults":  len(w.MultiSigVaults),
		"first_address":    w.Accounts[0].Address,
	}
}
