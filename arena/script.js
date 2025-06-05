$(function () {
  const deck = newDeck(true);

  let hand = [];

  $("#allCards").html(renderCards(deck));

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

  function sequenceCodeToArray(sequenceCode) {
    return sequenceCode.split(":").filter((c) => c !== "");
  }

  function renderDiscardPreview() {
    const discard = sequenceCodeToArray($("#discard").val());

    $("#discardPreview").html(renderCards(discard, { size: "xs" }));
  }

  function render() {
    $("#hand").html(renderCards(hand));

    $("#sequenceCode").val(hand.join(":"));

    renderDiscardPreview();

    $("#round").val(hand.length);
  }

  render();
});
