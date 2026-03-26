# Ejercicio 3

## Objetivo

Implementar y validar un echo server accesible por TCP, junto con una verificacion de salud que permita distinguir respuestas correctas e incorrectas.

## Solucion (resumen)

1. Se levanta el servidor en Docker y se expone el puerto para pruebas.
2. El script de validacion (`validar-echo-server.sh`) envia un mensaje de prueba y compara la respuesta.
3. Si el servidor responde exactamente lo esperado, se registra `action: test_echo_server | result: success`.
4. El servidor loguea recepcion de mensaje con `action: receive_message | result: success`.
5. Tambien se valida el comportamiento frente a un servidor mock no saludable, donde el healthcheck debe fallar.

## Aspectos clave de la implementacion

1. El chequeo funcional principal se apoya en netcat/script para validar ida y vuelta del mensaje (echo).
2. La observabilidad se hace por logs estructurados de resultado en cliente de prueba y servidor.
3. Se prueban ambos escenarios: servidor saludable (pass) y servidor no saludable (fail esperado).
