package contracts

// --- BRONZE LAYER ---

// ProductsRaw represents the ingested data in the bronze layer.
type ProductsRaw struct {
	SKU           string   `json:"sku"`
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	Currency      string   `json:"currency"`
	Price         string   `json:"price"` // kept as string due to upstream inconsistencies
	Vendor        string   `json:"vendor"`
	ReceivedAt    string   `json:"received_at"`
	IngestedAt    string   `json:"ingested_at"`
	MissingFields []string `json:"missing_fields,omitempty"`
	RawSource     string   `json:"source"`
}

// --- SILVER LAYER ---

// ProductsNormalized represents cleaned and standardized product records in the silver layer.
type ProductsNormalized struct {
	SKU          string  `json:"sku"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Currency     string  `json:"currency"`
	Price        float64 `json:"price"`
	Vendor       string  `json:"vendor"`
	Standardized bool    `json:"standardized"`
	ValidatedAt  string  `json:"validated_at"`
	IngestedAt   string  `json:"ingested_at"`
}

// --- GOLD LAYER ---

// ProductsCurated represents enriched, business-ready product records in the gold layer.
type ProductsCurated struct {
	SKU              string   `json:"sku"`
	FamilyID         string   `json:"family_id"`
	Name             string   `json:"name"`
	Description      string   `json:"description"`
	Category         string   `json:"category"`
	Currency         string   `json:"currency"`
	Price            float64  `json:"price"`
	LTVScore         float64  `json:"ltv_score"`
	Vendor           string   `json:"vendor"`
	Discontinued     bool     `json:"discontinued"`
	FirstAvailableAt string   `json:"first_available_at"`
	LastSeenAt       string   `json:"last_seen_at"`
	Tags             []string `json:"tags,omitempty"`
	EnrichedAt       string   `json:"enriched_at"`
}

// Each struct corresponds to a medallion stage contract — bronze, silver, and gold.
