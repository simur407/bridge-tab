(function () {
  var popup = document.getElementById("next-round-popup");
  if (!popup) {
    return;
  }
  var stay = document.getElementById("next-round-stay");
  var go = document.getElementById("next-round-go");
  var form = document.getElementById("advance-form");
  var timerLabel = document.getElementById("next-round-timer");
  var inactivityMs = 30000;
  var autoSeconds = 20;
  var idleTimer;
  var autoTimer;
  var remaining;

  function advance() {
    if (form) {
      form.submit();
    }
  }

  function updateTimer() {
    if (timerLabel) {
      timerLabel.textContent = String(remaining);
    }
  }

  function stopAuto() {
    if (autoTimer) {
      clearInterval(autoTimer);
      autoTimer = null;
    }
  }

  function hidePopup() {
    popup.hidden = true;
    stopAuto();
    startIdle();
  }

  function showPopup() {
    popup.hidden = false;
    if (idleTimer) {
      clearTimeout(idleTimer);
      idleTimer = null;
    }
    remaining = autoSeconds;
    updateTimer();
    stopAuto();
    autoTimer = setInterval(function () {
      remaining -= 1;
      updateTimer();
      if (remaining <= 0) {
        stopAuto();
        advance();
      }
    }, 1000);
  }

  function startIdle() {
    if (idleTimer) {
      clearTimeout(idleTimer);
    }
    idleTimer = setTimeout(showPopup, inactivityMs);
  }

  function bumpIdle() {
    if (!popup.hidden) {
      return;
    }
    startIdle();
  }

  if (stay) {
    stay.addEventListener("click", hidePopup);
  }
  if (go) {
    go.addEventListener("click", function (event) {
      event.preventDefault();
      advance();
    });
  }
  ["pointerdown", "keydown", "touchstart"].forEach(function (name) {
    document.addEventListener(name, bumpIdle, { passive: true });
  });

  showPopup();
})();
