package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/runink/finance-platform/contracts/cdm_trade/fdc3events"
	"github.com/segmentio/kafka-go"
	"github.com/snowflakedb/gosnowflake"
)

func getenv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Missing required environment variable: %s", key)
	}
	return val
}

func main() {
	ctx := context.Background()

	// Load Kafka config from env vars
	kafkaBrokers := strings.Split(getenv("CDM_KAFKA_BROKERS"), ",") // e.g., localhost:9092
	kafkaTopic := getenv("CDM_KAFKA_TOPIC")                         // e.g., topics.trade_events
	kafkaGroup := getenv("CDM_KAFKA_GROUP_ID")                      // e.g., runink-group

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: kafkaBrokers,
		Topic:   kafkaTopic,
		GroupID: kafkaGroup,
	})
	defer r.Close()

	// Load Snowflake config from env vars
	dsn, err := gosnowflake.DSN(&gosnowflake.Config{
		Account:  getenv("CDM_SF_ACCOUNT"),  // e.g., my_snowflake
		User:     getenv("CDM_SF_USER"),     // e.g., etl_user
		Password: getenv("CDM_SF_PASSWORD"), // e.g., secretpass
		Database: getenv("CDM_SF_DATABASE"), // e.g., FINANCE
		Schema:   getenv("CDM_SF_SCHEMA"),   // e.g., TRADE_PIPELINE
		Role:     getenv("CDM_SF_ROLE"),     // e.g., ETL_RUNNER
	})
	if err != nil {
		log.Fatalf("Failed to construct DSN: %v", err)
	}

	db, err := sql.Open("snowflake", dsn)
	if err != nil {
		log.Fatalf("Failed to open Snowflake connection: %v", err)
	}
	defer db.Close()

	for {
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			log.Printf("Kafka read error: %v", err)
			continue
		}

		var event fdc3events.CDMTradeEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("Unmarshal failed: %v", err)
			continue
		}

		decoded := fdc3events.DecodeCDMEvents(event)
		validated := fdc3events.ValidateLifecycle(decoded)

		if !validated.IsValid {
			log.Printf("Invalid trade %s: %s", validated.TradeID, validated.ValidationErr)
			_, err := db.ExecContext(ctx, `INSERT INTO CONTROL.INVALID_CDM_EVENTS (TRADE_ID, PRODUCT, EVENT_TYPE, ERROR) VALUES (?, ?, ?, ?)`,
				validated.TradeID, validated.Product, validated.EventType, validated.ValidationErr)
			if err != nil {
				log.Printf("DLQ insert failed: %v", err)
			}
			continue
		}

		tagged := fdc3events.TagWithFDC3Context(validated)
		log.Printf("Tagged valid trade: %s", tagged.TradeID)

		_, err = db.ExecContext(ctx, `INSERT INTO GOLD.TRADE_ENRICHED_FDC3 (TRADE_ID, PRODUCT, EVENT_TYPE, CONTEXT) VALUES (?, ?, ?, ?)`,
			tagged.TradeID, tagged.Product, tagged.EventType, tagged.FDC3Context)
		if err != nil {
			log.Printf("Snowflake insert failed: %v", err)
		}
	}
}
