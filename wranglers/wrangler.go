package wranglers

import (
	"strings"
	"time"

	"../contracts"
)

// -----------------------------
// Bronze Layer Transformations
// -----------------------------

func AddIngestMetadata(p contracts.ProductsRaw) contracts.ProductsRaw {
	p.IngestedAt = time.Now().Format(time.RFC3339)
	p.RawSource = "vendor_feed"
	return p
}

func TagMissingFields(p contracts.ProductsRaw) contracts.ProductsRaw {
	var missing []string
	if strings.TrimSpace(p.SKU) == "" {
		missing = append(missing, "sku")
	}
	if strings.TrimSpace(p.Name) == "" {
		missing = append(missing, "name")
	}
	if strings.TrimSpace(p.Price) == "" {
		missing = append(missing, "price")
	}
	p.MissingFields = missing
	return p
}

// -----------------------------
// Silver Layer Transformations
// -----------------------------

func TrimProductNames(p contracts.ProductsRaw) contracts.ProductsNormalized {
	return contracts.ProductsNormalized{
		SKU:         p.SKU,
		Name:        strings.TrimSpace(p.Name),
		Description: p.Description,
		Currency:    strings.ToUpper(p.Currency),
		Price:       0, // to be filled in StandardizeCurrency
		Vendor:      p.Vendor,
		IngestedAt:  p.IngestedAt,
	}
}

func StandardizeCurrency(p contracts.ProductsNormalized) contracts.ProductsNormalized {
	if p.Currency == "" {
		p.Currency = "USD"
	}
	p.Standardized = true
	return p
}

func FixEmptyDescriptions(p contracts.ProductsNormalized) contracts.ProductsNormalized {
	if strings.TrimSpace(p.Description) == "" {
		p.Description = "N/A"
	}
	p.ValidatedAt = time.Now().Format(time.RFC3339)
	return p
}

// -----------------------------
// Gold Layer Transformations
// -----------------------------

func GroupVariantsByFamily(p contracts.ProductsNormalized) contracts.ProductsCurated {
	return contracts.ProductsCurated{
		SKU:         p.SKU,
		FamilyID:    strings.Split(p.SKU, "-")[0],
		Name:        p.Name,
		Description: p.Description,
		Currency:    p.Currency,
		Price:       p.Price,
		Vendor:      p.Vendor,
		Tags:        []string{},
	}
}

func EnrichWithCategoryLTV(p contracts.ProductsCurated) contracts.ProductsCurated {
	switch strings.ToLower(p.Category) {
	case "electronics":
		p.LTVScore = 0.95
	case "apparel":
		p.LTVScore = 0.6
	default:
		p.LTVScore = 0.75
	}
	p.EnrichedAt = time.Now().Format(time.RFC3339)
	return p
}

func DetectDiscontinuedItems(p contracts.ProductsCurated) contracts.ProductsCurated {
	desc := strings.ToLower(p.Description)
	p.Discontinued = strings.Contains(desc, "discontinued") || strings.Contains(desc, "no longer available")
	return p
}
