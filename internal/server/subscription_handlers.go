package server

import (
	"librarease/internal/usecase"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Subscription struct {
	ID           string      `json:"id"`
	UserID       string      `json:"user_id"`
	MembershipID string      `json:"membership_id"`
	CreatedAt    string      `json:"created_at"`
	UpdatedAt    string      `json:"updated_at"`
	DeletedAt    *string     `json:"deleted_at,omitempty"`
	User         *User       `json:"user,omitempty"`
	Membership   *Membership `json:"membership,omitempty"`

	// Granfathering the membership
	ExpiresAt       string `json:"expires_at"`
	FinePerDay      int    `json:"fine_per_day"`
	LoanPeriod      int    `json:"loan_period"`
	ActiveLoanLimit int    `json:"active_loan_limit"`
}

type ListSubscriptionsRequest struct {
	Skip         int    `query:"skip"`
	Limit        int    `query:"limit" validate:"required,gte=1,lte=100"`
	UserID       string `query:"user_id" validate:"omitempty,uuid"`
	MembershipID string `query:"membership_id" validate:"omitempty,uuid"`
	LibraryID    string `query:"library_id" validate:"omitempty,uuid"`
}

func (s *Server) ListSubscriptions(ctx echo.Context) error {
	var req ListSubscriptionsRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(400, map[string]string{"error": err.Error()})
	}
	if err := s.validator.Struct(req); err != nil {
		return ctx.JSON(422, map[string]string{"error": err.Error()})
	}

	subs, _, err := s.server.ListSubscriptions(ctx.Request().Context(), usecase.ListSubscriptionsOption{
		Skip:         req.Skip,
		Limit:        req.Limit,
		UserID:       req.UserID,
		MembershipID: req.MembershipID,
		LibraryID:    req.LibraryID,
	})
	if err != nil {
		return ctx.JSON(500, map[string]string{"error": err.Error()})
	}
	list := make([]Subscription, 0, len(subs))

	for _, sub := range subs {
		var d *string
		if sub.DeletedAt != nil {
			tmp := sub.DeletedAt.String()
			d = &tmp
		}
		m := Subscription{
			ID:              sub.ID.String(),
			UserID:          sub.UserID.String(),
			MembershipID:    sub.MembershipID.String(),
			CreatedAt:       sub.CreatedAt.String(),
			UpdatedAt:       sub.UpdatedAt.String(),
			DeletedAt:       d,
			ExpiresAt:       sub.ExpiresAt.String(),
			FinePerDay:      sub.FinePerDay,
			LoanPeriod:      sub.LoanPeriod,
			ActiveLoanLimit: sub.ActiveLoanLimit,
		}
		if sub.User != nil {
			m.User = &User{
				ID:   sub.User.ID.String(),
				Name: sub.User.Name,
			}
		}
		if sub.Membership != nil {
			m.Membership = &Membership{
				ID:        sub.Membership.ID.String(),
				Name:      sub.Membership.Name,
				LibraryID: sub.Membership.LibraryID.String(),
			}

			if lib := sub.Membership.Library; lib != nil {
				m.Membership.Library = &Library{
					ID:   lib.ID.String(),
					Name: lib.Name,
				}
			}
		}
		list = append(list, m)
	}

	return ctx.JSON(200, list)
}

type GetSubscriptionByIDRequest struct {
	ID string `param:"id" validate:"required,uuid"`
}

func (s *Server) GetSubscriptionByID(ctx echo.Context) error {
	var req GetSubscriptionByIDRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(400, map[string]string{"error": err.Error()})
	}
	if err := s.validator.Struct(req); err != nil {
		return ctx.JSON(422, map[string]string{"error": err.Error()})
	}

	id, _ := uuid.Parse(req.ID)

	sub, err := s.server.GetSubscriptionByID(ctx.Request().Context(), id)
	if err != nil {
		return ctx.JSON(500, map[string]string{"error": err.Error()})
	}

	var d *string
	if sub.DeletedAt != nil {
		tmp := sub.DeletedAt.String()
		d = &tmp
	}
	m := Subscription{
		ID:              sub.ID.String(),
		UserID:          sub.UserID.String(),
		MembershipID:    sub.MembershipID.String(),
		CreatedAt:       sub.CreatedAt.String(),
		UpdatedAt:       sub.UpdatedAt.String(),
		DeletedAt:       d,
		ExpiresAt:       sub.ExpiresAt.String(),
		FinePerDay:      sub.FinePerDay,
		LoanPeriod:      sub.LoanPeriod,
		ActiveLoanLimit: sub.ActiveLoanLimit,
	}
	if sub.User != nil {
		m.User = &User{
			ID:   sub.User.ID.String(),
			Name: sub.User.Name,
		}
	}
	if sub.Membership != nil {
		m.Membership = &Membership{
			ID:              sub.Membership.ID.String(),
			Name:            sub.Membership.Name,
			LibraryID:       sub.Membership.LibraryID.String(),
			Duration:        sub.Membership.Duration,
			ActiveLoanLimit: sub.Membership.ActiveLoanLimit,
			LoanPeriod:      sub.Membership.LoanPeriod,
			FinePerDay:      sub.Membership.FinePerDay,
			CreatedAt:       sub.Membership.CreatedAt.String(),
			UpdatedAt:       sub.Membership.UpdatedAt.String(),
		}

		if lib := sub.Membership.Library; lib != nil {
			m.Membership.Library = &Library{
				ID:   lib.ID.String(),
				Name: lib.Name,
			}
		}
	}

	return ctx.JSON(200, m)
}

type CreateSubscriptionRequest struct {
	UserID       string `json:"user_id" validate:"required,uuid"`
	MembershipID string `json:"membership_id" validate:"required,uuid"`
}

func (s *Server) CreateSubscription(ctx echo.Context) error {
	var req CreateSubscriptionRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(400, map[string]string{"error": err.Error()})
	}
	if err := s.validator.Struct(req); err != nil {
		return ctx.JSON(422, map[string]string{"error": err.Error()})
	}

	userID, _ := uuid.Parse(req.UserID)
	membershipID, _ := uuid.Parse(req.MembershipID)

	id, err := s.server.CreateSubscription(ctx.Request().Context(), usecase.Subscription{
		UserID:       userID,
		MembershipID: membershipID,
	})
	if err != nil {
		return ctx.JSON(500, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(200, id)
}

type UpdateSubscriptionRequest struct {
	ID              string `param:"id" validate:"required,uuid"`
	UserID          string `json:"user_id" validate:"omitempty,uuid"`
	MembershipID    string `json:"membership_id" validate:"omitempty,uuid"`
	ExpiresAt       string `json:"expires_at" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	FinePerDay      int    `json:"fine_per_day" validate:"omitempty,number"`
	LoanPeriod      int    `json:"loan_period" validate:"omitempty,number"`
	ActiveLoanLimit int    `json:"active_loan_limit" validate:"omitempty,number"`
}

func (s *Server) UpdateSubscription(ctx echo.Context) error {
	var req UpdateSubscriptionRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(400, map[string]string{"error": err.Error()})
	}
	if err := s.validator.Struct(req); err != nil {
		return ctx.JSON(422, map[string]string{"error": err.Error()})
	}

	id, _ := uuid.Parse(req.ID)
	userID, _ := uuid.Parse(req.UserID)
	membershipID, _ := uuid.Parse(req.MembershipID)

	var (
		exp time.Time
		err error
	)
	if req.ExpiresAt != "" {
		exp, err = time.Parse(time.RFC3339, req.ExpiresAt)
		if err != nil {
			return ctx.JSON(422, map[string]string{"error": "invalid expires_at"})
		}
	}

	sub, err := s.server.UpdateSubscription(ctx.Request().Context(), usecase.Subscription{
		ID:              id,
		UserID:          userID,
		MembershipID:    membershipID,
		ExpiresAt:       exp,
		FinePerDay:      req.FinePerDay,
		LoanPeriod:      req.LoanPeriod,
		ActiveLoanLimit: req.ActiveLoanLimit,
	})
	if err != nil {
		return ctx.JSON(500, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(200, Subscription{
		ID:              sub.ID.String(),
		UserID:          sub.UserID.String(),
		MembershipID:    sub.MembershipID.String(),
		CreatedAt:       sub.CreatedAt.String(),
		UpdatedAt:       sub.UpdatedAt.String(),
		ExpiresAt:       sub.ExpiresAt.String(),
		FinePerDay:      sub.FinePerDay,
		LoanPeriod:      sub.LoanPeriod,
		ActiveLoanLimit: sub.ActiveLoanLimit,
	})
}
