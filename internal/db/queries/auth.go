package queries

import "context"

func (q *Queries) InsertStationKey(ctx context.Context, id, stationID, keyHash string) error {
	_, err := q.db.ExecContext(ctx,
		`INSERT INTO station_keys (id, station_id, key_hash) VALUES (?, ?, ?)`,
		id, stationID, keyHash)
	return err
}

func (q *Queries) GetStationIDByKey(ctx context.Context, keyHash string) (string, error) {
	var stationID string
	err := q.db.QueryRowContext(ctx,
		`SELECT station_id FROM station_keys WHERE key_hash = ? AND revoked = 0`, keyHash).Scan(&stationID)
	return stationID, err
}

func (q *Queries) RevokeAllKeys(ctx context.Context, stationID string) error {
	_, err := q.db.ExecContext(ctx,
		`UPDATE station_keys SET revoked = 1, revoked_at = NOW() WHERE station_id = ? AND revoked = 0`, stationID)
	return err
}

func (q *Queries) InsertAccessToken(ctx context.Context, id, stationID, tokenHash string) error {
	_, err := q.db.ExecContext(ctx,
		`INSERT INTO access_tokens (id, station_id, token_hash) VALUES (?, ?, ?)`,
		id, stationID, tokenHash)
	return err
}

func (q *Queries) GetStationIDByToken(ctx context.Context, tokenHash string) (string, error) {
	var stationID string
	err := q.db.QueryRowContext(ctx,
		`SELECT station_id FROM access_tokens WHERE token_hash = ? AND revoked = 0`, tokenHash).Scan(&stationID)
	return stationID, err
}

func (q *Queries) RevokeAllTokens(ctx context.Context, stationID string) error {
	_, err := q.db.ExecContext(ctx,
		`UPDATE access_tokens SET revoked = 1, revoked_at = NOW() WHERE station_id = ? AND revoked = 0`, stationID)
	return err
}

func (q *Queries) GetLatestStationKey(ctx context.Context, stationID string) (string, error) {
	var hash string
	err := q.db.QueryRowContext(ctx,
		`SELECT key_hash FROM station_keys WHERE station_id = ? AND revoked = 0 ORDER BY created_at DESC LIMIT 1`,
		stationID).Scan(&hash)
	return hash, err
}

func (q *Queries) GetLatestAccessToken(ctx context.Context, stationID string) (string, error) {
	var hash string
	err := q.db.QueryRowContext(ctx,
		`SELECT token_hash FROM access_tokens WHERE station_id = ? AND revoked = 0 ORDER BY created_at DESC LIMIT 1`,
		stationID).Scan(&hash)
	return hash, err
}
