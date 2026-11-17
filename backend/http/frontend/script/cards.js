(function () {
  var form = document.getElementById("cards-form");
  if (!form) {
    return;
  }

  var seats = ["N", "E", "S", "W"];
  var seatNames = { N: "North", E: "East", S: "South", W: "West" };
  var suits = [
    { code: "S", symbol: "♠", cls: "club", handCls: "hand-suit-S" },
    { code: "H", symbol: "♥", cls: "red", handCls: "hand-suit-H" },
    { code: "D", symbol: "♦", cls: "red", handCls: "hand-suit-D" },
    { code: "C", symbol: "♣", cls: "club", handCls: "hand-suit-C" }
  ];
  var ranks = ["A", "K", "Q", "J", "T", "9", "8", "7", "6", "5", "4", "3", "2"];

  var phases = [{ type: "intro" }];
  var seatIndex;
  var suitIndex;
  for (seatIndex = 0; seatIndex < seats.length; seatIndex++) {
    if (seatIndex > 0) {
      phases.push({ type: "next", seatIndex: seatIndex });
    }
    for (suitIndex = 0; suitIndex < 4; suitIndex++) {
      phases.push({ type: "suit", seatIndex: seatIndex, suitIndex: suitIndex });
    }
  }
  phases.push({ type: "summary" });

  var stepIndex = 0;
  var westReady = false;
  var selections = { N: [[], [], [], []], E: [[], [], [], []], S: [[], [], [], []], W: [[], [], [], []] };

  var titleEl = document.getElementById("cards-title");
  var countEl = document.getElementById("cards-count");
  var gridEl = document.getElementById("cards-grid");
  var backBtn = document.getElementById("cards-back");
  var nextBtn = document.getElementById("cards-next");

  function cardKey(suitIndex, rank) {
    return suitIndex + ":" + rank;
  }

  function usedCardsBefore(seatIndex) {
    var used = {};
    var s;
    var suit;
    var i;
    for (s = 0; s < seatIndex; s++) {
      for (suit = 0; suit < 4; suit++) {
        for (i = 0; i < selections[seats[s]][suit].length; i++) {
          used[cardKey(suit, selections[seats[s]][suit][i])] = true;
        }
      }
    }
    return used;
  }

  function seatCardCount(seat) {
    var total = 0;
    var suit;
    for (suit = 0; suit < 4; suit++) {
      total += selections[seat][suit].length;
    }
    return total;
  }

  function remainingForWest() {
    var used = usedCardsBefore(3);
    var remaining = [[], [], [], []];
    var suit;
    var r;
    for (suit = 0; suit < 4; suit++) {
      for (r = 0; r < ranks.length; r++) {
        if (!used[cardKey(suit, ranks[r])]) {
          remaining[suit].push(ranks[r]);
        }
      }
    }
    return remaining;
  }

  function ensureWestSelection() {
    if (westReady) {
      return;
    }
    selections.W = remainingForWest();
    westReady = true;
  }

  function clearWestSelection() {
    westReady = false;
    selections.W = [[], [], [], []];
  }

  function shouldKeepWest(phase) {
    if (!phase) {
      return false;
    }
    if (phase.type === "summary") {
      return true;
    }
    return phase.type === "suit" && seats[phase.seatIndex] === "W";
  }

  function setHidden(el, on) {
    if (!el) {
      return;
    }
    if (on) {
      el.setAttribute("hidden", "hidden");
    } else {
      el.removeAttribute("hidden");
    }
  }

  function sortRanks(list) {
    list.sort(function (a, b) {
      return ranks.indexOf(a) - ranks.indexOf(b);
    });
  }

  function displayRank(rank) {
    return rank === "T" ? "10" : rank;
  }

  function toggleRank(suitIndex, rank, checked) {
    var phase = phases[stepIndex];
    if (!phase || phase.type !== "suit") {
      return;
    }
    var seat = seats[phase.seatIndex];
    var list = selections[seat][suitIndex];
    var idx = list.indexOf(rank);
    if (checked && idx < 0) {
      list.push(rank);
      sortRanks(list);
    }
    if (!checked && idx >= 0) {
      list.splice(idx, 1);
    }
    updateCount();
  }

  function seatsProgressHtml(currentSeatIndex) {
    var seatsHtml = "";
    var i;
    for (i = 0; i < seats.length; i++) {
      var seatClass = "cards-seat";
      if (i < currentSeatIndex) {
        seatClass += " cards-seat-done";
      } else if (i === currentSeatIndex) {
        seatClass += " cards-seat-current";
      } else {
        seatClass += " cards-seat-todo";
      }
      seatsHtml += '<span class="' + seatClass + '">' + seatNames[seats[i]] + "</span>";
    }
    return seatsHtml;
  }

  function updateCount() {
    var phase = phases[stepIndex];
    if (!phase || phase.type !== "suit") {
      countEl.innerHTML = "";
      return;
    }
    var seat = seats[phase.seatIndex];
    var suitCount = selections[seat][phase.suitIndex].length;
    var total = seatCardCount(seat);
    var text = "Kolor: " + suitCount + " · ręka: " + total + " / 13";
    if (seat === "W") {
      text += " · karty dobrane automatycznie";
    }
    countEl.innerHTML = text;
  }

  function formatHandHtml(seat) {
    var parts = [];
    var suit;
    for (suit = 0; suit < 4; suit++) {
      var list = selections[seat][suit];
      if (!list.length) {
        continue;
      }
      var ranksHtml = [];
      var i;
      for (i = 0; i < list.length; i++) {
        ranksHtml.push(displayRank(list[i]));
      }
      parts.push(
        '<span class="hand-suit ' + suits[suit].handCls + '">' +
          suits[suit].symbol +
          " " +
          ranksHtml.join(" ") +
          "</span>"
      );
    }
    return parts.join(" ");
  }

  function renderInfoPanel(messageHtml) {
    gridEl.innerHTML = '<div class="cards-info">' + messageHtml + "</div>";
  }

  function renderIntro() {
    titleEl.innerHTML =
      '<span class="cards-heading-label">Zapis kart</span>' +
      '<span class="cards-heading-seats">' +
      seatsProgressHtml(0) +
      "</span>";
    countEl.innerHTML = "";
    renderInfoPanel(
      "<p>Zapisujemy karty rozdania.</p>" +
        "<p class=\"cards-info-emphasis\">Zaczynamy od <strong>North</strong>.</p>" +
        "<p class=\"cards-info-hint\">Jeśli w danym kolorze nie ma kart, kliknij <strong>Dalej</strong>, żeby go pominąć.</p>"
    );
    setHidden(backBtn, true);
    nextBtn.innerHTML = "Dalej";
  }

  function renderNext(phase) {
    var seat = seats[phase.seatIndex];
    var name = seatNames[seat];
    var detail =
      seat === "W"
        ? "<p>Karty West są już wstępnie zaznaczone — to pozostałe karty po innych rękach. Sprawdź je i popraw, jeśli trzeba.</p>"
        : "<p>Teraz wpisz karty " + name + ".</p>";
    titleEl.innerHTML =
      '<span class="cards-heading-label">Kolejna ręka</span>' +
      '<span class="cards-heading-seats">' +
      seatsProgressHtml(phase.seatIndex) +
      "</span>";
    countEl.innerHTML = "";
    renderInfoPanel(
      "<p class=\"cards-info-emphasis\">Następny: <strong>" + name + "</strong></p>" + detail
    );
    setHidden(backBtn, false);
    nextBtn.innerHTML = "Dalej";
  }

  function renderSummary() {
    ensureWestSelection();
    var html = '<div class="hand-summary">';
    var i;
    for (i = 0; i < seats.length; i++) {
      var seat = seats[i];
      html +=
        '<div class="hand-summary-seat">' +
        '<div class="hand-summary-name">' +
        seatNames[seat] +
        ":</div>" +
        '<div class="hand-summary-cards">' +
        formatHandHtml(seat) +
        "</div></div>";
    }
    html += "</div>";
    titleEl.innerHTML =
      '<span class="cards-heading-label">Podsumowanie</span>' +
      '<span class="cards-heading-seats">' +
      seatsProgressHtml(seats.length) +
      "</span>";
    countEl.innerHTML = "Sprawdź karty przed zapisem.";
    gridEl.innerHTML = html;
    setHidden(backBtn, false);
    nextBtn.innerHTML = "Zapisz";
  }

  function renderSuit(phase) {
    var seat = seats[phase.seatIndex];
    var suit = suits[phase.suitIndex];
    var used = usedCardsBefore(phase.seatIndex);
    var selected = {};
    var i;

    if (seat === "W") {
      ensureWestSelection();
    }

    for (i = 0; i < selections[seat][phase.suitIndex].length; i++) {
      selected[selections[seat][phase.suitIndex][i]] = true;
    }

    titleEl.innerHTML =
      '<span class="cards-heading-label">Karty <span class="suit ' +
      suit.cls +
      '">' +
      suit.symbol +
      "</span></span>" +
      '<span class="cards-heading-seats">' +
      seatsProgressHtml(phase.seatIndex) +
      "</span>";
    gridEl.innerHTML = "";

    var rows = [ranks.slice(0, 7), ranks.slice(7)];
    var rowIndex;
    for (rowIndex = 0; rowIndex < rows.length; rowIndex++) {
      var row = document.createElement("div");
      row.className = "card-row";
      for (i = 0; i < rows[rowIndex].length; i++) {
        var rank = rows[rowIndex][i];
        var key = cardKey(phase.suitIndex, rank);
        var id = "card-" + seat + "-" + suit.code + "-" + rank;
        var wrap = document.createElement("div");
        wrap.className = "card-choice";

        var input = document.createElement("input");
        input.type = "checkbox";
        input.id = id;
        input.value = rank;
        input.checked = !!selected[rank];
        if (used[key]) {
          input.disabled = true;
          wrap.className = "card-choice card-used";
        }

        var label = document.createElement("label");
        label.htmlFor = id;
        label.className = suit.cls;
        label.innerHTML =
          displayRank(rank) + '<span class="suit ' + suit.cls + '">' + suit.symbol + "</span>";

        (function (sIdx, rk, inp) {
          input.onchange = function () {
            toggleRank(sIdx, rk, inp.checked);
          };
        })(phase.suitIndex, rank, input);

        wrap.appendChild(input);
        wrap.appendChild(label);
        row.appendChild(wrap);
      }
      gridEl.appendChild(row);
    }
    setHidden(backBtn, false);
    nextBtn.innerHTML = "Dalej";
    updateCount();
  }

  function renderStep() {
    var phase = phases[stepIndex];
    if (phase.type === "intro") {
      renderIntro();
      return;
    }
    if (phase.type === "next") {
      renderNext(phase);
      return;
    }
    if (phase.type === "summary") {
      renderSummary();
      return;
    }
    renderSuit(phase);
  }

  function serializeHand(seat) {
    var parts = [];
    var suit;
    for (suit = 0; suit < 4; suit++) {
      parts.push(selections[seat][suit].join(""));
    }
    return parts.join(".");
  }

  function finish() {
    ensureWestSelection();
    document.getElementById("hand-north").value = serializeHand("N");
    document.getElementById("hand-east").value = serializeHand("E");
    document.getElementById("hand-south").value = serializeHand("S");
    document.getElementById("hand-west").value = serializeHand("W");
    form.submit();
  }

  backBtn.onclick = function () {
    if (stepIndex === 0) {
      return;
    }
    stepIndex -= 1;
    if (!shouldKeepWest(phases[stepIndex])) {
      clearWestSelection();
    }
    renderStep();
  };

  nextBtn.onclick = function () {
    var phase = phases[stepIndex];
    if (phase.type === "suit" && phase.suitIndex === 3) {
      var seat = seats[phase.seatIndex];
      if (seatCardCount(seat) !== 13) {
        window.alert(
          "Ręka " + seatNames[seat] + " musi mieć dokładnie 13 kart (teraz: " + seatCardCount(seat) + ")."
        );
        return;
      }
    }
    if (phase.type === "summary") {
      finish();
      return;
    }
    stepIndex += 1;
    renderStep();
  };

  renderStep();
})();
