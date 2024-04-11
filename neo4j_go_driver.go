package main

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log"
)

func main() {
	ctx := context.Background()
	dbUri := "neo4j://localhost:7687"
	dbUser := "neo4j"
	dbPassword := "secretgraph" // Replace with your new password after changing it

	driver, err := neo4j.NewDriverWithContext(dbUri, neo4j.BasicAuth(dbUser, dbPassword, ""), func(c *neo4j.Config) {
		c.Log = neo4j.ConsoleLogger(neo4j.ERROR) // Adjust the logging level as needed
	})
	if err != nil {
		log.Fatalf("Error creating driver: %v", err)
	}
	defer driver.Close(ctx)

	// Verify connectivity
	if err = driver.VerifyConnectivity(ctx); err != nil {
		log.Fatalf("Error verifying connectivity: %v", err)
	}
	log.Println("Connected to Neo4j!")

	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	_, err = session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		result, err := tx.Run(ctx, "CREATE (n:Person {name: $name}) RETURN n.name", map[string]interface{}{
			"name": "Alice",
		})
		if err != nil {
			return nil, err
		}
		if result.Next(ctx) {
			return result.Record().Values[0], nil
		}
		return nil, result.Err()
	})

	if err != nil {
		log.Fatalf("Error running query: %v", err)
	} else {
		log.Println("Created a person node named Alice")
	}
}
