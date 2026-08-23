(function () {
  var bar = document.getElementById("now-playing");
  var inner = document.getElementById("now-playing-inner");
  if (!bar || !inner) return;

  var ENDPOINT = "https://api.srin.cc/now-playing";
  var POLL_MS = 30000;

  function render(data) {
    inner.textContent = "";

    if (!data || !data.playing || !data.track) {
      bar.hidden = true;
      return;
    }

    inner.appendChild(document.createTextNode("now listening: "));

    var label = data.artist ? data.track + " — " + data.artist : data.track;
    if (data.url) {
      var link = document.createElement("a");
      link.href = data.url;
      link.target = "_blank";
      link.rel = "noopener";
      link.textContent = label;
      inner.appendChild(link);
    } else {
      inner.appendChild(document.createTextNode(label));
    }

    bar.hidden = false;
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
