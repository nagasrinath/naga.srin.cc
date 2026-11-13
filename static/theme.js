(function () {
  var btn = document.getElementById("theme-toggle");
  if (!btn) return;

  var root = document.documentElement;
  var BG = { light: "#f0eee9", dark: "#0d0f18" };
  // the other theme's mark and portrait, fetched ahead of a switch
  var art = [];

  function current() {
    var set = root.getAttribute("data-theme");
    if (set === "light" || set === "dark") return set;
    return window.matchMedia("(prefers-color-scheme: dark)").matches ? "dark" : "light";
  }

  function other() {
    return current() === "dark" ? "light" : "dark";
  }

  function label() {
    btn.setAttribute("aria-label", "switch to " + other() + " theme");
  }

  function apply(next) {
    root.setAttribute("data-theme", next);
    try { localStorage.setItem("theme", next); } catch (e) {}
    // the theme-color metas only follow the system preference, so pin them
    var metas = document.querySelectorAll('meta[name="theme-color"]');
    for (var i = 0; i < metas.length; i++) metas[i].setAttribute("content", BG[next]);
    label();
  }

  // the stylesheet names each theme's art -light.png / -dark.png, so what this
  // page would show in the other theme is one word away from what it shows now
  function preload() {
    var from = "-" + current() + ".";
    var to = "-" + other() + ".";
    var seen = {};
    var drawn = document.querySelectorAll(".mark, .home-avatar");
    for (var i = 0; i < drawn.length; i++) {
      var url = /url\(["']?([^"')]+)["']?\)/.exec(getComputedStyle(drawn[i]).backgroundImage);
      if (!url || url[1].indexOf(from) < 0 || seen[url[1]]) continue;
      seen[url[1]] = true;
      var img = new Image();
      img.src = url[1].replace(from, to);
      art.push(img);
    }
  }

  // settles once that art can be painted, or after `ms` regardless
  function ready(ms) {
    var waits = [];
    for (var i = 0; i < art.length; i++) {
      if (art[i].decode) waits.push(art[i].decode().catch(function () {}));
    }
    return Promise.race([
      Promise.all(waits),
      new Promise(function (done) { setTimeout(done, ms); })
    ]);
  }

  btn.addEventListener("click", function () {
    var next = other();
    if (!document.startViewTransition) {
      apply(next);
      return;
    }
    // the new theme opens out in a circle from the toggle (theme-reveal in the stylesheet)
    var box = btn.getBoundingClientRect();
    root.style.setProperty("--tx", box.left + box.width / 2 + "px");
    root.style.setProperty("--ty", box.top + box.height / 2 + "px");
    try {
      document.startViewTransition(function () {
        apply(next);
        return ready(80);
      });
    } catch (e) {
      apply(next);
    }
  });

  label();
  if (document.readyState === "complete") preload();
  else window.addEventListener("load", preload);
})();
