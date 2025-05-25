package postgres

import (
	"context"
	"database/sql"

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
			b.id, b.house_id, b.name, b.price, b.description, b.images,
			f.id, f.name, f.image, f.description, f.price
		FROM bathhouses b
		LEFT JOIN bathhouse_fill_options f ON f.bathhouse_id = b.id
		ORDER BY b.id, f.id
	`)
	if err != nil {
		return nil, errorspkg.NewErrRepoFailed("pool.Query", method, err)
	}
	defer rows.Close()

	bathhouseMap := make(map[int]*entities.Bathhouse)

	for rows.Next() {
		var (
			bathhouseID, houseID, fillID                   sql.NullInt64
			name, description, fillName, fillImg, fillDesc sql.NullString
			price, fillPrice                               sql.NullFloat64
			images                                         []string
		)
		err = rows.Scan(&bathhouseID, &houseID, &name, &price, &description, &images,
			&fillID, &fillName, &fillImg, &fillDesc, &fillPrice,
		)
		if err != nil {
			return nil, errorspkg.NewErrRepoFailed("rows.Scan", method, err)
		}

		bh, exists := bathhouseMap[int(bathhouseID.Int64)]
		if !exists {
			bh = &entities.Bathhouse{
				ID:          int(bathhouseID.Int64),
				HouseID:     int(houseID.Int64),
				Name:        name.String,
				Price:       int(price.Float64),
				Description: description.String,
				Images:      images,
				FillOptions: []entities.BathhouseFillOption{},
			}
			bathhouseMap[bh.ID] = bh
		}

		if fillID.Valid {
			bh.FillOptions = append(bh.FillOptions, entities.BathhouseFillOption{
				ID:          int(fillID.Int64),
				BathhouseID: int(bathhouseID.Int64),
				Name:        fillName.String,
				Image:       fillImg.String,
				Description: fillDesc.String,
				Price:       int(fillPrice.Float64),
			})
		}
	}

	result := make([]entities.Bathhouse, 0, len(bathhouseMap))
	for _, bh := range bathhouseMap {
		result = append(result, *bh)
	}
	return result, nil
}

func (r *BathhousesRepo) GetByHouse(ctx context.Context, houseID int) ([]entities.Bathhouse, error) {
	const method = "BathhousesRepo.GetByHouse"

	rows, err := r.pool.Query(ctx, `
		SELECT
			b.id, b.house_id, b.name, b.price, b.description, b.images,
			f.id, f.name, f.image, f.description, f.price
		FROM bathhouses b
		LEFT JOIN bathhouse_fill_options f ON f.bathhouse_id = b.id
		WHERE b.house_id = $1
		ORDER BY b.id, f.id
	`, houseID)
	if err != nil {
		return nil, errorspkg.NewErrRepoFailed("pool.Query", method, err)
	}
	defer rows.Close()

	bathhouseMap := make(map[int]*entities.Bathhouse)

	for rows.Next() {
		var (
			bathhouseID, houseIDDb, fillID                 sql.NullInt64
			name, description, fillName, fillImg, fillDesc sql.NullString
			price, fillPrice                               sql.NullFloat64
			images                                         []string
		)
		err = rows.Scan(&bathhouseID, &houseIDDb, &name, &price, &description, &images,
			&fillID, &fillName, &fillImg, &fillDesc, &fillPrice,
		)
		if err != nil {
			return nil, errorspkg.NewErrRepoFailed("rows.Scan", method, err)
		}

		bh, exists := bathhouseMap[int(bathhouseID.Int64)]
		if !exists {
			bh = &entities.Bathhouse{
				ID:          int(bathhouseID.Int64),
				HouseID:     int(houseIDDb.Int64),
				Name:        name.String,
				Price:       int(price.Float64),
				Description: description.String,
				Images:      images,
				FillOptions: []entities.BathhouseFillOption{},
			}
			bathhouseMap[bh.ID] = bh
		}

		if fillID.Valid {
			bh.FillOptions = append(bh.FillOptions, entities.BathhouseFillOption{
				ID:          int(fillID.Int64),
				BathhouseID: int(bathhouseID.Int64),
				Name:        fillName.String,
				Image:       fillImg.String,
				Description: fillDesc.String,
				Price:       int(fillPrice.Float64),
			})
		}
	}

	result := make([]entities.Bathhouse, 0, len(bathhouseMap))
	for _, bh := range bathhouseMap {
		result = append(result, *bh)
	}
	return result, nil
}

func (r *BathhousesRepo) Add(ctx context.Context, bathhouses []entities.Bathhouse) error {
	const method = "BathhousesRepo.Add"

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return errorspkg.NewErrRepoFailed("Begin", method, err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	for _, bh := range bathhouses {
		var id int
		err = tx.QueryRow(ctx, `
			INSERT INTO bathhouses (house_id, name, price, description, images)
			VALUES ($1, $2, $3, $4, $5) RETURNING id
		`, bh.HouseID, bh.Name, bh.Price, bh.Description, bh.Images).Scan(&id)
		if err != nil {
			return errorspkg.NewErrRepoFailed("Insert bathhouse", method, err)
		}

		for _, f := range bh.FillOptions {
			_, err = tx.Exec(ctx, `
				INSERT INTO bathhouse_fill_options (bathhouse_id, name, image, description, price)
				VALUES ($1, $2, $3, $4, $5)
			`, id, f.Name, f.Image, f.Description, f.Price)
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
	defer func() { _ = tx.Rollback(ctx) }()

	_, err = tx.Exec(ctx, `
		UPDATE bathhouses SET name=$1, price=$2, description=$3, images=$4 WHERE id=$5
	`, bh.Name, bh.Price, bh.Description, bh.Images, bh.ID)
	if err != nil {
		return errorspkg.NewErrRepoFailed("Update bathhouse", method, err)
	}

	_, err = tx.Exec(ctx, `DELETE FROM bathhouse_fill_options WHERE bathhouse_id=$1`, bh.ID)
	if err != nil {
		return errorspkg.NewErrRepoFailed("Clear fill_options", method, err)
	}
	for _, f := range bh.FillOptions {
		_, err = tx.Exec(ctx, `
			INSERT INTO bathhouse_fill_options (bathhouse_id, name, image, description, price)
			VALUES ($1, $2, $3, $4, $5)
		`, bh.ID, f.Name, f.Image, f.Description, f.Price)
		if err != nil {
			return errorspkg.NewErrRepoFailed("Insert fill_option", method, err)
		}
	}

	return tx.Commit(ctx)
}

func (r *BathhousesRepo) Delete(ctx context.Context, id int) error {
	const method = "BathhousesRepo.Delete"
	_, err := r.pool.Exec(ctx, `DELETE FROM bathhouses WHERE id = $1`, id)
	if err != nil {
		return errorspkg.NewErrRepoFailed("Delete bathhouse", method, err)
	}
	return nil
}
