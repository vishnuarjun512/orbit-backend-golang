package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (r *Repository) GetProfile(ctx context.Context, userID string) (*UserProfile, error) {
	const q = `
		SELECT 
			p.user_id,
			u.public_id, 
			p.display_name, 
			p.role_title, 
			p.pronouns, 
			p.bio, 
			p.avatar_url, 
			p.availability, 
			p.status_message, 
			p.updated_at
		FROM user_profiles p 
		JOIN users u USING (user_id) 
		WHERE p.user_id=$1 AND u.status='active'
	`

	p := &UserProfile{}
	err := r.db.QueryRow(ctx, q, userID).Scan(&p.UserID, &p.PublicID, &p.DisplayName, &p.RoleTitle, &p.Pronouns, &p.Bio, &p.AvatarURL, &p.Availability, &p.StatusMessage, &p.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUserNotRegistered
	}
	if err != nil {
		return nil, fmt.Errorf("Get user profile: %w", err)
	}
	return p, nil
}

func (r *Repository) GetProfileByPublicID(ctx context.Context, publicID string) (*UserProfile, error) {
	const q = `
		SELECT 
			p.user_id,
			u.public_id, 
			p.display_name, 
			p.role_title, 
			p.pronouns, 
			p.bio, 
			p.avatar_url, 
			p.availability, 
			p.status_message, 
			p.updated_at
		FROM user_profiles p 
		JOIN users u USING (user_id) 
		WHERE u.public_id=$1 AND u.status='active'
	`

	p := &UserProfile{}
	err := r.db.QueryRow(ctx, q, publicID).Scan(&p.UserID, &p.PublicID, &p.DisplayName, &p.RoleTitle, &p.Pronouns, &p.Bio, &p.AvatarURL, &p.Availability, &p.StatusMessage, &p.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUserNotRegistered
	}
	if err != nil {
		return nil, fmt.Errorf("Get user profile by public id: %w", err)
	}
	return p, nil
}

func (r *Repository) UpdateProfile(ctx context.Context, userID string, req ProfileUpdateRequest) (*UserProfile, error) {
	const q = `
		UPDATE user_profiles 
		SET
			display_name=COALESCE($2, display_name), 
			role_title=COALESCE($3, role_title),
			pronouns=COALESCE($4, pronouns), 
			bio=COALESCE($5, bio), 
			avatar_url=COALESCE($6, avatar_url),
			availability=COALESCE($7, availability), 
			status_message=COALESCE($8, status_message), 
			updated_at=NOW()
		WHERE user_id=$1 
		RETURNING 
			user_id, display_name, role_title, pronouns, bio, avatar_url, availability, status_message, updated_at`

	p := &UserProfile{}
	err := r.db.QueryRow(ctx, q, userID, req.DisplayName, req.RoleTitle, req.Pronouns, req.Bio, req.AvatarURL, req.Availability, req.StatusMessage).Scan(&p.UserID, &p.DisplayName, &p.RoleTitle, &p.Pronouns, &p.Bio, &p.AvatarURL, &p.Availability, &p.StatusMessage, &p.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUserNotRegistered
	}

	if err != nil {
		return nil, fmt.Errorf("update user profile: %w", err)
	}
	return p, nil
}

func (r *Repository) GetPreferences(ctx context.Context, userID string) (*UserPreferences, error) {
	const q = `
		SELECT 
			user_id, 
			active_workspace_id, 
			email_notifications, 
			task_updates, 
			mentions, 
			weekly_digest, 
			timezone, 
			week_starts_on, 
			theme, 
			updated_at 
		FROM user_preferences 
		WHERE user_id=$1
	`

	p := &UserPreferences{}
	err := r.db.QueryRow(ctx, q, userID).Scan(&p.UserID, &p.ActiveWorkspaceID, &p.EmailNotifications, &p.TaskUpdates, &p.Mentions, &p.WeeklyDigest, &p.Timezone, &p.WeekStartsOn, &p.Theme, &p.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUserNotRegistered
	}
	if err != nil {
		return nil, fmt.Errorf("get user preferences: %w", err)
	}
	return p, nil
}

func (r *Repository) UpdatePreferences(ctx context.Context, userID string, req PreferencesUpdateRequest) (*UserPreferences, error) {
	if req.ActiveWorkspaceID != nil && *req.ActiveWorkspaceID != "" {
		var member bool
		err := r.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM workspace_members WHERE user_id=$1 AND workspace_id=$2 AND status='active')`, userID, *req.ActiveWorkspaceID).Scan(&member)
		if err != nil {
			return nil, fmt.Errorf("validate active workspace: %w", err)
		}
		if !member {
			return nil, ErrInvalidPreferences
		}
	}
	const q = `
		UPDATE user_preferences 
		SET
			active_workspace_id=CASE WHEN $2::boolean THEN $3::uuid ELSE active_workspace_id END,
			email_notifications=COALESCE($4,email_notifications), 
			task_updates=COALESCE($5,task_updates),
			mentions=COALESCE($6,mentions), 
			weekly_digest=COALESCE($7,weekly_digest),
			timezone=COALESCE($8,timezone), 
			week_starts_on=COALESCE($9,week_starts_on), 
			theme=COALESCE($10,theme), 
			updated_at=NOW()
		WHERE user_id=$1 
		RETURNING 
			user_id, active_workspace_id, email_notifications, task_updates, mentions, weekly_digest, timezone, week_starts_on, theme, updated_at
	`

	clearWorkspace := req.ActiveWorkspaceID != nil && *req.ActiveWorkspaceID == ""
	var workspace any
	if req.ActiveWorkspaceID != nil && !clearWorkspace {
		workspace = *req.ActiveWorkspaceID
	}
	p := &UserPreferences{}
	err := r.db.QueryRow(ctx, q, userID, req.ActiveWorkspaceID != nil, workspace, req.EmailNotifications, req.TaskUpdates, req.Mentions, req.WeeklyDigest, req.Timezone, req.WeekStartsOn, req.Theme).Scan(&p.UserID, &p.ActiveWorkspaceID, &p.EmailNotifications, &p.TaskUpdates, &p.Mentions, &p.WeeklyDigest, &p.Timezone, &p.WeekStartsOn, &p.Theme, &p.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUserNotRegistered
	}
	if err != nil {
		return nil, fmt.Errorf("update user preferences: %w", err)
	}
	return p, nil
}
