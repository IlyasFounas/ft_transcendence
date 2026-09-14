package repositories

import (
	"time"
	"errors"
	"context"
	"database/sql"
	"github.com/uptrace/bun"
	"backend/internal/models/games"
)

type LobbyRepository struct {
	db *bun.DB
}

func NewLobbyRepository(db *bun.DB) *LobbyRepository {
	return &LobbyRepository{db: db}
}

// Specific functions
func (r *LobbyRepository) CreateLobby(ctx context.Context,
	lobby *games.Lobby) (*games.Lobby, error){
	_, err := r.db.NewInsert().
		Model(lobby).
		Exec(ctx)
	if err != nil {
		return nil, err
	}
	return lobby, nil
}

func (r *LobbyRepository) JoinLobby(ctx context.Context,
	lobby *games.Lobby) error{
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
		// Relation("Game", func(q *bun.SelectQuery) *bun.SelectQuery {
		// 	return q.Where("finished = false")
		// }).
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

func (r *LobbyRepository) IsCreator(ctx context.Context,
	userId int64, gameId int64) (bool, error) {
	lobby := new(games.Lobby)
	err := r.db.NewSelect().
		Model(lobby).
		Column("is_creator").
		Where("user_id = ?", userId).
		Where("game_id = ?", gameId).
		Scan(ctx)
	if err != nil {
		return false, err
	}
	return lobby.IsCreator, nil
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
	userId int64) (int, error){
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

func (r *LobbyRepository) Historic(ctx context.Context,
	userId int64) ([]*games.Game, error){
	var games []*games.Game
	err := r.db.NewSelect().
		Model(games).
		Column("id", "").
		Relation("Game", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Where("finished_at IS NOT NULL")
		}).
		Where("user_id = ?", userId).
		Where("left_at >= game.finished_at").
		Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("You don't have any games")
		}
		return nil, err
	}
	return games, nil
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