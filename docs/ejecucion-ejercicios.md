# Ejecucion de Ejercicios

Este documento describe como ejecutar manualmente cada ejercicio, considerando que el flujo cambia entre ramas.

## Requisitos

1. Docker y Docker Compose instalados.
2. Make disponible en el entorno.
3. Ubicarse en la raiz de `tp0-base`.

## Ejercicio 1

Objetivo: generar el archivo compose con cantidad variable de clientes.

1. `./generar-compose.sh docker-compose-dev.yaml 1`
2. `./generar-compose.sh docker-compose-dev.yaml 3`
3. Revisar que el archivo generado tenga `client1`, `client2`, `client3` segun corresponda.

## Ejercicio 2

Objetivo: inyectar configuraciones por volumen sin reconstruir imagen.

1. `./generar-compose.sh docker-compose-dev.yaml 1`
2. `make docker-compose-up`
3. Modificar `client/config.yaml` o `server/config.ini`.
4. `make docker-compose-down`
5. `make docker-compose-up`
6. Verificar en logs que se toman los nuevos valores sin rebuild.

## Ejercicio 3

Objetivo: validar echo server con script dedicado.

1. `./generar-compose.sh docker-compose-dev.yaml 1`
2. `make docker-compose-up`
3. `./validar-echo-server.sh`
4. Confirmar salida `action: test_echo_server | result: success`.
5. `make docker-compose-down`

## Ejercicio 4

Objetivo: cierre graceful con SIGTERM.

1. `./generar-compose.sh docker-compose-dev.yaml 1`
2. `make docker-compose-up`
3. `make docker-compose-down`
4. Verificar en logs mensajes de shutdown y cierre correcto de recursos.

## Ejercicio 5

Objetivo: enviar una apuesta y persistirla en servidor.

1. `./generar-compose.sh docker-compose-dev.yaml 1`
2. Configurar variables de apuesta para cliente en compose o entorno.
3. `make docker-compose-up`
4. Verificar logs:
5. Cliente: `action: apuesta_enviada | result: success | ...`
6. Servidor: `action: apuesta_almacenada | result: success | ...`
7. `make docker-compose-down`

## Ejercicio 6

Objetivo: envio por batches desde dataset por agencia.

1. Asegurar archivo `.data/agency-1.csv`.
2. `./generar-compose.sh docker-compose-dev.yaml 1`
3. Ajustar `batch.maxAmount` en `client/config.yaml`.
4. `make docker-compose-up`
5. Verificar logs de servidor `action: apuesta_recibida | result: success | cantidad: X`.
6. `make docker-compose-down`

## Ejercicio 7

Objetivo: coordinar sorteo y consulta de ganadores por agencia.

1. Preparar datasets `.data/agency-N.csv` para agencias activas.
2. `./generar-compose.sh docker-compose-dev.yaml <N_CLIENTES>`
3. `make docker-compose-up`
4. Verificar logs:
5. Servidor: `action: sorteo | result: success` cuando finalizan las agencias esperadas.
6. Cliente: `action: consulta_ganadores | result: success | cant_ganadores: X`.
7. `make docker-compose-down`

## Ejercicio 8

Objetivo: procesar conexiones y mensajes en paralelo en servidor.

1. Repetir el flujo de ejecucion de Ejercicio 7.
2. Ejecutar con mas de un cliente simultaneo (`N_CLIENTES > 1`).
3. Verificar que el servidor mantiene comportamiento funcional de Ej7, con procesamiento concurrente y sin inconsistencias en persistencia.
4. `make docker-compose-down`

## Nota sobre documentacion

1. `docs/implementacion-general.md` contiene protocolo y arquitectura general final.
2. `docs/ej8.md` resume decisiones especificas del ejercicio 8.
