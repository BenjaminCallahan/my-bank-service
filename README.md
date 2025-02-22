## Запуск приложения
### Зависимости
- go 1.17.5
- docker

Сборка образа
```
docker build -t my-bank-service .
```
Запуск контейнера
```
docker run --name my-bank-service_container \
    -p 8080:8080 \
    -e GIN_MODE=release \
    -v external-folder_sqlite:/usr/local/bin/app/sqlite_database/ \
    my-bank-service 
```



## Задача
Написать сервис личного банка, взаимодействие с которым реализуется по REST API.
Необходимо предоставить работающий сервис и описание API (желательно Postman коллекцией с тестами на указанные кейсы).

## Условия
Возможности банка и условия его работы:
1. Вы там являетесь _единственным_ вкладчиком и у вас там уже открыт счёт.
2. Валюта счёта - суверенный лысый попугай SBP (sovereign bald parrot). Дробная часть SBP - это **2х значное** суверенное пёрышко - spf (sovereign parrot feather). Курс SBP2RUB обеспечен рабским трудом фрилансеров на галерах и статичен много лет, составляя **0,7523**.
3. Все операции со счётом производятся в валюте счёта.
4. Банк даёт возможность пополнять счёт любыми суммами.
5. Банк обеспечивает накопления в размере 6% _от суммы на счёте_ сразу же после _каждого_ пополнения. Сумма дохода складывается на этот же счёт.
6. Банк позволяет узнать баланс вашего счёта как в валюте счета (по-умолчанию), так и в RUB.
7. Банк предоставляет возможность снятия денег со счёта, но _не более_ 70% от суммы на счёте _за раз_.

## Требования
Банк у нас прогрессивный, работает на sqlite базе. Первый запуск сервиса должен создать базу и проинициализировать её.

В репозитории находится [интерфейс](interface.go), который должен реализовывать объект счёта.

Формат входных-выходных данных: JSON

Других ограничений или требований при реализации не предусматривается. 

## Проверка
API будет проверено следующими тест-кейсами:
1. Успешное пополнение изначального нулевого баланса на 72.00 SBP
2. Запрос баланса в SBP. Результат должен быть равен **76.32 SBP**
3. Успешное пополнение на 37.50 SBP
4. Запрос баланса в SBP. Результат должен быть равен **120.65 SBP**
5. Успешное пополнение на 10.20 SBP
6. Запрос баланса в SBP. Результат должен быть равен **138.70SBP**
7. *Неуспешный* вывод 127.60 SBP
8. Запрос баланса в SBP. Результат должен быть равен **138.70SBP**
9. Успешный вывод 30.00 SBP
10. Запрос баланса в SBP. Результат должен быть равен **108.70SBP**
11. Запрос баланса в RUB. Результат должен быть равен **81.78RUB**

Последовательность действий и проверка результатов следует оформить в Postman тестах (проверяем значение баланса) для ускорения проверки.

Код должен быть оформлен в виде форка данного репозитория.
