package repositories

import (
	"backend/internal/models/games"
	"context"
	"database/sql"
	"errors"
	"github.com/uptrace/bun"
	"time"
)

type LobbyRepository struct {
	db *bun.DB
}

func NewLobbyRepository(db *bun.DB) *LobbyRepository {
	return &LobbyRepository{db: db}
}

// Specific functions
func (r *LobbyRepository) CreateLobby(ctx context.Context,
	lobby *games.Lobby) (*games.Lobby, error) {
	_, err := r.db.NewInsert().
		Model(lobby).
		Exec(ctx)
	if err != nil {
		return nil, err
	}
	return lobby, nil
}

func (r *LobbyRepository) JoinLobby(ctx context.Context,
	lobby *games.Lobby) error {
	_, err := r.db.NewInsert().
		Model(lobby).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

// return the actual lobby of the userid and gameid
func (r *LobbyRepository) ActualLobbyGameId(ctx context.Context,
	userId int64, gameId int64) (*games.Lobby, error) {
	lobby := new(games.Lobby)
	err := r.db.NewSelect().
		Model(lobby).
		Where("user_id = ?", userId).
		Where("game_id = ?", gameId).
		Where("left_at IS NULL").
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return lobby, nil
}

// return the actual lobby / game the user is in
// WARNING : if the user isn't in a lobby, it return nil
// so you must check if the lobby is nil after calling this function
func (r *LobbyRepository) ActualLobby(ctx context.Context,
	userId int64) (*games.Lobby, error) {
	lobby := new(games.Lobby)
	err := r.db.NewSelect().
		Model(lobby).
		Where("lo.user_id = ?", userId).
		Where("lo.left_at IS NULL").
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return lobby, nil
}

// Getters
func (r *LobbyRepository) Lobby(ctx context.Context,
	joinCode string) (*games.Lobby, error) {
	lobby := new(games.Lobby)
	err := r.db.NewSelect().
		Model(lobby).
		Where("join_code = ?", &joinCode).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return lobby, nil
}

func (r *LobbyRepository) GameIdLobby(ctx context.Context,
	joinCode string) (int64, error) {
	lobby := new(games.Lobby)
	err := r.db.NewSelect().
		Model(lobby).
		Column("game_id").
		Where("join_code = ?", &joinCode).
		Scan(ctx)
	if err != nil {
		return 0, err
	}
	return lobby.GameId, nil
}

// Get the number of players in the lobby
func (r *LobbyRepository) LobbyPlayerNb(ctx context.Context,
	gameId int64) (int, error) {
	count, err := r.db.NewSelect().
		Model((*games.Lobby)(nil)).
		Where("game_id = ?", gameId).
		Where("left_at IS NULL").
		Count(ctx)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *LobbyRepository) IsCreator(ctx context.Context, userId int64, gameId int64) (bool, error) {
	return r.db.NewSelect().
		Model((*games.Lobby)(nil)).
		Where("user_id = ?", userId).
		Where("game_id = ?", gameId).
		Where("is_creator = ?", true).
		Exists(ctx)
}

func (r *LobbyRepository) AllPlayers(ctx context.Context,
	gameId int64) ([]*games.Lobby, error) {
	var lobbies []*games.Lobby
	err := r.db.NewSelect().
		Model(&lobbies).
		Where("game_id = ?", gameId).
		Where("left_at IS NULL").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return lobbies, nil
}

func (r *LobbyRepository) NbGames(ctx context.Context,
	userId int64) (int, error) {
	count, err := r.db.NewSelect().
		Model((*games.Lobby)(nil)).
		Relation("Game", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Where("finished_at IS NOT NULL")
		}).
		Where("user_id = ?", userId).
		Where("left_at >= game.finished_at").
		Count(ctx)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *LobbyRepository) historicSubquery(userId int64) *bun.SelectQuery {
	return r.db.NewSelect().
		Model((*games.Lobby)(nil)).
		Column("lo.game_id").
		Join("JOIN game AS game ON game.id = lo.game_id").
		Where("lo.user_id = ?", userId).
		Where("lo.kicked IS FALSE").
		Where("game.finished_at IS NOT NULL").
		Where("lo.left_at IS NOT NULL").
		DistinctOn("lo.game_id").
		Order("lo.game_id", "lo.left_at DESC")
}

func (r *LobbyRepository) Historic(ctx context.Context,
	userId int64) ([]*games.HistoricGames, error) {
	var gameHistoric []*games.HistoricGames
	err := r.db.NewSelect().
		Model((*games.Lobby)(nil)).
		ColumnExpr("g.id as id").
		ColumnExpr("COALESCE(SUM(gu.base_score), 0) as xp").
		ColumnExpr("g.finished_at AS date").
		With("sub", r.historicSubquery(userId)).
		ModelTableExpr("sub AS lo").
		Join("JOIN game AS g on g.id = lo.game_id").
		Join("JOIN rounds AS r on r.game_id = g.id").
		Join("LEFT JOIN guesses AS gu ON gu.round_id = r.id AND gu.user_id = ?", userId).
		GroupExpr("g.id, g.finished_at").
		Order("g.finished_at DESC").
		Limit(20).
		Scan(ctx, &gameHistoric)
	if err != nil {
		return nil, err
	}
	return gameHistoric, nil
}

func (r *LobbyRepository) JoinCode(ctx context.Context,
	gameId int64) (*string, error) {
	lobby := new(games.Lobby)
	err := r.db.NewSelect().
		Model(lobby).
		Column("join_code").
		Where("game_id = ?", gameId).
		Where("join_code IS NOT NULL").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return lobby.JoinCode, nil
}

// Setters
func (r *LobbyRepository) SetGameId(ctx context.Context, lobby_id int64,
	game_id int64) error {
	lobby := &games.Lobby{GameId: game_id}
	_, err := r.db.NewUpdate().
		Model(lobby).
		Column("game_id").
		Where("id = ?", lobby_id).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *LobbyRepository) SetLeftAt(ctx context.Context,
	lobby_id int64) error {
	actualTime := time.Now()
	lobby := &games.Lobby{LeftAt: &actualTime}
	_, err := r.db.NewUpdate().
		Model(lobby).
		Column("left_at").
		Where("id = ?", lobby_id).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

// Kick removes a player from the lobby and remembers it: a kicked player can't join this game again (see WasKicked)
func (r *LobbyRepository) Kick(ctx context.Context,
	lobby_id int64) error {
	actualTime := time.Now()
	lobby := &games.Lobby{LeftAt: &actualTime, Kicked: true}
	_, err := r.db.NewUpdate().
		Model(lobby).
		Column("left_at", "kicked").
		Where("id = ?", lobby_id).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

// WasKicked tells whether userId has been kicked out of the game gameId
func (r *LobbyRepository) WasKicked(ctx context.Context,
	userId int64, gameId int64) (bool, error) {
	return r.db.NewSelect().
		Model((*games.Lobby)(nil)).
		Where("user_id = ?", userId).
		Where("game_id = ?", gameId).
		Where("kicked = ?", true).
		Exists(ctx)
}

func (r *LobbyRepository) SetIsCreator(ctx context.Context,
	lobby_id int64, state bool) error {
	lobby := &games.Lobby{IsCreator: state}
	_, err := r.db.NewUpdate().
		Model(lobby).
		Column("is_creator").
		Where("id = ?", lobby_id).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *LobbyRepository) SetJoinCode(ctx context.Context,
	lobby_id int64, joinCode *string) error {
	lobby := &games.Lobby{JoinCode: joinCode}
	_, err := r.db.NewUpdate().
		Model(lobby).
		Column("join_code").
		Where("id = ?", lobby_id).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

// Delete
func (r *LobbyRepository) RemoveLobby(ctx context.Context,
	lobby *games.Lobby) error {
	res, err := r.db.NewDelete().
		Model(lobby).
		WherePK().
		Exec(ctx)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("You are not in a lobby")
	}
	return nil
}
