// Starts the entrance animations once nothing is left that could move the text
// while it comes into focus: the web fonts, and on the homepage the song line.
// The inline script in <head> put .hold on <html> before first paint, which
// keeps every entrance at its first frame; taking it off lets them all run.
(function () {
  var root = document.documentElement;
  if (!root.classList.contains("hold")) return;

  // longest the page waits before starting anyway
  var CAP_MS = 800;

  var fonts = false;
  var song = false;

  function release() {
    root.classList.remove("hold");
  }

  function check() {
    if (fonts && song) release();
  }

  function fontsIn() {
    fonts = true;
    check();
  }

  function songIn() {
    song = true;
    check();
  }

  setTimeout(release, CAP_MS);

  // the two families every page is set in, named by the stylesheet's own tokens
  try {
    var tokens = getComputedStyle(root);
    Promise.all(["--sans", "--mono"].map(function (name) {
      return document.fonts.load("1em " + tokens.getPropertyValue(name));
    })).then(fontsIn, fontsIn);
  } catch (e) {
    fontsIn();
  }

  // the song line has settled once spotify.js has swapped its skeleton for a
  // track, or hidden it
  var line = document.getElementById("now-playing");

  function settled() {
    return line.hidden || !line.querySelector(".now-playing-skeleton");
  }

  if (!line || settled()) {
    songIn();
  } else {
    var watcher = new MutationObserver(function () {
      if (!settled()) return;
      watcher.disconnect();
      songIn();
    });
    watcher.observe(line, { childList: true, attributes: true, attributeFilter: ["hidden"] });
  }
})();
