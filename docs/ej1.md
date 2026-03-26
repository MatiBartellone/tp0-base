# Ejercicio 1

## Objetivo

Levantar una topologia cliente-servidor con cantidad variable de clientes y validar que cada cliente complete su loop de mensajes sin fallos.

## Solucion (resumen)

1. Se genera `docker-compose-dev.yaml` con N clientes usando el script de compose.
2. Se configura la cantidad de iteraciones (`loop_amount`) que cada cliente debe ejecutar.
3. Se construyen imagenes y se levantan servicios.
4. Cada cliente intercambia mensajes con el servidor y registra recepciones en logs.
5. Al finalizar, cada cliente debe completar su checklist de mensajes esperados y salir en success.

## Aspectos clave de la implementacion

1. El script `generar-compose.sh` habilita escalado de clientes sin duplicar definiciones manuales.
2. La validacion se basa en logs estructurados (`receive_message`, `loop_finished`, `exit`).
3. El servidor permanece activo aunque los clientes terminen su ejecucion.
