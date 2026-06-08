package seed

import (
	"log"
	"time"

	"community-forum/backend/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {
	var userCount int64
	db.Model(&models.User{}).Count(&userCount)
	if userCount > 0 {
		return
	}

	log.Println("No data found — seeding database...")

	hash := func(pw string) string {
		h, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.MinCost)
		if err != nil {
			log.Fatalf("seed: failed to hash password: %v", err)
		}
		return string(h)
	}

	now := time.Date(2026, 4, 20, 10, 0, 0, 0, time.UTC)

	users := []models.User{
		{Username: "architect", Email: "architect@midnight-forge.io", Password: hash("password123"), Role: "admin", Bio: "Protocol architect. Building the next generation of decentralized infrastructure.", Stacks: `["Rust", "Solidity", "Go"]`, Avatar: "", CreatedAt: now, UpdatedAt: now},
		{Username: "cypher", Email: "cypher@midnight-forge.io", Password: hash("password123"), Role: "user", Bio: "Cryptography researcher. Breaking and building secure systems.", Stacks: `["Python", "C++", "Haskell"]`, Avatar: "", CreatedAt: now.Add(-24 * time.Hour), UpdatedAt: now.Add(-24 * time.Hour)},
		{Username: "nexus", Email: "nexus@midnight-forge.io", Password: hash("password123"), Role: "user", Bio: "Network engineer. Routing packets through the digital underworld.", Stacks: `["Go", "Rust", "C"]`, Avatar: "", CreatedAt: now.Add(-48 * time.Hour), UpdatedAt: now.Add(-48 * time.Hour)},
		{Username: "oracle", Email: "oracle@midnight-forge.io", Password: hash("password123"), Role: "user", Bio: "Data scientist. Finding patterns in the noise.", Stacks: `["Python", "R", "Julia"]`, Avatar: "", CreatedAt: now.Add(-72 * time.Hour), UpdatedAt: now.Add(-72 * time.Hour)},
	}

	for i := range users {
		if err := db.Create(&users[i]).Error; err != nil {
			log.Fatalf("seed: failed to create user %s: %v", users[i].Username, err)
		}
	}

	tags := []models.Tag{
		{Name: "architecture", Color: "#f97316"},
		{Name: "security", Color: "#ef4444"},
		{Name: "distributed-systems", Color: "#8b5cf6"},
		{Name: "cryptography", Color: "#06b6d4"},
		{Name: "protocol-design", Color: "#10b981"},
		{Name: "systems", Color: "#6366f1"},
	}

	for i := range tags {
		if err := db.Create(&tags[i]).Error; err != nil {
			log.Fatalf("seed: failed to create tag %s: %v", tags[i].Name, err)
		}
	}

	type threadDef struct {
		Title   string
		Slug    string
		Content string
		Author  int
		Tags    []int
		Status  string
		Views   uint
		Created time.Time
	}

	threadDefs := []threadDef{
		{
			Title: "Architectural Shift: Migrating from Monolith to Mesh",
			Slug:  "architectural-shift",
			Content: `In this post, I want to discuss the architectural shift from traditional monolithic systems to a mesh-based service architecture.

We've been running a monolith for the past three years, and it's reaching its limits. The deployment pipeline is bottlenecked, scaling is expensive, and developer velocity is dropping.

## The Problem

Our current monolith handles authentication, billing, notifications, and core business logic all in one process. A single bug in the billing module can bring down the entire system.

## Proposed Solution

We're looking at a service mesh approach using:

1. **Sidecar proxies** for observability and traffic management
2. **gRPC** for inter-service communication
3. **Event sourcing** for data consistency across services

I'd love to hear from teams that have made this transition. What were your biggest challenges?`,
			Author: 0, Tags: []int{0, 2, 4}, Status: "published", Views: 342, Created: now.Add(-24 * time.Hour),
		},
		{
			Title: "Deep Dive: Post-Quantum Cryptography Standards",
			Slug:  "post-quantum-crypto",
			Content: `The NIST post-quantum cryptography standardization process has been making significant progress. Here's my analysis of the leading candidates and what they mean for our industry.

## CRYSTALS-Kyber

Selected for general encryption, Kyber offers:

- Compact key sizes (800 bytes for public key)
- Efficient performance comparable to ECDH
- Strong security proofs based on Module-LWE

## CRYSTALS-Dilithium

Selected for digital signatures, Dilithium provides:

- Fast verification
- Reasonable signature sizes (~2.5KB)
- Resistance to side-channel attacks

## Migration Timeline

Industry experts suggest starting migration planning now, with full transitions expected by 2030.`,
			Author: 1, Tags: []int{1, 3}, Status: "published", Views: 567, Created: now.Add(-48 * time.Hour),
		},
		{
			Title: "Building a Distributed Consensus Protocol from Scratch",
			Slug:  "distributed-consensus-protocol",
			Content: `After months of research and prototyping, I'd like to share my experience building a custom consensus protocol for a closed-network system.

## Why Not Raft or Paxos?

Existing protocols assume open participation and adversarial environments. Our use case is different:

- All nodes are authenticated
- Network topology is known
- Latency is predictable

## Our Approach

We designed a lightweight consensus mechanism based on:

1. **Leader election** with round-robin rotation
2. **Write-ahead logging** with fsync guarantees
3. **Quorum-based commits** with configurable thresholds

The result is a protocol that achieves sub-100ms finality with 7 nodes.`,
			Author: 2, Tags: []int{2, 4, 5}, Status: "published", Views: 891, Created: now.Add(-72 * time.Hour),
		},
		{
			Title: "Zero-Knowledge Proofs: A Practical Introduction",
			Slug:  "zero-knowledge-proofs",
			Content: `Zero-knowledge proofs are often described as "magic" but the underlying mathematics is accessible. Let me break down how zk-SNARKs work in practice.

## Core Concepts

A ZKP allows a prover to convince a verifier of a statement's truth without revealing any information beyond the validity of the statement.

## Real-World Applications

1. **Private transactions** on blockchain networks
2. **Identity verification** without revealing personal data
3. **Scalability** through proof aggregation

## Code Example

The key insight is that any computation can be represented as an arithmetic circuit, and the prover can generate a proof that the circuit was executed correctly.`,
			Author: 1, Tags: []int{1, 3, 4}, Status: "published", Views: 234, Created: now.Add(-96 * time.Hour),
		},
		{
			Title: "The State of WebAssembly Beyond the Browser",
			Slug:  "wasm-beyond-browser",
			Content: `WebAssembly is expanding far beyond client-side web development. Here's a survey of the most exciting WASM use cases in 2026.

## Server-Side WASM

WASI (WebAssembly System Interface) enables running WASM modules on servers with:

- Near-native performance
- Strong sandboxing guarantees
- Language-agnostic plugin systems

## Edge Computing

Platforms like Cloudflare Workers and Fastly Compute are leveraging WASM for:

- Sub-millisecond cold starts
- Portable deployments across CDN nodes
- Secure multi-tenant execution

## Plugin Ecosystems

Several projects now use WASM as their plugin runtime, including Envoy proxy and various database engines.`,
			Author: 0, Tags: []int{5, 0}, Status: "published", Views: 178, Created: now.Add(-120 * time.Hour),
		},
		{
			Title: "Incident Report: Database Outage Analysis",
			Slug:  "incident-report-db-outage",
			Content: `Last week we experienced a 47-minute database outage. Here's the full post-mortem.

## Timeline

- 14:32 — P95 latency spikes to 12s
- 14:35 — Connection pool exhaustion
- 14:38 — Read replicas also degrade
- 14:42 — Primary goes read-only
- 14:47 — Engineering paged
- 15:19 — Full recovery

## Root Cause

A migration script ran an unindexed ALTER TABLE on a 50M-row table, causing a full table rewrite that blocked all queries.

## Action Items

1. Implement query timeout middleware
2. Add migration review checklist
3. Deploy read-only failover automation`,
			Author: 3, Tags: []int{5, 2}, Status: "published", Views: 1456, Created: now.Add(-36 * time.Hour),
		},
		{
			Title: "Optimizing Go Garbage Collection for Low-Latency Systems",
			Slug:  "go-gc-optimization",
			Content: `Go's garbage collector has improved dramatically, but for low-latency systems you still need to understand its behavior.

## GC Fundamentals

Go uses a concurrent, tri-color mark-and-sweep GC. Key parameters:

- **GOGC** — triggers GC when heap grows by this percentage (default: 100)
- **GOMEMLIMIT** — soft memory limit (Go 1.19+)

## Optimization Strategies

1. Reduce pointer density in hot paths
2. Pre-allocate slices with known capacity
3. Use object pools (sync.Pool) for frequently allocated types
4. Consider GOMEMLIMIT for predictable latency

## Benchmarks

With proper tuning, we reduced P99 GC pause from 2ms to under 200μs.`,
			Author: 2, Tags: []int{5}, Status: "published", Views: 723, Created: now.Add(-60 * time.Hour),
		},
		{
			Title: "Designing Rate-Limiting for Distributed APIs",
			Slug:  "rate-limiting-distributed",
			Content: `Rate limiting is critical for API stability, but designing it for distributed systems introduces unique challenges.

## Algorithms Compared

1. **Token Bucket** — Simple but allows bursts
2. **Sliding Window** — More accurate, higher memory
3. **Leaky Bucket** — Smooths traffic, adds latency

## Distributed Considerations

- **Consistency** — Centralized Redis works but is a SPOF
- **Performance** — Local rate limiting is faster but less accurate
- **Fairness** — Per-tenant limits need careful design

## Our Implementation

We settled on a hybrid approach: local token buckets synced periodically with a Redis cluster. P99 overhead is under 5μs per request.`,
			Author: 0, Tags: []int{0, 2, 5}, Status: "published", Views: 445, Created: now.Add(-84 * time.Hour),
		},
		{
			Title: "Hardware Security Modules: A Field Guide",
			Slug:  "hsm-field-guide",
			Content: `Hardware Security Modules (HSMs) are essential for protecting cryptographic keys at scale. Here's what I've learned deploying them.

## HSM Types

1. **Cloud HSMs** (AWS CloudHSM, Azure Dedicated HSM) — managed, expensive
2. **On-prem HSMs** — full control, high upfront cost
3. **Software HSMs** — flexible, less secure

## Key Management Best Practices

- Use separate HSMs for development and production
- Implement key rotation with overlap periods
- Audit all key access through HSM logs

## Common Pitfalls

- Underestimating latency (network round trips add up)
- Forgetting backup procedures
- Inadequate testing of failover scenarios`,
			Author: 3, Tags: []int{1, 5}, Status: "published", Views: 312, Created: now.Add(-108 * time.Hour),
		},
		{
			Title: "Structured Logging: Moving Beyond fmt.Println",
			Slug:  "structured-logging",
			Content: `After years of debugging with ad-hoc log statements, I'm convinced that structured logging is the single highest-ROI observability investment you can make.

## Why Structured Logging?

- Machine-parseable output
- Built-in context propagation
- Easy integration with log aggregation systems

## Implementation

Using zerolog or slog (Go 1.21+), you can add structured logging incrementally:

- Start with HTTP middleware that logs request IDs
- Add structured error types with metadata
- Integrate with OpenTelemetry for distributed tracing`,
			Author: 1, Tags: []int{5, 0}, Status: "published", Views: 198, Created: now.Add(-132 * time.Hour),
		},
	}

	var threads []models.Thread
	for _, t := range threadDefs {
		thread := models.Thread{
			Title:     t.Title,
			Slug:      t.Slug,
			Content:   t.Content,
			AuthorID:  users[t.Author].ID,
			Status:    t.Status,
			ViewCount: t.Views,
			CreatedAt: t.Created,
			UpdatedAt: t.Created,
		}
		if err := db.Create(&thread).Error; err != nil {
			log.Fatalf("seed: failed to create thread %s: %v", t.Slug, err)
		}
		for _, ti := range t.Tags {
			if err := db.Model(&thread).Association("Tags").Append(&tags[ti]); err != nil {
				log.Fatalf("seed: failed to associate tag: %v", err)
			}
		}
		threads = append(threads, thread)
	}

	comments := []struct {
		Content string
		Author  int
		Thread  int
		Created time.Time
	}{
		{Content: "Great analysis. We made a similar migration last year and the biggest pain point was data consistency. Event sourcing helped but added significant complexity.", Author: 1, Thread: 0, Created: now.Add(-12 * time.Hour)},
		{Content: "Have you considered an incremental approach? Start by extracting the billing service first since it has the clearest boundaries.", Author: 2, Thread: 0, Created: now.Add(-10 * time.Hour)},
		{Content: "Kyber's performance characteristics make it a clear winner for most applications. The real challenge will be key rotation at scale.", Author: 0, Thread: 1, Created: now.Add(-36 * time.Hour)},
		{Content: "We started our PQC migration last quarter. The biggest issue isn't the algorithms themselves but the ecosystem readiness.", Author: 3, Thread: 1, Created: now.Add(-30 * time.Hour)},
		{Content: "Interesting approach. Have you benchmarked this against Raft with reduced election timeouts?", Author: 0, Thread: 2, Created: now.Add(-60 * time.Hour)},
		{Content: "Yes, we found that our approach performs better in the constrained environment but lacks the battle-tested guarantees of Raft.", Author: 2, Thread: 2, Created: now.Add(-55 * time.Hour)},
		{Content: "The arithmetic circuit analogy is the key insight that made ZKPs click for me. Great explanation.", Author: 3, Thread: 3, Created: now.Add(-80 * time.Hour)},
		{Content: "For those wanting to dive deeper, I recommend the 'MoonMath Manual' — it's freely available and covers the math behind zk-SNARKs.", Author: 0, Thread: 3, Created: now.Add(-75 * time.Hour)},
		{Content: "We hit the exact same issue last month. Now we run all migrations through a automated review system that flags unindexed operations.", Author: 1, Thread: 5, Created: now.Add(-24 * time.Hour)},
		{Content: "Adding query timeouts saved us multiple times. I'd recommend setting statement_timeout at the PostgreSQL level as a safety net.", Author: 2, Thread: 5, Created: now.Add(-20 * time.Hour)},
		{Content: "sync.Pool is great but be careful with cleanup. We had a memory leak because pooled objects held references to large buffers.", Author: 1, Thread: 6, Created: now.Add(-48 * time.Hour)},
		{Content: "Have you looked at the hybrid approach from Google's paper on sliding window counters? It uses very little memory and is surprisingly accurate.", Author: 3, Thread: 7, Created: now.Add(-72 * time.Hour)},
	}

	for i, c := range comments {
		comment := models.Comment{
			Content:   c.Content,
			AuthorID:  users[c.Author].ID,
			ThreadID:  threads[c.Thread].ID,
			CreatedAt: c.Created,
			UpdatedAt: c.Created,
		}
		if err := db.Create(&comment).Error; err != nil {
			log.Fatalf("seed: failed to create comment %d: %v", i, err)
		}
	}

	votes := []struct {
		User   int
		Thread int
		Value  int8
	}{
		{0, 0, 1}, {1, 0, 1}, {2, 0, 1},
		{0, 1, 1}, {2, 1, 1}, {3, 1, -1},
		{1, 2, 1}, {3, 2, 1},
		{0, 3, 1}, {2, 3, 1},
		{1, 5, 1}, {3, 5, 1},
	}

	for _, v := range votes {
		vote := models.Vote{
			UserID:   users[v.User].ID,
			ThreadID: &threads[v.Thread].ID,
			Value:    v.Value,
		}
		if err := db.Create(&vote).Error; err != nil {
			log.Fatalf("seed: failed to create vote: %v", err)
		}
	}

	log.Println("Seed complete — 4 users, 6 tags, 10 threads, 12 comments, 11 votes created")
}
