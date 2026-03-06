# Viking Chess

Hnefatafl is a Scandinavian board game somewhat similar to Chess. It
has some unusual and intriguing features:

* Sides having differing objectives (escape vs capture)
* Sides have differing numbers of pieces
* There is only one king
* Sides have different piece types

# Berserk

The berserk variant adds new pieces and rules, most notably chaining
successive captures together in a single turn (Berserk moves). See
this for the complete rules and history of the variant:

https://aagenielsen.dk/berserk_rules.php

# Monte Carlo Tree Search

This implementation has a very rudimentary implementation of MCTS that
can be used to come up with opponent moves. To get half-decent move
suggestions you must run hundreds of random games per move, which
makes it very slow.

I have noted that because of Hnefatafl's asymmetry, purely random
moves heavily favors the defender. Perhaps that's not surprising --
escaping requires moving a single piece while capturing requires
coordination of 2 to 4 pieces.

If, instead, captures are made when possible then the outcomes balance
out. However, this almost certainly biases the search toward moves
that may not be optimal.
