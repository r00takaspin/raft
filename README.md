[Алогоритма консенсуса RAFT](https://raft.github.io)
---

Исплементация алгоритма нахождения консенсуса Raft на языке Go. Ссылка на [описание алгоритма](https://raft.github.io/raft.pdf).

На данный момент реализован механизм избрания лидера. 

Зависимости:
   
   * docker
   * docker-compose

Сборка:

    $ docker build . -t node
    $ docker-compose build
    

Запуск:
    
    $ docker-compose up

Переизбрание нового лидера можно инициировать убийством одной из нод:

    $ docker ps
    $ docker kill <container_id>
    
    
    
 