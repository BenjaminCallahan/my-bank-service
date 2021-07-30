Задача написать сервис личного банка, взаимодействие с которым реализуется по REST API.

Возможности банка и условия его работы:
1. Вы там являетесь _единственным_ вкладчиком и у вас там уже открыт счёт.
2. Валюта счёта - SBP (sovereign bald parrot). Курс SBP2RUB обеспечен рабским трудом фрилансеров на галерах и статичен много лет, составляя **0,7523**.
3. Банк даёт возможность пополнять счёт любыми суммами в валюте счёта.
4. Банк обеспечивает накопления в размере 6% от суммы _каждого_ пополнения сразу же. Сумма накопления пополняет этот же счёт.
5. Банк позволяет узнать баланс вашего счёта. Баланс можно узнать как в валюте счета, так и в RUB.
6. Банк предоставляет возможность снятия денег со счёта, но _не более_ 30% от суммы вклада _за раз_.

Банк работает на sqlite базе. Первый запуск сервиса должен создать базу и проинициировать её.

Необходимо предоставить работающий сервис и описание API (желательно Postman коллекцией).

API будет проверено следующей последовательностью действий:
1. Пополнение на 72.00 SBP
2. Пополнение на 37.50 SBP
3. Пополнение на 10.20 SBP
4. Запрос баланса в SBP
5. Вывод 127.60 SBP
6. Вывод 30.00 SBP
7. Запрос баланса в RUB
8. Запрос баланса в SBP

Результат двух последних действия должен составить 81.75 RUB и 108.67 SBP, соответственно
