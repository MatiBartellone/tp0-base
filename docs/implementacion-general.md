# Implementacion General Final

Esta seccion resume la arquitectura final alcanzada y sirve como base comun para todas las ramas. Luego, cada rama puede complementar este documento con sus particularidades por ejercicio.

## Protocolo Final

Se implemento un protocolo binario propio sobre TCP, con serializacion compacta y lecturas/escrituras robustas para evitar short read y short write.

### Tipos y tamanos base

1. u8: 1 byte
2. u16: 2 bytes
3. u32: 4 bytes
4. Endianness: big-endian para campos multibyte

### Tipos de mensaje cliente -> servidor

1. `TYPE_BATCH (1)`
2. `TYPE_FINISH (2)`
3. `TYPE_WINNERS_QUERY (3)`

### Estructura de TYPE_BATCH

Header de batch (3 bytes):

1. type: u8 (1 byte) -> valor 1
2. cantidad: u8 (1 byte) -> cantidad de apuestas del batch
3. agencia: u8 (1 byte) -> id de agencia

Cuerpo de batch: repeticion de N apuestas serializadas, donde N es el campo cantidad.

Estructura de una apuesta serializada:

1. first_name_len: u16 (2 bytes)
2. first_name: bytes UTF-8 (first_name_len bytes)
3. last_name_len: u16 (2 bytes)
4. last_name: bytes UTF-8 (last_name_len bytes)
5. birth_year: u16 (2 bytes)
6. birth_month: u8 (1 byte)
7. birth_day: u8 (1 byte)
8. document: u32 (4 bytes)
9. number: u16 (2 bytes)

Tamano fijo minimo por apuesta (sin strings): 14 bytes.

### Estructura de TYPE_FINISH

Mensaje de 2 bytes:

1. type: u8 (1 byte) -> valor 2
2. agencia: u8 (1 byte)

### Estructura de TYPE_WINNERS_QUERY

Mensaje de 2 bytes:

1. type: u8 (1 byte) -> valor 3
2. agencia: u8 (1 byte)

### Respuestas servidor -> cliente

ACK simple (usado para TYPE_BATCH y TYPE_FINISH):

1. ack: u8 (1 byte)
2. ack = 1: exito
3. ack = 0: error

Respuesta a TYPE_WINNERS_QUERY:

1. status: u8 (1 byte)
2. winners_count: u16 (2 bytes)
3. winners: lista de winners_count documentos, cada uno como u32 (4 bytes)

Estados de status:

1. status = 1: success
2. status = 2: pending (sorteo aun no habilitado)
3. status = 0: failure

Formula de tamano de respuesta de ganadores:

1. 1 byte (status) + 2 bytes (count) + 4 * winners_count bytes (documentos)

### Nota de robustez del transporte

El protocolo se implementa con loops de lectura y escritura exacta para garantizar recepcion/envio completo de cada campo, incluso cuando el socket entrega fragmentos parciales.

La construccion de mensajes en cliente se centraliza en un `ProtocolBuilder` para evitar duplicacion de constantes y reducir errores de framing.

## Logica de Cliente

El cliente sigue un flujo en etapas:

1. Lee su archivo `.data/agency-N.csv`.
2. Agrupa apuestas en batches segun `batch.maxAmount` y limite de tamano.
3. Envia cada batch y espera ACK.
4. Notifica finalizacion de carga (`FINISH`).
5. Consulta ganadores por agencia (`WINNERS_QUERY`) hasta obtener resultado final.

La implementacion se fue modularizando en archivos por responsabilidad (orquestacion, envio por lotes, consulta de ganadores, protocolo, serializacion y logs), buscando bajo acoplamiento y funciones pequeñas.

## Logica de Servidor

El servidor:

1. Acepta conexiones y procesa mensajes por tipo.
2. Persiste apuestas recibidas por batch.
3. Registra notificaciones de fin de agencias.
4. Cuando se alcanza la cantidad esperada de agencias, habilita el sorteo.
5. Responde consultas de ganadores por agencia con datos filtrados.

Para concurrencia se agrego procesamiento paralelo de conexiones y sincronizacion explicita sobre estado compartido y acceso a persistencia.

## Evolucion Entre Ramas

El desarrollo se hizo de manera incremental, ejercicio por ejercicio, conservando compatibilidad con los requisitos de cada etapa.

1. Se partio de una base simple cliente-servidor.
2. Se incorporo protocolo de negocio y serializacion.
3. Se agrego batching y control de ACK.
4. Se sumo coordinacion de sorteo y consulta de ganadores.
5. Se llego a una version mas madura en `ej8`, con mayor modularizacion, mejor separacion de responsabilidades, limpieza de logs y soporte de concurrencia.

En otras palabras, la version de `ej8` representa la forma mas limpia y modular alcanzada dentro de la evolucion del TP, y sirve como referencia para retropropagar mejoras de diseño al resto de las ramas.
