package zeldb

import (
	"fmt"
	"github.com/dgraph-io/badger/v4"
	"log"
	"time"
)

func NewDb(olddir, datadir string) *ZelDB {
	dir := fmt.Sprintf("%s/gzel/%s", datadir, olddir)
	opts := badger.DefaultOptions(dir)

	opts = opts.WithSyncWrites(true)           // safer but slower, set false if you want faster writes at risk of crash loss
	opts = opts.WithValueThreshold(32)         // inline values <= 32 bytes (tune smaller for many small values)
	opts = opts.WithValueLogFileSize(64 << 20) // 64MB vlog files; smaller files allow faster GC but more files overall
	opts = opts.WithNumMemtables(3)            // number of memtables to hold before flush; more means higher mem usage but smoother writes

	opts = opts.WithNumLevelZeroTables(10)      // controls level 0 table count before compaction triggers
	opts = opts.WithNumLevelZeroTablesStall(20) // controls when writes stall due to too many level 0 tables
	opts = opts.WithMaxLevels(7)                // number of LSM levels, 7 is default and balanced
	opts = opts.WithCompactL0OnClose(true)      // compact level 0 on DB close to reduce fragmentation

	// Resource tuning
	opts = opts.WithNumCompactors(2) // number of compaction workers; increase for SSD, decrease for HDD
	opts = opts.WithLogger(nil)      // disable default logging or set your own logger

	db, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		for {
			time.Sleep(10 * time.Second)
			db.RunValueLogGC(0.3)

		}
	}()

	//os.Exit(12)
	return &ZelDB{
		DB: db,
	}

}
