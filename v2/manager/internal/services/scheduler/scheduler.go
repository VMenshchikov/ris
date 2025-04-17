package scheduler

import (
	"hash_manager/internal/domain/model"
	"hash_manager/internal/domain/repo"
	"hash_manager/internal/infra/kafka"
	"log"
	"math"
	"sync/atomic"
	"time"
)

const (
	DEFAULT_TASKS_COUNT = 2
	RESEND_DURATION     = time.Minute
	RESEND_CHECK        = time.Minute
)

type Scheduler struct {
	tasksChannel chan Task

	writer *kafka.Sender
	repo   repo.ManagerRepository

	currentTasks atomic.Int64
	maxTasks     atomic.Int64
}

type Task struct {
	OrderId     uint64
	TargetHash  [16]byte
	BlockSize   uint
	BlockNumber uint
	MaxLen      uint
}

func CreateScheduler(repo repo.ManagerRepository, kafkaWriter *kafka.Sender) *Scheduler {
	ret := Scheduler{
		tasksChannel: make(chan Task),
		writer:       kafkaWriter,
		repo:         repo,
		currentTasks: atomic.Int64{},
		maxTasks:     atomic.Int64{},
	}
	ret.maxTasks.Store(DEFAULT_TASKS_COUNT)
	return &ret
}

func (s *Scheduler) UpdateMaxTasks(delta int64) {
	s.maxTasks.Add(delta)
}

func (s *Scheduler) FinishTask() {
	for {
		current := s.currentTasks.Load()
		if current <= DEFAULT_TASKS_COUNT {
			log.Println("!!!!!!!!!!!!1")
			break
		}

		if s.currentTasks.CompareAndSwap(current, current-1) {
			continue
		}
	}
	log.Printf("Updated task count: %d/%d", s.currentTasks.Load(), s.maxTasks.Load())
}

func (s *Scheduler) Run() {
	go s.sendTasksToKafka()
	go s.resendTasks()
	go s.createTasks()
	go s.closeTimeoutOrders()
}

func (s *Scheduler) resendTasks() {
	for {
		time.Sleep(RESEND_CHECK)

		tasks, err := s.repo.FindOutdatedSendedTasks(time.Now().Add(-RESEND_DURATION))
		if err != nil {
			log.Println(err)
			continue
		}

		for _, v := range tasks {
			order, err := s.repo.FindOrder(v.OrderId)
			if err != nil {
				log.Println(err)
			}

			s.tasksChannel <- s.createTask(
				v.OrderId,
				order.TargetHash,
				uint(order.BlockSize),
				v.BlockNumber,
				order.MaxLen,
			)
		}

	}
}

func (s *Scheduler) sendTasksToKafka() {
	for {
		log.Println("Current:", s.currentTasks.Load(), "Max:", s.maxTasks.Load())
		s.waitFreeSlot()

		task := <-s.tasksChannel
		msg := kafka.MessageTask{
			OrderId:     task.OrderId,
			TargetHash:  task.TargetHash,
			BlockSize:   task.BlockSize,
			BlockNumber: task.BlockNumber,
			MaxLen:      task.MaxLen,
		}
		if err := s.writer.Send(msg); err != nil {
			log.Println("Ошибка отправки в кафку, задача не отправлена ")
			continue
		}
		log.Println("Task sent to Kafka successfully")

		domainTask, err := s.repo.FindTask(task.OrderId, uint64(task.BlockNumber))
		if err != nil {
			log.Println(err)
		}

		domainTask.Status = model.SENDED
		s.repo.UpdateTask(*domainTask)
		log.Println("Task status change to SENDED")

		s.currentTasks.Add(1)
		log.Printf("Updated task count: %d/%d", s.currentTasks.Load(), s.maxTasks.Load())
	}
}

func (s *Scheduler) waitFreeSlot() {
	for {
		if s.currentTasks.Load() < s.maxTasks.Load() {
			return
		}
		time.Sleep(time.Millisecond * 50)
	}
}

func (s *Scheduler) closeTimeoutOrders() {
	for {
		time.Sleep(10 * time.Second)
		s.repo.CloseTimeoutOrders()
	}
}

func (s *Scheduler) createTasks() {
	for {
		oModel, err := s.repo.FindOrderForExecution()
		if err != nil {
			log.Println(err)
			time.Sleep(time.Second) //todo wakeup? неопрравданное усложнение
			continue
		}
		if oModel.Timeout.Before(time.Now()) {
			oModel.Status = model.ERROR
			s.repo.UpdateOrder(oModel)
			continue
		}

		orderTasks, err := s.repo.FindTasksByOrderId(oModel.Id)
		if err != nil {
			log.Println("Ошибка при поиске задач заказа:", err)
			return
		}
		log.Println(oModel.MaxLen, oModel.BlockSize)
		allTasksCount := uint(math.Ceil(math.Pow(float64(36), float64(oModel.MaxLen)) / float64(oModel.BlockSize)))
		orderTasksCount := len(orderTasks)

		for i := s.maxTasks.Load() - s.currentTasks.Load(); i > 0; i-- {
			log.Printf("Шедулер создает задачу %d/%d для заказа %d", orderTasksCount+1, allTasksCount, oModel.Id)

			s.tasksChannel <- s.createTask(
				oModel.Id,
				oModel.TargetHash,
				uint(oModel.BlockSize),
				uint(orderTasksCount),
				oModel.MaxLen,
			)
			orderTasksCount++

			if orderTasksCount == int(allTasksCount) {
				oModel.Status = model.ALL_TASKS
				s.repo.UpdateOrder(oModel)
				log.Println("Созданы все задачи для ")
				break
			}
		}

	}
}

func (s *Scheduler) createTask(orderId uint64, targetHash [16]byte, blockSize, blockNumber, maxLen uint) Task {
	task := Task{
		OrderId:     orderId,
		TargetHash:  targetHash,
		BlockSize:   blockSize,
		BlockNumber: blockNumber,
		MaxLen:      maxLen,
	}

	err := s.repo.AddTask(model.TaskInfo{
		OrderId:     task.OrderId,
		Status:      model.CREATED,
		BlockNumber: task.BlockNumber,
		UpdatedAt:   time.Now(),
	})
	if err != nil {
		log.Println(err)
	}
	return task
}
