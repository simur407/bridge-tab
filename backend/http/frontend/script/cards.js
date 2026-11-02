(function () {
  var form = document.getElementById("cards-form");
  if (!form) {
    return;
  }

  var seats = ["N", "E", "S", "W"];
  var seatNames = { N: "North", E: "East", S: "South", W: "West" };
  var suits = [
    { code: "S", symbol: "♠", cls: "club" },
    { code: "H", symbol: "♥", cls: "red" },
    { code: "D", symbol: "♦", cls: "red" },
    { code: "C", symbol: "♣", cls: "club" }
  ];
  var ranks = ["A", "K", "Q", "J", "T", "9", "8", "7", "6", "5", "4", "3", "2"];

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

  function toggleRank(suitIndex, rank, checked) {
    var seat = seats[Math.floor(stepIndex / 4)];
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

  function updateCount() {
    var seat = seats[Math.floor(stepIndex / 4)];
    var suitIndex = stepIndex % 4;
    var suitCount = selections[seat][suitIndex].length;
    var total = seatCardCount(seat);
    countEl.innerHTML = "Kolor: " + suitCount + " · ręka: " + total + " / 13";
  }

  function renderStep() {
    var seatIndex = Math.floor(stepIndex / 4);
    var suitIndex = stepIndex % 4;
    var seat = seats[seatIndex];
    var suit = suits[suitIndex];
    var used = usedCardsBefore(seatIndex);
    var selected = {};
    var i;

    if (seat === "W") {
      ensureWestSelection();
    }

    for (i = 0; i < selections[seat][suitIndex].length; i++) {
      selected[selections[seat][suitIndex][i]] = true;
    }

    var seatsHtml = "";
    for (i = 0; i < seats.length; i++) {
      var seatClass = "cards-seat";
      if (i < seatIndex) {
        seatClass += " cards-seat-done";
      } else if (i === seatIndex) {
        seatClass += " cards-seat-current";
      } else {
        seatClass += " cards-seat-todo";
      }
      seatsHtml += '<span class="' + seatClass + '">' + seatNames[seats[i]] + "</span>";
    }
    titleEl.innerHTML =
      '<span class="cards-heading-label">Karty <span class="suit ' + suit.cls + '">' + suit.symbol + "</span></span>" +
      '<span class="cards-heading-seats">' + seatsHtml + "</span>";
    gridEl.innerHTML = "";

    var rows = [ranks.slice(0, 7), ranks.slice(7)];
    var rowIndex;
    for (rowIndex = 0; rowIndex < rows.length; rowIndex++) {
      var row = document.createElement("div");
      row.className = "card-row";
      for (i = 0; i < rows[rowIndex].length; i++) {
        var rank = rows[rowIndex][i];
        var key = cardKey(suitIndex, rank);
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
        label.innerHTML = (rank === "T" ? "10" : rank) + '<span class="suit ' + suit.cls + '">' + suit.symbol + "</span>";

        (function (sIdx, rk, inp) {
          input.onchange = function () {
            toggleRank(sIdx, rk, inp.checked);
          };
        })(suitIndex, rank, input);

        wrap.appendChild(input);
        wrap.appendChild(label);
        row.appendChild(wrap);
      }
      gridEl.appendChild(row);
    }
    setHidden(backBtn, stepIndex === 0);
    nextBtn.innerHTML = stepIndex === 15 ? "Zapisz" : "Dalej";
    updateCount();
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
    if (stepIndex < 12) {
      westReady = false;
      selections.W = [[], [], [], []];
    }
    renderStep();
  };

  nextBtn.onclick = function () {
    var seatIndex = Math.floor(stepIndex / 4);
    var suitIndex = stepIndex % 4;
    var seat = seats[seatIndex];
    if (suitIndex === 3 && seatCardCount(seat) !== 13) {
      window.alert("Ręka " + seatNames[seat] + " musi mieć dokładnie 13 kart (teraz: " + seatCardCount(seat) + ").");
      return;
    }
    if (stepIndex === 15) {
      finish();
      return;
    }
    stepIndex += 1;
    renderStep();
  };

  renderStep();
})();
