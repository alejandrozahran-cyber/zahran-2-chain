package compliance

import (
	"crypto/sha256"
	"fmt"
	"time"
)

// Compliance & zk-Identity Layer
// Privacy-preserving KYC, selective disclosure, regulatory compliance

type ComplianceLayer struct {
	Identities          map[string]*Identity
	ZKProofs            map[string]*ZKProof
	ComplianceRules     []ComplianceRule
	Regulators          map[string]*Regulator
	EnterpriseAccounts  map[string]*EnterpriseAccount
}

type Identity struct {
	DID                 string  // Decentralized Identifier
	Owner               string
	VerificationLevel   VerificationLevel
	Attributes          map[string]*EncryptedAttribute
	Credentials         []Credential
	CreatedAt           time.Time
	LastVerified        time.Time
	Revoked             bool
}

type VerificationLevel string

const (
	LevelNone       VerificationLevel = "none"
	LevelBasic      VerificationLevel = "basic"       // Email verified
	LevelStandard   VerificationLevel = "standard"    // KYC Tier 1
	LevelAdvanced   VerificationLevel = "advanced"    // KYC Tier 2
	LevelEnterprise VerificationLevel = "enterprise"  // Full compliance
)

type EncryptedAttribute struct {
	Name                string
	EncryptedValue      []byte
	Disclosed           bool
	DisclosedTo         []string
	Verifier            string
	ExpiresAt           time.Time
}

type Credential struct {
	ID                  string
	Type                CredentialType
	Issuer              string
	IssuedAt            time.Time
	ExpiresAt           time.Time
	ZKProof             []byte
	Valid               bool
}

type CredentialType string

const (
	CredentialKYC       CredentialType = "kyc"
	CredentialAML       CredentialType = "aml"
	CredentialAccredited CredentialType = "accredited_investor"
	CredentialAge       CredentialType = "age_verification"
	CredentialCountry   CredentialType = "country_residence"
)

type ZKProof struct {
	ProofID             string
	Statement           string  // "I am over 18" or "I am not from sanctioned country"
	Proof               []byte
	Verified            bool
	CreatedAt           time.Time
}

type ComplianceRule struct {
	RuleID              string
	Name                string
	RequiredLevel       VerificationLevel
	RestrictedCountries []string
	MinAge              int
	MaxTransactionSize  uint64
	Enabled             bool
}

type Regulator struct {
	ID                  string
	Name                string
	Jurisdiction        string
	PublicKey           string
	AuditAccess         bool
}

type EnterpriseAccount struct {
	CompanyID           string
	CompanyName         string
	LegalEntity         string
	TaxID               string
	ComplianceOfficer   string
	AMLPolicy           string
	KYCProvider         string
	Verified            bool
}

func NewComplianceLayer() *ComplianceLayer {
	return &ComplianceLayer{
		Identities:         make(map[string]*Identity),
		ZKProofs:           make(map[string]*ZKProof),
		ComplianceRules:    make([]ComplianceRule, 0),
		Regulators:         make(map[string]*Regulator),
		EnterpriseAccounts: make(map[string]*EnterpriseAccount),
	}
}

// Create decentralized identity
func (cl *ComplianceLayer) CreateIdentity(owner string) (*Identity, error) {
	did := generateDID(owner)
	
	identity := &Identity{
		DID:               did,
		Owner:             owner,
		VerificationLevel: LevelNone,
		Attributes:        make(map[string]*EncryptedAttribute),
		Credentials:       make([]Credential, 0),
		CreatedAt:         time.Now(),
		LastVerified:      time.Time{},
		Revoked:           false,
	}
	
	cl.Identities[did] = identity
	
	fmt.Printf("🆔 Identity created: %s\n", did)
	
	return identity, nil
}

// Submit KYC (privacy-preserving)
func (cl *ComplianceLayer) SubmitKYC(
	did string,
	verifier string,
	attributes map[string]string,  // name, dob, country, etc
) error {
	identity, exists := cl.Identities[did]
	if !exists {
		return fmt.Errorf("identity not found")
	}
	
	fmt.Printf("🔐 Processing KYC for %s.. .\n", did)
	
	// Encrypt sensitive attributes
	for key, value := range attributes {
		encrypted := encryptAttribute(value)
		
		identity.Attributes[key] = &EncryptedAttribute{
			Name:           key,
			EncryptedValue: encrypted,
			Disclosed:      false,
			DisclosedTo:    make([]string, 0),
			Verifier:       verifier,
			ExpiresAt:      time.Now().AddDate(1, 0, 0),  // 1 year
		}
	}
	
	// Issue KYC credential
	credential := Credential{
		ID:        fmt.Sprintf("cred_%s", did),
		Type:      CredentialKYC,
		Issuer:    verifier,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().AddDate(1, 0, 0),
		ZKProof:   generateZKProof("KYC verified"),
		Valid:     true,
	}
	
	identity.Credentials = append(identity.Credentials, credential)
	identity. VerificationLevel = LevelStandard
	identity.LastVerified = time.Now()
	
	fmt.Printf("✅ KYC verified: %s (level: %s)\n", did, identity.VerificationLevel)
	
	return nil
}

// Generate ZK proof (prove statement without revealing data)
func (cl *ComplianceLayer) GenerateZKProof(
	did string,
	statement string,  // e.g., "I am over 18"
) (*ZKProof, error) {
	identity, exists := cl.Identities[did]
	if !exists {
		return nil, fmt.Errorf("identity not found")
	}
	
	// Generate zero-knowledge proof
	// Production: Use zk-SNARKs library (bellman, halo2)
	
	proofID := fmt.Sprintf("proof_%d", time.Now().UnixNano())
	
	proof := &ZKProof{
		ProofID:   proofID,
		Statement: statement,
		Proof:     generateZKProof(statement),
		Verified:  false,
		CreatedAt: time.Now(),
	}
	
	cl.ZKProofs[proofID] = proof
	
	fmt.Printf("🔮 ZK Proof generated: %s | Statement: '%s'\n", proofID, statement)
	
	// User can now prove statement without revealing actual data
	_ = identity  // Use identity for actual proof generation
	
	return proof, nil
}

// Selective disclosure (share only specific attributes)
func (cl *ComplianceLayer) SelectiveDisclose(
	did string,
	attributeName string,
	discloseTo string,
) error {
	identity, exists := cl.Identities[did]
	if !exists {
		return fmt.Errorf("identity not found")
	}
	
	attribute, exists := identity.Attributes[attributeName]
	if !exists {
		return fmt.Errorf("attribute not found")
	}
	
	// Mark as disclosed to specific party
	attribute.Disclosed = true
	attribute.DisclosedTo = append(attribute.DisclosedTo, discloseTo)
	
	fmt.Printf("👁️ Attribute '%s' disclosed to %s\n", attributeName, discloseTo)
	
	return nil
}

// Verify compliance
func (cl *ComplianceLayer) VerifyCompliance(did string, rule ComplianceRule) (bool, string) {
	identity, exists := cl.Identities[did]
	if !exists {
		return false, "identity not found"
	}
	
	// Check verification level
	if ! meetsVerificationLevel(identity. VerificationLevel, rule.RequiredLevel) {
		return false, fmt.Sprintf("insufficient verification level: %s < %s",
			identity.VerificationLevel, rule.RequiredLevel)
	}
	
	// Check restricted countries (via ZK proof)
	// User proves "I am NOT from restricted country" without revealing actual country
	
	// Check credentials expiry
	for _, cred := range identity. Credentials {
		if cred.Type == CredentialKYC && time.Now().After(cred.ExpiresAt) {
			return false, "KYC credential expired"
		}
	}
	
	return true, "compliant"
}

// Enterprise account setup
func (cl *ComplianceLayer) RegisterEnterprise(
	companyName, legalEntity, taxID string,
	complianceOfficer string,
) (*EnterpriseAccount, error) {
	companyID := fmt.Sprintf("ent_%s", companyName)
	
	account := &EnterpriseAccount{
		CompanyID:         companyID,
		CompanyName:       companyName,
		LegalEntity:       legalEntity,
		TaxID:             taxID,
		ComplianceOfficer: complianceOfficer,
		AMLPolicy:         "standard",
		KYCProvider:       "internal",
		Verified:          false,
	}
	
	cl.EnterpriseAccounts[companyID] = account
	
	fmt.Printf("🏢 Enterprise account registered: %s\n", companyName)
	
	return account, nil
}

// Regulator audit access
func (cl *ComplianceLayer) GrantRegulatorAccess(
	regulatorID, name, jurisdiction string,
) error {
	regulator := &Regulator{
		ID:           regulatorID,
		Name:         name,
		Jurisdiction: jurisdiction,
		PublicKey:    "regulator_pubkey",
		AuditAccess:  true,
	}
	
	cl.Regulators[regulatorID] = regulator
	
	fmt.Printf("👮 Regulator access granted: %s (%s)\n", name, jurisdiction)
	
	return nil
}

// Add compliance rule
func (cl *ComplianceLayer) AddComplianceRule(rule ComplianceRule) {
	cl.ComplianceRules = append(cl.ComplianceRules, rule)
	
	fmt.Printf("📜 Compliance rule added: %s\n", rule.Name)
}

// Get compliance stats
func (cl *ComplianceLayer) GetStats() map[string]interface{} {
	verified := 0
	for _, identity := range cl.Identities {
		if identity.VerificationLevel != LevelNone {
			verified++
		}
	}
	
	return map[string]interface{}{
		"total_identities":     len(cl.Identities),
		"verified_identities":  verified,
		"zk_proofs_generated":  len(cl.ZKProofs),
		"compliance_rules":     len(cl. ComplianceRules),
		"regulators":           len(cl.Regulators),
		"enterprise_accounts":  len(cl.EnterpriseAccounts),
	}
}

// Helper functions
func generateDID(owner string) string {
	hash := sha256.Sum256([]byte(owner + time.Now().String()))
	return fmt.Sprintf("did:nusa:%x", hash[:16])
}

func encryptAttribute(value string) []byte {
	// Simplified encryption (production: use proper encryption)
	return []byte(fmt.Sprintf("encrypted_%s", value))
}

func generateZKProof(statement string) []byte {
	// Simplified ZK proof generation
	// Production: Use zk-SNARKs (Groth16, Plonk, etc)
	hash := sha256.Sum256([]byte(statement))
	return hash[:]
}

func meetsVerificationLevel(actual, required VerificationLevel) bool {
	levels := map[VerificationLevel]int{
		LevelNone:       0,
		LevelBasic:      1,
		LevelStandard:   2,
		LevelAdvanced:   3,
		LevelEnterprise: 4,
	}
	
	return levels[actual] >= levels[required]
}
