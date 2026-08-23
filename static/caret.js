(function () {
  var caret = document.querySelector(".caret");
  if (!caret) return;
  var typed = document.getElementById("typed");
  var editing = false;
  var visible = true;
  setInterval(function () {
    if (editing) return;
    visible = !visible;
    caret.style.visibility = visible ? "visible" : "hidden";
  }, 560);
  if (typed) {
    typed.addEventListener("focus", function () {
      editing = true;
      caret.style.visibility = "hidden";
    });
    typed.addEventListener("blur", function () {
      editing = false;
      visible = true;
      caret.style.visibility = "visible";
    });
  }
})();
