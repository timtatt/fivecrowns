$(function () {
  const seed = localStorage.getItem("seed") || "123";

  $("#seed").val(seed);

  $("#reseed").on("click", function () {
    localStorage.setItem("seed", $("#seed").val());

    location.reload();
  });

  $("#randomSeed").on("click", function () {
    const seed = Date.now();
    localStorage.setItem("seed", seed);

    location.reload();
  });

  const round = 10;

  const stacks = {
    1: [],
    2: [],
    3: [],
    4: [],
    5: [],
    6: [],
    player: [],
    deck: shuffleDeck(cards),
    discard: [],
  };

  function render() {
    for (let i = 1; i <= 6; i++) {
      console.log("player", bestScore(stacks[i]));
      $(`#hand${i}`).html(renderCards(stacks[i]));
    }

    // $(`#player`).html(renderCards(stacks.player));
    // $("#deck").html(renderCards(stacks.deck.toReversed()));
    // $("#discard").html(renderCards(stacks.discard.toReversed()));
  }

  function deal() {
    for (let j = 0; j < round; j++) {
      for (let i = 1; i <= 6; i++) {
        stacks[i].push(stacks.deck.pop());
      }
    }
  }

  deal();
  render();

  function bestScore(cards) {
    console.log(cards);

    const cardIdx = {};

    // convert hand into an index
    for (const card of cards) {
      if (!cardIdx[card]) {
        cardIdx[card] = 1;
      } else {
        cardIdx[card] += 1;
      }
    }

    // calculate possible sequences
    const best = calculateBest(0, cardIdx, []);
  }

  function calculateBest(depth, cardIdx, parentSequences) {
    const sequences = getSequences(cardIdx);

    if (sequences.length === 0 || depth == 2) {
      // no more sequences left, now we have to try and get rid of some numbers

      const remaining = [];
      let score = 0;

      for (const [card, count] of Object.entries(cardIdx)) {
        if (count > 0) {
          score += decodeCard(card).number * count;
          for (let i = 0; i < count; i++) {
            remaining.push(card);
          }
        }
      }

      out = {
        sequences: parentSequences,
        remaining: remaining,
        score: score,
      };

      console.log(out);

      return out;
    }

    for (const sequence of sequences) {
      const newCardIdx = { ...cardIdx };

      for (const card of sequence.split(":")) {
        newCardIdx[card]--;
      }

      const best = calculateBest(depth + 1, newCardIdx, [
        ...parentSequences,
        sequence,
      ]);
    }
  }

  function getSequences(cardIdx) {
    const sequences = {};

    for (const card of Object.keys(cardIdx)) {
      if (cardIdx[card] === 0) {
        // dont calculate anything with 0 cards
        break;
      }

      if (card === "*") {
        // no need to calculate joker
        continue;
      }

      cardDec = decodeCard(card);

      const sequence = [card];

      // available optimisation: existing sequences with this card

      // calculate the largest same suite sequence for this card

      // traverse left of the current card
      for (let i = cardDec.number - 1; i >= 3; i--) {
        const siblingCard = `${i}-${cardDec.suite}`;

        if (cardIdx[siblingCard] && cardIdx[siblingCard] > 0) {
          sequence.unshift(siblingCard);
        } else {
          break;
        }
      }

      // traverse right of current card
      for (let i = cardDec.number + 1; i <= 13; i++) {
        const siblingCard = `${i}-${cardDec.suite}`;

        if (cardIdx[siblingCard] && cardIdx[siblingCard] > 0) {
          sequence.push(siblingCard);
        } else {
          break;
        }
      }

      if (sequence.length >= 3) {
        sequences[sequence.join(":")] = 1;
      }

      // calculate the largest number set possible
      const cardSet = [];
      for (const suite of suites) {
        const siblingCard = `${cardDec.number}-${suite}`;
        const count = cardIdx[siblingCard] ?? 0;

        for (let i = 0; i < count; i++) {
          cardSet.push(siblingCard);
        }
      }

      if (cardSet.length >= 3) {
        sequences[sortCards(cardSet).join(":")] = 1;
      }
    }

    // permute the sequences
    for (const sequence of Object.keys(sequences)) {
      const permutations = permuteSequence(sequence);

      for (const permutation of permutations) {
        sequences[permutation] = 1;
      }
    }

    return Object.keys(sequences);
  }

  function permuteSequence(sequence) {
    const cards = sequence.split(":");

    if (cards.length === 3) {
      return [sequence];
    }

    const permutations = [];

    if (isSet(cards)) {
      return permuteSet(cards).values();
    } else {
      // this is a run of cards
      for (let l = 3; l <= cards.length; l++) {
        for (let o = 0; o <= cards.length - l; o++) {
          permutations.push(cards.slice(o, o + l));
        }
      }

      return permutations;
    }
  }

  function permuteSet(cards) {
    let permutations = new Set([cards.join(":")]);

    if (cards.length > 3) {
      for (let i = 0; i < cards.length; i++) {
        permutations = permutations.union(
          permuteSet([...cards].toSpliced(i, 1)),
        );
      }
    }

    return permutations;
  }

  function isSet(cards) {
    return decodeCard(cards[0]).number === decodeCard(cards[1]).number;
  }

  $("#drawDeck").on("click", function () {
    const card = stacks.deck.pop();
    stacks.player.push(card);

    stacks.discard.pop();

    render();
  });

  $("#player").on("click", ".pcard", function (e) {
    const cardIdx = $(e.currentTarget).index();

    stacks.player.splice(cardIdx, 1);

    // stacks.discard.push(discarded[0]);

    render();
  });

  $("#discardTop").on("click", ".pcard", function () {
    const card = stacks.discard.pop();

    stacks.player.push(card);

    // get a new deck choice
    stacks.deck.pop();

    render();
  });
});
