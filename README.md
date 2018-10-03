# [Имплементация алогоритма консенсуса RAFT](https://raft.github.io)

Сборка:

    $ docker build . -t raft
    
Запуск:
    
    $ docker run -p 24816:24816 raft
    
Тестовый запрос:
    
    $ grpc_cli call 127.0.0.1:24816 SetValue "value: 1" --protofiles proto/raft-grpc.proto
    
 