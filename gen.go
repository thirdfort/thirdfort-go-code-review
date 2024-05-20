package main

import _ "go.uber.org/mock/mockgen/model"

//go:generate mockgen -destination=./internal/service/mock_service.go -package=service github.com/thirdfort/thirdfort-go-code-review/internal/service MainService
//go:generate mockgen -destination=./internal/service/mock_ops.go -package=service -source=./internal/service/task_ops.go github.com/thirdfort/thirdfort-go-code-review/internal/service MockOp[DataType]
//go:generate mockgen -destination=./internal/clients/pa/mock_pa.go -package=pa github.com/thirdfort/thirdfort-go-code-review/internal/clients/pa Client
//go:generate mockgen -destination=./internal/repositories/mock_repository.go -package=repositories github.com/thirdfort/thirdfort-go-code-review/internal/repositories DataStore
//go:generate mockgen -destination=./internal/observability/mock_observability.go -package=observability github.com/thirdfort/thirdfort-go-code-review/internal/observability MetricsInterface
//go:generate mockgen -destination=./internal/models/mock_models.go -package=models -source=./internal/models/address.go github.com/thirdfort/thirdfort-go-code-review/internal/models Item[DataType]
//go:generate mockgen -destination=./internal/cache/mock_cache.go -package=cache github.com/thirdfort/thirdfort-go-code-review/internal/cache Cache
//go:generate mockgen -destination=./internal/repositories/mock_crud.go -package=repositories -source=./internal/repositories/crud.go github.com/thirdfort/thirdfort-go-code-review/internal/repositories MockOp[DataType]
