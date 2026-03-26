# Ejercicio 8

## Objetivo

Permitir que el servidor acepte conexiones y procese mensajes en paralelo, manteniendo la misma semantica funcional de Ej7.

## Solucion (resumen)

1. El servidor crea un thread por conexion entrante para desacoplar aceptacion y procesamiento.
2. Se conserva el protocolo de mensajes de Ej7 (`BATCH`, `FINISH`, `WINNERS_QUERY`).
3. Se agregan mecanismos de sincronizacion para proteger estado compartido y persistencia.

## Mecanismos de sincronizacion utilizados

1. Lock para estado de sorteo (agencias finalizadas y bandera de sorteo realizado).
2. Lock para operaciones sobre persistencia compartida (escrituras de apuestas y lecturas para calcular ganadores).