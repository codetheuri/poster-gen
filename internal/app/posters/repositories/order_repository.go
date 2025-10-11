package repositories

import (
    "context"
    "github.com/codetheuri/poster-gen/internal/app/posters/models"
    "github.com/codetheuri/poster-gen/pkg/logger"
    "gorm.io/gorm"
)

// OrderSubRepository interface for order operations
type OrderSubRepository interface {
    CreateOrder(ctx context.Context, order *models.Order) error
    GetOrderByID(ctx context.Context, id uint) (*models.Order, error)
    GetOrderByNumber(ctx context.Context, orderNumber string) (*models.Order, error)
    UpdateOrder(ctx context.Context, order *models.Order) error
    DeleteOrder(ctx context.Context, id uint) error
}

type orderSubRepository struct {
    db  *gorm.DB
    log logger.Logger
}

// NewOrderSubRepository constructor
func NewOrderSubRepository(db *gorm.DB, log logger.Logger) OrderSubRepository {
    return &orderSubRepository{
        db:  db,
        log: log,
    }
}

func (r *orderSubRepository) CreateOrder(ctx context.Context, order *models.Order) error {
    r.log.Info("Creating order", "order_number", order.OrderNumber, "user_id", order.UserID)
    return r.db.WithContext(ctx).Create(order).Error
}

func (r *orderSubRepository) GetOrderByID(ctx context.Context, id uint) (*models.Order, error) {
    r.log.Info("Getting order by ID", "id", id)
    var order models.Order
    err := r.db.WithContext(ctx).First(&order, id).Error
    if err != nil {
        r.log.Error("Failed to get order by ID", err, "id", id)
        return nil, err
    }
    return &order, nil
}

func (r *orderSubRepository) GetOrderByNumber(ctx context.Context, orderNumber string) (*models.Order, error) {
    r.log.Info("Getting order by number", "order_number", orderNumber)
    var order models.Order
    err := r.db.WithContext(ctx).Where("order_number = ?", orderNumber).First(&order).Error
    if err != nil {
        r.log.Error("Failed to get order by number", err, "order_number", orderNumber)
        return nil, err
    }
    return &order, nil
}

func (r *orderSubRepository) UpdateOrder(ctx context.Context, order *models.Order) error {
    r.log.Info("Updating order", "id", order.ID)
    return r.db.WithContext(ctx).Save(order).Error
}

func (r *orderSubRepository) DeleteOrder(ctx context.Context, id uint) error {
    r.log.Info("Deleting order", "id", id)
    return r.db.WithContext(ctx).Delete(&models.Order{}, id).Error
}