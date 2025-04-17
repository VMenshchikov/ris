#!/bin/bash

set -e

echo "Запуск docker-compose..."
docker-compose up -d --build

echo "Ожидаем, брокера на 9092"
until docker exec kafka-broker bash -c "echo > /dev/tcp/kafka-broker/9092" 2>/dev/null; do
  sleep 1
done

echo "Kafka доступна, создаём топики..."

docker exec kafka-broker kafka-topics.sh \
  --bootstrap-server localhost:9092 \
  --topic __consumer_offsets \
  --partitions 1 \
  --replication-factor 1 \
  --create \
  --if-not-exists

docker exec kafka-broker kafka-topics.sh \
  --bootstrap-server localhost:9092 \
  --create --if-not-exists \
  --topic manager-to-workers \
  --partitions 3 \
  --replication-factor 1

docker exec kafka-broker kafka-topics.sh \
  --bootstrap-server localhost:9092 \
  --create --if-not-exists \
  --topic workers-to-manager \
  --partitions 1 \
  --replication-factor 1

docker start v2-manager-1

docker-compose logs -f
