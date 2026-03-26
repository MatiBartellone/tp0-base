from .utils import load_bets, has_won


class WinnersService:
    def find_winner_documents_by_agency(self, agency):
        winners = []
        for bet in load_bets():
            if str(bet.agency) != agency:
                continue
            if not has_won(bet):
                continue
            winners.append(int(bet.document))

        return winners
