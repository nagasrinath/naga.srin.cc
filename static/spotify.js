(function () {
  var bar = document.getElementById("now-playing");
  var inner = document.getElementById("now-playing-inner");
  if (!bar || !inner) return;

  var ENDPOINT = "https://api.srin.cc/now-playing";
  var POLL_MS = 30000;

  function label(text) {
    var span = document.createElement("span");
    span.className = "now-playing-label";
    span.textContent = text;
    return span;
  }

  function render(data) {
    inner.textContent = "";

    if (!data || !data.track) {
      inner.appendChild(label("♪ no signal"));
      return;
    }

    inner.appendChild(label(data.playing ? "♪ now listening: " : "♪ was listening: "));

    var text = data.artist ? data.track + " — " + data.artist : data.track;
    if (data.url) {
      var link = document.createElement("a");
      link.href = data.url;
      link.target = "_blank";
      link.rel = "noopener";
      link.textContent = text;
      inner.appendChild(link);
    } else {
      inner.appendChild(document.createTextNode(text));
    }
  }

  function refresh() {
    fetch(ENDPOINT)
      .then(function (res) { return res.json(); })
      .then(render)
      .catch(function () { render(null); });
  }

  refresh();
  setInterval(refresh, POLL_MS);
})();
