package storage

import (
	"crypto/sha256"
	"fmt"
	"time"
)

// Distributed Storage Layer - IPFS/Arweave style
// For: NFT metadata, GameFi assets, dApp files, AI datasets

type DistributedStorage struct {
	Files        map[string]*StoredFile
	Nodes        map[string]*StorageNode
	Replicas     int    // Minimum replicas
	ChunkSize    int    // 256KB chunks
	TotalStorage uint64 // Total bytes stored
}

type StoredFile struct {
	CID          string   // Content Identifier (hash)
	FileName     string
	FileSize     uint64
	ContentType  string
	Owner        string
	UploadedAt   time.Time
	Chunks       []string // Chunk CIDs
	Replicas     []string // Node addresses storing this file
	Pinned       bool     // Permanent storage
	AccessCount  uint64
	StorageCost  uint64 // NUSA cost per GB per year
}

type StorageNode struct {
	Address       string
	TotalSpace    uint64 // GB
	UsedSpace     uint64
	Bandwidth     uint64 // MB/s
	Uptime        float64
	FilesStored   int
	EarningsTotal uint64
	Online        bool
}

func NewDistributedStorage() *DistributedStorage {
	return &DistributedStorage{
		Files:        make(map[string]*StoredFile),
		Nodes:        make(map[string]*StorageNode),
		Replicas:     3, // Store on 3 nodes minimum
		ChunkSize:    256 * 1024, // 256KB
		TotalStorage: 0,
	}
}

// Upload file to distributed storage
func (ds *DistributedStorage) UploadFile(
	fileName string,
	content []byte,
	owner string,
	pinned bool,
) (string, error) {
	// Calculate content ID (CID)
	hash := sha256.Sum256(content)
	cid := fmt. Sprintf("Qm%x", hash) // IPFS-style CID

	// Check if file already exists (deduplication)
	if existing, exists := ds.Files[cid]; exists {
		fmt.Printf("✅ File already exists: %s (deduplicated)\n", cid)
		return existing.CID, nil
	}

	// Split into chunks
	chunks := ds.chunkFile(content)
	chunkCIDs := make([]string, len(chunks))

	for i, chunk := range chunks {
		chunkHash := sha256.Sum256(chunk)
		chunkCID := fmt.Sprintf("Qm%x", chunkHash)
		chunkCIDs[i] = chunkCID
	}

	// Select storage nodes (most available space + uptime)
	selectedNodes := ds.selectStorageNodes(ds.Replicas)

	if len(selectedNodes) < ds. Replicas {
		return "", fmt.Errorf("insufficient storage nodes available")
	}

	// Calculate storage cost (0.1 NUSA per GB per year)
	fileSizeGB := float64(len(content)) / (1024 * 1024 * 1024)
	storageCost := uint64(fileSizeGB * 0.1 * 1e8) // Convert to smallest unit

	file := &StoredFile{
		CID:         cid,
		FileName:    fileName,
		FileSize:    uint64(len(content)),
		ContentType: detectContentType(fileName),
		Owner:       owner,
		UploadedAt:  time.Now(),
		Chunks:      chunkCIDs,
		Replicas:    selectedNodes,
		Pinned:      pinned,
		AccessCount: 0,
		StorageCost: storageCost,
	}

	ds.Files[cid] = file
	ds.TotalStorage += file.FileSize

	// Update node usage
	for _, nodeAddr := range selectedNodes {
		if node, exists := ds.Nodes[nodeAddr]; exists {
			node.UsedSpace += file.FileSize
			node.FilesStored++
		}
	}

	fmt.Printf("📦 File uploaded: %s (%d bytes) on %d nodes | Cost: %d NUSA/year\n",
		cid, len(content), len(selectedNodes), storageCost)

	return cid, nil
}

// Retrieve file from distributed storage
func (ds *DistributedStorage) RetrieveFile(cid string) (*StoredFile, error) {
	file, exists := ds.Files[cid]
	if !exists {
		return nil, fmt. Errorf("file not found: %s", cid)
	}

	// Check if at least one replica is available
	availableNodes := 0
	for _, nodeAddr := range file.Replicas {
		if node, exists := ds.Nodes[nodeAddr]; exists && node.Online {
			availableNodes++
		}
	}

	if availableNodes == 0 {
		return nil, fmt.Errorf("no available nodes for file: %s", cid)
	}

	file.AccessCount++

	fmt.Printf("📥 File retrieved: %s (accessed %d times)\n", cid, file.AccessCount)

	return file, nil
}

// Pin file for permanent storage
func (ds *DistributedStorage) PinFile(cid string, owner string) bool {
	file, exists := ds.Files[cid]
	if !exists {
		return false
	}

	if file.Owner != owner {
		fmt. Println("❌ Not file owner")
		return false
	}

	file.Pinned = true

	fmt.Printf("📌 File pinned permanently: %s\n", cid)

	return true
}

// Register storage node
func (ds *DistributedStorage) RegisterNode(
	address string,
	totalSpace uint64,
	bandwidth uint64,
) bool {
	if _, exists := ds.Nodes[address]; exists {
		return false
	}

	node := &StorageNode{
		Address:       address,
		TotalSpace:    totalSpace,
		UsedSpace:     0,
		Bandwidth:     bandwidth,
		Uptime:        100.0,
		FilesStored:   0,
		EarningsTotal: 0,
		Online:        true,
	}

	ds.Nodes[address] = node

	fmt. Printf("🖥️ Storage node registered: %s (%d GB capacity)\n", address, totalSpace/(1024*1024*1024))

	return true
}

// Garbage collection - Remove unpinned, unused files
func (ds *DistributedStorage) GarbageCollect(maxAge time.Duration) int {
	removed := 0

	for cid, file := range ds.Files {
		// Don't remove pinned files
		if file. Pinned {
			continue
		}

		// Remove files not accessed in maxAge
		if time.Since(file.UploadedAt) > maxAge && file.AccessCount < 10 {
			ds.TotalStorage -= file.FileSize

			// Update node usage
			for _, nodeAddr := range file.Replicas {
				if node, exists := ds.Nodes[nodeAddr]; exists {
					node.UsedSpace -= file.FileSize
					node.FilesStored--
				}
			}

			delete(ds.Files, cid)
			removed++
		}
	}

	fmt.Printf("🗑️ Garbage collection: removed %d files\n", removed)

	return removed
}

// Chunk file into smaller pieces
func (ds *DistributedStorage) chunkFile(content []byte) [][]byte {
	chunks := make([][]byte, 0)

	for i := 0; i < len(content); i += ds.ChunkSize {
		end := i + ds.ChunkSize
		if end > len(content) {
			end = len(content)
		}
		chunks = append(chunks, content[i:end])
	}

	return chunks
}

// Select best storage nodes
func (ds *DistributedStorage) selectStorageNodes(count int) []string {
	selected := make([]string, 0)

	// Sort nodes by: available space + uptime
	for addr, node := range ds.Nodes {
		if ! node.Online {
			continue
		}

		availableSpace := node.TotalSpace - node.UsedSpace
		if availableSpace > 1024*1024*1024 { // Min 1GB free
			selected = append(selected, addr)
		}

		if len(selected) >= count {
			break
		}
	}

	return selected
}

// Detect content type from filename
func detectContentType(fileName string) string {
	// Simplified detection
	if len(fileName) > 4 {
		ext := fileName[len(fileName)-4:]
		switch ext {
		case ".jpg", ".png", ".gif":
			return "image"
		case ".mp4", ".avi":
			return "video"
		case ".mp3", ".wav":
			return "audio"
		case ".pdf":
			return "document"
		case "json":
			return "json"
		}
	}
	return "unknown"
}

// Get storage stats
func (ds *DistributedStorage) GetStats() map[string]interface{} {
	totalNodes := len(ds.Nodes)
	onlineNodes := 0
	totalCapacity := uint64(0)
	usedCapacity := uint64(0)

	for _, node := range ds.Nodes {
		if node.Online {
			onlineNodes++
		}
		totalCapacity += node.TotalSpace
		usedCapacity += node.UsedSpace
	}

	return map[string]interface{}{
		"total_files":    len(ds.Files),
		"total_storage":  ds.TotalStorage,
		"total_nodes":    totalNodes,
		"online_nodes":   onlineNodes,
		"total_capacity": totalCapacity,
		"used_capacity":  usedCapacity,
		"utilization":    float64(usedCapacity) / float64(totalCapacity) * 100,
	}
}
