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

  // The markup ships a resolved "no signal" state so it reads sensibly
  // without JS; once we're running, show the skeleton until the first
  // response lands.
  function showSkeleton() {
    inner.textContent = "";
    inner.appendChild(label("♪"));
    inner.appendChild(document.createTextNode(" "));
    var skel = document.createElement("span");
    skel.className = "now-playing-skeleton";
    skel.setAttribute("aria-hidden", "true");
    inner.appendChild(skel);
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
    // No point polling a tab nobody is looking at; visibilitychange below
    // catches it up the moment it comes back.
    if (document.hidden) return;
    fetch(ENDPOINT)
      .then(function (res) {
        if (!res.ok) throw new Error("status " + res.status);
        return res.json();
      })
      .then(render)
      .catch(function () { render(null); });
  }

  showSkeleton();
  refresh();
  setInterval(refresh, POLL_MS);
  document.addEventListener("visibilitychange", function () {
    if (!document.hidden) refresh();
  });
})();
