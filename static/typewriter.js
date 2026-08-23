(function () {
  var el = document.getElementById("typed");
  if (!el) return;
  var text = el.value;
  var n = text.length;
  var STEP_MS = 45;
  el.style.width = "0ch";
  var i = 0;
  var timer = setInterval(function () {
    i++;
    el.style.width = i + "ch";
    if (i >= n) clearInterval(timer);
  }, STEP_MS);

  document.addEventListener("keydown", function (e) {
    if (document.activeElement === el) return;
    var active = document.activeElement;
    if (active && (active.tagName === "INPUT" || active.tagName === "TEXTAREA" || active.isContentEditable)) return;
    if (e.ctrlKey || e.metaKey || e.altKey) return;

    if (e.key === "Backspace") {
      e.preventDefault();
      el.value = el.value.slice(0, -1);
      el.focus();
      el.setSelectionRange(el.value.length, el.value.length);
    } else if (e.key.length === 1) {
      e.preventDefault();
      el.value += e.key;
      el.focus();
      el.setSelectionRange(el.value.length, el.value.length);
    }
  });
})();
