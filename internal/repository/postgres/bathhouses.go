package postgres

import (
	"context"
	"database/sql"
	"github.com/jackc/pgx/v5"

	"github.com/calyrexx/QuietGrooveBackend/internal/entities"
	"github.com/calyrexx/QuietGrooveBackend/internal/pkg/errorspkg"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BathhousesRepo struct {
	pool *pgxpool.Pool
}

func NewBathhousesRepo(pool *pgxpool.Pool) *BathhousesRepo {
	return &BathhousesRepo{pool: pool}
}

func (r *BathhousesRepo) GetAll(ctx context.Context) ([]entities.Bathhouse, error) {
	const method = "BathhousesRepo.GetAll"

	rows, err := r.pool.Query(ctx, `
		SELECT
			bt.id,
			bt.name,
			bt.default_price,
			bt.created_at,
			
			hb.id,
			hb.house_id,
			hb.price,
			hb.description,
			hb.images,
			
			bfo.id,
			bfo.name,
			bfo.image,
			bfo.description,
			bfo.price
		FROM bathhouses_types bt
		LEFT JOIN house_bathhouses hb ON hb.bathhouse_type_id = bt.id
		LEFT JOIN bathhouse_fill_options bfo ON bfo.bathhouse_type_id = bt.id
		ORDER BY bt.id, hb.id, bfo.id
	`)
	if err != nil {
		return nil, errorspkg.NewErrRepoFailed("pool.Query", method, err)
	}
	defer rows.Close()

	bathhouseMap := make(map[int]*entities.Bathhouse)
	houseSet := make(map[[2]int]struct{}) // [TypeID, HouseBathhouseID]
	fillSet := make(map[[2]int]struct{})  // [TypeID, FillOptionID]

	for rows.Next() {
		var (
			btID, hbID, houseID, bfoID       sql.NullInt64
			btName, hbDesc, bfoName, bfoDesc sql.NullString
			btCreated, bfoImg                sql.NullString
			btDefPrice, hbPrice, bfoPrice    sql.NullInt32
			hbImgs                           []string
		)
		err = rows.Scan(&btID, &btName, &btDefPrice, &btCreated,
			&hbID, &houseID, &hbPrice, &hbDesc, &hbImgs,
			&bfoID, &bfoName, &bfoImg, &bfoDesc, &bfoPrice,
		)
		if err != nil {
			return nil, errorspkg.NewErrRepoFailed("rows.Scan", method, err)
		}
		typeID := int(btID.Int64)
		bh, ok := bathhouseMap[typeID]
		if !ok {
			bh = &entities.Bathhouse{
				ID:           typeID,
				Name:         btName.String,
				DefaultPrice: int(btDefPrice.Int32),
				Houses:       []entities.HouseBathhouse{},
				FillOptions:  []entities.BathhouseFill{},
			}
			bathhouseMap[typeID] = bh
		}
		if hbID.Valid {
			houseKey := [2]int{typeID, int(hbID.Int64)}
			if _, exists := houseSet[houseKey]; !exists {
				bh.Houses = append(bh.Houses, entities.HouseBathhouse{
					ID:          int(hbID.Int64),
					HouseID:     int(houseID.Int64),
					Price:       int(hbPrice.Int32),
					Description: hbDesc.String,
					Images:      hbImgs,
				})
				houseSet[houseKey] = struct{}{}
			}
		}
		// Add FillOption
		if bfoID.Valid {
			fillKey := [2]int{typeID, int(bfoID.Int64)}
			if _, exists := fillSet[fillKey]; !exists {
				bh.FillOptions = append(bh.FillOptions, entities.BathhouseFill{
					ID:          int(bfoID.Int64),
					Name:        bfoName.String,
					Image:       bfoImg.String,
					Description: bfoDesc.String,
					Price:       int(bfoPrice.Int32),
				})
				fillSet[fillKey] = struct{}{}
			}
		}
	}
	res := make([]entities.Bathhouse, 0, len(bathhouseMap))
	for _, bh := range bathhouseMap {
		res = append(res, *bh)
	}
	return res, nil
}

func (r *BathhousesRepo) Add(ctx context.Context, bhs []entities.Bathhouse) error {
	const method = "BathhousesRepo.Add"
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return errorspkg.NewErrRepoFailed("Begin", method, err)
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		_ = tx.Rollback(ctx)
	}(tx, ctx)

	for _, bh := range bhs {
		var id int
		err = tx.QueryRow(ctx, `
			INSERT INTO bathhouses_types (name, default_price) VALUES ($1, $2)
			RETURNING id
		`, bh.Name, bh.DefaultPrice).Scan(&id)
		if err != nil {
			return errorspkg.NewErrRepoFailed("Insert bathhouse type", method, err)
		}
		for _, h := range bh.Houses {
			_, err = tx.Exec(ctx, `
				INSERT INTO house_bathhouses (
					house_id,
				    bathhouse_type_id,                          
				    price,
				    description,
				    images
				)
				VALUES (
					$1,
				    $2,
				    $3,
				    $4, 
				    $5
				)
			`,
				h.HouseID,
				id,
				h.Price,
				h.Description,
				h.Images,
			)
			if err != nil {
				return errorspkg.NewErrRepoFailed("Insert house_bathhouses", method, err)
			}
		}
		for _, f := range bh.FillOptions {
			_, err = tx.Exec(ctx, `
				INSERT INTO bathhouse_fill_options (
					bathhouse_type_id,
				    name,
				    image,
				    description,
				    price
				)
				VALUES (
				    $1,
				    $2,
				    $3,
				    $4,
				    $5
				)
			`,
				id,
				f.Name,
				f.Image,
				f.Description,
				f.Price,
			)
			if err != nil {
				return errorspkg.NewErrRepoFailed("Insert fill_option", method, err)
			}
		}
	}

	return tx.Commit(ctx)
}

func (r *BathhousesRepo) Update(ctx context.Context, bh entities.Bathhouse) error {
	const method = "BathhousesRepo.Update"

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return errorspkg.NewErrRepoFailed("Begin", method, err)
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		_ = tx.Rollback(ctx)
	}(tx, ctx)

	_, err = tx.Exec(ctx, `
		UPDATE bathhouses_types SET name=$1, default_price=$2 WHERE id=$3
	`, bh.Name, bh.DefaultPrice, bh.ID)
	if err != nil {
		return errorspkg.NewErrRepoFailed("Update bathhouse_type", method, err)
	}

	_, err = tx.Exec(ctx, `DELETE FROM house_bathhouses WHERE bathhouse_type_id=$1`, bh.ID)
	if err != nil {
		return errorspkg.NewErrRepoFailed("Clear house_bathhouses", method, err)
	}
	for _, h := range bh.Houses {
		_, err = tx.Exec(ctx, `
			INSERT INTO house_bathhouses (house_id, bathhouse_type_id, price, description, images)
			VALUES ($1, $2, $3, $4, $5)
		`, h.HouseID, bh.ID, h.Price, h.Description, h.Images)
		if err != nil {
			return errorspkg.NewErrRepoFailed("Insert house_bathhouses (upd)", method, err)
		}
	}
	_, err = tx.Exec(ctx, `DELETE FROM bathhouse_fill_options WHERE bathhouse_type_id=$1`, bh.ID)
	if err != nil {
		return errorspkg.NewErrRepoFailed("Clear fill_options", method, err)
	}
	for _, f := range bh.FillOptions {
		_, err = tx.Exec(ctx, `
			INSERT INTO bathhouse_fill_options (bathhouse_type_id, name, image, description, price)
			VALUES ($1, $2, $3, $4, $5)
		`, bh.ID, f.Name, f.Image, f.Description, f.Price)
		if err != nil {
			return errorspkg.NewErrRepoFailed("Insert fill_option (upd)", method, err)
		}
	}

	return tx.Commit(ctx)
}

func (r *BathhousesRepo) Delete(ctx context.Context, id int) error {
	const method = "BathhousesRepo.Delete"
	_, err := r.pool.Exec(ctx, `DELETE FROM bathhouses_types WHERE id = $1`, id)
	if err != nil {
		return errorspkg.NewErrRepoFailed("Delete bathhouse_type", method, err)
	}
	return nil
}
