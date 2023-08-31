
# Webitel logger

Webitel logger


## Config

#### Setting up

To start up the logger you should parse program flags.


Flags available (THE DEFAULT VALUE ONLY FOR [id] flag!):

| Flag   | Type     | Description                                                                                                                            |
|:-------| :------- |:---------------------------------------------------------------------------------------------------------------------------------------|
| `amqp` | `string` | **Required**. Connection string to the rabbit client                                                                                   |
| `data_source` | `string` | **Required**. Connection string to the database                                                                                        |
| `consul`    | `string` | **Required**. Address of the consul                                                                                                    |
| `id`        | `string` | **NOT Required**. The tag of the registered service. Default 'logger'                                                                  |
| `grpc_addr` | `string` | **Required**. The address at which the grpc will start up and also this address will be send to the consul as the main access address. |

FOR DEVS!

To use client in your services you should import package client from pkg/client. To create a client you should use a NewClient function and open a connection with Open() method on client. After that call Rabbit().{Create|Update|Delete}Action method on the client to send the logs. Then you should call {One|Many} method to pass the record id and it's new state in json format. And in the end you call SendContext to send message you built.
The example of usage:

https://github.com/Dtsnko/test_logger/blob/master/main.go


The client itself knows when logs should be send, so all you have to do is call these function on every API.

Client firstly call logger service api to know if config for that system object enabled and logs should be saved. If yes logs will be send to the rabbit. If the service is not active or an error occured while sending grpc request, the error will be returned.

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `address` | `string` | **Required**. Address of the consul |
| `id` | `string` | **Required**. The tag of the registered service |
| `publicAddress` | `string` | **Required**. The address at which the grpc will start up and also this address will be send to the consul as the main access address. |


