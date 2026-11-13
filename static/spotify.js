(function () {
  var line = document.getElementById("now-playing");
  if (!line) return;

  var ENDPOINT = "https://api.srin.cc/now-playing";
  var POLL_MS = 30000;

  function bars() {
    var wrap = document.createElement("span");
    wrap.className = "eq";
    wrap.setAttribute("aria-hidden", "true");
    for (var i = 0; i < 3; i++) wrap.appendChild(document.createElement("span"));
    return wrap;
  }

  // holds the line's place until the first response lands
  function skeleton() {
    var bar = document.createElement("span");
    bar.className = "now-playing-skeleton";
    bar.setAttribute("aria-hidden", "true");
    return bar;
  }

  function text(s) {
    return document.createTextNode(s);
  }

  function plural(n, unit) {
    return n + " " + unit + (n === 1 ? "" : "s") + " ago";
  }

  function ago(iso) {
    var then = Date.parse(iso);
    if (isNaN(then)) return "";
    var mins = Math.floor((Date.now() - then) / 60000);
    if (mins < 1) return "just now";
    if (mins < 60) return plural(mins, "minute");
    var hours = Math.floor(mins / 60);
    if (hours < 24) return plural(hours, "hour");
    return plural(Math.floor(hours / 24), "day");
  }

  function render(data) {
    line.textContent = "";

    if (!data || !data.track) {
      line.hidden = true;
      return;
    }

    if (data.playing) {
      line.appendChild(bars());
      line.appendChild(text("listening to "));
    } else {
      line.appendChild(text("was listening to "));
    }

    if (data.url) {
      var link = document.createElement("a");
      link.href = data.url;
      link.target = "_blank";
      link.rel = "noopener";
      link.textContent = data.track;
      line.appendChild(link);
    } else {
      line.appendChild(text(data.track));
    }

    var tail = (data.artist ? " by " + data.artist : "") + " on spotify";
    if (data.playing) {
      tail += ", right now.";
    } else {
      var when = data.played_at ? ago(data.played_at) : "";
      tail += when ? ", " + when + "." : ".";
    }
    line.appendChild(text(tail));
    line.hidden = false;
  }

  function refresh() {
    fetch(ENDPOINT)
      .then(function (res) { return res.json(); })
      .then(render)
      .catch(function () { render(null); });
  }

  line.appendChild(skeleton());
  line.hidden = false;

  refresh();
  setInterval(refresh, POLL_MS);
})();
