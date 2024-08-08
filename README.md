![build](https://github.com/dzhordano/avito-aug2024/actions/workflows/default.yml/badge.svg)
![Coverage](https://img.shields.io/badge/Coverage-7.4%25-red)

## Problem & Solutions
    1. Отсутствие в API номера квартиры при вызове endpoint'a /flat/create.
    Решение: Решил добавить отдельное поле FlatNumber (flat_number в бд), которое это делает.
    
    2. В API указано, что при создании квартиры ей выдается статус on moderate (moderating),
        а в github по условию - created.
    Решение: Статус указывается как created при создании квартиры. 
        OnModerate устанавливается при запросе на /flat/update пока не будет завершен запрос,
        иначе (при отмене контекста или ошибке) возвращает изначальное состояние квартиры.

## Build & Run:
- Для запуска с Docker `make run`.
- Иначе последовательно `make init-db`, `make run-no-docker`.
- Использовать `make lint` для запуска линтера.
- Генерация документации swagger `make swag`.
---
### Testing
    Не знаю как исправить 'connection reset by peer', поэтому последовательно:
    Для интеграционных тестов:    
        1. make init-db-test
        2. make test.integration
    Для unit тестов:
        make test
---
- Go version `1.22.5`