# Ejercicio 7

## Objetivo

Coordinar el cierre de carga de apuestas entre agencias, ejecutar el sorteo una vez completada la etapa de carga y permitir la consulta de ganadores por agencia.

## Solucion (resumen)

1. El cliente envia apuestas por batch (heredado de Ej6).
2. Al finalizar el envio, cada cliente notifica al servidor con un mensaje `FINISH`.
3. Luego, cada cliente consulta ganadores con `WINNERS_QUERY`.
4. El servidor responde `pending` hasta que se alcanza la cantidad esperada de agencias finalizadas.
5. Cuando todas las agencias esperadas notifican fin, el servidor registra `action: sorteo | result: success`.
6. Desde ese momento, responde a cada agencia solo con sus ganadores.

## Aspectos clave de la implementacion

1. Se reutiliza un protocolo binario con tipos de mensaje para batch, fin de carga y consulta de ganadores.
2. Se usa estado de sorteo para controlar cuando habilitar respuestas definitivas.
3. La consulta de ganadores se filtra por agencia, evitando broadcast global.
