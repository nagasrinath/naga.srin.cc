(function () {
  var caret = document.querySelector(".caret");
  if (!caret) return;
  var visible = true;
  setInterval(function () {
    visible = !visible;
    caret.style.visibility = visible ? "visible" : "hidden";
  }, 560);
})();
