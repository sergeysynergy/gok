#!/bin/bash
SCRIPT_DIR="$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
cd "$SCRIPT_DIR" || exit 1

LOCAL_CONF="./docker-compose.local.yml"
POSTGRES_USER="gok"
POSTGRES_USER_PASSWORD="Passw0rd33"

# Если кто-то это читает и считает, что в опенсорсе всё должно быть на eng: нет проблем, присылайте пул реквест.
help() {
  cat << EOF
Скрипт для управления контейнерами программного комплекса (далее ПК)

Использование:
  ./manage.sh [команда]

Доступные команды:
  help        - помощь
  build all   - сборка всех образов контейнеров ПК
  init        - инициализация ПК
  launch      - сборка, инициализация и запуск контейнерной группировки
  start       - запуск контейнерной группировки
  stop        - остановка контейнерной группировки
  restart     - перезапуск контейнеров
  purge       - остановка и удаление контейнеров

EOF
}

init() {
  echo "Инициализируем программный комплекс: режим разработки"

  if [[ ! -f $LOCAL_CONF ]]; then
    echo "Скопируем конфиг для локальной разработки"
    cp ./docker-compose.yml $LOCAL_CONF
  fi
}

launch() {
  buildAll
  init
  start
}

start() {
  if [[ ! -f $LOCAL_CONF ]]; then
    echo "Скопируем конфиг для локальной разработки"
    cp ./docker-compose.yml $LOCAL_CONF
  fi

  echo "!!! Запуск контейнерной группировки МФСБ в режиме локальной разработки !!!"
  docker-compose -f $LOCAL_CONF up -d

  sleep 1 && docker ps --format "table {{.ID}}\t{{.Names}}\t{{.Ports}}"
}

stop() {
  echo "Останавливаем контейнеры"
  docker-compose stop
  docker-compose -f $LOCAL_CONF stop
}

purge() {
  echo "Удаляем все контейнеры ПК и их хранилища"
  docker-compose stop
  docker-compose -f $LOCAL_CONF stop
  docker-compose kill
  docker-compose -f $LOCAL_CONF kill
  docker-compose rm -vf
  docker-compose -f $LOCAL_CONF rm -vf
}

########################################################################################################################

buildPostgres() {
  echo "Пересобираем образ Postgres"
  docker build --no-cache \
    -t gok-dev-postgres \
    --build-arg USER="$POSTGRES_USER" \
    --build-arg USER_PASSWORD="$POSTGRES_USER_PASSWORD" \
    ./postgres/
}

buildAll() {
  buildPostgres
}

########################################################################################################################

if [[ $1 = "build" ]]; then
  case $2 in
    all)
      buildAll
    ;;
    postgres)
      buildPostgres
    ;;
    *)
      cat << EOF
Требуется второй дополнительный аргумент. Возможные варианты:
  all - собрать образы для всех контейнеров программного комплекса;
  postgres - собрать образ MariaDB;
EOF
    ;;
  esac
  exit 0
fi

########################################################################################################################

if [[ $1 = "init" ]]; then
  init
  exit 0
fi

if [[ $1 = "launch" ]]; then
  launch
  exit 0
fi

if [[ $1 = "start" ]]; then
  start
  exit 0
fi

if [[ $1 = "stop" ]]; then
  stop
  exit 0
fi

if [[ $1 = "restart" ]]; then
  stop
  start
  exit 0
fi

if [[ $1 = "purge" ]]; then
  purge
  exit 0
fi

########################################################################################################################

help
