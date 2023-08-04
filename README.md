
# Webitel logger

Webitel logger


## Config

#### Setting up

To start up the logger you should do a config file in json format.

The config file consists of 3 sections: "rabbit", "database", "consul".

Example of config.json:
```json
{
    "rabbit": {
        "url": "amqp://admin:admin@10.9.8.111:5672"
    },
    "database": {
        "url": "postgres://postgres:postgres@10.9.8.111:5432/postgres"
    },
    "consul": {
        "address": "10.9.8.111:8500",
        "id" : "logger",
        "publicAddress": "10.10.10.162:10001"
        
    }
}   
```

Rabbit object:

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `url` | `string` | **Required**. Connection string to the rabbit client |

Database object:

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `url` | `string` | **Required**. Connection string to the database |

Database object:

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `address` | `string` | **Required**. Address of the consul |
| `id` | `string` | **Required**. The tag of the registered service |
| `publicAddress` | `string` | **Required**. The address at which the grpc will start up and also this address will be send to the consul as the main access address. |

When you created your config file the next step will be run logger.

There is only one flag in the application and this is the flag that allows pass the path to the config file:

| Flag       | Type     | Description                                                                                               |
|:-----------| :------- |:----------------------------------------------------------------------------------------------------------|
| `--config` | `string` | **NOT Required**. Path to the config file. Default path points to the root directory config/config.json . |