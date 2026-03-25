import struct
import datetime

from .utils import Bet

TYPE_BATCH = 1

ACK_SUCCESS = 1
ACK_FAILURE = 0

U8_BYTE_SIZE = 1
U16_BYTE_SIZE = 2
U32_BYTE_SIZE = 4

EMPTY_READ_SIZE = 0

U16_FMT = ">H"
U32_FMT = ">I"

STRING_ENCODING = "utf-8"

BIRTH_YEAR_MIN = 1
BIRTH_MONTH_MIN = 1
BIRTH_DAY_MIN = 1


class ServerProtocol:
    def __init__(self, peer):
        self._peer = peer

    def recv_batch(self):
        msg_type = self._recv_u8()
        if msg_type != TYPE_BATCH:
            raise ValueError(f"unexpected message type: {msg_type}")

        batch_count = self._recv_u8()
        agency = str(self._recv_u8())

        bets = []
        for _ in range(batch_count):
            bets.append(self._recv_bet(agency))
        return bets

    def _recv_bet(self, agency):
        first_name = self._recv_string()
        last_name = self._recv_string()
        year = self._recv_u16()
        month = self._recv_u8()
        day = self._recv_u8()
        birthdate = datetime.date(year, month, day).isoformat()
        document = str(self._recv_u32())
        number = str(self._recv_u16())
        return Bet(agency, first_name, last_name, document, birthdate, number)

    def send_ack(self, ok):
        code = bytes([ACK_SUCCESS]) if ok else bytes([ACK_FAILURE])
        self._peer.sendall(code)

    def _recv_exact(self, size):
        data = bytearray()
        while len(data) < size:
            chunk = self._peer.recv(size - len(data))
            if not chunk:
                if len(data) == EMPTY_READ_SIZE:
                    raise EOFError("connection closed")
                raise ConnectionError("connection closed while receiving data")
            data.extend(chunk)
        return bytes(data)

    def _recv_u8(self):
        return self._recv_exact(U8_BYTE_SIZE)[ACK_FAILURE]

    def _recv_u16(self):
        return struct.unpack(U16_FMT, self._recv_exact(U16_BYTE_SIZE))[ACK_FAILURE]

    def _recv_u32(self):
        return struct.unpack(U32_FMT, self._recv_exact(U32_BYTE_SIZE))[ACK_FAILURE]

    def _recv_string(self):
        raw_size = self._recv_exact(U16_BYTE_SIZE)
        size = struct.unpack(U16_FMT, raw_size)[ACK_FAILURE]
        raw_value = self._recv_exact(size)
        return raw_value.decode(STRING_ENCODING)


def send_ack(peer, ok):
    ServerProtocol(peer).send_ack(ok)
