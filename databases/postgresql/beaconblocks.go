package postgresql

import (
	"github.com/duyhtq/incognito-data-sync/models"
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

// BeaconBlockStore is Beacon Block postgresql store
type BeaconBlockStore struct {
	PGStoreAbs
}

// NewBeaconBlockStore initialises BeaconBlockStore with pg connection
func NewBeaconBlockStore(db *sqlx.DB) (*BeaconBlockStore, error) {
	return &BeaconBlockStore{PGStoreAbs{DB: db}}, nil
}

func (st *BeaconBlockStore) GetLatestProcessedBCHeight() (uint64, error) {
	sqlStr := `
		SELECT block_height FROM beacon_blocks
		ORDER BY block_height DESC
		LIMIT 1
	`
	bcHeights := []uint64{}
	err := st.DB.Select(&bcHeights, sqlStr)
	if err != nil {
		return 0, err
	}
	if len(bcHeights) == 0 {
		return 0, nil
	}
	return bcHeights[0], nil
}

func (st *BeaconBlockStore) StoreBeaconBlock(beaconBlockModel *models.BeaconBlock) error {
	sqlStr := `
		INSERT INTO beacon_blocks (block_hash, block_height, data, instructions, pre_block, next_block, created_time, block_version, epoch, round)
		VALUES (:block_hash, :block_height, :data, :instructions, :pre_block, :next_block, :created_time, :block_version, :epoch, :round)
		RETURNING id
	`
	tx := st.DB.MustBegin()
	defer tx.Commit()
	_, err := tx.NamedQuery(sqlStr, beaconBlockModel)
	return err
}
