document.addEventListener("keyup", (event) => {
  if (event.key === "/") {
    const element = document.querySelector("input[type=search]");

    if (element && element instanceof HTMLInputElement && document.activeElement !== element) {
      event.preventDefault();

      element.focus();
      element.setSelectionRange(0, element.value.length);
    }
  }
});
