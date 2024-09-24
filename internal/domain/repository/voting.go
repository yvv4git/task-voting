package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yvv4git/task-voting/internal/infrastructure"

	sq "github.com/Masterminds/squirrel"
)

const (
	tbVoting           = "voting"
	tbVotingInvariance = "voting_invariance"
	tbVotingResults    = "voting_results"
)

type Voting struct {
	db *pgxpool.Pool
}

func NewVoting(db *pgxpool.Pool) *Voting {
	return &Voting{
		db: db,
	}
}

type (
	ListVotingRequest struct {
		Limit  int64
		Offset int64
	}

	InvarianceScore struct {
		ID    uuid.UUID `db:"invariance_id"`
		Name  string    `db:"invariance_name"`
		Score int64     `db:"invariance_score"`
	}

	VotingItem struct {
		ID          uuid.UUID         `db:"voting_id"`
		Name        string            `db:"voting_name"`
		Description string            `db:"voting_desc"`
		CreatedAt   time.Time         `db:"created_at"`
		StartAt     time.Time         `db:"started_at"`
		EndAt       time.Time         `db:"ended_at"`
		Invariance  []InvarianceScore `db:"invarianceItem"`
	}

	ListVotingResponse struct {
		Items []VotingItem
	}
)

func (v *Voting) List(ctx context.Context, r *ListVotingRequest) (*ListVotingResponse, error) {
	var response ListVotingResponse

	listBuilder := sq.Select(
		"v.id", "v.name", "v.description", "v.created_at", "v.started_at", "v.ended_at",
		"i.id", "i.name", "count(r.id)",
	).
		From("voting v").
		LeftJoin("voting_invariance i ON v.id = i.voting_id").
		LeftJoin("voting_results r ON i.id = r.invariant_id").
		Where(sq.Eq{"v.deleted_at": nil}).
		GroupBy("v.id", "i.id").
		OrderBy("v.name", "i.name")

	if r.Limit > 0 {
		listBuilder = listBuilder.Limit(uint64(r.Limit))
	}

	if r.Offset > 0 {
		listBuilder = listBuilder.Offset(uint64(r.Offset))
	}

	stmt, args, err := listBuilder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	rows, err := v.db.Query(ctx, stmt, args...)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			votingItem      VotingItem
			invarianceID    sql.NullString
			invarianceName  sql.NullString
			invarianceScore sql.NullInt64
		)
		if err = rows.Scan(
			&votingItem.ID,
			&votingItem.Name,
			&votingItem.Description,
			&votingItem.CreatedAt,
			&votingItem.StartAt,
			&votingItem.EndAt,
			&invarianceID,
			&invarianceName,
			&invarianceScore,
		); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}

		lastItemIDx := len(response.Items) - 1
		if len(response.Items) > 0 && votingItem.ID == response.Items[lastItemIDx].ID {
			if !invarianceID.Valid {
				continue
			}
			response.Items[lastItemIDx].Invariance = append(response.Items[lastItemIDx].Invariance, InvarianceScore{
				ID:    uuid.MustParse(invarianceID.String),
				Name:  invarianceName.String,
				Score: invarianceScore.Int64,
			})
		} else {
			votingItem.Invariance = []InvarianceScore{
				{
					ID:    uuid.MustParse(invarianceID.String),
					Name:  invarianceName.String,
					Score: invarianceScore.Int64,
				},
			}
			response.Items = append(response.Items, votingItem)
		}
	}

	return &response, nil
}

type (
	CreateVotingParams struct {
		Name        string
		Description string
		StartAt     time.Time
		EndAt       time.Time
		Invariance  []string
	}

	CreateVotingResult struct {
		ID uuid.UUID
	}
)

type ResultID struct {
	ID uuid.UUID `db:"id"`
}

func (v *Voting) CreateVoting(ctx context.Context, p *CreateVotingParams) (*CreateVotingResult, error) {
	/*
		Зачем использовать defer tx.Rollback(ctx)?
		1. Гарантия отката транзакции.
		2. Обработка исключений.
		Если в процессе выполнения транзакции возникает паника (panic), defer-функция все равно будет выполнена,
		что позволяет корректно завершить транзакцию.
		3. Чистота кода.
		Использование defer делает код более чистым и понятным, так как не нужно вручную откатывать транзакцию в каждом месте,
		где может возникнуть ошибка.

		Когда вы вызываете tx.Rollback(ctx) в блоке defer, это не обязательно означает, что будет выполнен запрос к базе данных.
		Фактический запрос к БД будет выполнен только в том случае, если транзакция действительно активна и не была зафиксирована (committed).
		Если транзакция уже была зафиксирована, вызов tx.Rollback(ctx) просто ничего не сделает.
	*/
	tx, err := v.db.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	resultID, err := v.createVoting(ctx, tx, p)
	if err != nil {
		return nil, fmt.Errorf("create voting: %w", err)
	}

	if err = v.addInvariance(ctx, tx, addInvarianceParams{
		votingID:       resultID.ID,
		invarianceItem: p.Invariance,
	}); err != nil {
		return nil, fmt.Errorf("add invarianceItem: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return &CreateVotingResult{
		ID: resultID.ID,
	}, nil
}

func (v *Voting) createVoting(ctx context.Context, tx pgx.Tx, p *CreateVotingParams) (*ResultID, error) {
	updateBuilder := sq.Insert(tbVoting).
		Columns("name", "description", "started_at", "ended_at").
		Values(p.Name, p.Description, p.StartAt, p.EndAt).
		Suffix("RETURNING id")

	stmt, args, err := updateBuilder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("build statement: %w", err)
	}

	resultIDPtr, err := infrastructure.FetchRow[ResultID](ctx, tx, stmt, args...)
	if err != nil {
		return nil, err
	}

	return resultIDPtr, nil
}

type addInvarianceParams struct {
	votingID       uuid.UUID
	invarianceItem []string
}

func (v *Voting) addInvariance(ctx context.Context, tx pgx.Tx, p addInvarianceParams) error {
	updateBuilder := sq.Insert(tbVotingInvariance).
		Columns("voting_id", "name")

	for _, VotingItem := range p.invarianceItem {
		updateBuilder = updateBuilder.Values(p.votingID, VotingItem)
	}

	stmt, args, err := updateBuilder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("build statement: %w", err)
	}

	if err = infrastructure.Execute(ctx, tx, stmt, args...); err != nil {
		return err
	}

	return nil
}

type UpdateVotingParams struct {
	ID          uuid.UUID
	Name        *string
	Description *string
	StartAt     *time.Time
	EndAt       *time.Time
	Invariance  []string
}

func (v *Voting) UpdateVoting(ctx context.Context, p *UpdateVotingParams) error {
	tx, err := v.db.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	if err = v.validateUpdate(ctx, tx, p); err != nil {
		return fmt.Errorf("validate update voting: %w", err)
	}

	if err = v.updateVoting(ctx, tx, p); err != nil {
		return fmt.Errorf("update voting: %w", err)
	}

	if err = v.updateInvariance(ctx, tx, p); err != nil {
		return fmt.Errorf("update invarianceItem: %w", err)
	}

	return tx.Commit(ctx)
}

func (v *Voting) validateUpdate(ctx context.Context, tx pgx.Tx, p *UpdateVotingParams) error {
	return v.checkVotingExists(ctx, tx, p.ID)
}

func (v *Voting) updateVoting(ctx context.Context, tx pgx.Tx, p *UpdateVotingParams) error {
	updateBuilder := sq.Update(tbVoting).
		Where(sq.Eq{"id": p.ID})

	if p.Name != nil && *p.Name != "" {
		updateBuilder = updateBuilder.Set("name", p.Name)
	}

	if p.Description != nil && *p.Description != "" {
		updateBuilder = updateBuilder.Set("description", p.Description)
	}

	if p.StartAt != nil && !p.StartAt.IsZero() {
		updateBuilder = updateBuilder.Set("started_at", p.StartAt)
	}

	if p.EndAt != nil && !p.EndAt.IsZero() {
		updateBuilder = updateBuilder.Set("ended_at", p.EndAt)
	}

	stmt, args, err := updateBuilder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("build statement: %w", err)
	}

	if err = infrastructure.Execute(ctx, tx, stmt, args...); err != nil {
		return err
	}

	return nil
}

func (v *Voting) updateInvariance(ctx context.Context, tx pgx.Tx, p *UpdateVotingParams) error {
	if len(p.Invariance) == 0 {
		return nil
	}

	updateBuilder := sq.Delete(tbVotingInvariance).
		Where(sq.Eq{"voting_id": p.ID})

	stmt, args, err := updateBuilder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("build statement: %w", err)
	}

	if err = infrastructure.Execute(ctx, tx, stmt, args...); err != nil {
		return err
	}

	return v.addInvariance(ctx, tx, addInvarianceParams{
		votingID:       p.ID,
		invarianceItem: p.Invariance,
	})
}

type СheckExistsReult struct {
	Status bool `db:"status"`
}

func (v *Voting) checkVotingExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	updateBuilder := sq.Select("1").
		From(tbVoting).
		Where(sq.Eq{"id": id}).
		Suffix("FOR UPDATE").Prefix("SELECT EXISTS(").
		Suffix(") AS status")

	stmt, args, err := updateBuilder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("build statement: %w", err)
	}

	status, err := infrastructure.FetchRow[СheckExistsReult](ctx, tx, stmt, args...)
	if err != nil {
		return err
	}
	if !status.Status {
		return fmt.Errorf("voting not found")
	}

	return nil
}

type DeleteVotingParams struct {
	ID uuid.UUID
}

func (v *Voting) DeleteVoting(ctx context.Context, params *DeleteVotingParams) error {
	tx, err := v.db.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	if err = v.checkVotingExists(ctx, tx, params.ID); err != nil {
		return err
	}

	if err = v.deleteVoting(ctx, tx, params.ID); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (v *Voting) deleteVoting(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	deleteStmt, args, err := sq.Delete(tbVoting).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("build delete statement: %w", err)
	}

	if err := infrastructure.Execute(ctx, tx, deleteStmt, args...); err != nil {
		return err
	}

	// Voting invarianceItem are deleted cascadingly at the database level.
	return nil
}

type MakeChoiceParams struct {
	InvarianceID uuid.UUID
	UserID       uuid.UUID
}

func (v *Voting) MakeChoice(ctx context.Context, p *MakeChoiceParams) error {
	tx, err := v.db.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := v.validateBeforeMakeChoice(ctx, tx, p); err != nil {
		return err
	}

	if err = v.makeChoice(ctx, tx, p); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (v *Voting) validateBeforeMakeChoice(ctx context.Context, tx pgx.Tx, p *MakeChoiceParams) error {
	if err := v.isFinished(ctx, tx, p); err != nil {
		return err
	}

	return v.isUserVoted(ctx, tx, p)
}

func (v *Voting) isFinished(ctx context.Context, tx pgx.Tx, p *MakeChoiceParams) error {
	internalBuilder := sq.Select("1").
		From("voting_invariance i").
		Join("voting v ON v.id = i.voting_id").
		Where(
			sq.And{
				sq.Eq{"i.id": p.InvarianceID},
				sq.Expr("v.ended_at <= current_timestamp"),
			},
		)

	mainBuilder := internalBuilder.Prefix("SELECT EXISTS(").Suffix(") AS status")

	stmt, args, err := mainBuilder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("build statement: %w", err)
	}

	status, err := infrastructure.FetchRow[СheckExistsReult](ctx, tx, stmt, args...)
	if err != nil {
		return err
	}
	if status.Status {
		return fmt.Errorf("voting is finished")
	}

	return nil
}

func (v *Voting) isUserVoted(ctx context.Context, tx pgx.Tx, p *MakeChoiceParams) error {
	cteBlock := sq.Select("i2.id").
		From("voting_invariance i").
		Join("voting_invariance i2 ON i.voting_id = i2.voting_id").
		Where(sq.Eq{"i.id": p.InvarianceID})

	cteBlock = cteBlock.Prefix("WITH invariances AS (").Suffix(")")

	innerSelect := sq.Select("1").
		From("voting_results r").
		Join("invariances i on i.id = r.invariant_id").
		Where(sq.Eq{"r.user_id": p.UserID}).
		PrefixExpr(cteBlock)

	mainQuery := innerSelect.Prefix("SELECT EXISTS(").Suffix(") AS status")

	stmt, args, err := mainQuery.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("build statement: %w", err)
	}

	status, err := infrastructure.FetchRow[СheckExistsReult](ctx, tx, stmt, args...)
	if err != nil {
		return err
	}
	if status.Status {
		return fmt.Errorf("user already voted")
	}

	return nil
}

func (v *Voting) makeChoice(ctx context.Context, tx pgx.Tx, p *MakeChoiceParams) error {
	insertBuilder := sq.Insert(tbVotingResults).
		Columns("invariant_id", "user_id").
		Values(p.InvarianceID, p.UserID)

	stmt, args, err := insertBuilder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("build statement: %w", err)
	}

	if err = infrastructure.Execute(ctx, tx, stmt, args...); err != nil {
		return err
	}

	return nil
}
