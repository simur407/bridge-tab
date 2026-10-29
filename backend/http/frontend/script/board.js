(function () {
  var form = document.getElementById("board-form");
  if (!form) {
    return;
  }
  var steps = ["step-lead", "step-contract", "step-double", "step-declarer", "step-tricks"];
  var ranks = ["A", "K", "Q", "J", "10", "9", "8", "7", "6", "5", "4", "3", "2"];
  var leadSuits = [
    { code: "S", symbol: "♠", cls: "club" },
    { code: "H", symbol: "♥", cls: "red" },
    { code: "D", symbol: "♦", cls: "red" },
    { code: "C", symbol: "♣", cls: "club" }
  ];
  var contractSuits = [
    { code: "NT", symbol: "NT", cls: "" },
    { code: "S", symbol: "♠", cls: "club" },
    { code: "H", symbol: "♥", cls: "red" },
    { code: "D", symbol: "♦", cls: "red" },
    { code: "C", symbol: "♣", cls: "club" }
  ];

  function setHidden(el, on) {
    if (on) {
      el.setAttribute("hidden", "hidden");
    } else {
      el.removeAttribute("hidden");
    }
  }

  function suitSpan(code) {
    var i;
    for (i = 0; i < contractSuits.length; i++) {
      if (contractSuits[i].code === code) {
        if (!contractSuits[i].cls) {
          return contractSuits[i].symbol;
        }
        return '<span class="suit ' + contractSuits[i].cls + '">' + contractSuits[i].symbol + "</span>";
      }
    }
    return "";
  }

  function enteredItem(label, html) {
    return '<span class="entered-item"><span class="entered-label">' + label + "</span>" + html + "</span>";
  }

  function renderEntered(stepId) {
    var el = document.getElementById("entered");
    var parts = [];
    var rank = document.getElementById("openingLeadFigure").value;
    var leadSuit = document.getElementById("openingLeadSuit").value;
    var level = document.getElementById("contractLevel").value;
    var contractSuit = document.getElementById("contractSuit").value;
    var modifier = document.getElementById("contractModifier").value;
    var declarer = document.getElementById("declarer").value;
    if (stepId !== "step-lead" && rank && leadSuit) {
      parts.push(enteredItem("Wist", rank + suitSpan(leadSuit)));
    }
    if ((stepId === "step-double" || stepId === "step-declarer" || stepId === "step-tricks") && level && contractSuit) {
      if (stepId === "step-double") {
        modifier = "";
      }
      parts.push(enteredItem("Kontrakt", level + suitSpan(contractSuit) + modifier));
    }
    if (stepId === "step-tricks" && declarer) {
      parts.push(enteredItem("Rozgrywa", declarer));
    }
    if (!parts.length) {
      el.innerHTML = "";
      setHidden(el, true);
      return;
    }
    el.innerHTML = parts.join("");
    setHidden(el, false);
  }

  function showStep(id) {
    steps.forEach(function (stepId) {
      setHidden(document.getElementById(stepId), stepId !== id);
    });
    renderEntered(id);
  }

  function selected(name) {
    var el = form.querySelector('input[name="' + name + '"]:checked');
    return el ? el.value : "";
  }

  function radioRow(container, name, items, options) {
    var opts = options || {};
    var row = document.createElement("div");
    var n;
    row.className = opts.rowClass || "combo-row";
    items.forEach(function (item) {
      var id = name + "-" + item.id;
      var input = document.createElement("input");
      input.type = "radio";
      input.name = name;
      input.id = id;
      input.value = item.value;
      var label = document.createElement("label");
      label.htmlFor = id;
      label.innerHTML = item.label;
      if (item.cls) {
        label.classList.add(item.cls);
      }
      row.appendChild(input);
      row.appendChild(label);
    });
    if (opts.pad && items.length < opts.pad) {
      for (n = items.length; n < opts.pad; n++) {
        var spacer = document.createElement("span");
        spacer.className = "combo-spacer";
        row.appendChild(spacer);
      }
    }
    container.appendChild(row);
  }

  function comboLabel(container, text) {
    var label = document.createElement("p");
    label.className = "combo-label";
    label.appendChild(document.createTextNode(text));
    container.appendChild(label);
  }

  var leadGrid = document.getElementById("lead-grid");
  comboLabel(leadGrid, "Karta");
  [ranks.slice(0, 7), ranks.slice(7)].forEach(function (row) {
    radioRow(leadGrid, "leadRank", row.map(function (rank) {
      return { id: rank, value: rank, label: rank };
    }), { pad: 7 });
  });
  comboLabel(leadGrid, "Kolor");
  radioRow(leadGrid, "leadSuit", leadSuits.map(function (suit) {
    return {
      id: suit.code,
      value: suit.code,
      label: '<span class="suit ' + suit.cls + '">' + suit.symbol + "</span>",
      cls: suit.cls
    };
  }), { rowClass: "combo-row choice-row" });

  var contractGrid = document.getElementById("contract-grid");
  contractSuits.forEach(function (suit) {
    radioRow(contractGrid, "contractCall", [1, 2, 3, 4, 5, 6, 7].map(function (level) {
      var symbol = suit.cls
        ? '<span class="suit ' + suit.cls + '">' + suit.symbol + "</span>"
        : suit.symbol;
      return {
        id: level + suit.code,
        value: level + "|" + suit.code,
        label: level + symbol,
        cls: suit.cls
      };
    }));
  });

  var tricksGrid = document.getElementById("tricks-grid");
  [[0, 1, 2, 3, 4, 5, 6], [7, 8, 9, 10, 11, 12, 13]].forEach(function (row) {
    radioRow(tricksGrid, "tricksChoice", row.map(function (n) {
      return { id: String(n), value: String(n), label: String(n) };
    }));
  });

  function fillFromLead() {
    var rank = selected("leadRank");
    var suit = selected("leadSuit");
    if (!rank || !suit) {
      return false;
    }
    document.getElementById("openingLeadFigure").value = rank;
    document.getElementById("openingLeadSuit").value = suit;
    return true;
  }

  function fillFromContract() {
    var value = selected("contractCall");
    if (!value) {
      return false;
    }
    var parts = value.split("|");
    document.getElementById("contractLevel").value = parts[0];
    document.getElementById("contractSuit").value = parts[1];
    return true;
  }

  function fillFromDouble() {
    var choice = form.querySelector("input[name='doubleChoice']:checked");
    if (!choice) {
      return false;
    }
    document.getElementById("contractModifier").value = choice.value;
    return true;
  }

  function fillFromDeclarer() {
    var value = selected("declarerChoice");
    if (!value) {
      return false;
    }
    document.getElementById("declarer").value = value;
    return true;
  }

  function fillFromTricks() {
    var value = selected("tricksChoice");
    if (value === "") {
      return false;
    }
    document.getElementById("tricks").value = value;
    return true;
  }

  function submitVerify() {
    if (!window.htmx) {
      return;
    }
    htmx.ajax("POST", form.getAttribute("data-verify-url"), {
      source: form,
      target: "#dialog",
      swap: "innerHTML"
    });
  }

  document.getElementById("lead-next").addEventListener("click", function () {
    if (!fillFromLead()) {
      alert("Wybierz kartę i kolor");
      return;
    }
    document.getElementById("contractModifier").value = "";
    showStep("step-contract");
  });

  document.getElementById("contract-back").addEventListener("click", function () {
    showStep("step-lead");
  });
  document.getElementById("contract-next").addEventListener("click", function () {
    if (!fillFromContract()) {
      alert("Wybierz kontrakt");
      return;
    }
    showStep("step-double");
  });

  document.getElementById("double-back").addEventListener("click", function () {
    showStep("step-contract");
  });
  document.getElementById("double-next").addEventListener("click", function () {
    if (!fillFromDouble()) {
      alert("Wybierz kontrę");
      return;
    }
    showStep("step-declarer");
  });

  document.getElementById("declarer-back").addEventListener("click", function () {
    showStep("step-double");
  });
  document.getElementById("declarer-next").addEventListener("click", function () {
    if (!fillFromDeclarer()) {
      alert("Wybierz rozgrywającego");
      return;
    }
    showStep("step-tricks");
  });

  document.getElementById("tricks-back").addEventListener("click", function () {
    showStep("step-declarer");
  });
  document.getElementById("tricks-next").addEventListener("click", function () {
    if (!fillFromLead() || !fillFromContract() || !fillFromDouble() || !fillFromDeclarer() || !fillFromTricks()) {
      alert("Uzupełnij wszystkie pola");
      return;
    }
    submitVerify();
  });

  document.getElementById("pass-button").addEventListener("click", function () {
    document.getElementById("openingLeadFigure").value = "";
    document.getElementById("openingLeadSuit").value = "";
    document.getElementById("contractLevel").value = "";
    document.getElementById("contractSuit").value = "";
    document.getElementById("contractModifier").value = "Pass";
    document.getElementById("declarer").value = "";
    document.getElementById("tricks").value = "0";
    submitVerify();
  });

  form.addEventListener("submit", function (ev) {
    ev.preventDefault();
  });
})();
