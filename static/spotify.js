(function () {
  var el = document.getElementById("now-playing");
  if (!el) return;

  var ENDPOINT = "https://api.srin.cc/now-playing";
  var POLL_MS = 30000;

  function refresh() {
    fetch(ENDPOINT)
      .then(function (res) { return res.json(); })
      .then(function (data) {
        if (data && data.playing && data.track) {
          el.textContent = "♪ " + data.track + " — " + data.artist;
        } else {
          el.textContent = "";
        }
      })
      .catch(function () {
        el.textContent = "";
      });
  }

  refresh();
  setInterval(refresh, POLL_MS);
})();
