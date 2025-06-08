$(function () {
  const deck = newDeck(true);

  let hand = [];

  const cachedHand = localStorage.getItem("hand");

  if (cachedHand) {
    hand = sequenceCodeToArray(cachedHand);
  }

  const action = localStorage.getItem("action");
  $("#action").val(action !== "" ? action : "score");
  $("#discard").val(localStorage.getItem("discard"));

  $("#allCards").html(renderCards(deck));

  calculate();

  $("#shuffle").on("click", async function () {
    const round = $("#round").val();

    const newHand = [];

    for (let i = 0; i < round; i++) {
      const randIdx = Math.floor(Math.random() * deck.length);

      newHand.push(deck[randIdx]);
    }

    hand = newHand;

    render();

    calculate();
  });

  $("#allCards").on("click", ".pcard", function (e) {
    const cardEl = $(e.currentTarget);

    const card = cardEl.data("card");

    hand.push(card);

    render();
  });

  $("#hand").on("click", ".pcard", function (e) {
    const cardEl = $(e.currentTarget);

    const card = cardEl.data("card");

    const cardIdx = hand.indexOf(card);
    hand.splice(cardIdx, 1);

    render();
  });

  $("#discard").on("change", function (e) {
    renderDiscardPreview();
  });

  $("#copySequenceCode").on("click", function (e) {
    const sequenceCode = $("#sequenceCode").val();
    navigator.clipboard.writeText(sequenceCode);
  });

  $("#sequenceCode").on("change", function (e) {
    hand = sequenceCodeToArray($(e.currentTarget).val());

    render();
  });

  async function calculate() {
    const action = $("#action").val();
    localStorage.setItem("action", action);

    const discard = $("#discard").val();
    localStorage.setItem("discard", discard);

    const response = await fetch("http://localhost:3000/bots/grugbot", {
      method: "POST",
      body: JSON.stringify({
        action,
        discard: discard.split(":"),
        playerCount: 1,
        lastTurn: false,
        newestCard: "",
        hand,
        round: hand.length - (action === "discard" ? 1 : 0),
      }),
    });

    const botResponse = await response.json();

    console.log("received bot response", botResponse);

    renderBotResponse(botResponse);
  }

  $("#calculate").on("click", async function () {
    await calculate();
  });

  function sequenceCodeToArray(sequenceCode) {
    return sequenceCode.split(":").filter((c) => c !== "");
  }

  function renderDiscardPreview() {
    const discard = sequenceCodeToArray($("#discard").val());

    $("#discardPreview").html(renderCards(discard, { size: "xs" }));
  }

  function renderBotResponse(res) {
    let summary = "";
    let sequences = "";

    let score = 0;

    if (res.sequences) {
      for (const seq of res.sequences) {
        if (seq.length < 3) {
          score += scoreSequence(seq);
        }

        sequences += `<li class="list-group-item"><div class="hand">${renderCards(seq, { size: "sm" })}</div></li>`;
      }
    }

    if (res.action === "draw") {
      summary = `Stack: ${res.stack}`;
    } else if (res.action === "score") {
      summary = `Score: ${score}, Flop: ${res.flop}`;
    } else if (res.action === "discard") {
      summary = `Score: ${score}, Flop: ${res.flop}, Discarded: ${res.card}`;
    }

    $("#botResponseSummary").html(summary);
    $("#botResponseSequences").html(sequences);
  }

  renderBotResponse({
    sequences: [["10-R", "9-R"], ["*"], ["7-Y"]],
    action: "draw",
    stack: "deck",
  });

  function render() {
    console.log(hand);
    $("#hand").html(renderCards(hand));

    $("#sequenceCode").val(hand.join(":"));

    renderDiscardPreview();

    $("#round").val(hand.length);

    // save the hand to local storage
    localStorage.setItem("hand", hand.join(":"));
  }

  render();
});
