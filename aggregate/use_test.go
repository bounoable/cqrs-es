package aggregate_test

// func TestUse(t *testing.T) {
// 	cases := map[string]struct {
// 		AggregateType cqrs.AggregateType
// 		AggregateID   uuid.UUID
// 		UseFn         func(context.Context, cqrs.Aggregate) error
// 		ExpectedErr   error
// 	}{
// 		"aggregate not found": {

// 		},
// 	}

// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	ctx := context.Background()
// 	repo := mock_cqrs.NewMockAggregateRepository(ctrl)
// 	id := uuid.New()

// 	err := aggregate.Use(ctx, repo, "testagg", id, func(ctx context.Context, agg cqrs.Aggregate) error {

// 	})
// }
