package syncer

type BlockSyncerAPI interface {
	// Start star block syncer
	Start() error
	// Stop stop block syncer
	Stop()
}
