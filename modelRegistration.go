package DatabaseFlow

import (
	"context"
	"sync"
)

var (
	models     []interface{}
	modelsSet  bool
	modelsOnce sync.Once
)

// RegisterModels allows external projects to register their models for migration.
func RegisterModels(ctx context.Context, newModels ...interface{}) {
	ctx, span := tracer.Start(ctx, "RegisterModels")
	defer span.End()

	// Check if models have already been set
	if modelsSet {
		logger.Error(ctx, "Models have already been set")
		return
	}

	// Set the models only once
	modelsOnce.Do(func() {
		models = newModels
		modelsSet = true
	})
}
