package postgresql

import (
	"github.com/duyhtq/incognito-data-sync/models"
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

// PDEStatePGStore is PDE State postgresql store
type PDEStatePGStore struct {
	PGStoreAbs
}

// NewPDEStatePGStore initialises PDEStatePGStore with pg connection
func NewPDEStatePGStore(db *sqlx.DB) (*PDEStatePGStore, error) {
	return &PDEStatePGStore{PGStoreAbs{DB: db}}, nil
}

func (st *PDEStatePGStore) EnsurePDEStatePGStore() {
	// TODO: implement it
}

func (st *PDEStatePGStore) GetLatestProcessedBCHeight() (uint64, error) {
	sqlStr := `
		SELECT beacon_height FROM pde_pool_pairs
		ORDER BY beacon_height DESC
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

func (st *PDEStatePGStore) StorePDEPoolPair(pdeStateModel *models.PDEPoolPair) error {
	sqlStr := `
		INSERT INTO pde_pool_pairs (token1_id_str, token1_pool_value, token2_id_str, token2_pool_value, beacon_height, beacon_time_stamp)
		VALUES (:token1_id_str, :token1_pool_value, :token2_id_str, :token2_pool_value, :beacon_height, :beacon_time_stamp)
		RETURNING id
	`
	tx := st.DB.MustBegin()
	defer tx.Commit()
	_, err := tx.NamedQuery(sqlStr, pdeStateModel)
	return err
}
