# Игра жизнь
Знакомы ли вы с игрой «Жизнь»? Сыграть в неё не получится, зато за ней можно наблюдать:

Игра без игроков

«Жизнь» развивается по строго предписанным правилам. 
Она не требует никаких действий со стороны игрока — разве что запустить себя сама она не может.

## Описание игры

Игра проходит на двумерной сетке. Клетки в этой сетке либо мертвы, либо живы. 
По правилам на каждом шаге клетки оживают, умирают или никак не меняют свой статус.

## Правила игры

Каждый шаг игра проверяет, где на поле находятся клетки:

живая клетка с двумя или тремя соседями продолжает жить
мёртвая клетка с тремя соседями оживает
живая клетка погибает, а мёртвая остаётся мёртвой, если у них больше трёх или меньше двух соседей Например:

Несмотря на простые правила, в игре регулярно возникают новые, непохожие друг на друга состояния клеток.

Некоторые из них остаются стабильными и никак не меняются. Эти узоры называют «натюрмортами». 
Проще всего узнать «натюрморт» в квадрате два на два. 
Если он появился с начала игры или возник в процессе эволюции, то будет существовать всю дальнейшую игру без изменений:


Планер (или глайдер) — другая наиболее известная фигура из клеток, которая, в отличие от квадрата, постоянно движется:


Вариантов фигур оказалось так много, что теперь их разделяют на категории. Вот некоторые из них:

- устойчивые (как квадрат)
- двигающиеся
- периодические (время от времени повторяют своё состояние)
- ружья (время от времени «стреляют» планерами)
- паровозы (оставляют след в движении)
- пожиратели (не разрушаются при столкновении с другими фигурами)

У каждой клетки на каждом шаге игры есть некоторое состояние (или state), 
а вся игра — это переход от одного состояния к другому. 
Её можно назвать стейт-машиной, то есть системой набора состояний, событий и переходов из одного состояния в другое. 
В теории программирования стейт-машину называют конечным автоматом (или бесконечным, если поле не ограничено).

## Запуск программы командой в терминале:
```json
go run cmd/life/main.go
```

Установка сторонних библиотек:
```
go get -u go.uber.org/zap
go get -u go.uber.org/zap/zapcore
```

## Запросы пользователя
Метод GET 
```
http://localhost:8081/nextstate
```
Успешное выполнение запроса **200 OK**
```
[[true,false,false,false,false,false,false,false,false,false],
[false,false,false,false,false,false,false,false,false,false],
[false,false,false,false,false,false,false,false,false,false],
[false,false,false,false,false,false,false,false,false,false],
[false,false,false,false,true,false,true,true,true,false],
[true,false,true,true,false,false,false,false,false,true],
[false,false,true,false,false,false,false,false,false,true],
[true,false,false,false,false,false,false,false,false,false],
[false,false,false,false,false,false,true,false,false,false],
[false,true,true,true,false,false,false,true,true,false]]
```
Метод POST
```
http://localhost:8081/setstate
```
body:
```
{
    "fill": 5
}
```
Успешное выполнение запроса **200 OK**
```
[[false,true,true,true,false,false,false,false,false,false],
[false,true,true,true,false,false,false,false,false,false],
[false,false,false,false,false,false,false,false,false,false],
[false,false,false,false,false,false,false,false,false,false],
[false,false,false,false,false,false,false,false,false,false],
[false,false,false,false,false,false,false,false,false,false],
[false,false,false,false,false,false,false,false,false,false],
[false,false,false,false,false,false,false,false,false,false],
[false,false,false,false,false,false,false,false,false,false],
[false,true,true,true,false,false,false,false,false,false]]
```
Метод PUT
```
http://localhost:8081/reset
```
Успешное выполнение запроса **200 OK**
```
{
"fill":5
}
```