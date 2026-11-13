(function () {
  var el = document.getElementById("ist");
  if (!el) return;

  function pad(n) { return n < 10 ? "0" + n : "" + n; }

  function tick() {
    // IST is UTC+5:30, no daylight saving
    var ist = new Date(Date.now() + 19800000);
    el.textContent = pad(ist.getUTCHours()) + ":" + pad(ist.getUTCMinutes());
  }

  tick();
  setInterval(tick, 20000);
})();
