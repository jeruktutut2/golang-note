## note web server
httprouter will conflict if create endpoint like this: /:patvariable/link and /link
fiber cannot get cancelation signal from client
servemux cannot override not found handler and method not allowed
gin cannot override method not allowed

## note pdf
github.com/johnfercher/maroto

#note google oauth2
go get golang.org/x/oauth2

## rabbitmq
https://hub.docker.com/_/rabbitmq
docker pull rabbitmq
docker run -d --hostname my-rabbit --name rabbitmq-note -e RABBITMQ_DEFAULT_USER=user -e RABBITMQ_DEFAULT_PASS=password -e RABBITMQ_DEFAULT_VHOST=my_vhost -p 15672:15672 -p 5672:5672 rabbitmq:3-management
docker run -d --hostname my-rabbit --name rabbitmq-note -e RABBITMQ_DEFAULT_USER=user -e RABBITMQ_DEFAULT_PASS=password -e RABBITMQ_DEFAULT_VHOST=my_vhost -p 5672:5672 rabbitmq:3-management
use this, no initial vhost couse error Exception (403) Reason: "no access to this vhost", use default vhost instead: docker run -d --hostname my-rabbit --name rabbitmq-note -e RABBITMQ_DEFAULT_USER=user -e RABBITMQ_DEFAULT_PASS=password -p 15672:15672 -p 5672:5672 rabbitmq:3-management
go get github.com/rabbitmq/amqp091-go
http://localhost:15672/#/
web management: localhost:15672
rabbitmq port: localhost: 5672
create exchange (notification) - create queue (email) - create bindings

## kafka
go get -u github.com/lovoo/goka
docker-compose up -d
docker-compose down
create email topic: kafka-topics --bootstrap-server localhost:9092 --create --topic email
create text message topic: kafka-topics --bootstrap-server localhost:9092 --create --topic text-message
show list topic: kafka-topics --bootstrap-server localhost:9092 --list
producer email: kafka-console-producer --broker-list localhost:9092 --topic email
producer text-message: kafka-console-producer --broker-list localhost:9092 --topic text-message
consumer email: kafka-console-consumer --bootstrap-server localhost:9092 --topic email
consumer text-message: kafka-console-consumer --bootstrap-server localhost:9092 --topic text-message
consumer group email: kafka-console-consumer --bootstrap-server localhost:9092 --topic email --group email-consumer-group
consumer group text-message: kafka-console-consumer --bootstrap-server localhost:9092 --topic email --group text-message-consumer-group
https://www.conduktor.io/kafka/how-to-start-kafka-using-docker/
https://github.com/conduktor/kafka-stack-docker-compose
docker-compose -f zk-single-kafka-single.yml up -d
show list topic: kafka-topics --bootstrap-server kafka1:9092 --list
create email topic: kafka-topics --bootstrap-server kafka1:9092 --create --topic email
create text message topic: kafka-topics --bootstrap-server kafka1:9092 --create --topic text-message