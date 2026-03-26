import socket
import logging
import signal
import threading

from .utils import store_bets
from .draw_state import DrawState
from .winners_service import WinnersService
from .protocol import (
    ServerProtocol,
    TYPE_BATCH,
    TYPE_FINISH,
    TYPE_WINNERS_QUERY,
    ACK_SUCCESS,
)

FAIL_BATCH_COUNT = 0
CLIENT_THREAD_JOIN_TIMEOUT_SECONDS = 1


class Server:
    def __init__(self, port, listen_backlog, expected_agencies):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self._running = True
        self._draw_state = DrawState(expected_agencies)
        self._winners_service = WinnersService()
        self._storage_lock = threading.Lock()
        self._threads_lock = threading.Lock()
        self._client_threads = set()
        self.__register_signals()

    def __register_signals(self):
        signal.signal(signal.SIGTERM, self.shutdown)
        signal.signal(signal.SIGINT, self.shutdown)

    def shutdown(self, _signum=None, _frame=None):
        if not self._running:
            return

        self._running = False
        logging.info('action: shutdown | result: in_progress | resource: server_socket')
        try:
            self._server_socket.close()
            logging.info('action: shutdown | result: success | resource: server_socket')
        except OSError as e:
            logging.error(f'action: shutdown | result: fail | resource: server_socket | error: {e}')

        self.__join_client_threads()

    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """
        while self._running:
            client_sock = self.__accept_new_connection()
            if client_sock is None:
                continue

            thread = threading.Thread(target=self.__run_client_connection, args=(client_sock,), daemon=True)
            self.__register_client_thread(thread)
            thread.start()

    def __run_client_connection(self, client_sock):
        try:
            self.__handle_client_connection(client_sock)
        finally:
            self.__unregister_current_thread()

    def __register_client_thread(self, thread):
        with self._threads_lock:
            self._client_threads.add(thread)

    def __unregister_current_thread(self):
        current_thread = threading.current_thread()
        with self._threads_lock:
            self._client_threads.discard(current_thread)

    def __join_client_threads(self):
        with self._threads_lock:
            threads_to_join = list(self._client_threads)

        for thread in threads_to_join:
            thread.join(CLIENT_THREAD_JOIN_TIMEOUT_SECONDS)

    def __handle_client_connection(self, client_sock):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        protocol = ServerProtocol(client_sock)
        try:
            while self._running:
                try:
                    msg_type = protocol.recv_message_type()
                except EOFError:
                    break
                except (OSError, ValueError, ConnectionError) as e:
                    if self._running:
                        logging.error(f'action: apuesta_recibida | result: fail | cantidad: {FAIL_BATCH_COUNT} | error: {e}')
                    try:
                        protocol.send_ack(False)
                    except OSError:
                        pass
                    break

                if msg_type == TYPE_BATCH:
                    if not self.__handle_batch_message(protocol):
                        break
                    continue

                if msg_type == TYPE_FINISH:
                    if not self.__handle_finish_message(protocol):
                        break
                    continue

                if msg_type == TYPE_WINNERS_QUERY:
                    if not self.__handle_winners_query(protocol):
                        break
                    continue

                try:
                    protocol.send_ack(False)
                except OSError:
                    pass
                break
        finally:
            try:
                client_sock.close()
            except OSError:
                pass

    def __handle_batch_message(self, protocol):
        try:
            batch_count, agency = protocol.recv_batch_payload_header()
            bets = protocol.recv_batch_payload(batch_count, agency)
            with self._storage_lock:
                store_bets(bets)
            logging.info(f"action: apuesta_recibida | result: success | cantidad: {len(bets)}")
            protocol.send_ack(True)
            return True
        except (OSError, ValueError, ConnectionError) as e:
            logging.error(f'action: apuesta_recibida | result: fail | cantidad: {FAIL_BATCH_COUNT} | error: {e}')
            try:
                protocol.send_ack(False)
            except OSError:
                pass
            return False

    def __handle_finish_message(self, protocol):
        try:
            agency = protocol.recv_agency()
            protocol.send_ack(True)

            if self._draw_state.mark_agency_finished(agency):
                logging.info('action: sorteo | result: success')
            return True
        except (OSError, ValueError, ConnectionError):
            try:
                protocol.send_ack(False)
            except OSError:
                pass
            return False

    def __handle_winners_query(self, protocol):
        try:
            agency = protocol.recv_agency()
            if not self._draw_state.is_draw_done():
                protocol.send_pending_response()
                return True

            with self._storage_lock:
                winners = self._winners_service.find_winner_documents_by_agency(agency)
            protocol.send_winners_response(ACK_SUCCESS, winners)
            return True
        except (OSError, ValueError, ConnectionError):
            return False

    def __accept_new_connection(self):
        """
        Accept new connections

        Function blocks until a connection to a client is made.
        Then connection created is printed and returned
        """

        # Connection arrived
        if not self._running:
            return None

        logging.info('action: accept_connections | result: in_progress')
        try:
            c, addr = self._server_socket.accept()
            logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
            return c
        except OSError:
            return None
