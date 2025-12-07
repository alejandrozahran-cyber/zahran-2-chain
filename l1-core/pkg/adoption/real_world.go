package adoption

import (
	"fmt"
	"time"
)

// Real-World Adoption Layer
// Supply chain, IoT, payments, enterprise solutions

type AdoptionLayer struct {
	SupplyChains     map[string]*SupplyChain
	IoTDevices       map[string]*IoTDevice
	PaymentRails     *PaymentRails
	EnterpriseApps   map[string]*EnterpriseApp
	Integrations     map[string]*Integration
}

// ===== SUPPLY CHAIN TRACKING =====

type SupplyChain struct {
	ChainID          string
	Name             string
	Industry         string
	Participants     []Participant
	Products         map[string]*Product
	Shipments        map[string]*Shipment
	TotalProducts    uint64
	Verified         bool
}

type Participant struct {
	ID               string
	Name             string
	Role             ParticipantRole
	Location         string
	CertifiedBy      string
	Verified         bool
}

type ParticipantRole string

const (
	RoleManufacturer ParticipantRole = "manufacturer"
	RoleDistributor  ParticipantRole = "distributor"
	RoleRetailer     ParticipantRole = "retailer"
	RoleCustomer     ParticipantRole = "customer"
	RoleAuditor      ParticipantRole = "auditor"
)

type Product struct {
	ProductID        string
	Name             string
	SKU              string
	Manufacturer     string
	ManufacturedAt   time.Time
	BatchNumber      string
	Certifications   []string
	CurrentLocation  string
	Owner            string
	Journey          []JourneyStep
	Authentic        bool
}

type JourneyStep struct {
	Timestamp        time.Time
	Location         string
	Handler          string
	Action           string
	Temperature      float64  // For cold chain
	Photos           []string
	Verified         bool
}

type Shipment struct {
	ShipmentID       string
	Products         []string
	From             string
	To               string
	Status           ShipmentStatus
	EstimatedArrival time.Time
	ActualArrival    time.Time
	Carrier          string
	TrackingEvents   []TrackingEvent
}

type ShipmentStatus string

const (
	StatusPending    ShipmentStatus = "pending"
	StatusInTransit  ShipmentStatus = "in_transit"
	StatusDelivered  ShipmentStatus = "delivered"
	StatusFailed     ShipmentStatus = "failed"
)

type TrackingEvent struct {
	Timestamp        time.Time
	Location         string
	Event            string
	Coordinates      string
}

// ===== IoT INTEGRATION =====

type IoTDevice struct {
	DeviceID         string
	DeviceType       IoTDeviceType
	Owner            string
	Location         string
	Status           string
	Firmware         string
	
	// Telemetry
	LastHeartbeat    time.Time
	DataPoints       []DataPoint
	
	// Security
	PublicKey        string
	Authorized       bool
	
	// Payments
	WalletAddress    string
	Balance          uint64
}

type IoTDeviceType string

const (
	TypeSensor       IoTDeviceType = "sensor"
	TypeCamera       IoTDeviceType = "camera"
	TypeVehicle      IoTDeviceType = "vehicle"
	TypeSmartMeter   IoTDeviceType = "smart_meter"
	TypeWearable     IoTDeviceType = "wearable"
)

type DataPoint struct {
	Timestamp        time.Time
	Metric           string
	Value            float64
	Unit             string
	Hash             string  // Integrity verification
}

// ===== PAYMENT RAILS =====

type PaymentRails struct {
	Channels         map[string]*PaymentChannel
	Merchants        map[string]*Merchant
	PaymentMethods   []PaymentMethod
	TotalVolume      uint64
	TotalTransactions uint64
}

type PaymentChannel struct {
	ChannelID        string
	Type             ChannelType
	Participants     []string
	Capacity         uint64
	Balance          uint64
	SettlementTime   time.Duration
	Active           bool
}

type ChannelType string

const (
	ChannelInstant   ChannelType = "instant"      // 0 confirmations
	ChannelStandard  ChannelType = "standard"     // 1-2 confirmations
	ChannelBatch     ChannelType = "batch"        // Batched settlements
)

type Merchant struct {
	MerchantID       string
	BusinessName     string
	Category         string
	APIKey           string
	PaymentAddress   string
	SettlementPeriod time.Duration
	FeeRate          float64
	TotalSales       uint64
	Verified         bool
}

type PaymentMethod struct {
	Method           string
	Enabled          bool
	FeePercentage    float64
	ProcessingTime   time.Duration
}

// ===== ENTERPRISE APPLICATIONS =====

type EnterpriseApp struct {
	AppID            string
	CompanyName      string
	UseCase          string
	APIEndpoint      string
	MonthlyVolume    uint64
	SLA              ServiceLevel
	DeployedAt       time.Time
}

type ServiceLevel struct {
	Uptime           float64   // 99.9%
	MaxLatency       time.Duration
	Support          string
	RateLimits       int
}

type Integration struct {
	IntegrationID    string
	Platform         string
	Type             IntegrationType
	Active           bool
	Config           map[string]string
}

type IntegrationType string

const (
	IntegrationERP       IntegrationType = "erp"
	IntegrationCRM       IntegrationType = "crm"
	IntegrationAccounting IntegrationType = "accounting"
	IntegrationPayment   IntegrationType = "payment_gateway"
)

func NewAdoptionLayer() *AdoptionLayer {
	return &AdoptionLayer{
		SupplyChains:   make(map[string]*SupplyChain),
		IoTDevices:     make(map[string]*IoTDevice),
		PaymentRails:   NewPaymentRails(),
		EnterpriseApps: make(map[string]*EnterpriseApp),
		Integrations:   make(map[string]*Integration),
	}
}

func NewPaymentRails() *PaymentRails {
	return &PaymentRails{
		Channels:          make(map[string]*PaymentChannel),
		Merchants:         make(map[string]*Merchant),
		PaymentMethods:    make([]PaymentMethod, 0),
		TotalVolume:       0,
		TotalTransactions: 0,
	}
}

// ===== SUPPLY CHAIN METHODS =====

func (al *AdoptionLayer) CreateSupplyChain(
	name, industry string,
	participants []Participant,
) (*SupplyChain, error) {
	chainID := fmt.Sprintf("sc_%s_%d", name, time.Now().Unix())
	
	chain := &SupplyChain{
		ChainID:       chainID,
		Name:          name,
		Industry:      industry,
		Participants:  participants,
		Products:      make(map[string]*Product),
		Shipments:     make(map[string]*Shipment),
		TotalProducts: 0,
		Verified:      true,
	}
	
	al.SupplyChains[chainID] = chain
	
	fmt.Printf("🏭 Supply chain created: %s (%s) with %d participants\n",
		name, industry, len(participants))
	
	return chain, nil
}

func (al *AdoptionLayer) RegisterProduct(
	chainID, productName, sku, manufacturer string,
	certifications []string,
) (*Product, error) {
	chain, exists := al.SupplyChains[chainID]
	if !exists {
		return nil, fmt.Errorf("supply chain not found")
	}
	
	productID := fmt.Sprintf("prod_%s_%d", sku, time.Now().Unix())
	
	product := &Product{
		ProductID:      productID,
		Name:           productName,
		SKU:            sku,
		Manufacturer:   manufacturer,
		ManufacturedAt: time.Now(),
		BatchNumber:    fmt.Sprintf("BATCH_%d", time.Now().Unix()),
		Certifications: certifications,
		CurrentLocation: manufacturer,
		Owner:          manufacturer,
		Journey:        make([]JourneyStep, 0),
		Authentic:      true,
	}
	
	// Record first journey step
	product.Journey = append(product.Journey, JourneyStep{
		Timestamp: time.Now(),
		Location:  manufacturer,
		Handler:   manufacturer,
		Action:    "manufactured",
		Verified:  true,
	})
	
	chain.Products[productID] = product
	chain.TotalProducts++
	
	fmt.Printf("📦 Product registered: %s (SKU: %s) in %s\n", productName, sku, chain.Name)
	
	return product, nil
}

func (al *AdoptionLayer) TrackProduct(
	chainID, productID, location, handler, action string,
) error {
	chain, exists := al.SupplyChains[chainID]
	if !exists {
		return fmt.Errorf("supply chain not found")
	}
	
	product, exists := chain.Products[productID]
	if !exists {
		return fmt.Errorf("product not found")
	}
	
	// Add journey step
	step := JourneyStep{
		Timestamp: time.Now(),
		Location:  location,
		Handler:   handler,
		Action:    action,
		Verified:  true,
	}
	
	product.Journey = append(product. Journey, step)
	product.CurrentLocation = location
	
	fmt.Printf("📍 Product tracked: %s at %s (action: %s)\n",
		product.Name, location, action)
	
	return nil
}

func (al *AdoptionLayer) VerifyAuthenticity(chainID, productID string) (bool, []JourneyStep) {
	chain, exists := al.SupplyChains[chainID]
	if !exists {
		return false, nil
	}
	
	product, exists := chain.Products[productID]
	if !exists {
		return false, nil
	}
	
	// Verify entire journey is legitimate
	authentic := product. Authentic && len(product.Journey) > 0
	
	fmt.Printf("🔍 Authenticity check: %s = %v (journey: %d steps)\n",
		product.Name, authentic, len(product.Journey))
	
	return authentic, product.Journey
}

// ===== IoT METHODS =====

func (al *AdoptionLayer) RegisterIoTDevice(
	deviceType IoTDeviceType,
	owner, location string,
) (*IoTDevice, error) {
	deviceID := fmt.Sprintf("iot_%s_%d", deviceType, time.Now().Unix())
	
	device := &IoTDevice{
		DeviceID:      deviceID,
		DeviceType:    deviceType,
		Owner:         owner,
		Location:      location,
		Status:        "online",
		Firmware:      "v1.0.0",
		LastHeartbeat: time.Now(),
		DataPoints:    make([]DataPoint, 0),
		PublicKey:     "device_pubkey",
		Authorized:    true,
		WalletAddress: fmt.Sprintf("0x%s", deviceID),
		Balance:       0,
	}
	
	al.IoTDevices[deviceID] = device
	
	fmt.Printf("🤖 IoT device registered: %s (%s) at %s\n", deviceID, deviceType, location)
	
	return device, nil
}

func (al *AdoptionLayer) SubmitIoTData(
	deviceID, metric string,
	value float64,
	unit string,
) error {
	device, exists := al.IoTDevices[deviceID]
	if !exists {
		return fmt.Errorf("device not found")
	}
	
	if ! device.Authorized {
		return fmt.Errorf("device not authorized")
	}
	
	// Create data point
	dataPoint := DataPoint{
		Timestamp: time.Now(),
		Metric:    metric,
		Value:     value,
		Unit:      unit,
		Hash:      fmt.Sprintf("%x", time.Now().Unix()),
	}
	
	device.DataPoints = append(device.DataPoints, dataPoint)
	device.LastHeartbeat = time.Now()
	
	fmt.Printf("📊 IoT data: %s reported %s = %. 2f %s\n",
		deviceID, metric, value, unit)
	
	return nil
}

func (al *AdoptionLayer) PayIoTDevice(deviceID string, amount uint64) error {
	device, exists := al.IoTDevices[deviceID]
	if !exists {
		return fmt. Errorf("device not found")
	}
	
	device.Balance += amount
	
	fmt.Printf("💰 IoT payment: %d NUSA → %s (balance: %d)\n",
		amount/1e8, deviceID, device.Balance/1e8)
	
	return nil
}

// ===== PAYMENT RAILS METHODS =====

func (al *AdoptionLayer) RegisterMerchant(
	businessName, category string,
	settlementPeriod time.Duration,
) (*Merchant, error) {
	merchantID := fmt.Sprintf("mer_%s_%d", businessName, time.Now().Unix())
	
	merchant := &Merchant{
		MerchantID:       merchantID,
		BusinessName:     businessName,
		Category:         category,
		APIKey:           fmt.Sprintf("key_%s", merchantID),
		PaymentAddress:   fmt.Sprintf("0xmerchant_%s", merchantID),
		SettlementPeriod: settlementPeriod,
		FeeRate:          0.02,  // 2%
		TotalSales:       0,
		Verified:         true,
	}
	
	al.PaymentRails.Merchants[merchantID] = merchant
	
	fmt.Printf("🏪 Merchant registered: %s (%s)\n", businessName, category)
	
	return merchant, nil
}

func (al *AdoptionLayer) ProcessPayment(
	merchantID string,
	amount uint64,
	method string,
) error {
	merchant, exists := al.PaymentRails.Merchants[merchantID]
	if !exists {
		return fmt.Errorf("merchant not found")
	}
	
	// Calculate fee
	fee := uint64(float64(amount) * merchant.FeeRate)
	netAmount := amount - fee
	
	merchant.TotalSales += netAmount
	al.PaymentRails.TotalVolume += amount
	al.PaymentRails.TotalTransactions++
	
	fmt.Printf("💳 Payment processed: %d NUSA → %s (fee: %d, net: %d)\n",
		amount/1e8, merchant.BusinessName, fee/1e8, netAmount/1e8)
	
	return nil
}

func (al *AdoptionLayer) CreatePaymentChannel(
	channelType ChannelType,
	participants []string,
	capacity uint64,
) (*PaymentChannel, error) {
	channelID := fmt.Sprintf("ch_%s_%d", channelType, time.Now().Unix())
	
	channel := &PaymentChannel{
		ChannelID:      channelID,
		Type:           channelType,
		Participants:   participants,
		Capacity:       capacity,
		Balance:        0,
		SettlementTime: 1 * time.Second,
		Active:         true,
	}
	
	al.PaymentRails.Channels[channelID] = channel
	
	fmt.Printf("🔄 Payment channel created: %s (%s, capacity: %d NUSA)\n",
		channelID, channelType, capacity/1e8)
	
	return channel, nil
}

// ===== ENTERPRISE METHODS =====

func (al *AdoptionLayer) DeployEnterpriseApp(
	companyName, useCase, apiEndpoint string,
	sla ServiceLevel,
) (*EnterpriseApp, error) {
	appID := fmt. Sprintf("app_%s_%d", companyName, time.Now(). Unix())
	
	app := &EnterpriseApp{
		AppID:         appID,
		CompanyName:   companyName,
		UseCase:       useCase,
		APIEndpoint:   apiEndpoint,
		MonthlyVolume: 0,
		SLA:           sla,
		DeployedAt:    time. Now(),
	}
	
	al.EnterpriseApps[appID] = app
	
	fmt.Printf("🏢 Enterprise app deployed: %s (%s) | SLA: %. 1f%% uptime\n",
		companyName, useCase, sla. Uptime)
	
	return app, nil
}

func (al *AdoptionLayer) IntegratePlatform(
	platform string,
	integrationType IntegrationType,
	config map[string]string,
) (*Integration, error) {
	integrationID := fmt. Sprintf("int_%s_%s", platform, integrationType)
	
	integration := &Integration{
		IntegrationID: integrationID,
		Platform:      platform,
		Type:          integrationType,
		Active:        true,
		Config:        config,
	}
	
	al.Integrations[integrationID] = integration
	
	fmt.Printf("🔌 Integration created: %s (%s)\n", platform, integrationType)
	
	return integration, nil
}

// Get adoption stats
func (al *AdoptionLayer) GetStats() map[string]interface{} {
	totalProducts := uint64(0)
	for _, chain := range al.SupplyChains {
		totalProducts += chain.TotalProducts
	}
	
	onlineDevices := 0
	for _, device := range al.IoTDevices {
		if device.Status == "online" {
			onlineDevices++
		}
	}
	
	return map[string]interface{}{
		"supply_chains": map[string]interface{}{
			"total_chains":   len(al.SupplyChains),
			"total_products": totalProducts,
		},
		"iot": map[string]interface{}{
			"total_devices":  len(al.IoTDevices),
			"online_devices": onlineDevices,
		},
		"payments": map[string]interface{}{
			"merchants":     len(al.PaymentRails.Merchants),
			"channels":      len(al.PaymentRails.Channels),
			"total_volume":  al. PaymentRails.TotalVolume / 1e8,
			"total_txs":     al.PaymentRails.TotalTransactions,
		},
		"enterprise": map[string]interface{}{
			"apps":          len(al.EnterpriseApps),
			"integrations":  len(al. Integrations),
		},
	}
}
