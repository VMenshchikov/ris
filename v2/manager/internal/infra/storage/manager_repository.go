package storage

import (
	"context"
	"fmt"
	"hash_manager/internal/domain/model"
	"hash_manager/internal/domain/repo"
	"time"

	"github.com/ztrue/tracerr"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	_ repo.ManagerRepository = (*ManagerRepository)(nil)
)

type ManagerRepository struct {
	db *mongo.Database
}

func New(database *mongo.Database) *ManagerRepository {
	return &ManagerRepository{
		db: database,
	}
}

func (m *ManagerRepository) AddOrder(targetHash [16]byte, maxLen uint, timeout time.Time, blockSize uint) (uint64, error) {
	id, err := m.getAndIncrementSeq()
	if err != nil {
		return 0, tracerr.Wrap(err)
	}

	order := &model.OrderInfo{
		Id:         id,
		Status:     model.IN_PROGRESS,
		TargetHash: targetHash,
		MaxLen:     maxLen,
		Timeout:    timeout,
		BlockSize:  int64(blockSize),
		Results:    []string{},
		CreatedAt:  time.Now(),
	}

	_, err = m.db.Collection("orders").InsertOne(context.Background(), order)
	if err != nil {
		return 0, tracerr.Wrap(err)
	}

	return id, nil
}

func (m *ManagerRepository) UpdateOrder(order *model.OrderInfo) error {
	filter := bson.M{"id": order.Id}
	update := bson.M{
		"$set": bson.M{
			"status":      order.Status,
			"target_hash": order.TargetHash,
			"max_len":     order.MaxLen,
			"timeout":     order.Timeout,
			"block_size":  order.BlockSize,
			"results":     order.Results,
		},
	}

	_, err := m.db.Collection("orders").UpdateOne(context.Background(), filter, update)
	if err != nil {
		return tracerr.Wrap(err)
	}

	return nil
}

func (m *ManagerRepository) FindOrder(id uint64) (*model.OrderInfo, error) {
	filter := bson.M{"id": id}

	var order model.OrderInfo
	if err := m.db.Collection("orders").FindOne(context.Background(), filter).Decode(&order); err != nil {
		return nil, tracerr.Wrap(err)
	}

	return &order, nil
}

func (m *ManagerRepository) FindOrderForExecution() (*model.OrderInfo, error) {
	filter := bson.M{
		"status": model.IN_PROGRESS,
	}

	opts := options.FindOne().
		SetSort(bson.D{{Key: "created_at", Value: 1}}) // находим самый ранний

	var order model.OrderInfo
	err := m.db.Collection("orders").FindOne(context.Background(), filter, opts).Decode(&order)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	return &order, nil
}

func (m *ManagerRepository) CloseTimeoutOrders() error {
	now := time.Now()

	filter := bson.M{
		"timeout": bson.M{"$lt": now},
		"status":  bson.M{"$eq": model.IN_PROGRESS},
	}

	update := bson.M{
		"$set": bson.M{
			"status": model.ERROR,
		},
	}

	_, err := m.db.Collection("orders").UpdateMany(context.Background(), filter, update)
	if err != nil {
		return tracerr.Wrap(err)
	}
	return nil
}

func (m *ManagerRepository) AddTask(task model.TaskInfo) error {
	_, err := m.db.Collection("tasks").InsertOne(context.Background(), task)
	if err != nil {
		return tracerr.Wrap(err)
	}
	return nil
}

func (m *ManagerRepository) UpdateTask(task model.TaskInfo) error {
	// Создаем фильтр для поиска задачи по ее order_id и block_number
	filter := bson.M{
		"order_id":     task.OrderId,
		"block_number": task.BlockNumber,
	}

	// Создаем обновление для задачи, обновляем только статус и результаты
	update := bson.M{
		"$set": bson.M{
			"status":     task.Status,
			"updated_at": time.Now(),
		},
	}

	// Выполняем обновление в коллекции
	_, err := m.db.Collection("tasks").UpdateOne(context.Background(), filter, update)
	if err != nil {
		return tracerr.Wrap(err)
	}

	return nil
}

func (m *ManagerRepository) FindTask(orderId, blockNumber uint64) (*model.TaskInfo, error) {
	filter := bson.M{
		"order_id":     orderId,
		"block_number": blockNumber,
	}

	var task model.TaskInfo
	err := m.db.Collection("tasks").FindOne(context.Background(), filter).Decode(&task)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	return &task, nil
}

func (m *ManagerRepository) FindTasksByOrderId(orderId uint64) ([]model.TaskInfo, error) {
	filter := bson.M{"order_id": orderId}

	cursor, err := m.db.Collection("tasks").Find(context.Background(), filter)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	defer cursor.Close(context.Background())

	var tasks []model.TaskInfo
	if err := cursor.All(context.Background(), &tasks); err != nil {
		return nil, tracerr.Wrap(err)
	}

	return tasks, nil
}

func (m *ManagerRepository) CountCompletedTasksByOrderId(orderId uint64) (int64, error) {
	filter := bson.M{
		"order_id": orderId,
		"status":   model.COMPLETED,
	}

	count, err := m.db.Collection("tasks").CountDocuments(context.Background(), filter)
	if err != nil {
		return 0, tracerr.Wrap(err)
	}

	return count, nil
}

func (m *ManagerRepository) FindOutdatedSendedTasks(before time.Time) ([]model.TaskInfo, error) {
	filter := bson.M{
		"status":     model.CREATED,
		"updated_at": bson.M{"$lt": before},
	}

	cursor, err := m.db.Collection("tasks").Find(context.Background(), filter)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	defer cursor.Close(context.Background())

	var tasks []model.TaskInfo
	if err := cursor.All(context.Background(), &tasks); err != nil {
		return nil, tracerr.Wrap(err)
	}

	return tasks, nil
}

func (m *ManagerRepository) AddWorker(worker model.Worker) error {
	_, err := m.db.Collection("workers").InsertOne(context.Background(), worker)
	if err != nil {
		return tracerr.Wrap(err)
	}
	return nil
}

func (m *ManagerRepository) UpdateWorker(worker model.Worker) error {
	filter := bson.M{
		"worker_id": worker.Id,
	}

	update := bson.M{
		"$set": bson.M{
			"max_tasks":   worker.MaxTasks,
			"last_action": worker.LastAction,
		},
	}

	_, err := m.db.Collection("workers").UpdateOne(context.Background(), filter, update)
	if err != nil {
		return tracerr.Wrap(err)
	}

	return nil
}

func (m *ManagerRepository) FindWorker(id int64) (model.Worker, error) {
	filter := bson.M{
		"worker_id": id,
	}

	var worker model.Worker
	err := m.db.Collection("workers").FindOne(context.Background(), filter).Decode(&worker)
	if err != nil {
		return model.Worker{}, tracerr.Wrap(err)
	}

	return worker, nil
}

func (m *ManagerRepository) getAndIncrementSeq() (uint64, error) {
	seqCollection := m.db.Collection("seq")

	filter := bson.D{{Key: "name", Value: "order"}}

	update := bson.D{
		{Key: "$inc", Value: bson.D{{Key: "value", Value: 1}}},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var result model.Sequence

	err := seqCollection.FindOneAndUpdate(context.Background(), filter, update, opts).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			_, err := seqCollection.InsertOne(context.Background(), model.Sequence{Name: "order", Value: 1})
			if err != nil {
				return 0, tracerr.New(fmt.Sprintf("error inserting new sequence: %v", err))
			}

			err = seqCollection.FindOne(context.Background(), filter).Decode(&result)
			if err != nil {
				return 0, tracerr.New(fmt.Sprintf("error retrieving new sequence: %v", err))
			}
		} else {
			return 0, tracerr.New(fmt.Sprintf("error incrementing sequence: %v", err))
		}
	}

	return result.Value, nil
}
