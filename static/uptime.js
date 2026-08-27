(function () {
  var el = document.getElementById("uptime");
  if (!el) return;
  var dob = parseInt(el.getAttribute("data-dob"), 10);
  var days = Math.floor((Date.now() / 1000 - dob) / 86400);
  el.textContent = "~" + days + "d";
})();
