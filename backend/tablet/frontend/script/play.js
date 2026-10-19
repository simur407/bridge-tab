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

  function setHidden(el, on) {
    if (on) {
      el.setAttribute("hidden", "hidden");
    } else {
      el.removeAttribute("hidden");
    }
  }

  function hidePopup() {
    setHidden(popup, true);
    stopAuto();
    startIdle();
  }

  function showPopup() {
    setHidden(popup, false);
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
    if (!popup.hasAttribute("hidden")) {
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
  document.addEventListener("mousedown", bumpIdle, false);
  document.addEventListener("keydown", bumpIdle, false);
  document.addEventListener("touchstart", bumpIdle, false);

  showPopup();
})();
