# General :

- Optimize / improve the engine.

# Specifics :

- Fix the remaining move count on the transposition table version (v1.4).

> For example, the position "61612776255435661442454452" is returning 5 remaining moves (incorrect value), but when you play the best move suggested (column 1), the remaining moves becomes 12 (which is the correct value).
