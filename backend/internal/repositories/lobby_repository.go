package repositories

import (
	"context"
	"github.com/uptrace/bun"
	"backend/internal/models/games"
)

type LobbyRepository struct {
	db *bun.DB
}

func NewLobbyRepository(db *bun.DB) *LobbyRepository {
	return &LobbyRepository{db: db}
}

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

func (r *LobbyRepository) Lobby(ctx context.Context,
	joinCode string) (*games.Lobby, error) {
	lobby := new(games.Lobby)
	err := r.db.NewSelect().
		Model(lobby).
		Where("join_code = ?", joinCode).
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
		Where("join_code = ?", joinCode).
		Scan(ctx)
	if err != nil {
		return 0, err
	}
	return lobby.GameId, nil
}

func (r *LobbyRepository) LobbyPlayerNb(ctx context.Context,
	gameId int64) (int, error) {
	count, err := r.db.NewSelect().
		Model((*games.Lobby)(nil)).
		Where("game_id = ?", gameId).
		Count(ctx)
	if err != nil {
		return 0, err
	}
	return count, nil
}

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
