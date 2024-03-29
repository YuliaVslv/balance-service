# Микросервис для работы с балансом пользователя

RESTful микросервис, предоставляющий возможность взаимодействия с балансом пользователя (зачисление средств, резервирование средств на отдельном счете, списывание средств при подтверждении операции, отмена оплаты). 
Также есть возможность сформировать отчет об оплаченных услугах за указанный месяц в формате csv и получить полную историю транзакций для указанного пользователя (реализована пагинация и сортировка для списка транзакций).

Взаимодействие с сервисом происходит через HTTP-запросы с телом запроса в формате JSON.

## Запуск

    docker-compose build  
    docker-compose up

## Поддерживаемые запросы
Запросы принимаются на localhost:7000
Все запросы возвращают ответ в формате JSON (данные или сообщения).

### GET /balance/{id}
      Пример запроса: localhost:7000/balance/1

Получение баланса пользователя с заданным id.  
Возращает данные счета в формате JSON, если пользователь существует. В противном случае возвращает код 400 (BadRequest) и сообщение в формате JSON.

### POST /credit
      Пример запроса: localhost:7000/credit  

      {  
          "user_id": 1,  
          "value": 1000  
      } 

Зачисление средств на счет пользователя.  
Возвращает код 200 (OK) в случае успеха и сообщение в формате JSON.

### POST /reserve
    Пример запроса: localhost:7000/reserve

    {  
      "user_id": 1,  
      "service_id": 200,  
      "order_id": 250,  
      "value": 1000  
    }  

Резервирование средств на отдельном счете.  
Возвращает код 200 (OK) в случае успеха и сообщение в формате JSON.

### POST /withdraw
    Пример запроса: localhost:7000/withdraw  

    {  
      "user_id": 1,  
      "service_id": 200,  
      "order_id": 250,  
      "value": 1000  
    }  

Подтверждение выручки: списывание средств с резервного счета.  
Возвращает код 200 (OK) в случае успеха и сообщение в формате JSON.  

### POST /refund
    Пример запроса: localhost:7000/refund  

    {  
      "user_id": 1,  
      "service_id": 200,  
      "order_id": 250,  
      "value": 1000  
    }

Отмена оплаты: возврат средств с резервного счета на счет пользователя.  
Возвращает код 200 (OK) в случае успеха и сообщение в формате JSON.

### GET /report/{year}/{month}
    Пример запроса: localhost:7000/report/2022/11  

+ year - год: число от 2007 до номера текущего года
+ month - месяц: число от 1 до 12  

Запрос отчета за заданный год и месяц. Формируется отчет, записывается в CSV файл.  
Возвращает код 200 (OK) в случае успеха и сообщение в формате JSON с указанием пути до файла отчета.

### GET /history/{id}/{page}/{sort_field}/{order}
    Примеры запросов:  
    localhost:7000/history/1/1/value/desc - транзакции пользователя с id=1, отсортированные по сумме по убыванию (страница 1)
    localhost:7000/history/3/2/date - транзакции пользователя с id=3, отсортированные по сумме по возрастанию (страница 2)
    localhost:7000/history/5/10 - транзакции пользователя с id=5 без сортировки (страница 10)
    localhost:7000/history/7 - транзакции пользователя с id=7 без сортировки (страница 1) 
    

+ id - идентификатор пользователя  
+ page - номер страницы (необязательно, по умолчанию 1)  
+ sort_field - по какому полю сортировка (необязательно, по умолчанию сортируется по id транзакции)  
  1. value - по сумме 
  2. date - по дате
+ order - порядок сортировки (необязательно, по умолчанию по возрастанию)
  1. asc - в порядке возрастания
  2. desc - в порядке убывания

Запрос истории транзакций для заданного пользователя. Реализована пагинация и сортировка по сумме или дате.  
Возвращает код 200 (OK) в случае успеха и список транзакций в формате JSON (по 10 транзакций на странице).
