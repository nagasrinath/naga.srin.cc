(function () {
  var el = document.getElementById("typed");
  var out = document.getElementById("term-output");
  if (!el || !out) return;

  var history = [];
  var historyIndex = 0;

  function readList(selector) {
    return Array.prototype.map.call(document.querySelectorAll(selector), function (n) {
      return n.textContent;
    });
  }

  var COMMANDS = {
    help: {
      desc: "show this list",
      run: function () {
        return Object.keys(COMMANDS).sort().map(function (name) {
          return name + " - " + COMMANDS[name].desc;
        });
      },
    },
    whoami: { desc: "who you're talking to", run: function () { return "naga"; } },
    pwd: { desc: "current directory", run: function () { return "/world/india/banglore/naga"; } },
    ls: {
      desc: "list files",
      run: function () { return "bio.txt  interests.txt  links.txt"; },
    },
    cat: {
      desc: "print a file (try: cat bio.txt)",
      run: function (args) {
        var file = (args || "").trim();
        if (file === "bio.txt") {
          return readList(".bio p");
        }
        if (file === "interests.txt") {
          return readList(".tags-row .tag").join(", ");
        }
        if (file === "links.txt") {
          return Array.prototype.map.call(document.querySelectorAll("nav.links a"), function (a) {
            return a.textContent + ": " + a.getAttribute("href");
          });
        }
        if (!file) return "usage: cat FILE";
        return "cat: " + file + ": No such file or directory";
      },
    },
    neofetch: {
      desc: "system summary",
      run: function () {
        var uptimeEl = document.getElementById("uptime");
        var location = document.body.getAttribute("data-location") || "";
        var interests = readList(".tags-row .tag").join(", ");
        return [
          "naga@srin.cc",
          "-------------",
          "OS: Nix OS",
          "Uptime: " + (uptimeEl ? uptimeEl.textContent : "?"),
          "Location: " + location,
          "Interests: " + interests,
        ];
      },
    },
    uptime: {
      desc: "how long this site has been up",
      run: function () {
        var uptimeEl = document.getElementById("uptime");
        var n = uptimeEl ? uptimeEl.textContent : "?";
        return "up " + n + " · still self-hosting";
      },
    },
    date: { desc: "current date and time", run: function () { return new Date().toString(); } },
    echo: { desc: "print TEXT back", run: function (args) { return args || ""; } },
    history: {
      desc: "show command history",
      run: function () {
        if (!history.length) return "";
        return history.map(function (cmd, i) { return "  " + (i + 1) + "  " + cmd; });
      },
    },
    man: {
      desc: "one-line description of CMD",
      run: function (args) {
        var name = (args || "").trim().split(/\s+/)[0];
        if (name && COMMANDS[name]) return name + " - " + COMMANDS[name].desc;
        return "No manual entry for " + (name || "");
      },
    },
    uname: {
      desc: "kernel info",
      run: function () { return "Linux srin.cc 6.1.0-naga #1 SMP PREEMPT x86_64 GNU/Linux"; },
    },
    hostname: { desc: "this machine's hostname", run: function () { return "srin.cc"; } },
    sudo: {
      desc: "try it",
      run: function () { return "naga is not in the sudoers file.  This incident will be reported."; },
    },
    clear: { desc: "clear the screen", run: function () { return null; } },
  };

  function appendLine(text, cls) {
    var div = document.createElement("div");
    if (cls) div.className = cls;
    div.textContent = text;
    out.appendChild(div);
  }

  function runCommand(raw) {
    var trimmed = raw.trim();
    if (!trimmed) return;

    var spaceIdx = trimmed.indexOf(" ");
    var name = spaceIdx === -1 ? trimmed : trimmed.slice(0, spaceIdx);
    var args = spaceIdx === -1 ? "" : trimmed.slice(spaceIdx + 1);

    history.push(trimmed);
    historyIndex = history.length;

    if (name === "clear") {
      out.textContent = "";
      return;
    }

    appendLine("$ " + trimmed, "term-cmd");

    var result;
    if (COMMANDS[name]) {
      result = COMMANDS[name].run(args);
    } else {
      result = "bash: " + name + ": command not found";
    }

    var lines = Array.isArray(result) ? result : [result];
    lines.forEach(function (line) {
      if (line !== "" && line != null) appendLine(line, "term-out");
    });

    out.scrollTop = out.scrollHeight;
  }

  el.addEventListener("keydown", function (e) {
    if (e.key === "Enter") {
      e.preventDefault();
      runCommand(el.value);
      el.value = "";
    } else if (e.key === "ArrowUp") {
      if (!history.length) return;
      e.preventDefault();
      historyIndex = Math.max(0, historyIndex - 1);
      el.value = history[historyIndex] || "";
      el.setSelectionRange(el.value.length, el.value.length);
    } else if (e.key === "ArrowDown") {
      if (!history.length) return;
      e.preventDefault();
      historyIndex = Math.min(history.length, historyIndex + 1);
      el.value = history[historyIndex] || "";
      el.setSelectionRange(el.value.length, el.value.length);
    }
  });
})();
