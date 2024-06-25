package cards

import (
	"github.com/KBingsoo/entities/pkg/models"
	"github.com/literalog/go-wise/wise"
)

type Repository wise.MongoRepository[models.Card]
