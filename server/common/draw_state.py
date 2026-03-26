import threading


class DrawState:
    def __init__(self, expected_agencies):
        self._expected_agencies = expected_agencies
        self._finished_agencies = set()
        self._draw_done = False
        self._lock = threading.Lock()

    def mark_agency_finished(self, agency):
        with self._lock:
            self._finished_agencies.add(agency)
            if self._draw_done:
                return False
            if len(self._finished_agencies) < self._expected_agencies:
                return False

            self._draw_done = True
            return True

    def is_draw_done(self):
        with self._lock:
            return self._draw_done
