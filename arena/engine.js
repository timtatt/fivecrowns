const suites = ["B", "R", "Y", "X", "G"];
const numbers = [3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13];

function newDeck(unique = false) {
  const cards = [];

  for (const suite of suites) {
    for (const number of numbers) {
      for (let i = 0; i < (unique ? 1 : 2); i++) {
        cards.push(`${number}-${suite}`);
      }
    }
  }

  for (let i = 0; i < (unique ? 1 : 6); i++) {
    cards.push("*");
  }

  return cards;
}

function decodeCard(card) {
  if (card === "*") {
    return {
      joker: true,
    };
  }

  const [number, suite] = card.split("-");
  return {
    joker: false,
    number: parseInt(number),
    suite,
  };
}

function scoreSequence(sequence) {
  let score = 0;

  console.log("scoring sequence", sequence);
  for (const card of sequence) {
    if (card === "*") {
      score += 25;
      continue;
    }

    const cardDec = decodeCard(card);

    score += cardDec.number;
  }

  return score;
}

function shuffleDeck(cards, seed) {
  let random = new Math.seedrandom(seed);

  const deck = [...cards];

  deck.sort(function () {
    return random() > 0.5 ? 1 : -1;
  });

  return deck;
}

function renderCard(card, { size = "md" } = {}) {
  const decodedCard = decodeCard(card);
  if (decodedCard.joker) {
    return `<button class="pcard pcard-${size} pcard-joker" data-card="${card}"></button>`;
  } else {
    return `<button class="pcard pcard-${size} pcard-${decodedCard.number} pcard-${decodedCard.suite.toLowerCase()}" data-card="${card}"></button>`;
  }
}

function renderCards(cards, opts) {
  let html = "";
  for (const card of cards) {
    html += renderCard(card, opts);
  }

  return html;
}
