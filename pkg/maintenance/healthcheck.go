package maintenance

import "context"

type HealthCheck func(ctx context.Context) (bool, error)

type HealthChecker interface {
	HealthCheck(ctx context.Context) error
}
