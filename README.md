# Kafka Toolkit

Libreria para crear Sarama consumers/producers y reducir la cantidad de codigo base al momento de escribir una aplicacion que consume o produce en topicos Kafka, sin embargo al utilizar interfaces genericas podria utilizar otro driver como confluent kafka en el futuro.

## Interfaces provistas

### Message Handler

[Message Handler](message_handler.go), es una interfaz que nos permite implementar nuestra propia logica al momento de recibir un mensaje en un topic Kafka, un Message Handler de ejemplo se encuentra en nuestro codigo base: [Logging Message Handler](logging_message_handler.go).

### Consumer Error Handler
[Consumer Error Handler](consumer_error_handler.go), es una interfaz que nos permite implementar nuestra propia logica para el manejo de errores al momento de leer un mensaje del topico, un Consumer Error Handler se encuentra en nuestro codigo base: [Logging Consumer Error Handler](logging_consumer_error_handler.go).


## Como inicializar consumer
Importar package `"gitlab.falabella.com/fif/integracion/forthehorde/commons/kafka-toolkit"` en proyecto e inicializar configuracion, en el ejemplo a continuacion creamos un objeto de tipo **ConsumerGroupInput** que almacenara nuestra configuracion:

Primero inicializamos nuestra configuracion naga.

```go
return ConfigEntries() []naga.ConfigEntry{
		{
			VariableName: "kafka_brokers",
			Description:  "brokers kafka a conectar",
			DefaultValue: "localhost:9091",
		},
		{
			VariableName: "client_id_prefix",
			Description:  "Prefijo de Kafka Client ID, el clientID resultante sera <prefijo>@<PID>",
			DefaultValue: "ino-prefix",
		},
		{
			VariableName: "consumer_group",
			Description:  "consumer group Kafka",
			DefaultValue: ":8080",
		},
		{
			VariableName: "version",
			Description:  "Kafka version",
			DefaultValue: "2.1.1",
		},
		{
			VariableName: "input_topic",
			Description:  "Topicos de entrada (raw)",
			DefaultValue: "TEST",
		},
		{
			VariableName: "strategy",
			Description:  "Estrategia de consumer group (range, roundrobin)",
			DefaultValue: "roundrobin",
		},
		{
			VariableName: "earliest",
			Description:  "Indica si utiliza el offset mas nuevo inicialmente",
			Shortcut:     "e",
			DefaultValue: false,
		},
		{
			VariableName: "security_enabled",
			Description:  "Seguridad de Kafka habilitada",
			DefaultValue: false,
		},
		{
			VariableName: "kafka_cert_part",
			Description:  "Path de CA Cert",
			DefaultValue: "key.pem",
		},
		{
			VariableName: "kafka_username",
			Description:  "SASL Username",
			DefaultValue: "username",
		},
		{
			VariableName: "kafka_password",
			Description:  "SASL Password",
			DefaultValue: "password",
		},
	}
}
```

Posterior a eso inicializamos la configuracion Kafka
```go
func InitKafkaConfig() kafka.ConsumerGroupInput {
	clientID := fmt.Sprintf("%v-%v", options["client_id_prefix"].(string), os.Getpid())

	return kafka.ConsumerGroupInput{
		Brokers:         options["kafka_brokers"].(string),
		Topic:           options["input_topic"].(string),
		Group:           options["consumer_group"].(string),
		ClientID:        clientID,
		BalanceStrategy: options["strategy"].(string),
		Oldest:          options["earliest"].(bool),
		Version:         options["version"].(string),
		Security:        options["security_enabled"].(bool),
		CaFile:          options["kafka_cert_part"].(string),
		Username:        options["kafka_username"].(string),
		Password:        options["kafka_password"].(string),
	}
}
```

Luego llamamos a la inicializacion de la configuracion Kafka en nuestra funcion init en main.go:

```go
func main() {
	inputConf = kafka.InitKafkaConfig()
	...
}
```
Luego de haber inicializado la configuracion podemos empezar a implementar la funcion main, donde debemos inicializar un **Message Handler** y un **Error Handler**

Ahora incializamos el consumo del topic con nuestro message y error handlers, en este caso utilizamos las implementaciones que generan logs tanto de mensaje como para los errores, pero se pueden crear otras implementaciones de ser necesario.
```go
    msgHandler := kafka.NewBaseHandler(ep, // Endpoint de tipo gokit endpoint
			decodeFunc, // Decode func debe ser de tipo DecodeConsumerMessageFunc 
	errorHandler := kafka.NewLoggingConsumerErrorHandler()

	consumer, err := kafka.MakeSaramaConsumerBuilder(inputConf, msgHandler).WithErrorHandler(errorHandler).Build()
	if err != nil {
		log.Panicf("Error creando consumer: %v", err)
	}

    if err := consumer.Start() ; err != nil {
        log.Panicf("Error inicializando consumer: %v", err)
    }
}
```
- ep es una variable de tipo go-kit endpoint, en general este endpoint puede ser quien llame a un service
- decodeFunc es la funcion que convierte el objecto de tipo `ConsumerMessage` y lo convierte al tipo de dato a ser procesado en el endpoint.

ver [Encode y Decode Funcs](encode_decode.go)

Listo, con esa configuracion debiesemos estar listos para empezar a consumir mensajes Kafka.

## Como inicializar un producer
Para inicializar un producer necesitamos crear un nuevo simple sync producer, kafka-toolkit nos provee una funcion para inicializar este producer:

```go
	producerConfInput := kafka.BaseProducerConfigInput{
		Brokers: brokers,
		Ack:     -1,
		Retries: 3,
		Version: version,
		Topic:   outTopic,
	}

	producer, err := kafka.NewSimpleSyncProducer(producerConfInput)

	if err != nil {
        log.Panicf("Error inicializando producer: %v", err)
    }
```
Para enviar un mensaje debemos crear un message producer endpoint.

se debe crear un `encode func` que recibe un objeto (interface) y lo convierte a el tipo `ProducerMessage`

Una vez implementada el encode func podemos crear el endpoint

```go
	producerEndpoint := kafka.MakeMessageProducerEndpoint(producer, encodeOK)
```

Para luego enviar nuestro mensaje

```go
_, err := producerEndpoint(ctx, someInterface)
if err != nil {
	return nil, err
}
```

## Como crear un streamer
Un Streamer es un tipo de consumer que ademas de consumir un mensaje desde un topico, posterior a procesar, enviara un mensaje a otro topico.
Para crear un nuevo streamer debemos primero que nada crear un `StreamProcessor`, el cual se encarga de "decodear" el mensaje de entrada, llamar a un service endpoint y "encodear" el mensaje de salida.

```go
	streamProcessor := kafka.MakeStreamProcessor(ep, // Endpoint de tipo gokit endpoint
			decodeFunc, // Decode func debe ser de tipo DecodeConsumerMessageFunc
			encodeFunc) // Encode func debe ser de tipo EncodeProducerMessageFunc
```

*Nota*: decodeFunc y encodeFunc son el mismo tipo de funcion mencionados en consumer y producer kafka respectivamente.

- ep es una variable de tipo go-kit endpoint, en general este endpoint puede ser quien llame a un service
- decodeFunc es la funcion que convierte el objecto de tipo `ConsumerMessage` y lo convierte al tipo de dato a ser procesado en el endpoint.
- encodeFunc recibe un objeto, normalmente la response del endpoint y lo convierte a el tipo `ProducerMessage`

Luego podemos crear un streamer del siguiente modo:


```go
	streamerBuilder ,err := kafka.MakeStreamerBuilder(consumerConf, streamProcessor, producerConf)
	if err != nil {
		log.Panicf("Error creando streamer builder: %v", err)
	}
	
	streamer, err := streamerBuilder.Build()
	if err != nil {
		log.Panicf("Error creando streamer: %v", err)
	}

    if err := streamer.Start() ; err != nil {
        log.Panicf("Error inicializando consumer: %v", err)
    }

```

### StreamProcessor middlewares
Tambien se ha agregado la opcion de crear middlewares para el `StreamProcessor`, de modo de poder decorar la funcionalidad por defecto, a continuacion algunos de los middlewares incluidos en la libreria.

#### HeaderPassThruStreamProcessor
Middleware que toma los headers del mensaje de entrada y los deja en el header de salida.

```go
	streamProcessor := kafka.MakeStreamProcessor(ep, // Endpoint de tipo gokit endpoint
			decodeFunc, // Decode func debe ser de tipo DecodeConsumerMessageFunc
			encodeFunc) // Encode func debe ser de tipo EncodeProducerMessageFunc
	
	streamProcessor = kafka.MakeHeaderPassThruStreamProcessorMiddleware()(streamProcessor)
```

#### OpenTracingStreamProcessor
Middleware que utilizando un tracer de opentracing y un nombre de operacion:

* Toma Headers de mensaje kafka de entrada, crea un span y lo agrega al context.
* Toma Span desde context y agregar Headers de mensaje de salida de Kafka.


```go
	streamProcessor := kafka.MakeStreamProcessor(ep, // Endpoint de tipo gokit endpoint
			decodeFunc, // Decode func debe ser de tipo DecodeConsumerMessageFunc
			encodeFunc) // Encode func debe ser de tipo EncodeProducerMessageFunc
	
	streamProcessor = kafka.MakeOpenTracingStreamProcessorMiddleware("my_operation", tracer)(streamProcessor)
```

### KeyFilterStreamMiddleware
Permite filtrar mensajes de streams con keys especificas

```go
	streamProcessor := kafka.MakeStreamProcessor(ep, // Endpoint de tipo gokit endpoint
			decodeFunc, // Decode func debe ser de tipo DecodeConsumerMessageFunc
			encodeFunc) // Encode func debe ser de tipo EncodeProducerMessageFunc
	
	allowedKeys := []string{"my_key"}

	streamProcessor = kafka.MakeKeyFilterMessageHandlerMiddleware(allowedKeys)(streamProcessor)
```


## Como crear un health Check
Se puede usar la función **Health** definida en la interfaz **HealthCheck**, este se utiliza de la siguiente forma:

```go

	healthCheck := kafka.NewHealthCheck([]string{brokers}, sarama.NewConfig())
	if err := healthCheck.Health(); err != nil {
		log.Fatal(err)
	}

```

Esto retornara un error en caso de no poder conectarse a algún broker. Para conexiones con seguridad puede hacerlo de la siguiente forma:
```go

	config, err  := kafka.SaramaSASLConfig(KafkaSASLSecurity{
			Username: username,
			Password: password,
			CaFile:   aaFile,
		})
	if err != nil {
		log.Println(err)
	}

	healthCheck := kafka.NewHealthCheck([]string{brokers},config)
	if err := healthCheck.Health(); err != nil {
		log.Fatal(err)
	}
```

### Endpoint HTTP de Health (y metrics)
Para crear endpoint http de health se debe inicializar el builder de HTTP handler.

```go
	kafkaMetrics := kafka.MakeDefaultKafkaMetrics("my_service_name") // nombre de service o de app, utilizado para las metricas

	httpHandler := kafka.MakeHealthHandlerBuilder(
		logger, // logger de go-kit
		healthChek, // service de healthCheck mencionado en punto anterior
	).WithMetrics(kafkaMetrics).Build()
```

Este handler provee de 2 enpoints "/healthz" y "/metrics" los cuales respectivamente indican la salud y las metricas del service.

*NOTA:* Para mas info en como utilizar el HTTP Handler builder ver [go microservices commons](https://github.com/validatecl/go-microservices-commons)

### Decoracion de metricas en endpoint para consumer o streamer kafka
Simplemente se debe crear endpoint de consumer o streamer y decorar con los middleware de metrics de la libreria commons:

```go
	endpoint := endpoint.MakeMyEndpoint(service)
	endpoint = commons.EndpointSuccesAndFailureMetricsMiddleware(kafkaMetrics.SuccessfulRequests, kafkaMetrics.FailedRequests)(endpoint)
	endpoint = commons.EndpointTimeTakenMetricsMiddleware("my_operation", kafkaMetrcis.RequestDuration)(endpoint)
```

## Kafka Middleware functions
Se han agregado funciones middleware para decorar distintas funcionalidades de consumer, producer y modifier Kafka, esto permite crear nuevos decoradores personalizadas como tracing, instrumentation, logging, etc...
A continuacion los 4 nuevos middlewares provistos:

- MessageHandlerMiddleware
- MessageProducerMiddleware
- DecodeConsumerMessageFuncMiddleware
- EncodeProducerMessageFuncMiddleware

ver [Middleware functions](middleware.go) y [Encode y Decode middlewares](encode_decode.go)

adicional a esto se implemento un handler logging middleware que puede ser utilizado en caso de requerir loggear el mensaje y error de un message handler.

```go
//myMessageHandler MessageHandler personalizado
myMessageHandler := kafka.MakeLoggingMessageHandlerMiddleware()(myMessageHandler)
```

ver [Tests de logging handler middleware](logging_message_handler_middleware_test.go)

## Como correr tests

Para correr los tests simplemente se debe hacer `make test`

