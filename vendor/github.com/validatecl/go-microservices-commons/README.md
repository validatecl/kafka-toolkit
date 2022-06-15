# Go Microservices commons

![pipeline](https://github.com/validatecl/commons/go-microservices-commons/badges/develop/pipeline.svg)![coverage](https://github.com/validatecl/commons/go-microservices-commons/badges/develop/coverage.svg)

Componentes comunes de microservicios desarrolados con gokit

## Errores comunes y su manejo
Microservices commons trae variados errores comunes de los servicios que desarrollamos, estos distintos tipos de errores son manejados a nivel de transporte http por [Error2Wrapper](error_wrapper.go) el cual permite manejar estos errores comunes y retornar un struct con un status code correspondiente al error, ejemplo:

    {
        "code": "400",
        "message": "Datos de entrada no válidos.",
        "docs": "https://developerportal.apigee.io/apis/{apidoc}"
    }

| Error | Status | Obs |
| ---   |    --- |  --- |
| InputError | 400 | 
| ServiceError | 500 |
| GatewayError | 503 | 
| BackendCodedError | 502 | Errores de Backend con codigos de error |

### Personalizacion de DOCS URL
Para personalizar contenido de docs en response de error debe modificarse la variable `DOCSURL` del package, ejemplo:

    DOCSURL = "https://developerportal.apigee.io/apis/{apidoc}"

## HTTP Transport response encoder
Se agrego tambien una funcion default para Encoding de response http que en combinacion con el `Error2Wrapper` mencionado en el punto anterior. Esta funcion `EncodeHTTPResponseFunc` posterior a ser creada con `MakeEncodeHTTPResponseFunc` puede ser usada como argumento como funcion de EncodeResponse al momento de crear un http server. para referencias ver [tests](http_transport_test.go)

## Default Endpoint y sus middlewares
Se provee de una funcion middleware `MakeDefaultEntryEndpoint` para crear endpoint por defecto al cual se decora con los siguientes middlewares:

| Middleware | Obs |
| --- | --- |
| EndpointLogMiddleware | Loguea operacion, tiempo de inicio, tiempo que tomo la request y errores |

## Creacion de http handler
Se provee de un builder para inicializar un http handler el cual inicializa endpoints y aplica mapping de errores, middlewares de tracing y logging al server y a los endpoints provistos.

Antes de inicializar el service es necesario crear la configuracion de nuestro endpoint indicando path, nombre de operacion, endpoint, request decoder y response encoder.

    endpointCfgs := []commons.EndpointConfig{
		commons.GET("/drink/beer", "DRINK_BEER", beerEndpoint, decoder, nil),
	}
**Nota 1**: En el ejemplo el ultimo argumento "response encoder" esta nil, de este modo utiliza el response encoder por defecto, este argumento solo debe ser utilizado en caso de requerir personalizar response.

**Nota 2**: En el ejemplo se llama a la funcion `GET` pero tambien estan disponibles `POST`, `PUT`, `PATCH` y `DELETE`.


Luego simplemente creamos nuestro http handler con el builder:

    commons.MakeHTTPHandlerBuilder(logger, endpointCfgs).Build()

Si queremos un manejo de errores personalizado podemos crear un middleware o una funcion personalizada para manejo de errores, para esto se provee de un modificador para el builder `WithCustomErrorWrapper`:

    commons.MakeHTTPHandlerBuilder(logger, endpointCfgs).WithCustomErrorWrapper(errWrapperFunc).Build()

Para habilitar el tracing es necesario utilizar el modificador `WithTracer` para setear un tracer del tipo *Tracer* de "github.com/opentracing/opentracing-go"

    commons.MakeHTTPHandlerBuilder(logger, endpointCfgs).WithTracer(tracer).Build()

Para habilitar instrumentacion de similar modo se puede agregar generando metricas por defecto y utilizando el modificador ``

    metrics := commons.MakeDefaultEndpointMetrics("integration", "beer_service")

    handler := commons.MakeHTTPHandlerBuilder(logger, endpointCfgs).WithMetrics(metrics).Build()

Para personalizar el context a partir de la request se puede usar un `RequestFuncMiddleware` primero se crea un middleware func

    func makeCustomRequestFuncMiddleware() RequestFuncMiddleware {
        return func(next httptransport.RequestFunc) httptransport.RequestFunc {
            return func(ctx context.Context, req *http.Request) context.Context {
                superHeaderValue := req.Header.Get("super-value")
                ctx = context.WithValue(ctx, "special-key", superHeaderValue)
                
                return next(ctx, req)
            }
        }
    }

Para luego setearse en el builder

    handler := commons.MakeHTTPHandlerBuilder(logger, endpointCfgs).WithCustomRequestFuncMiddleware(makeCustomRequestFuncMiddleware()).Build()

Obviamente estas configuraciones se pueden combinar

    commons.MakeHTTPHandlerBuilder(logger, endpointCfgs).WithTracer(tracer).WithCustomErrorWrapper(errWrapperFunc).Build()

## Creacion de Logger
Para crear un nuevo logger se debe llamar a la funcion `ConfigureLogger` indicando el nivel de logger, los niveles disponibles son `debug`, `info`, `warn` y `error`.

    logger := c.ConfigureLogger("info")

## Creacion de Server
Para crea un server (group) se debe utilizar la funcion `CreateServer` pasando como argumento un http.Handler, port y logger.

    g := c.CreateServer(handler, ":8080", logger)
    logger.Log("exit", g.Run())

## Middlewares utilitarios

### Cache
Se ha provisto de una implementacion de endpoint cache middleware utilizando gcache. Para decorar un endpoint con este middleware se debe crear indicando una funcion keyResolver el cual pueda crear una key a partir de la interfaz que se utilice como argumento, cache size y un ttl.

    type Person struct{
        Name string
        Age string
    }

    func resolvePersonKey(in interface{}) string {
        person := in.(*Person)

        return fmt.Sprintf(%s-%d, person.Name, person.Age)
    }

    ...
    ep = commons.MakeCacheEndpointMiddleware(resolvePersonKey, 20, time.Minute * 15)(serviceEndpoint)

## HTTP client
Para crea un nuevo http client se debe utilizar el http client builder.

    commons.MakeHTTPClientBuilder("http://servicehost/requestpath",
			"POST",
			time.Second,
			encodeRequest,
			decodeResponse,
			log.NewNopLogger()).
			Build()

Adicional a esto el builder tiene un metodo para crear un client que utilice opentracing

    commons.MakeHTTPClientBuilder("http://servicehost/requestpath",
			"POST",
			time.Second,
			encodeRequest,
			decodeResponse,
			log.NewNopLogger()).
            WithTracer(tracer).
			Build()

### Client: Default Request encoder
Se puede utilizar un default request encoder para client que convierte cualquier struct a json.

    commons.DefaultRequestEncode

### Client: Default Response decoder
Se puede utilizar un default response decoder que simplemente requiere inicializar el puntero a la estructura.

    func CustomDecode(ctx context.Context, r *http.Response) (interface{}, error) {
		return commons.DefaultDecodeResponse(ctx, r, new(entity.MyClientResponseStruct))
	}

### Client Credentials OAUTH client y Request encode middleware
De ser necesario de autenticarse por OAuth se provee de dos componentes.

Auth Client (o endpoint) permite llamar a una api de autenticacion oauth con grant type `client_credentials` de modo de obtener un token de acceso.

    authClient := commons.MakeClientCredentialsOAUTHClient("http://authme.plz", "clientID", "clientSecret", time.Second, logger)

Donde authClient es un endpoint que llamara a la API de autenticacion y retornara un `TokenResponse`.

Adicional a esto se ha provisto de un middleware para clients http, el cual permite inyectar el header de Authorization utilizando este AuthClient.

    myEncodeRequestFunc := commons.MakeOAuthClientCredentialsRequestEncodeMiddleware(authClient)(commons.DefaultRequestEncode)

## Paginacion
Para Paginacion solo se proveen las estructuras necesarias para implementar paginacion con el fin de ser reutilizadas para implementar paginacion. 

Ver [paging.go](paging.go)

## Validate Middleware
Validar struct, se está validando por los siguientes Tag
- required (Los siguientes campos son requeridos: X-channel)
- oneof (Los siguientes campos no hacen match con el enumerado: X-channel)
- min (Los siguientes campos no cumplen con el mínimo de caracteres: X-channel)
- default "Cualquier otra validate que no tenemos controlado"(Los siguientes campos son inválidos: X-channel) 

Implementar el método de la siguiente forma dentro del endpoint
    myEndpoint = commons.ValidateMiddleware()(myEndpoint)